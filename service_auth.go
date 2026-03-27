package main

import (
	"context"
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-webauthn/webauthn/protocol"
	"github.com/go-webauthn/webauthn/webauthn"
)

type AuthService struct {
	webauthn *webauthn.WebAuthn
	db       *sql.DB
}

func NewAuthService(db *sql.DB) *AuthService {
	wauthn, err := webauthn.New(&webauthn.Config{
		RPDisplayName: "My Cloud Storage",
		RPID:          "localhost",                       // Adjust when deploying
		RPOrigins:     []string{"http://localhost:8080"}, // Adjust when deploying
	})
	if err != nil {
		log.Fatalf("Failed to initialize WebAuthn: %v", err)
	}
	return &AuthService{webauthn: wauthn, db: db}
}

func (a *AuthService) BeginRegistration(user *User) (*protocol.CredentialCreation, *webauthn.SessionData, error) {
	options, sessionData, err := a.webauthn.BeginRegistration(user)
	if err != nil {
		return nil, nil, err
	}
	return options, sessionData, nil
}

func (a *AuthService) FinishRegistration(user *User, sessionData *webauthn.SessionData, r *http.Request) error {
	credential, err := a.webauthn.FinishRegistration(user, *sessionData, r)
	if err != nil {
		return err
	}
	return savePasskey(a.db, user.ID, credential)
}

func (a *AuthService) BeginLoginWithoutUser() (*protocol.CredentialAssertion, *webauthn.SessionData, error) {
	options, sessionData, err := a.webauthn.BeginDiscoverableLogin()
	if err != nil {
		return nil, nil, err
	}
	return options, sessionData, nil
}

func (a *AuthService) FinishLoginWithoutUser(r *http.Request, sessionData *webauthn.SessionData) (*User, *webauthn.Credential, error) {
	parsedResponse, err := protocol.ParseCredentialRequestResponseBody(r.Body)
	if err != nil {
		return nil, nil, err
	}

	user, credential, err := a.webauthn.ValidatePasskeyLogin(
		func(rawID []byte, userHandle []byte) (webauthn.User, error) {
			return getUserByCredentialID(a.db, rawID)
		},
		*sessionData,
		parsedResponse,
	)
	if err != nil {
		return nil, nil, err
	}

	appUser, ok := user.(*User)
	if !ok {
		return nil, nil, fmt.Errorf("unexpected user type returned from webauthn")
	}
	return appUser, credential, err
}

func exchangeCodeForToken(ctx context.Context, code string) (string, error) {
	reqBody := strings.NewReader(fmt.Sprintf("client_id=%s&client_secret=%s&code=%s", GithubClientID, GithubSecret, code))

	req, err := http.NewRequestWithContext(ctx, "POST", "https://github.com/login/oauth/access_token", reqBody)
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer func() {
		err := resp.Body.Close()
		if err != nil {
			log.Printf("Failed to close response body: %v", err)
		}
	}()

	var result struct {
		AccessToken string `json:"access_token"`
		TokenType   string `json:"token_type"`
		Scope       string `json:"scope"`
		Error       string `json:"error"`
		ErrorDesc   string `json:"error_description"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	if result.Error != "" {
		return "", fmt.Errorf("GitHub OAuth error: %s - %s", result.Error, result.ErrorDesc)
	}

	return result.AccessToken, nil
}

func getGithubUserID(ctx context.Context, accessToken string) (string, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", "https://api.github.com/user", nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Accept", "application/vnd.github+json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var user struct {
		ID int64 `json:"id"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		return "", err
	}
	return strconv.FormatInt(user.ID, 10), nil
}

func generateToken() (string, error) {
	bytes := make([]byte, 16)

	if _, err := rand.Read(bytes); err != nil {
		return "", fmt.Errorf("failed to generate random token: %w", err)
	}
	return hex.EncodeToString(bytes), nil
}
