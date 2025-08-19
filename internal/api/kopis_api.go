package api

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/SOMTHING-ITPL/ITPL-server/config"
	"github.com/clbanning/mxj/v2"
)

// 그냥 config에 넣기 귀찮아 아니 이게 오히려 깔끔한 것 같기도 하고
const KopisBaseURL = "http://kopis.or.kr/openApi/restful/pblprfr"

type PerformanceListRequest struct {
	StartDate string // 공연시작일, 필수, 8
	EndDate   string // 공연종료일, 필수, 8 (최대 31일)
	CPage     string // 현재페이지, 필수, 3
	Rows      string // 페이지당 목록수, 필수, 3
	Running   bool
	AfterDate *string // 해당일자 이후, 8
}

func GetPerformanceList(req PerformanceListRequest) (map[string]interface{}, error) {
	params := url.Values{}
	params.Set("service", config.KopisCfg.SecretKey)
	params.Set("stdate", req.StartDate)
	params.Set("eddate", req.EndDate)
	params.Set("cpage", req.CPage)
	params.Set("rows", req.Rows)
	params.Set("shcate", "CCCD") // 대중음악

	if req.Running {
		params.Set("prfstate", "02")
	} else {
		params.Set("prfstate", "01")
	}

	if req.AfterDate != nil {
		params.Set("afterdate", *req.AfterDate)
	}

	fullURL := fmt.Sprintf("%s?%s", KopisBaseURL, params.Encode())

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get(fullURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	//걍 xml -> json 형태로
	mv, err := mxj.NewMapXml(body)
	if err != nil {
		return nil, err
	}

	return mv, nil
}

func GetDetailPerformance(performanceId string) (map[string]interface{}, error) {

	fullURL := fmt.Sprintf("%s/%s", KopisBaseURL, performanceId)

	params := url.Values{}
	params.Set("service", config.KopisCfg.SecretKey)

	fullURL = fmt.Sprintf("%s?%s", fullURL, params.Encode())

	client := &http.Client{Timeout: 10 * time.Second}

	resp, err := client.Get(fullURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	//걍 xml -> json 형태로
	mv, err := mxj.NewMapXml(body)
	if err != nil {
		return nil, err
	}

	return mv, nil
}

func GetDetailFacility(facilityID string) (map[string]interface{}, error) {

	fullURL := fmt.Sprintf("%s/%s", KopisBaseURL, facilityID)

	params := url.Values{}
	params.Set("service", config.KopisCfg.SecretKey)

	fullURL = fmt.Sprintf("%s?%s", fullURL, params.Encode())

	client := &http.Client{Timeout: 10 * time.Second}

	resp, err := client.Get(fullURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	//걍 xml -> json 형태로
	mv, err := mxj.NewMapXml(body)
	if err != nil {
		return nil, err
	}

	return mv, nil
}
