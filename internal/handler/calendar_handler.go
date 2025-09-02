package handler

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/SOMTHING-ITPL/ITPL-server/calendar"
	"github.com/SOMTHING-ITPL/ITPL-server/performance"
	"github.com/gin-gonic/gin"
)

func NewCalendarHandler(calRepo *calendar.Repository, perfRepo *performance.Repository) *CalendarHandler {
	return &CalendarHandler{
		calendarRepo:    calRepo,
		performanceRepo: perfRepo,
	}
}

func (ch *CalendarHandler) CreateCalendarData() gin.HandlerFunc {
	type req struct {
		Date          string `json:"date" binding:"required"`
		PerformanceID uint   `json:"performance_id" binding:"required"`
	}

	return func(c *gin.Context) {
		var req req
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid query params"})
			return
		}

		parsedDate, err := time.Parse("20060102", req.Date)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "date must be in format yyyymmdd"})
			return
		}

		userIDStr, _ := c.Get("userID")

		userID, ok := userIDStr.(uint)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid user id type"})
			return
		}

		_, err = ch.performanceRepo.GetPerformanceById(req.PerformanceID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "ther is no performance id"})
			return
		}

		if ch.calendarRepo.CreateCalendar(req.PerformanceID, uint(userID), parsedDate) != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "fail to create Calendar record"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	}
}

func (ch *CalendarHandler) DeleteCalendarData() gin.HandlerFunc {
	return func(c *gin.Context) {
		calendarIDStr := c.Param("id")
		calendarID, err := strconv.ParseUint(calendarIDStr, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "path parameter type is not uint"})
			return
		}

		if ch.calendarRepo.DeleteCalendar(uint(calendarID)) != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "fail to delete calendar data "})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "success"})
	}
}

func (ch *CalendarHandler) GetCalendarData() gin.HandlerFunc {
	type res struct {
		Performance performanceShort `json:"performances"`
		CalendarID  uint             `json:"calendar_id"`
	}

	return func(c *gin.Context) {
		userIDStr, _ := c.Get("userID")

		userID, ok := userIDStr.(uint)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid user id type"})
			return
		}
		yearStr := c.Query("year")
		monthStr := c.Query("month")

		if yearStr == "" || monthStr == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "year and month are required"})
			return
		}

		year, err := strconv.Atoi(yearStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid year"})
			return
		}

		month, err := strconv.Atoi(monthStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid month"})
			return
		}

		data, err := ch.calendarRepo.GetCalendar(userID, month, year)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "fail to get data!"})
			return
		}

		result := make(map[string][]res)
		for _, cal := range data {
			short := ToPerformanceShort(cal.Performance)

			dayKey := fmt.Sprintf("%d", cal.Day)
			result[dayKey] = append(result[dayKey], res{
				Performance: short,
				CalendarID:  cal.ID,
			})
		}
		c.JSON(http.StatusOK, CommonRes{
			Message: "success",
			Data:    result,
		})
	}
}
