package handler

import (
	"net/http"
	"strconv"

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
		Genre   string `form:"genre"`   //optional
		Region  string `form:"region"`  // optional
		Keyword string `form:"keyword"` // optional
	}

	return func(c *gin.Context) {
		var req PerformanceListRequest
		if err := c.ShouldBindQuery(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid query params"})
			return
		}

		//조회
		performances, err := p.performanceRepo.FindPerformances(req.Page, req.Limit, req.Genre, req.Region, req.Keyword)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "fail to findPerformances"})
			return
		}

		c.JSON(http.StatusOK, PerformanceListRes{
			Performances: ToPerformanceShortList(performances),
			Count:        len(performances),
		})
	}
}

// 공연 상세 조회
func (p *PerformanceHandler) GetPerformanceDetail() gin.HandlerFunc {
	return func(c *gin.Context) {
		idStr := c.Param("id")
		id, err := strconv.ParseUint(idStr, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
			return
		}

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

		c.JSON(http.StatusOK, gin.H{
			"message":     "success",
			"performance": perfRes,
			"facility":    facilityRes,
		})

	}
}

// 공연 시설 목록 조회
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

		facilities, err := p.performanceRepo.FindFacilities(req.Page, req.Limit, req.Region)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"fail to get facility": err})
		}

		c.JSON(http.StatusOK, FacilityListRes{
			Facilities: ToFacilityShortList(facilities),
			Count:      len(facilities),
		})
	}
}

// 공연 시설 상세 조회
func (p *PerformanceHandler) GetFacilityDetail() gin.HandlerFunc {
	return func(c *gin.Context) {
		idStr := c.Param("id")
		id, err := strconv.ParseUint(idStr, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
			return
		}

		//get facility
		facility, err := p.performanceRepo.GetFacilityById(uint(id))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"fail to get facility": err})
		}

		c.JSON(http.StatusOK, gin.H{
			"message":  "success",
			"facility": ToFacilityDetail(facility),
		})
	}
}

// 최근 본 공연 목록 조회
func (p *PerformanceHandler) GetRecentViewPerformance() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("userID")
		if !exists {
			c.JSON(401, gin.H{"error": "unauthorized"})
			return
		}

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

		c.JSON(http.StatusOK, PerformanceListRes{
			Performances: ToPerformanceShortList(performances),
			Count:        len(performances),
		})
	}
}

// top N 공연 목록 조회 + score 표시해줘야 함?
func (p *PerformanceHandler) GetTopPerformances() gin.HandlerFunc {
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

		perfIDs := make([]uint, len(perfScores))
		for i, score := range perfScores {
			perfIDs[i] = score.ID
		}

		performances, err := p.performanceRepo.GetPerformancesByIDs(perfIDs)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "fail to get performances"})
			return
		}

		c.JSON(http.StatusOK, PerformanceListRes{
			Performances: ToPerformanceShortList(performances),
			Count:        len(performances),
		})
	}
}

// 공연 좋아요 생성
func (p *PerformanceHandler) CreatePerformanceLike() gin.HandlerFunc {
	return func(c *gin.Context) {
		userIDStr, exists := c.Get("userID")
		if !exists {
			c.JSON(401, gin.H{"error": "unauthorized"})
			return
		}

		//middleware로 빼버릴까..
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
		userIDStr, exists := c.Get("userID")
		if !exists {
			c.JSON(401, gin.H{"error": "unauthorized"})
			return
		}

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

		c.JSON(http.StatusOK, PerformanceListRes{
			Performances: ToPerformanceShortList(performances),
			Count:        len(performances),
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
