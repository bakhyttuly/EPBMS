package main

import (
	"epbms/config"
	"epbms/routes"
	"path/filepath"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

func main() {
	config.ConnectDB()

	r := gin.Default()

	store := cookie.NewStore([]byte("super-secret-session-key"))
	r.Use(sessions.Sessions("epbms_session", store))

	frontendPath := filepath.Join("..", "frontend")
	templatesPath := filepath.Join(frontendPath, "templates", "*")
	staticPath := filepath.Join(frontendPath, "static")

	r.LoadHTMLGlob(templatesPath)
	r.Static("/static", staticPath)

	routes.SetupRoutes(r)

	r.Run(":8080")
}
