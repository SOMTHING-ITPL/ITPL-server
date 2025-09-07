package scheduler

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/SOMTHING-ITPL/ITPL-server/internal/api"
	"github.com/SOMTHING-ITPL/ITPL-server/performance"
)

type PerformanceScheduler struct {
	PerformanceRepo *performance.Repository
}

func (s *PerformanceScheduler) BuilderFacility(facility *api.FacilityDetailRes, region string) (*performance.Facility, error) {
	lo, err := strconv.ParseFloat(facility.Longitude, 64)
	if err != nil {
		return nil, err
	}

	la, err := strconv.ParseFloat(facility.Latitude, 64)
	if err != nil {
		return nil, err
	}

	return &performance.Facility{
		KopisFacilityKey: facility.ID,
		Name:             facility.Name,
		SeatCount:        facility.SeatScale,
		OpenedYear:       facility.OpenYear,
		Phone:            facility.TelNo,
		Homepage:         facility.URL,
		Address:          facility.Address,
		Latitude:         la,
		Longitude:        lo,
		Store:            facility.Store,
		Restaurant:       facility.Restaurant,
		Cafe:             facility.Cafe,
		ParkingLot:       facility.ParkingLot,
		Region:           region,
	}, nil
}

func (s *PerformanceScheduler) BuilderPerformanceImages(perfID uint, urls []string) []performance.PerformanceImage {
	sites := make([]performance.PerformanceImage, len(urls))
	for i, url := range urls {
		sites[i] = performance.PerformanceImage{
			PerformanceID: perfID,
			URL:           url,
		}
	}
	return sites
}
func (s *PerformanceScheduler) BuilderPerformanceTicketSite(perfID uint, urls []api.Relate) []performance.PerformanceTicketSite {
	sites := make([]performance.PerformanceTicketSite, len(urls))

	for i, url := range urls {
		sites[i] = performance.PerformanceTicketSite{
			PerformanceID: perfID,
			URL:           url.URL,
			Name:          url.Name,
		}
	}
	return sites
}

func (s *PerformanceScheduler) BuilderPerformance(res *api.PerformanceDetailRes, gptRes *GPTResponse, facilityId uint) (*performance.Performance, error) {
	layout := "2006.01.02"

	fromTime, err := time.Parse(layout, res.EndDate)
	if err != nil {
		fmt.Println("from 변환 실패:", err)
		fmt.Println("layout:", layout)
	}

	toTime, err := time.Parse(layout, res.StartDate)
	if err != nil {
		fmt.Println("to 변환 실패:", err)
		fmt.Println("layout:", layout)
	}
	fmt.Println("layout:", layout)

	layout = "2006-01-02 15:04:05"

	lastTime, err := time.Parse(layout, res.UpdateDate)
	if err != nil {
		fmt.Println("last 변환 실패:", err)
	}
	if strings.TrimSpace(res.Cast) == "" {
		res.Cast = gptRes.Cast
	}

	return &performance.Performance{
		Title:               res.Name,
		StartDate:           toTime,
		EndDate:             fromTime,
		Cast:                &res.Cast,
		Crew:                &res.Crew,
		KopisPerformanceKey: res.ID,
		KopisFacilityKey:    res.FacilityID,
		FacilityID:          facilityId,
		Runtime:             &res.Runtime,
		AgeRating:           &res.Age,
		FacilityName:        res.Facility,
		// Producer: res.,
		TicketPrice:  &res.Price,
		PosterURL:    &res.Poster,
		Region:       &res.Area,
		Status:       res.State,
		IsForeign:    res.Visit,
		LastModified: lastTime,
		Story:        &res.Story,
		DateGuidance: &res.DateGuidance,
		Genre:        gptRes.Genre,
		Keyword:      gptRes.Keyword,
	}, nil
}

// 공연 목록 조회 -> 공연 상세 조회 / LLM 추가 정보 수집 -> 공연 시설 조회
func (s *PerformanceScheduler) PutPerformanceList(startDate string, endDate string, isRunnung bool, afterDay *string) error {
	pge, row := 1, 100

	for {
		req := api.PerformanceListRequest{
			StartDate: startDate,
			EndDate:   endDate,
			CPage:     strconv.Itoa(pge),
			Rows:      strconv.Itoa(row),
			AfterDate: afterDay,
			Running:   isRunnung,
		}

		performanceList, err := api.GetPerformanceList(req)

		if err != nil {
			return err
		}
		if len(performanceList) == 0 {
			fmt.Print("no data ")
			return nil
		}

		for _, performance := range performanceList {
			performanceRes, err := api.GetDetailPerformance(performance.ID)
			if err != nil {
				return fmt.Errorf("Scheduler: Get Performance fail: %w", err)
			}

			facilityId, err := s.PutFacilityDetail(performanceRes.FacilityID, performanceRes.Area)
			if err != nil {
				return err
			}

			_, err = s.PutPerformanceDetail(performanceRes, performance.ID, facilityId)
			if err != nil {
				return err
			}
		}
		pge++
	}
}

// 아 애매하네 .. 이것도 캐싱형태로 해야 하나? 으으음 계속 날리는 형태로 ? 애매한데 ...
func (s *PerformanceScheduler) PutFacilityDetail(id string, region string) (uint, error) {
	//is already in db? 존재하거나 / error 일 경우 ->
	data, err := s.PerformanceRepo.GetFacilityByKopisID(id)
	if err != nil {
		return 0, err
	}
	if data != nil {
		return data.ID, nil //already exist
	}

	facilityRes, err := api.GetDetailFacility(id)
	if err != nil {
		return 0, fmt.Errorf("Scheduler: GetFacility fail: %w", err)
	}
	facility, err := s.BuilderFacility(facilityRes, region)
	if err != nil {
		return 0, fmt.Errorf("Scheduler: building Facility fail: %w", err)
	}
	facilityId, err := s.PerformanceRepo.CreateFacility(facility)
	if err != nil {
		return 0, fmt.Errorf("Scheduler: Create Facility fail: %w", err)
	}
	return facilityId, nil
}

func (s *PerformanceScheduler) UpdatePerformance(updatePerf *api.PerformanceDetailRes, originalPerf *performance.Performance, facilityID uint) (uint, error) {
	var layout = "2006-01-02 15:04:05"

	fromTime, err := time.Parse(layout, updatePerf.EndDate)
	if err != nil {
		fmt.Println("from 변환 실패:", err)
	}

	toTime, err := time.Parse(layout, updatePerf.StartDate)
	if err != nil {
		fmt.Println("to 변환 실패:", err)
	}

	//poster 나 예매처 변경된거에 대해서도 반영을 해줘야 하는부분 ..? 귀찮은데 ..
	originalPerf.AgeRating = &updatePerf.Age

	if strings.TrimSpace(updatePerf.Cast) != "" {
		originalPerf.Cast = &updatePerf.Cast
	}
	originalPerf.Crew = &updatePerf.Crew
	originalPerf.UpdatedAt = time.Now()
	originalPerf.DateGuidance = &updatePerf.DateGuidance
	originalPerf.EndDate = fromTime
	originalPerf.StartDate = toTime
	originalPerf.FacilityID = facilityID
	originalPerf.FacilityName = updatePerf.Facility

	if err := s.PerformanceRepo.UpdatePerformance(originalPerf); err != nil {
		return 0, fmt.Errorf("Scheduler: updating performance fail: %w", err)
	}

	return originalPerf.ID, nil
}

func (s *PerformanceScheduler) PutPerformanceDetail(res *api.PerformanceDetailRes, id string, facilityID uint) (uint, error) {
	// performanceRes, err := api.GetDetailPerformance(id)
	data, err := s.PerformanceRepo.GetPerformanceByKopisID(id)
	if err != nil {
		return 0, err
	}
	if data != nil {
		s.UpdatePerformance(res, data, facilityID)
		return data.ID, nil
	}

	//preprocess
	gptRes, err := PreProcessPerformance(res.Name, res.Cast)
	if err != nil {
		fmt.Printf("Scheduler: preprocess performance fail: %w", err) //log는 따로 남겨놔야 할 듯?
		gptRes = &GPTResponse{                                        //일단 임의의 값을 채워넣는 형태
			Keyword: res.Name,
			Genre:   13,
			Cast:    res.Cast,
		}
	}

	performance, err := s.BuilderPerformance(res, gptRes, facilityID)
	if err != nil {
		return 0, fmt.Errorf("Scheduler: building performance fail: %w", err)
	}

	performanceID, err := s.PerformanceRepo.CreatePerformance(performance)
	if err != nil {
		return 0, fmt.Errorf("Scheduler: Creating performance fail: %w", err)
	}

	imageList := s.BuilderPerformanceImages(performanceID, res.StyUrls)

	for _, image := range imageList {
		if err := s.PerformanceRepo.CreatePerformanceImage(&image); err != nil {
			return 0, fmt.Errorf("Scheduler: creating site fail: %w", err)
		}
	}

	ticketList := s.BuilderPerformanceTicketSite(performanceID, res.Relates)

	for _, ticket := range ticketList {
		if err := s.PerformanceRepo.CreatePerformanceTicketSite(&ticket); err != nil {
			return 0, fmt.Errorf("Scheduler: creating site fail: %w", err)
		}
	}

	return performanceID, nil
}
