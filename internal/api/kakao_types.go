package api

type KakaoAccountProfile struct {
	Nickname        *string `json:"nickname"`
	ProfileImageURL *string `json:"profile_image_url"`
}

type KakaoAccount struct {
	Email   *string             `json:"email"`
	Profile KakaoAccountProfile `json:"profile"`
}

type KakaoUserResponse struct {
	ID           int64        `json:"id"`
	KakaoAccount KakaoAccount `json:"kakao_account"`
}
