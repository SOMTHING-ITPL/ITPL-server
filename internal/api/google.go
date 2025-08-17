package api

import (
	"encoding/json"
	"log"
	"net/http"
	"net/url"

	"github.com/SOMTHING-ITPL/ITPL-server/config"
)

type GoogleUserResponse struct {
	ID            string `json:"id"`
	Email         string `json:"email"`
	VerifiedEmail bool   `json:"verified_email"`
	Name          string `json:"name"`
	GivenName     string `json:"given_name"`
	FamilyName    string `json:"family_name"`
	Picture       string `json:"picture"`
	Locale        string `json:"locale"`
}

type GoogleClient struct {
	cfg    *config.GoogleConfig
	client *http.Client
}

func NewGoogleClient() *GoogleClient {
	return &GoogleClient{
		cfg:    config.GoogleCfg,
		client: &http.Client{},
	}
}

// Authorization Code -> AccessToken
func (g *GoogleClient) MakeAccessTokenForm(code string) url.Values {
	data := url.Values{}
	data.Set("grant_type", "authorization_code")
	data.Set("client_id", g.cfg.ClientId)
	data.Set("client_secret", g.cfg.ClientSecret)
	data.Set("redirect_uri", g.cfg.RedirectURI)
	data.Set("code", code)
	return data
}

func (g *GoogleClient) Login(code string) (OAuthUserInfo, error) {
	res, err := g.getAccessToken(code)
	if err != nil {
		return OAuthUserInfo{}, err
	}

	user, err := g.getUserInfo(res.AccessToken)
	if err != nil {
		return OAuthUserInfo{}, err
	}
	return user, nil
}

func (g *GoogleClient) getAccessToken(code string) (OAuthTokenResponse, error) {
	var tokenResp OAuthTokenResponse

	tokenURL := "https://oauth2.googleapis.com/token"
	res, err := g.client.PostForm(tokenURL, g.MakeAccessTokenForm(code))
	if err != nil {
		log.Printf("error occurred: %v", err)
		return tokenResp, err
	}
	defer res.Body.Close()

	if err := json.NewDecoder(res.Body).Decode(&tokenResp); err != nil {
		return tokenResp, err
	}
	return tokenResp, nil
}

func (g *GoogleClient) getUserInfo(accessToken string) (OAuthUserInfo, error) {
	var userInfo OAuthUserInfo
	var googleRes GoogleUserResponse

	URL := "https://www.googleapis.com/oauth2/v2/userinfo"

	req, _ := http.NewRequest("GET", URL, nil)
	req.Header.Set("Authorization", "Bearer "+accessToken)

	res, err := g.client.Do(req)
	if err != nil {
		return userInfo, err
	}
	defer res.Body.Close()

	if err := json.NewDecoder(res.Body).Decode(&googleRes); err != nil {
		return userInfo, err
	}

	userInfo.ID = googleRes.ID
	userInfo.Email = &googleRes.Email

	return userInfo, nil
}
