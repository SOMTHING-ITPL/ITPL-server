package user

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"

	"github.com/SOMTHING-ITPL/ITPL-server/internal/api"
)

func GetProviderClient(p string) (api.OAuthClient, error) {
	if p == string(ProviderGoogle) {
		return api.NewKakaoClient(), nil //this should be google client
	}
	if p == string(ProviderKakao) {
		return api.NewKakaoClient(), nil
	}

	return nil, fmt.Errorf("unsupported provider: %s", p)
}

func GenerateNanoIDNickname() string {
	bytes := make([]byte, 8)
	rand.Read(bytes)
	id := hex.EncodeToString(bytes)
	return "user-" + id
}
