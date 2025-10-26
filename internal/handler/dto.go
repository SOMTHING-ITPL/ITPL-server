package handler

import (
	"fmt"
	"time"

	"github.com/SOMTHING-ITPL/ITPL-server/aws/s3"
	"github.com/SOMTHING-ITPL/ITPL-server/chat"
	"github.com/SOMTHING-ITPL/ITPL-server/course"
	"github.com/SOMTHING-ITPL/ITPL-server/performance"
	"github.com/SOMTHING-ITPL/ITPL-server/place"
	"github.com/SOMTHING-ITPL/ITPL-server/user"
	"github.com/aws/aws-sdk-go-v2/aws"
	"gorm.io/gorm"
)

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
		IntroImageURL: IntroImageUrlToString(p.PerformanceImages), // []PerformanceImage
		TicketSite:    p.TicketSites,                              // []PerformanceTicketSite
		LastModified:  p.LastModified,
	}
}

func IntroImageUrlToString(urls []performance.PerformanceImage) []string {
	res := make([]string, 0, len(urls))
	for _, u := range urls {
		res = append(res, u.URL)
	}
	return res
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

func ToCourseInfo(course course.Course) CourseInfoResponse {
	return CourseInfoResponse{
		ID:          course.ID,
		CreatedAt:   course.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   course.UpdatedAt.Format(time.RFC3339),
		UserID:      course.UserID,
		Title:       course.Title,
		Description: course.Description,
		IsAICreated: course.IsAICreated,
		FacilityID:  course.FacilityID,
	}
}

func ToCourseDetails(db *gorm.DB, details []course.CourseDetail) ([]CourseDetailResponse, error) {
	var courseDetailResponse []CourseDetailResponse
	for _, detail := range details {
		res := CourseDetailResponse{
			ID:         detail.ID,
			CreatedAt:  detail.CreatedAt.Format(time.RFC3339),
			UpdatedAt:  detail.UpdatedAt.Format(time.RFC3339),
			CourseID:   detail.CourseID,
			Day:        detail.Day,
			Sequence:   detail.Sequence,
			PlaceID:    detail.PlaceID,
			PlaceTitle: detail.PlaceTitle,
			PlaceImage: place.GetImageByPlaceID(db, detail.PlaceID),
			Address:    detail.Address,
			Latitud:    detail.Latitud,
			Longitude:  detail.Longitude,
		}
		courseDetailResponse = append(courseDetailResponse, res)
	}
	return courseDetailResponse, nil
}

func ToChatRoomInfoResponse(cfg aws.Config, bucketName string, room chat.ChatRoom) (ChatRoomInfoResponse, error) {
	imageUrl, err := s3.GetPresignURL(cfg, bucketName, *room.ImageKey)
	currentMembers := len(room.Members)
	return ChatRoomInfoResponse{
		ID:             room.ID,
		Title:          room.Title,
		ImageUrl:       &imageUrl,
		PerformanceDay: room.PerformanceDay,
		MaxMembers:     room.MaxMembers,
		DepartureName:  room.DepartureName,
		CurrentMembers: currentMembers,
		ArrivalName:    room.ArrivalName,
	}, err
}

func ToChatRoomMemberInfoResponse(cfg aws.Config, bucketName string, DB *gorm.DB, userID uint) (ChatRoomMemberResponse, error) {
	var user user.User

	result := DB.First(&user, userID)
	if result.Error != nil {
		fmt.Printf("get user error : %s\n", result.Error)
		return ChatRoomMemberResponse{}, result.Error
	}
	var url string
	var err error
	if user.Photo != nil {
		url, err = s3.GetPresignURL(cfg, bucketName, aws.ToString(user.Photo))
	}

	if err != nil {
		return ChatRoomMemberResponse{}, err
	}
	return ChatRoomMemberResponse{
		UserID:   user.ID,
		Nickname: user.NickName,
		ImageUrl: url,
	}, nil
}
