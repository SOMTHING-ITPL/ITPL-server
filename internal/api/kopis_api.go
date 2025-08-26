package api

import (
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/SOMTHING-ITPL/ITPL-server/config"
)

// 그냥 config에 넣기 귀찮아 아니 이게 오히려 깔끔한 것 같기도 하고
const KopisBaseURL = "http://kopis.or.kr/openApi/restful"

type KopisResponse[T any] struct {
	Items []T `xml:"db" json:"items"`
}

type PerformanceListRes struct {
	ID        string `xml:"mt20id"`
	Name      string `xml:"prfnm"`
	StartDate string `xml:"prfpdfrom"`
	EndDate   string `xml:"prfpdto"`
	Facility  string `xml:"fcltynm"`
	Poster    string `xml:"poster"`
	Area      string `xml:"area"`
	Genre     string `xml:"genrenm"`
	OpenRun   string `xml:"openrun"`
	State     string `xml:"prfstate"`
}

type PerformanceListRequest struct {
	StartDate string // 공연시작일, 필수, 8
	EndDate   string // 공연종료일, 필수, 8 (최대 31일)
	CPage     string // 현재페이지, 필수, 3
	Rows      string // 페이지당 목록수, 필수, 3
	Running   bool
	AfterDate *string // 해당일자 이후, 8
}

type PerformanceDetailRes struct {
	ID        string `xml:"mt20id"       `
	Name      string `xml:"prfnm"        `
	StartDate string `xml:"prfpdfrom"    `
	EndDate   string `xml:"prfpdto"      `
	Facility  string `xml:"fcltynm"     ` //공연 시설명
	Cast      string `xml:"prfcast"      `
	Crew      string `xml:"prfcrew"      `
	Runtime   string `xml:"prfruntime"   `
	Age       string `xml:"prfage"       ` //관람 연령
	// Company       string   `xml:"entrpsnm"     json:"company"` //기획 제작사
	// CompanyP      string   `xml:"entrpsnmP"    json:"company_p"` //제작사
	// CompanyA      string   `xml:"entrpsnmA"    json:"company_a"`// 기획사
	// CompanyH string `xml:"entrpsnmH"    json:"company_h"` //주최
	// CompanyS string `xml:"entrpsnmS"    json:"company_s"` //주관
	Price  string `xml:"pcseguidance"` //티켓 가격
	Poster string `xml:"poster"`       // 포스터 이미지 경로
	Story  string `xml:"sty"`          //줄거리
	Area   string `xml:"area"`         // 지역
	// Genre    string `xml:"genrenm"      json:"genre"`     // 장르
	Visit string `xml:"visit"        ` //내한
	// Festival      string   `xml:"festival"     json:"festival"` //축제 여부
	UpdateDate   string   `xml:"updatedate"` //최근 수정일
	State        string   `xml:"prfstate"`
	StyUrls      []string `xml:"styurls>styurl"`
	FacilityID   string   `xml:"mt10id"`
	DateGuidance string   `xml:"dtguidance"`

	Relates []Relate `xml:"relates>relate"`
}

// 공연 시설 상세 정보
type FacilityDetailRes struct {
	Name string `xml:"fcltynm"     `
	ID   string `xml:"mt10id"      `
	// VenueCount int    `xml:"mt13cnt"     json:"venue_count"`
	Category  *string `xml:"fcltychartr"`
	OpenYear  *string `xml:"opende"    `
	SeatScale *string `xml:"seatscale" `
	TelNo     *string `xml:"telno"     `
	URL       *string `xml:"relateurl"   ` //homepage
	Address   string  `xml:"adres"       `
	Latitude  string  `xml:"la"          `
	Longitude string  `xml:"lo"          `

	Restaurant string `xml:"restaurant"  `
	Cafe       string `xml:"cafe"        `
	Store      string `xml:"store"       `
	// Nolibang    string `xml:"nolibang"    json:"nolibang"`
	// Suyu        string `xml:"suyu"        json:"suyu"`
	// ParkBarrier string `xml:"parkbarrier" json:"park_barrier"`
	// RestBarrier string `xml:"restbarrier" json:"rest_barrier"`
	// RunwBarrier string `xml:"runwbarrier" json:"runw_barrier"`
	// ElevBarrier string `xml:"elevbarrier" json:"elev_barrier"`
	ParkingLot string `xml:"parkinglot"  ` //주차시설
}

type Relate struct {
	Name string `xml:"relatenm"  json:"name,omitempty"`
	URL  string `xml:"relateurl" json:"url"`
}

func kopisGet[T any](endpoint string, params map[string]string) ([]T, error) {
	values := url.Values{}
	values.Set("service", config.KopisCfg.SecretKey)
	for k, v := range params {
		values.Set(k, v)
	}
	fullURL := fmt.Sprintf("%s?%s", endpoint, values.Encode())
	client := &http.Client{Timeout: 10 * time.Second}

	resp, err := client.Get(fullURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var res KopisResponse[T]
	if err := xml.Unmarshal(body, &res); err != nil {
		return nil, err
	}
	return res.Items, nil
}

func GetPerformanceList(req PerformanceListRequest) ([]PerformanceListRes, error) {
	params := map[string]string{
		"stdate":   req.StartDate,
		"eddate":   req.EndDate,
		"cpage":    req.CPage,
		"rows":     req.Rows,
		"shcate":   "CCCD",
		"prfstate": "01",
	}

	if req.Running {
		params["prfstate"] = "02"
	}
	if req.AfterDate != nil {
		params["afterdate"] = *req.AfterDate
	}
	url := fmt.Sprintf("%s/%s", KopisBaseURL, "pblprfr")

	return kopisGet[PerformanceListRes](url, params)
}

func GetDetailPerformance(performanceID string) (*PerformanceDetailRes, error) {
	items, err := kopisGet[PerformanceDetailRes](fmt.Sprintf("%s/%s/%s", KopisBaseURL, "pblprfr", performanceID), nil)
	if err != nil {
		return nil, err
	}
	if len(items) == 0 {
		return nil, fmt.Errorf("performance not found: %s", performanceID)
	}
	return &items[0], nil
}

func GetDetailFacility(facilityID string) (*FacilityDetailRes, error) {
	items, err := kopisGet[FacilityDetailRes](fmt.Sprintf("%s/%s/%s", KopisBaseURL, "prfplc", facilityID), nil)
	if err != nil {
		return nil, err
	}
	if len(items) == 0 {
		return nil, fmt.Errorf("facility not found: %s", facilityID)
	}
	return &items[0], nil
}
