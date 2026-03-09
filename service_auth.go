package main

type AuthService struct{}

func NewAuthService() *AuthService {
	return &AuthService{}
}

func (a *AuthService) Register(email, password string) (User, error) {
	return User{}, nil
}
