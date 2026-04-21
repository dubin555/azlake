package acl

import (
	"context"
	"time"

	"github.com/dubin555/azlake/pkg/auth"
)

func PolicyName(groupID string) string {
	return "ACL(" + groupID + ")"
}

func WriteGroupACL(ctx context.Context, svc auth.Service, groupID string, acl ACL, createdAt time.Time, overwrite bool) error {
	return nil // stub
}

type ACL struct {
	Permission string
}
