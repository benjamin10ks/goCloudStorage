package main

import (
	"strconv"

	"github.com/go-webauthn/webauthn/webauthn"
)

type User struct {
	ID          int64  `json:"id"`
	Username    string `json:"username"`
	DisplayName string `json:"display_name"`
}

func (u *User) WebAuthnID() []byte { return []byte(strconv.FormatInt(u.ID, 10)) }

func (u *User) WebAuthnName() string { return u.Username }

func (u *User) WebAuthnDisplayName() string { return u.DisplayName }

// Placeholder... Later get passkey with query
func (u *User) WebAuthnCredentials() []webauthn.Credential { return nil }
