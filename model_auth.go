package main

type Passkey struct {
	ID           int64  `json:"id"`
	UserID       int64  `json:"user_id"`
	CredentialID string `json:"credential_id"`
	PublicKey    string `json:"public_key"`
	SignCount    int64  `json:"sign_count"`
	DeviceName   string `json:"device_name"`
	Backedup     bool   `json:"backedup"`
	LastUsedAt   int64  `json:"last_used_at"`
}

type OAuthAccound struct {
	ID             int64  `json:"id"`
	UserID         int64  `json:"user_id"`
	Provider       string `json:"provider"`
	ProviderUserID string `json:"provider_user_id"`
	AccessToken    string `json:"access_token"`
	TokenExpiresAt int64  `json:"token_expires_at"`
}

type Session struct {
	ID         int64  `json:"id"`
	UserID     int64  `json:"user_id"`
	AuthMethod string `json:"auth_method"`
	ExpiresAt  int64  `json:"expires_at"`
}
