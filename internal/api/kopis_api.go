package api

import (
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

// 그냥 config에 넣기 귀찮아 아니 이게 오히려 깔끔한 것 같기도 하고
const KopisBaseURL = "http://kopis.or.kr/openApi/restful/pblprfr"

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
	ID        string `xml:"mt20id"       json:"id"`
	Name      string `xml:"prfnm"        json:"name"`
	StartDate string `xml:"prfpdfrom"    json:"start_date"`
	EndDate   string `xml:"prfpdto"      json:"end_date"`
	Facility  string `xml:"fcltynm"      json:"facility"` //공연 시설명
	Cast      string `xml:"prfcast"      json:"cast"`
	Crew      string `xml:"prfcrew"      json:"crew"`
	Runtime   string `xml:"prfruntime"   json:"runtime"`
	Age       string `xml:"prfage"       json:"age"` //관람 연령
	// Company       string   `xml:"entrpsnm"     json:"company"` //기획 제작사
	// CompanyP      string   `xml:"entrpsnmP"    json:"company_p"` //제작사
	// CompanyA      string   `xml:"entrpsnmA"    json:"company_a"`// 기획사
	// CompanyH string `xml:"entrpsnmH"    json:"company_h"` //주최
	// CompanyS string `xml:"entrpsnmS"    json:"company_s"` //주관
	Price  string `xml:"pcseguidance" json:"price"`  //티켓 가격
	Poster string `xml:"poster"       json:"poster"` // 포스터 이미지 경로
	Story  string `xml:"sty"          json:"story"`  //줄거리
	Area   string `xml:"area"         json:"area"`   // 지역
	// Genre    string `xml:"genrenm"      json:"genre"`     // 장르
	Visit string `xml:"visit"        json:"visit"` //내한
	// Festival      string   `xml:"festival"     json:"festival"` //축제 여부
	UpdateDate   string   `xml:"updatedate"   json:"update_date"`   //최근 수정일
	State        string   `xml:"prfstate"     json:"state"`         //공연 상태
	StyUrls      []string `xml:"styurls>styurl" json:"sty_urls"`    //소개 이미지 목록
	FacilityID   string   `xml:"mt10id"       json:"facility_id"`   //공연 시설 ID
	DateGuidance string   `xml:"dtguidance"   json:"date_guidance"` // 공연 시간

	Relates []Relate `xml:"relates>relate" json:"relates"` //예매처
}

// 공연 시설 상세 정보
type FacilityDetailRes struct {
	Name string `xml:"fcltynm"     json:"name"`
	ID   string `xml:"mt10id"      json:"id"`
	// VenueCount int    `xml:"mt13cnt"     json:"venue_count"`
	Category  *string `xml:"fcltychartr" json:"category"`
	OpenYear  *string `xml:"opende"      json:"open_year"`
	SeatScale *string `xml:"seatscale"   json:"seat_scale"`
	TelNo     *string `xml:"telno"       json:"tel"`
	URL       *string `xml:"relateurl"   json:"url"` //homepage
	Address   string  `xml:"adres"       json:"address"`
	Latitude  string  `xml:"la"          json:"latitude"`
	Longitude string  `xml:"lo"          json:"longitude"`

	Restaurant string `xml:"restaurant"  json:"restaurant"`
	Cafe       string `xml:"cafe"        json:"cafe"`
	Store      string `xml:"store"       json:"store"`
	// Nolibang    string `xml:"nolibang"    json:"nolibang"`
	// Suyu        string `xml:"suyu"        json:"suyu"`
	// ParkBarrier string `xml:"parkbarrier" json:"park_barrier"`
	// RestBarrier string `xml:"restbarrier" json:"rest_barrier"`
	// RunwBarrier string `xml:"runwbarrier" json:"runw_barrier"`
	// ElevBarrier string `xml:"elevbarrier" json:"elev_barrier"`
	ParkingLot string `xml:"parkinglot"  json:"parking_lot"` //주차시설
}

type Relate struct {
	Name string `xml:"relatenm"  json:"name,omitempty"`
	URL  string `xml:"relateurl" json:"url"`
}

func kopisGet[T any](endpoint string, params map[string]string) ([]T, error) {
	values := url.Values{}
	values.Set("service", "config.KopisCfg.SecretKey")
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

	return kopisGet[PerformanceListRes](KopisBaseURL, params)
}

func GetDetailPerformance(performanceID string) (*PerformanceDetailRes, error) {
	items, err := kopisGet[PerformanceDetailRes](fmt.Sprintf("%s/%s", KopisBaseURL, performanceID), nil)
	if err != nil {
		return nil, err
	}
	if len(items) == 0 {
		return nil, fmt.Errorf("performance not found: %s", performanceID)
	}
	return &items[0], nil
}

func GetDetailFacility(facilityID string) (*FacilityDetailRes, error) {
	items, err := kopisGet[FacilityDetailRes](fmt.Sprintf("%s/%s", KopisBaseURL, facilityID), nil)
	if err != nil {
		return nil, err
	}
	if len(items) == 0 {
		return nil, fmt.Errorf("facility not found: %s", facilityID)
	}
	return &items[0], nil
}
