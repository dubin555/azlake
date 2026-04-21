package setup

import (
	"context"
	"fmt"
	"time"

	"github.com/dubin555/azlake/pkg/auth"
	"github.com/dubin555/azlake/pkg/auth/model"
)

// CreateAdminUser creates the initial admin user and returns their credentials.
func CreateAdminUser(ctx context.Context, authService auth.Service, superuser *model.SuperuserConfiguration) (*model.Credential, error) {
	return AddAdminUser(ctx, authService, superuser)
}

// AddAdminUser creates a user and generates credentials for them.
func AddAdminUser(ctx context.Context, authService auth.Service, user *model.SuperuserConfiguration) (*model.Credential, error) {
	user.Source = "internal"
	_, err := authService.CreateUser(ctx, &user.User)
	if err != nil {
		return nil, fmt.Errorf("create user: %w", err)
	}

	var creds *model.Credential
	if user.AccessKeyID == "" {
		creds, err = authService.CreateCredentials(ctx, user.Username)
	} else {
		creds, err = authService.AddCredentials(ctx, user.Username, user.AccessKeyID, user.SecretAccessKey)
	}
	if err != nil {
		// Clean up user on credential failure
		_ = authService.DeleteUser(ctx, user.Username)
		return nil, fmt.Errorf("create credentials for %s: %w", user.Username, err)
	}
	return creds, nil
}

// CreateInitialAdminUser creates the first admin user with auto-generated credentials.
func CreateInitialAdminUser(ctx context.Context, authService auth.Service, username string) (*model.Credential, error) {
	return CreateInitialAdminUserWithKeys(ctx, authService, username, nil, nil)
}

// CreateInitialAdminUserWithKeys creates the first admin user, optionally with specific keys.
func CreateInitialAdminUserWithKeys(ctx context.Context, authService auth.Service, username string, accessKeyID, secretAccessKey *string) (*model.Credential, error) {
	adminUser := &model.SuperuserConfiguration{
		User: model.User{
			CreatedAt: time.Now(),
			Username:  username,
		},
	}
	if accessKeyID != nil && secretAccessKey != nil {
		adminUser.AccessKeyID = *accessKeyID
		adminUser.SecretAccessKey = *secretAccessKey
	}
	return CreateAdminUser(ctx, authService, adminUser)
}
