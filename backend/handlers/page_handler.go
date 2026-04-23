package handlers

import "github.com/gin-gonic/gin"

func ShowLoginPage(c *gin.Context) {
	c.HTML(200, "login.html", gin.H{
		"title": "Login",
	})
}

func ShowRegisterPage(c *gin.Context) {
	c.HTML(200, "register.html", gin.H{
		"title": "Register",
	})
}

func ShowDashboardPage(c *gin.Context) {
	c.HTML(200, "dashboard.html", gin.H{
		"title": "Dashboard",
	})
}

func ShowPerformersPage(c *gin.Context) {
	c.HTML(200, "performers.html", gin.H{
		"title": "Performers",
	})
}

func ShowBookingsPage(c *gin.Context) {
	c.HTML(200, "bookings.html", gin.H{
		"title": "Bookings",
	})
}

func ShowCalendarPage(c *gin.Context) {
	c.HTML(200, "calendar.html", gin.H{
		"title": "Calendar",
	})
}

func ShowMySchedulePage(c *gin.Context) {
	c.HTML(200, "my_schedule.html", gin.H{
		"title": "My Schedule",
	})
}
