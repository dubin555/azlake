package model

import (
	"time"

	"github.com/dubin555/azlake/pkg/auth/crypt"
	"github.com/dubin555/azlake/pkg/kv"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

const (
	StatementEffectAllow = "allow"
	StatementEffectDeny  = "deny"

	PartitionKey           = "auth"
	usersPrefix            = "users"
	usersCredentialsPrefix = "uCredentials" // #nosec G101 -- False positive: this is only a kv key prefix
	credentialsPrefix      = "credentials"
)

//nolint:gochecknoinits
func init() {
	kv.MustRegisterType("auth", "users", (&UserData{}).ProtoReflect().Type())
	kv.MustRegisterType("auth", kv.FormatPath("uCredentials", "*", "credentials"), (&CredentialData{}).ProtoReflect().Type())
}

// UserPath returns the KV key for a user.
func UserPath(userName string) []byte {
	return []byte(kv.FormatPath(usersPrefix, userName))
}

// CredentialPath returns the KV key for a credential.
func CredentialPath(userName string, accessKeyID string) []byte {
	return []byte(kv.FormatPath(usersCredentialsPrefix, userName, credentialsPrefix, accessKeyID))
}

type PaginationParams struct {
	Prefix string
	After  string
	Amount int
}

type Paginator struct {
	Amount        int
	NextPageToken string
}

type SuperuserConfiguration struct {
	User
	AccessKeyID    string
	SecretAccessKey string
}

type User struct {
	CreatedAt         time.Time `db:"created_at"`
	Username          string    `db:"display_name" json:"display_name"`
	FriendlyName      *string   `db:"friendly_name" json:"friendly_name"`
	Email             *string   `db:"email" json:"email"`
	EncryptedPassword []byte    `db:"encrypted_password" json:"encrypted_password"`
	Source            string    `db:"source" json:"source"`
	ExternalID        *string   `db:"external_id" json:"external_id"`
}

func (u *User) Committer() string {
	if u.Email != nil && *u.Email != "" {
		return *u.Email
	}
	return u.Username
}

type BaseCredential struct {
	AccessKeyID                   string    `db:"access_key_id"`
	SecretAccessKey               string    `db:"-" json:"-"`
	SecretAccessKeyEncryptedBytes []byte    `db:"secret_access_key" json:"-"`
	IssuedDate                    time.Time `db:"issued_date"`
}

type Credential struct {
	Username string
	BaseCredential
}

type CredentialKeys struct {
	AccessKeyID     string `json:"access_key_id"`
	SecretAccessKey string `json:"secret_access_key"`
}

type Policy struct {
	CreatedAt   time.Time  `db:"created_at"`
	DisplayName string     `db:"display_name" json:"display_name"`
	Statement   Statements `db:"statement"`
}

type Statement struct {
	Effect    string                         `json:"Effect"`
	Action    []string                       `json:"Action"`
	Resource  string                         `json:"Resource"`
	Condition map[string]map[string][]string `json:"Condition,omitempty"`
}

type Statements []Statement

// Proto conversion helpers

func UserFromProto(pb *UserData) *User {
	return &User{
		CreatedAt:         pb.CreatedAt.AsTime(),
		Username:          pb.Username,
		FriendlyName:      &pb.FriendlyName,
		Email:             &pb.Email,
		EncryptedPassword: pb.EncryptedPassword,
		Source:            pb.Source,
		ExternalID:        &pb.ExternalId,
	}
}

func ProtoFromUser(u *User) *UserData {
	fn := ""
	if u.FriendlyName != nil {
		fn = *u.FriendlyName
	}
	email := ""
	if u.Email != nil {
		email = *u.Email
	}
	extID := ""
	if u.ExternalID != nil {
		extID = *u.ExternalID
	}
	return &UserData{
		CreatedAt:         timestamppb.New(u.CreatedAt),
		Username:          u.Username,
		FriendlyName:      fn,
		Email:             email,
		EncryptedPassword: u.EncryptedPassword,
		Source:            u.Source,
		ExternalId:        extID,
	}
}

func CredentialFromProto(s crypt.SecretStore, pb *CredentialData) (*Credential, error) {
	secret, err := DecryptSecret(s, pb.SecretAccessKeyEncryptedBytes)
	if err != nil {
		return nil, err
	}
	return &Credential{
		Username: string(pb.UserId),
		BaseCredential: BaseCredential{
			AccessKeyID:                   pb.AccessKeyId,
			SecretAccessKey:               secret,
			SecretAccessKeyEncryptedBytes: pb.SecretAccessKeyEncryptedBytes,
			IssuedDate:                    pb.IssuedDate.AsTime(),
		},
	}, nil
}

func ProtoFromCredential(c *Credential) *CredentialData {
	return &CredentialData{
		AccessKeyId:                   c.AccessKeyID,
		SecretAccessKeyEncryptedBytes: c.SecretAccessKeyEncryptedBytes,
		IssuedDate:                    timestamppb.New(c.IssuedDate),
		UserId:                        []byte(c.Username),
	}
}

func ConvertUsersDataList(users []proto.Message) []*User {
	result := make([]*User, 0, len(users))
	for _, u := range users {
		result = append(result, UserFromProto(u.(*UserData)))
	}
	return result
}

func DecryptSecret(s crypt.SecretStore, value []byte) (string, error) {
	decrypted, err := s.Decrypt(value)
	if err != nil {
		return "", err
	}
	return string(decrypted), nil
}

func EncryptSecret(s crypt.SecretStore, secretAccessKey string) ([]byte, error) {
	return s.Encrypt([]byte(secretAccessKey))
}
