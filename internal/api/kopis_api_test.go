package api

import (
	"fmt"
	"testing"
)

func TestGetPerformanceList(t *testing.T) {
	req := PerformanceListRequest{
		StartDate: "20250801",
		EndDate:   "20250820",
		CPage:     "1",
		Rows:      "5",
		Running:   true,
	}

	list, err := GetPerformanceList(req)
	if err != nil {
		t.Fatalf("GetPerformanceList error: %v", err)
	}

	if len(list) == 0 {
		t.Fatal("No performances returned")
	}

	for _, p := range list {
		fmt.Printf("%s | %s | %s ~ %s\n", p.ID, p.Name, p.StartDate, p.EndDate)
	}
}

func TestGetDetailPerformance(t *testing.T) {
	perfID := "PF132236"
	detail, err := GetDetailPerformance(perfID)
	if err != nil {
		t.Fatalf("GetDetailPerformance error: %v", err)
	}

	fmt.Printf("ID: %s\nName: %s\nFacility: %s\nStartDate: %s\nEndDate: %s\n",
		detail.ID, detail.Name, detail.Facility, detail.StartDate, detail.EndDate)
}

func TestGetDetailFacility(t *testing.T) {
	facilityID := "FC001247"
	facility, err := GetDetailFacility(facilityID)
	if err != nil {
		t.Fatalf("GetDetailFacility error: %v", err)
	}

	fmt.Printf("Name: %s\nCategory: %s\nSeatScale: %s\nAddress: %s\n",
		facility.Name, facility.Category, facility.SeatScale, facility.Address)
}
