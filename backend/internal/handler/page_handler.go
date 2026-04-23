package handler

import (
	"net/http"
	"github.com/gin-gonic/gin"
)

// PageHandler serves the HTML templates for the frontend.
type PageHandler struct{}

func NewPageHandler() *PageHandler {
	return &PageHandler{}
}

func (h *PageHandler) ShowLoginPage(c *gin.Context) {
	c.HTML(http.StatusOK, "login.html", gin.H{"title": "Login"})
}

func (h *PageHandler) ShowRegisterPage(c *gin.Context) {
	c.HTML(http.StatusOK, "register.html", gin.H{"title": "Register"})
}

func (h *PageHandler) ShowDashboardPage(c *gin.Context) {
	c.HTML(http.StatusOK, "dashboard.html", gin.H{"title": "Dashboard"})
}

func (h *PageHandler) ShowPerformersPage(c *gin.Context) {
	c.HTML(http.StatusOK, "performers.html", gin.H{"title": "Performers"})
}

func (h *PageHandler) ShowBookingsPage(c *gin.Context) {
	c.HTML(http.StatusOK, "bookings.html", gin.H{"title": "Bookings"})
}

func (h *PageHandler) ShowCalendarPage(c *gin.Context) {
	c.HTML(http.StatusOK, "calendar.html", gin.H{"title": "Calendar"})
}

func (h *PageHandler) ShowMySchedulePage(c *gin.Context) {
	c.HTML(http.StatusOK, "my_schedule.html", gin.H{"title": "My Schedule"})
}
