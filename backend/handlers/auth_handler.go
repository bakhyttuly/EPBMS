package handlers

import (
	"epbms/config"
	"epbms/models"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type RegisterRequest struct {
	FullName string `json:"full_name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func Register(c *gin.Context) {
	var req RegisterRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	if req.Role != "admin" && req.Role != "organizer" && req.Role != "performer" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid role"})
		return
	}

	var existing models.User
	if err := config.DB.Where("email = ?", req.Email).First(&existing).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "email already exists"})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to hash password"})
		return
	}

	user := models.User{
		FullName: req.FullName,
		Email:    req.Email,
		Password: string(hashedPassword),
		Role:     req.Role,
	}

	if err := config.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create user"})
		return
	}

	if user.Role == "performer" {
		var existingPerformer models.Performer
		err := config.DB.Where("user_id = ?", user.ID).First(&existingPerformer).Error

		if err != nil {
			performer := models.Performer{
				UserID:             &user.ID,
				Name:               user.FullName,
				Category:           "Not set",
				Price:              0,
				Description:        "",
				AvailabilityStatus: "available",
			}

			if err := config.DB.Create(&performer).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create performer profile"})
				return
			}
		}
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "registered successfully",
	})
}

func Login(c *gin.Context) {
	var req LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	var user models.User
	if err := config.DB.Where("email = ?", req.Email).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid email or password"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid email or password"})
		return
	}

	session := sessions.Default(c)
	session.Set("user_id", user.ID)
	session.Set("role", user.Role)

	if err := session.Save(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save session"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "login successful",
		"role":    user.Role,
	})
}

func Logout(c *gin.Context) {
	session := sessions.Default(c)
	session.Clear()
	_ = session.Save()

	c.Redirect(http.StatusFound, "/")
}
