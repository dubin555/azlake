package auth

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/dubin555/azlake/pkg/auth/crypt"
	"github.com/dubin555/azlake/pkg/auth/keys"
	"github.com/dubin555/azlake/pkg/auth/model"
	"github.com/dubin555/azlake/pkg/kv"
	"github.com/dubin555/azlake/pkg/logging"
	"google.golang.org/protobuf/types/known/timestamppb"
)

const InvalidUserID = ""

// Service is the simplified auth service interface for API key authentication only.
type Service interface {
	CreateUser(ctx context.Context, user *model.User) (string, error)
	GetUser(ctx context.Context, username string) (*model.User, error)
	ListUsers(ctx context.Context, params *model.PaginationParams) ([]*model.User, *model.Paginator, error)
	DeleteUser(ctx context.Context, username string) error

	CreateCredentials(ctx context.Context, username string) (*model.Credential, error)
	AddCredentials(ctx context.Context, username, accessKeyID, secretAccessKey string) (*model.Credential, error)
	GetCredentials(ctx context.Context, accessKeyID string) (*model.Credential, error)
	GetCredentialsForUser(ctx context.Context, username, accessKeyID string) (*model.Credential, error)
	DeleteCredentials(ctx context.Context, username, accessKeyID string) error

	SecretStore() crypt.SecretStore
}

// KVAuthService implements Service backed by the KV store.
type KVAuthService struct {
	store       kv.Store
	secretStore crypt.SecretStore
	log         logging.Logger
}

func NewKVAuthService(store kv.Store, secretStore crypt.SecretStore, logger logging.Logger) *KVAuthService {
	logger.Info("initialized Auth service")
	return &KVAuthService{
		store:       store,
		secretStore: secretStore,
		log:         logger,
	}
}

func (s *KVAuthService) SecretStore() crypt.SecretStore {
	return s.secretStore
}

func (s *KVAuthService) CreateUser(ctx context.Context, user *model.User) (string, error) {
	if err := model.ValidateAuthEntityID(user.Username); err != nil {
		return InvalidUserID, err
	}
	userKey := model.UserPath(user.Username)
	err := kv.SetMsgIf(ctx, s.store, model.PartitionKey, userKey, model.ProtoFromUser(user), nil)
	if err != nil {
		if errors.Is(err, kv.ErrPredicateFailed) {
			return "", ErrAlreadyExists
		}
		return "", fmt.Errorf("create user: %w", err)
	}
	return user.Username, nil
}

func (s *KVAuthService) GetUser(ctx context.Context, username string) (*model.User, error) {
	userKey := model.UserPath(username)
	m := &model.UserData{}
	_, err := kv.GetMsg(ctx, s.store, model.PartitionKey, userKey, m)
	if err != nil {
		if errors.Is(err, kv.ErrNotFound) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("get user: %w", err)
	}
	return model.UserFromProto(m), nil
}

func (s *KVAuthService) ListUsers(ctx context.Context, params *model.PaginationParams) ([]*model.User, *model.Paginator, error) {
	var (
		it  *kv.PrimaryIterator
		err error
	)
	userKey := model.UserPath("")
	if params != nil && params.After != "" {
		it, err = kv.NewPrimaryIterator(ctx, s.store, (&model.UserData{}).ProtoReflect().Type(),
			model.PartitionKey, userKey, kv.IteratorOptionsAfter([]byte(params.After)))
	} else {
		it, err = kv.NewPrimaryIterator(ctx, s.store, (&model.UserData{}).ProtoReflect().Type(),
			model.PartitionKey, userKey, kv.IteratorOptionsAfter([]byte("")))
	}
	if err != nil {
		return nil, nil, fmt.Errorf("list users: %w", err)
	}
	defer it.Close()

	amount := 100
	if params != nil && params.Amount > 0 {
		amount = params.Amount
	}

	users := make([]*model.User, 0)
	for it.Next() {
		entry := it.Entry()
		user := model.UserFromProto(entry.Value.(*model.UserData))
		users = append(users, user)
		if len(users) >= amount {
			break
		}
	}
	if err := it.Err(); err != nil {
		return nil, nil, fmt.Errorf("list users: %w", err)
	}

	paginator := &model.Paginator{Amount: len(users)}
	if len(users) > 0 && it.Next() {
		paginator.NextPageToken = users[len(users)-1].Username
	}
	return users, paginator, nil
}

func (s *KVAuthService) DeleteUser(ctx context.Context, username string) error {
	userKey := model.UserPath(username)

	// Delete all credentials for this user first
	credKey := model.CredentialPath(username, "")
	credIt, err := kv.NewPrimaryIterator(ctx, s.store, (&model.CredentialData{}).ProtoReflect().Type(),
		model.PartitionKey, credKey, kv.IteratorOptionsAfter([]byte("")))
	if err != nil {
		return fmt.Errorf("list user credentials for delete: %w", err)
	}
	defer credIt.Close()
	for credIt.Next() {
		if err := s.store.Delete(ctx, []byte(model.PartitionKey), credIt.Entry().Key); err != nil {
			return fmt.Errorf("delete credential: %w", err)
		}
	}

	if err := s.store.Delete(ctx, []byte(model.PartitionKey), userKey); err != nil {
		if errors.Is(err, kv.ErrNotFound) {
			return ErrNotFound
		}
		return fmt.Errorf("delete user: %w", err)
	}
	return nil
}

func (s *KVAuthService) CreateCredentials(ctx context.Context, username string) (*model.Credential, error) {
	accessKeyID := keys.GenAccessKeyID()
	secretAccessKey := keys.GenSecretAccessKey()
	return s.addCredentials(ctx, username, accessKeyID, secretAccessKey)
}

func (s *KVAuthService) AddCredentials(ctx context.Context, username, accessKeyID, secretAccessKey string) (*model.Credential, error) {
	return s.addCredentials(ctx, username, accessKeyID, secretAccessKey)
}

func (s *KVAuthService) addCredentials(ctx context.Context, username, accessKeyID, secretAccessKey string) (*model.Credential, error) {
	// Verify user exists
	if _, err := s.GetUser(ctx, username); err != nil {
		return nil, fmt.Errorf("get user for credentials: %w", err)
	}

	encryptedKey, err := s.secretStore.Encrypt([]byte(secretAccessKey))
	if err != nil {
		return nil, fmt.Errorf("encrypt secret key: %w", err)
	}

	now := time.Now()
	c := &model.Credential{
		Username: username,
		BaseCredential: model.BaseCredential{
			AccessKeyID:                   accessKeyID,
			SecretAccessKey:               secretAccessKey,
			SecretAccessKeyEncryptedBytes: encryptedKey,
			IssuedDate:                    now,
		},
	}

	credKey := model.CredentialPath(username, accessKeyID)
	credData := &model.CredentialData{
		AccessKeyId:                   accessKeyID,
		SecretAccessKeyEncryptedBytes: encryptedKey,
		IssuedDate:                    timestamppb.New(now),
		UserId:                        []byte(username),
	}

	err = kv.SetMsgIf(ctx, s.store, model.PartitionKey, credKey, credData, nil)
	if err != nil {
		if errors.Is(err, kv.ErrPredicateFailed) {
			return nil, ErrAlreadyExists
		}
		return nil, fmt.Errorf("save credentials: %w", err)
	}
	return c, nil
}

func (s *KVAuthService) GetCredentials(ctx context.Context, accessKeyID string) (*model.Credential, error) {
	// Scan all users' credentials to find by access key ID
	credPrefix := model.CredentialPath("", "")
	it, err := kv.NewPrimaryIterator(ctx, s.store, (&model.CredentialData{}).ProtoReflect().Type(),
		model.PartitionKey, credPrefix, kv.IteratorOptionsAfter([]byte("")))
	if err != nil {
		return nil, fmt.Errorf("get credentials: %w", err)
	}
	defer it.Close()

	for it.Next() {
		entry := it.Entry()
		credData := entry.Value.(*model.CredentialData)
		if credData.AccessKeyId == accessKeyID {
			return model.CredentialFromProto(s.secretStore, credData)
		}
	}
	if err := it.Err(); err != nil {
		return nil, fmt.Errorf("get credentials: %w", err)
	}
	return nil, ErrNotFound
}

func (s *KVAuthService) GetCredentialsForUser(ctx context.Context, username, accessKeyID string) (*model.Credential, error) {
	credKey := model.CredentialPath(username, accessKeyID)
	m := &model.CredentialData{}
	_, err := kv.GetMsg(ctx, s.store, model.PartitionKey, credKey, m)
	if err != nil {
		if errors.Is(err, kv.ErrNotFound) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("get credentials for user: %w", err)
	}
	return model.CredentialFromProto(s.secretStore, m)
}

func (s *KVAuthService) DeleteCredentials(ctx context.Context, username, accessKeyID string) error {
	credKey := model.CredentialPath(username, accessKeyID)
	if err := s.store.Delete(ctx, []byte(model.PartitionKey), credKey); err != nil {
		if errors.Is(err, kv.ErrNotFound) {
			return ErrNotFound
		}
		return fmt.Errorf("delete credentials: %w", err)
	}
	return nil
}
