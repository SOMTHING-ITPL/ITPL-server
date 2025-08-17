package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"

	"github.com/SOMTHING-ITPL/ITPL-server/config"
)

// type KakaoAccountProfile struct {
// 	Nickname        *string `json:"nickname"`
// 	ProfileImageURL *string `json:"profile_image_url"`
// }

// type KakaoAccount struct {
// 	Email   *string             `json:"email"`
// 	Profile KakaoAccountProfile `json:"profile"`
// }

type KakaoUserResponse struct {
	ID int64 `json:"id"`
	// KakaoAccount KakaoAccount `json:"kakao_account"`
}

type KakaoClient struct {
	cfg    *config.KaKaoConfig
	client *http.Client
}

func NewKakaoClient() *KakaoClient {
	return &KakaoClient{
		cfg:    config.KakaoCfg,
		client: &http.Client{},
	}
}

func (k *KakaoClient) MakeAccessTokenForm(code string) url.Values {
	data := url.Values{}
	data.Set("grant_type", "authorization_code")
	data.Set("client_id", k.cfg.ClientId)
	data.Set("redirect_uri", k.cfg.RedirectURI)
	data.Set("client_secret", k.cfg.ClientSecret)
	data.Set("code", code)
	return data
}

func (k *KakaoClient) Login(code string) (OAuthUserInfo, error) {
	res, err := k.getAccessToken(code)
	if err != nil {
		return OAuthUserInfo{}, err
	}

	user, err := k.getUserInfo(res.AccessToken)
	if err != nil {
		return OAuthUserInfo{}, err
	}
	return user, nil
}

func (k *KakaoClient) getAccessToken(code string) (OAuthTokenResponse, error) {
	var tokenResp OAuthTokenResponse

	tokenURL := fmt.Sprintf("%s/oauth/token", k.cfg.Domain)
	res, err := k.client.PostForm(tokenURL, k.MakeAccessTokenForm(code))
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

func (k *KakaoClient) getUserInfo(accessToken string) (OAuthUserInfo, error) {
	var userInfo OAuthUserInfo
	var kakaoRes KakaoUserResponse

	URL := fmt.Sprintf("%s/v2/user/me", k.cfg.Domain)

	req, _ := http.NewRequest("GET", URL, nil)
	req.Header.Set("Authorization", "Bearer "+accessToken)

	res, err := k.client.Do(req)
	if err != nil {
		return userInfo, err
	}
	defer res.Body.Close()

	if err := json.NewDecoder(res.Body).Decode(&kakaoRes); err != nil {
		return userInfo, err
	}

	userInfo.ID = fmt.Sprintf("%d", kakaoRes.ID)
	// userInfo.Nickname = kakaoRes.KakaoAccount.Profile.Nickname
	// userInfo.Photo = kakaoRes.KakaoAccount.Profile.ProfileImageURL

	return userInfo, nil
}
