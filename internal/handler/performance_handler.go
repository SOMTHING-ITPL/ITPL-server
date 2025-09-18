package handler

import (
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/SOMTHING-ITPL/ITPL-server/performance"
	"github.com/gin-gonic/gin"
)

func NewPerformanceHandler(repo *performance.Repository) *PerformanceHandler {
	return &PerformanceHandler{performanceRepo: repo}
}

// performance List 검색 시에..
// query parameter : page, limit , 장르 , 지역 , 키워드 검색 가능하도록
func (p *PerformanceHandler) GetPerformanceShortList() gin.HandlerFunc {
	type PerformanceListRequest struct {
		Page    int    `form:"page" binding:"required"`
		Limit   int    `form:"limit" binding:"required"`
		Genre   int    `form:"genre"`   //optional
		Region  string `form:"region"`  // optional
		Keyword string `form:"keyword"` // optional
	}

	return func(c *gin.Context) {
		var req PerformanceListRequest
		if err := c.ShouldBindQuery(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid body"})
			return
		}

		performances, total, err := p.performanceRepo.FindPerformances(req.Page, req.Limit, req.Genre, req.Region, req.Keyword)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "fail to findPerformances"})
			return
		}

		c.JSON(http.StatusOK, CommonRes{
			Message: "success",
			Data: PerformanceListRes{
				Performances: ToPerformanceShortList(performances),
				Count:        int(total),
			},
		})
	}
}

// 공연 상세 조회
func (p *PerformanceHandler) GetPerformanceDetail() gin.HandlerFunc {
	type res struct {
		Performance PerformanceDetail `json:"performance"`
		Facility    FacilityDetail    `json:"facility"`
	}
	return func(c *gin.Context) {
		idStr := c.Param("id")
		id, err := strconv.ParseUint(idStr, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
			return
		}

		userIDStr, _ := c.Get("userID")
		userID, _ := userIDStr.(uint)

		perfDetail, err := p.performanceRepo.GetPerformanceWithTicketsAndImages(uint(id))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"Get detail performance data error": err})
			return
		}

		perfRes := ToPerformanceDetail(*perfDetail)

		facility, err := p.performanceRepo.GetFacilityById(perfDetail.FacilityID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"Get detail facility data error": err})
			return
		}
		facilityRes := ToFacilityDetail(facility)

		if err := p.performanceRepo.IncrementPerformanceScore(uint(id), 1, c.Request.Context()); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "fail to increase score"})
			return
		}
		if err := p.performanceRepo.CreateRecentView(userID, uint(id), c.Request.Context()); err != nil { //user 최근 공연 집계
			c.JSON(http.StatusInternalServerError, gin.H{"error": "fail to increase score"})
			return
		}

		c.JSON(http.StatusOK, CommonRes{
			Message: "success",
			Data: res{
				Performance: perfRes,
				Facility:    facilityRes,
			},
		})
	}
}

// 공연 시설 목록 조회 + 그것도 해야 하나 .. ?
func (p *PerformanceHandler) GetFacilityList() gin.HandlerFunc {
	type FacilityListRequest struct {
		Page   int    `form:"page" binding:"required"`
		Limit  int    `form:"limit" binding:"required"`
		Region string `form:"region"` //optional
	}

	return func(c *gin.Context) {
		var req FacilityListRequest
		if err := c.ShouldBindQuery(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid query params"})
			return
		}

		facilities, total, err := p.performanceRepo.FindFacilities(req.Page, req.Limit, req.Region)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"fail to get facility": err})
		}

		c.JSON(http.StatusOK, CommonRes{
			Message: "success",
			Data: FacilityListRes{
				Facilities: ToFacilityShortList(facilities),
				Count:      int(total), // 이러면 안되는게 맞긴 한데 .. ㅋㅋ 64비트 꽉 차겠음?
			},
		})
	}
}

// 공연 시설 상세 조회
func (p *PerformanceHandler) GetFacilityDetail() gin.HandlerFunc {
	type res struct {
		Facility     FacilityDetail     `json:"facility"`
		Performances []performanceShort `json:"related_performances"`
	}
	return func(c *gin.Context) {
		idStr := c.Param("id")
		id, err := strconv.ParseUint(idStr, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
			return
		}

		facility, performances, err := p.performanceRepo.GetFacilityByIdWithPerf(uint(id))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"fail to get facility": err})
		}

		c.JSON(http.StatusOK, CommonRes{
			Message: "success",
			Data: res{
				Facility:     ToFacilityDetail(facility),
				Performances: ToPerformanceShortList(performances),
			},
		})
	}
}

// 최근 본 공연 목록 조회
func (p *PerformanceHandler) GetRecentViewPerformance() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, _ := c.Get("userID")

		userIDUint, ok := userID.(uint)
		if !ok {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id type"})
			return
		}

		perfIds, err := p.performanceRepo.GetRecentViews(userIDUint, c.Request.Context())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "fail to get recent views"})
			return
		}
		//recentViews를 통해서 short 반환해야 함.

		performances, err := p.performanceRepo.GetPerformancesByIDs(perfIds)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "fail to get performances"})
			return
		}

		c.JSON(http.StatusOK, CommonRes{
			Message: "success",
			Data: PerformanceListRes{
				Performances: ToPerformanceShortList(performances),
				Count:        len(performances),
			},
		})
	}
}

// top N 공연 목록 조회 + score 표시해줘야 함?
func (p *PerformanceHandler) GetTopPerformances() gin.HandlerFunc {
	type res struct {
		Performance performanceShort `json:"performance"`
		Score       float64          `json:"score"`
	}
	return func(c *gin.Context) {
		numStr := c.Query("num")
		if numStr == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Top num is required"})
			return
		}

		num, err := strconv.ParseInt(numStr, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid num type"})
			return
		}

		perfScores, err := p.performanceRepo.GetTopPerformances(num, c.Request.Context())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "fail to get top performances"})
			return
		}

		//여기 잘못 짠 듯
		perfIDs := make([]uint, len(perfScores))
		for i, score := range perfScores {
			perfIDs[i] = score.ID
		}

		performances, err := p.performanceRepo.GetPerformancesByIDs(perfIDs)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "fail to get performances"})
			return
		}

		result := make([]res, 0, len(perfScores))
		performanceMap := make(map[uint]performance.Performance)
		for _, p := range performances {
			performanceMap[p.ID] = p
		}

		for _, score := range perfScores {
			if perf, ok := performanceMap[score.ID]; ok {
				result = append(result, res{
					Performance: ToPerformanceShort(perf),
					Score:       score.Score,
				})
			}
		}

		c.JSON(http.StatusOK, CommonRes{
			Message: "success",
			Data:    result,
		})
	}
}

// 공연 좋아요 생성
func (p *PerformanceHandler) CreatePerformanceLike() gin.HandlerFunc {
	return func(c *gin.Context) {
		userIDStr, _ := c.Get("userID")

		userID, ok := userIDStr.(uint)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid user id type"})
			return
		}

		perfIdStr := c.Param("id")
		perfid, err := strconv.ParseUint(perfIdStr, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "path parameter type is not uint"})
			return
		}

		if err = p.performanceRepo.CreateUserLike(uint(perfid), uint(userID)); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "fail to put data record"})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"message": "success",
		})
	}
}

// 공연 좋아요 목록 조회
func (p *PerformanceHandler) GetPerformanceLike() gin.HandlerFunc {
	return func(c *gin.Context) {
		userIDStr, _ := c.Get("userID")

		//middleware로 빼버릴까..
		userID, ok := userIDStr.(uint)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid user id type"})
			return
		}

		performances, err := p.performanceRepo.GetUserLike(userID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "fail to get userLike"})

		}

		c.JSON(http.StatusOK, CommonRes{
			Message: "success",
			Data: PerformanceListRes{
				Performances: ToPerformanceShortList(performances),
				Count:        len(performances),
			},
		})
	}
}

// 공연 좋아요 삭제
func (p *PerformanceHandler) DeletePerformanceLike() gin.HandlerFunc {
	return func(c *gin.Context) {
		userIDStr, exists := c.Get("userID")
		if !exists {
			c.JSON(401, gin.H{"error": "unauthorized"})
			return
		}

		userID, ok := userIDStr.(uint)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid user id type"})
			return
		}

		perfIdStr := c.Param("id")
		perfid, err := strconv.ParseUint(perfIdStr, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "path parameter type is not uint"})
			return
		}

		if err = p.performanceRepo.DeleteUserLike(uint(perfid), uint(userID)); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "fail to delete data record"})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"message": "success",
		})

	}
}

func (p *PerformanceHandler) AiRecommendation() gin.HandlerFunc {
	type res struct {
		Performance performanceShort `json:"performance"`
		Score       float64          `json:"score"`
	}

	return func(c *gin.Context) {
		// userIDStr, exists := c.Get("userID") //getUserId
		// userID, ok := userIDStr.(uint)

		perfIds := make([]uint, 3)
		for i := 0; i < 3; i++ {
			num := rand.Intn(101) + 480 // 0~100 + 450 = 450~550
			perfIds[i] = uint(num)
		}

		perf, err := p.performanceRepo.GetPerformancesByIDs(perfIds)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "gail to get perfs"})
			return
		}

		result := make([]res, 0, 3)
		performanceMap := make(map[uint]performance.Performance)
		for _, p := range perf {
			performanceMap[p.ID] = p
		}

		for idx, score := range perf {
			idx++

			if perf, ok := performanceMap[score.ID]; ok {
				result = append(result, res{
					Performance: ToPerformanceShort(perf),
					Score:       float64(idx),
				})
			}
		}

		c.JSON(http.StatusOK, CommonRes{
			Message: "success",
			Data:    result,
		})
	}
}

// func (p *PerformanceHandler) IncrementPerformanceView() gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		userIDStr, _ := c.Get("userID")

// 		userID, ok := userIDStr.(uint)
// 		if !ok {
// 			c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid user id type"})
// 			return
// 		}

// 		perfIdStr := c.Query("perfId")
// 		if perfIdStr == "" {
// 			c.JSON(http.StatusBadRequest, gin.H{"error": " is required"})
// 			return
// 		}

// 		perfId, err := strconv.ParseInt(perfIdStr, 10, 64)
// 		if err != nil {
// 			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid num type"})
// 			return
// 		}
// 		//performance 존재하는지 확인
// 		_, err = p.performanceRepo.GetPerformanceById(uint(perfId))
// 		if errors.Is(err, gorm.ErrRecordNotFound) {
// 			c.JSON(http.StatusBadRequest, gin.H{"error": "not exist that performance id "})
// 			return
// 		} else if err != nil {
// 			c.JSON(http.StatusInternalServerError, gin.H{"error": err})
// 			return
// 		}
// 		if err := p.performanceRepo.IncrementPerformanceScore(uint(perfId), 1, c.Request.Context()); err != nil {
// 			c.JSON(http.StatusInternalServerError, gin.H{"error": "fail to increase score"})
// 			return
// 		}
// 		if err := p.performanceRepo.CreateRecentView(userID, uint(perfId), c.Request.Context()); err != nil { //user 최근 공연 집계
// 			c.JSON(http.StatusInternalServerError, gin.H{"error": "fail to increase score"})
// 			return
// 		}

// 		c.JSON(http.StatusOK, gin.H{"message": "success"})
// 		return
// 	}
// }

// 다가오는 공연
func (p *PerformanceHandler) GetCommingPerformances() gin.HandlerFunc {
	type res struct {
		Performance performanceShort `json:"performance"`
	}
	return func(c *gin.Context) {
		numStr := c.Query("num")
		if numStr == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "num is required"})
			return
		}

		num, err := strconv.ParseInt(numStr, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid num type"})
			return
		}

		perfs, err := p.performanceRepo.GetRecentPerformance(time.Now(), int(num))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "fail to get performance data"})
			return
		}

		result := make([]res, len(perfs))
		for i, p := range perfs {
			result[i] = res{Performance: ToPerformanceShort(p)}
		}

		c.JSON(http.StatusOK, CommonRes{
			Message: "success",
			Data:    result,
		})
	}
}
