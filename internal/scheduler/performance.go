package scheduler

import (
	"fmt"
	"strconv"
	"time"

	"github.com/SOMTHING-ITPL/ITPL-server/internal/api"
	"github.com/SOMTHING-ITPL/ITPL-server/performance"
)

type PerformanceScheduler struct {
	PerformanceRepo *performance.Repository
}

func (s *PerformanceScheduler) BuilderFacility(facility *api.FacilityDetailRes) (*performance.Facility, error) {
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
	}, nil
}

func (s *PerformanceScheduler) BuildPerformanceTicketSites(perfID uint, urls []string) []performance.PerformanceTicketSite {
	sites := make([]performance.PerformanceTicketSite, len(urls))
	for i, url := range urls {
		sites[i] = performance.PerformanceTicketSite{
			PerformanceID: perfID,
			TicketSite:    url,
		}
	}
	return sites
}

func (s *PerformanceScheduler) BuilderPerformance(res *api.PerformanceDetailRes, facilityId uint) (*performance.Performance, error) {
	layout := "2006.01.02"

	fromTime, err := time.Parse(layout, res.EndDate)
	if err != nil {
		fmt.Println("from 변환 실패:", err)
	}

	toTime, err := time.Parse(layout, res.StartDate)
	if err != nil {
		fmt.Println("to 변환 실패:", err)
	}

	layout = "2006-01-02 15:04:05"

	lastTime, err := time.Parse(layout, res.UpdateDate)
	if err != nil {
		fmt.Println("last 변환 실패:", err)
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
		// Producer: res.,
		TicketPrice:  &res.Price,
		PosterURL:    &res.Poster,
		Region:       &res.Area,
		Status:       &res.State,
		IsForeign:    res.Visit,
		LastModified: lastTime,
		Story:        &res.Story,
		DateGuidance: &res.DateGuidance,
	}, nil
}

// 공연 목록 조회 -> 공연 상세 조회 / LLM 추가 정보 수집 -> 공연 시설 조회
func (s *PerformanceScheduler) PutPerformanceList(startDate string, endDate string, afterDay *string, isRunnung bool) error {
	pge, row := 0, 100

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
			return nil
		}

		for _, performance := range performanceList {
			facilityId, err := s.PutFacilityDetail(performance.Facility)
			if err != nil {
				return err
			}

			_, err = s.PutPerformanceDetail(performance.ID, facilityId)
			if err != nil {
				return err
			}
		}
		pge++
	}
}

// 아 애매하네 .. 이것도 캐싱형태로 해야 하나? 으으음 계속 날리는 형태로 ? 애매한데 ...
func (s *PerformanceScheduler) PutFacilityDetail(id string) (uint, error) {
	//is already in db?
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
	facility, err := s.BuilderFacility(facilityRes)
	if err != nil {
		return 0, fmt.Errorf("Scheduler: building Facility fail: %w", err)
	}
	facilityId, err := s.PerformanceRepo.CreateFacility(facility)
	if err != nil {
		return 0, fmt.Errorf("Scheduler: Create Facility fail: %w", err)
	}
	return facilityId, nil
}

func (s *PerformanceScheduler) PutPerformanceDetail(id string, facilityID uint) (uint, error) {
	performanceRes, err := api.GetDetailPerformance(id)
	if err != nil {
		return 0, fmt.Errorf("Scheduler: Get Performance fail: %w", err)
	}

	performance, err := s.BuilderPerformance(performanceRes, facilityID)
	if err != nil {
		return 0, fmt.Errorf("Scheduler: building performance fail: %w", err)
	}

	performanceID, err := s.PerformanceRepo.CreatePerformance(performance)
	if err != nil {
		return 0, fmt.Errorf("Scheduler: Creating performance fail: %w", err)
	}

	ticketList := s.BuildPerformanceTicketSites(performanceID, performanceRes.StyUrls)
	if err != nil {
		return 0, fmt.Errorf("Scheduler: building site fail: %w", err)
	}
	for _, ticket := range ticketList {
		if err := s.PerformanceRepo.CreatePerformanceTicketSite(&ticket); err != nil {
			return 0, fmt.Errorf("Scheduler: creating site fail: %w", err)
		}
	}
	return performanceID, nil
}

//스케줄러 예시
// func main() {
// 	ticker := time.NewTicker(24 * time.Hour) // 24시간마다 실행
// 	defer ticker.Stop()

// 	// 초기 실행
// 	doTask()

// 	for {
// 		select {
// 		case <-ticker.C:
// 			doTask()
// 		}
// 	}
// }

// func doTask() {
// 	fmt.Println("오늘 공연 데이터를 업데이트합니다:", time.Now())
// 	// 여기에 공연 데이터를 긁어오는 로직 작성
// }
