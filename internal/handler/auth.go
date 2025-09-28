
package handler

import (
	"errors"
	"net/http"
	"runrun/internal"
	"runrun/internal/protocol"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// AuthHandler handles user registration and login.
func AuthHandler(c *gin.Context) {
	var req protocol.AuthRequest
	// Bind the request body to the AuthRequest struct.
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "Invalid request body"})
		return
	}

	// Basic validation
	if req.Account == "" || req.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "Account and password are required"})
		return
	}

	var user internal.User
	// Check if the user account already exists.
	result := internal.DB.Where("account = ?", req.Account).First(&user)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		// Case 1: User does not exist, create a new user (register).
		newUser := internal.User{
			Account:  req.Account,
			Password: req.Password, // TODO: Hash the password before saving!
		}
		if createResult := internal.DB.Create(&newUser); createResult.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": "Failed to create user"})
			return
		}
		c.JSON(http.StatusCreated, gin.H{"code": 1, "msg": "登记成功"})
	} else if result.Error == nil {
		// Case 2: User exists, check password (login).
		if user.Password == req.Password { // TODO: Use a secure password comparison!
			// Password matches
			c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "登录成功"})
		} else {
			// Password does not match
			c.JSON(http.StatusUnauthorized, gin.H{"code": 3, "msg": "登录失败"})
		}
	} else {
		// Any other database error.
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": "Database error"})
	}
}
