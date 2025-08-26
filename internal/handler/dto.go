package handler

import "github.com/SOMTHING-ITPL/ITPL-server/performance"

// derefString safely dereferences a string pointer, returning empty string if nil
func derefString(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

func ToPerformanceShort(p performance.Performance) performanceShort {
	return performanceShort{
		Id:           p.ID,
		Title:        p.Title,
		State:        p.Status,
		PosterURL:    derefString(p.PosterURL),
		FacilityName: p.FacilityName,
		StartDate:    p.StartDate.Format("2006-01-02"),
		EndDate:      p.EndDate.Format("2006-01-02"),
	}
}

func ToPerformanceShortList(perfs []performance.Performance) []performanceShort {
	result := make([]performanceShort, len(perfs))
	for i, p := range perfs {
		result[i] = ToPerformanceShort(p)
	}
	return result
}
func ToPerformanceDetail(p performance.PerformanceWithTicketsAndImage) PerformanceDetail {
	return PerformanceDetail{
		Id:            p.ID,
		Title:         p.Title,
		StartDate:     p.StartDate,
		EndDate:       p.EndDate,
		FacilityID:    p.FacilityID,
		FacilityName:  p.FacilityName,
		AgeRating:     derefString(p.AgeRating), //null 너무 무  서  워
		TicketPrice:   derefString(p.TicketPrice),
		PosterURL:     derefString(p.PosterURL),
		Status:        p.Status,
		IsForeign:     p.IsForeign,
		DateGuidance:  p.DateGuidance,
		IntroImageURL: p.PerformanceImages, // []PerformanceImage
		TicketSite:    p.TicketSites,       // []PerformanceTicketSite
		LastModified:  p.LastModified,
	}
}

func ToFacilityShort(f *performance.Facility) FacilityShort {
	return FacilityShort{
		Id:        f.ID,
		Name:      f.Name,
		SeatCount: derefString(f.SeatCount),
	}
}

func ToFacilityDetail(f *performance.Facility) FacilityDetail {
	return FacilityDetail{
		Id:         f.ID,
		Name:       f.Name,
		OpenedYear: f.OpenedYear,
		SeatCount:  derefString(f.SeatCount),
		Phone:      f.Phone,
		Homepage:   f.Homepage,
		Address:    f.Address,
		Latitude:   f.Latitude,
		Longitude:  f.Longitude,
		Restaurant: f.Restaurant,
		Cafe:       f.Cafe,
		Store:      f.Store,
		ParkingLot: f.ParkingLot,
	}
}

func ToFacilityShortList(facilities []performance.Facility) []FacilityShort {
	details := make([]FacilityShort, len(facilities))
	for i, f := range facilities {
		details[i] = ToFacilityShort(&f)
	}
	return details
}
