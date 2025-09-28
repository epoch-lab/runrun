package handler

import (
	"errors"
	"log"
	"math/rand"
	"net/http"
	"runrun/internal"
	"runrun/internal/protocol"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// AuthHandler handles user registration and login.
func AuthHandler(c *gin.Context) {
	// 1. 从请求体中解析 account 和 password
	var req protocol.AuthRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "Invalid request body"})
		return
	}

	// 基础验证
	if req.Account == "" || req.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "Account and password are required"})
		return
	}

	// 2. 先调用一次Login，验证账号密码是否正确
	clientInfo := protocol.GenerateFakeClient()
	userInfo, err := protocol.Login(req.Account, req.Password, clientInfo)
	if err != nil {
		// Login失败
		c.JSON(http.StatusUnauthorized, gin.H{"code": 3, "msg": "账号或者密码错误"})
		return
	}

	// Login成功，继续处理数据库逻辑
	var user internal.User
	result := internal.DB.Where("account = ?", req.Account).First(&user)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		// Case 1: User does not exist, create a new user (register).
		newUser := internal.User{
			Account:  req.Account,
			Password: req.Password, // 直接存储明文密码
		}
		if createResult := internal.DB.Create(&newUser); createResult.Error != nil {
			log.Printf("Failed to create user: %v", createResult.Error)
			c.JSON(http.StatusInternalServerError, gin.H{"code": 5, "msg": "其他错误"})
			return
		}
		c.JSON(http.StatusCreated, gin.H{"code": 1, "msg": "登记成功"})
		
	} else if result.Error == nil {
		// 账号已存在：验证密码并执行跑步
		if user.Password == req.Password {
			// 密码匹配：调用submit立即执行一次跑步
			if err := executeRunForUser(*userInfo, clientInfo, req.Account); err != nil {
				log.Printf("Failed to execute run for user %s: %v", req.Account, err)
				c.JSON(http.StatusInternalServerError, gin.H{"code": 5, "msg": "其他错误"})
				return
			}
			c.JSON(http.StatusOK, gin.H{"code": 2, "msg": "已登记"})
		} else {
			// 密码不匹配
			c.JSON(http.StatusUnauthorized, gin.H{"code": 3, "msg": "账号或者密码错误"})
		}
		
	} else {
		// 其他数据库错误
		log.Printf("Database error in AuthHandler: %v", result.Error)
		c.JSON(http.StatusInternalServerError, gin.H{"code": 5, "msg": "其他错误"})
	}
}


// executeRunForUser 为用户执行一次跑步
func executeRunForUser(userInfo protocol.UserInfo, clientInfo protocol.ClientInfo, account string) error {
	// 生成随机跑步参数 (4-5km, 20-30分钟)
	distance := int64(4000 + rand.Intn(1000)) // 4000-4999米
	duration := int32(20 + rand.Intn(10))     // 20-29分钟
	
	// 提交跑步记录
	err := protocol.Submit(userInfo, clientInfo, duration, distance)
	if err != nil {
		return err
	}
	
	// 更新数据库中的跑步距离
	return updateUserRunningProgress(account, float64(distance)/1000.0) // 转换为公里
}

// updateUserRunningProgress 更新用户跑步进度
func updateUserRunningProgress(account string, distanceKm float64) error {
	var user internal.User
	// 通过账号查找用户
	result := internal.DB.Where("account = ?", account).First(&user)
	if result.Error != nil {
		return result.Error
	}
	
	// 更新当前距离
	user.CurrentDistance += distanceKm
	
	// 检查是否达到目标距离
	if user.CurrentDistance >= user.TargetDistance {
		user.IsRunningRequired = false
		log.Printf("User %s has reached target distance %.2f km", user.Account, user.TargetDistance)
	}
	
	// 保存更新
	return internal.DB.Save(&user).Error
}
