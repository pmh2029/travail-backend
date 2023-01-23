package config

import "golang.org/x/oauth2"

type OAuthConfig struct {
	GoogleLoginConfig oauth2.Config
}

type GoogleUserInfo struct {
	ID    string `json:"id"`
	Email string `json:"email"`
	Name  string `json:"name"`
}
