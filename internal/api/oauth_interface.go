package api

import (
	"net/url"
)

type OAuthClient interface {
	MakeAccessTokenForm(code string) url.Values
	Login(code string) (OAuthUserInfo, error)
	getAccessToken(code string) (OAuthTokenResponse, error)
	getUserInfo(accessToken string) (OAuthUserInfo, error)
}

type OAuthTokenResponse struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	RefreshToken string `json:"refresh_token,omitempty"`
	ExpiresIn    int    `json:"expires_in,omitempty"`
	Scope        string `json:"scope,omitempty"`
}

type OAuthUserInfo struct {
	ID    string  `json:"id"`
	Email *string `json:"email,omitempty"`
}
