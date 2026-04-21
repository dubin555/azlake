package model

import "time"

const (
	StatementEffectAllow = "allow"
	StatementEffectDeny  = "deny"
)

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
