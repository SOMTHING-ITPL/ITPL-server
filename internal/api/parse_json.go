package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

//조씨. 이기범이 reflect 찾아보래.

// Item 구조체: JSON "item" 배열 안의 객체 하나
type Item struct {
	ContentId   string `json:"contentid"`
	Title       string `json:"title"`
	Addr1       string `json:"addr1"`
	Addr2       string `json:"addr2"`
	Tel         string `json:"tel"`
	MapX        string `json:"mapx"`
	MapY        string `json:"mapy"`
	FirstImage  string `json:"firstimage"`
	CreatedTime string `json:"createdtime"`
	CategoryID  string `json:"contenttypeid"`
}

// Items 구조체: JSON "items" 필드
type Items struct {
	ItemList []Item `json:"item"`
}

// Body 구조체: JSON "body" 필드
type Body struct {
	Items Items `json:"items"`
}

// Header 구조체: JSON "header" 필드
type Header struct {
	ResultCode string `json:"resultCode"`
	ResultMsg  string `json:"resultMsg"`
}

// Response 구조체: JSON "response" 필드
type Response struct {
	Header Header `json:"header"`
	Body   Body   `json:"body"`
}

// Root 구조체: 최상위 JSON 구조
type Root struct {
	Response Response `json:"response"`
}

// fetchAndParseJSON 함수: URL로 GET 요청 보내고, JSON을 파싱해서 Item 리스트 반환
func FetchAndParseJSON(url string) ([]Item, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("HTTP 요청 실패: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("비정상 응답 코드: %d", resp.StatusCode)
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("응답 바디 읽기 실패: %w", err)
	}

	var root Root
	if err := json.Unmarshal(bodyBytes, &root); err != nil {
		return nil, fmt.Errorf("JSON 파싱 실패: %w", err)
	}

	// 결과 코드 확인
	if root.Response.Header.ResultCode != "0000" {
		return nil, fmt.Errorf("API 응답 실패: %s", root.Response.Header.ResultMsg)
	}

	return root.Response.Body.Items.ItemList, nil
}
