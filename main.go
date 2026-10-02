package main

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"ttcs/database"
	"ttcs/util"
	"ttcs/web"
	"ttcs/xray"
)

const version = "1.0.0"

func main() {
	dbPath := "/etc/ttcs/ttcs.db"
	port := "12345"

	if err := database.Init(dbPath); err != nil {
		println("Database error:", err.Error())
		os.Exit(1)
	}

	if err := database.AutoMigrate(); err != nil {
		println("Migration error:", err.Error())
		os.Exit(1)
	}

	r := web.NewServer()

	// =========================
	// LOGIN
	// =========================

	r.POST("/login", loginHandler)

	// =========================
	// HEALTH
	// =========================

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "ok",
			"version": version,
		})
	})

	// =========================
	// SYSTEM INFORMATION
	// =========================

	r.GET("/api/system", func(c *gin.Context) {
		info, err := util.GetSystemInfo()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, info)
	})
	r.GET("/api/inbounds", func(c *gin.Context) {
		var inbounds []database.Inbound

		if err := database.DB.Find(&inbounds).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"success":  true,
			"inbounds": inbounds,
		})
	})

	r.POST("/api/inbounds", func(c *gin.Context) {
		var inbound database.Inbound

		if err := c.ShouldBindJSON(&inbound); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "Invalid request",
			})
			return
		}

		if inbound.Remark == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "Remark cannot be empty",
			})
			return
		}

		if inbound.Protocol == "" {
			inbound.Protocol = "vmess"
		}

		if inbound.Port <= 0 {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "Port must be greater than 0",
			})
			return
		}

		if err := database.DB.Create(&inbound).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": err.Error(),
			})
			return
		}
		if err := xray.Apply(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "Inbound saved but failed to apply Xray config: " + err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "Inbound created",
			"inbound": inbound,
		})
	})

	r.PUT("/api/inbounds/:id", func(c *gin.Context) {
		var inbound database.Inbound

		if err := database.DB.First(&inbound, c.Param("id")).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{
				"success": false,
				"message": "Inbound not found",
			})
			return
		}

		var req database.Inbound

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "Invalid request",
			})
			return
		}

		req.ID = inbound.ID

		if err := database.DB.Save(&req).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": err.Error(),
			})
			return
		}
		if err := xray.Apply(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "Inbound updated but failed to apply Xray config: " + err.Error(),
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "Inbound updated",
			"inbound": req,
		})
	})

	r.DELETE("/api/inbounds/:id", func(c *gin.Context) {
		var inbound database.Inbound

		if err := database.DB.First(&inbound, c.Param("id")).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{
				"success": false,
				"message": "Inbound not found",
			})
			return
		}

		if err := database.DB.Delete(&inbound).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": err.Error(),
			})
			return
		}
		if err := xray.Apply(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "Inbound deleted but failed to apply Xray config: " + err.Error(),
			})
			return
		}
		if err := xray.Apply(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "Inbound deleted but failed to apply Xray config: " + err.Error(),
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "Inbound deleted",
		})
	})

	r.POST("/api/inbounds/:id/toggle", func(c *gin.Context) {
		var inbound database.Inbound

		if err := database.DB.First(&inbound, c.Param("id")).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{
				"success": false,
				"message": "Inbound not found",
			})
			return
		}

		inbound.Enabled = !inbound.Enabled

		if err := database.DB.Save(&inbound).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": err.Error(),
			})
			return
		}
		if err := xray.Apply(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "Inbound toggled but failed to apply Xray config: " + err.Error(),
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"enabled": inbound.Enabled,
		})
	})
	// =========================
	// SETTINGS - GET ALL
	// =========================

	r.GET("/api/settings/all", func(c *gin.Context) {
		var settings []database.Setting

		if err := database.DB.Find(&settings).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		result := make(map[string]string)

		for _, setting := range settings {
			result[setting.Key] = setting.Value
		}

		c.JSON(http.StatusOK, result)
	})

	// =========================
	// SETTINGS - UPDATE
	// =========================

	r.POST("/api/settings/update", func(c *gin.Context) {
		var req map[string]string

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "Invalid request",
			})
			return
		}

		for key, value := range req {
			var setting database.Setting

			result := database.DB.Where("key = ?", key).First(&setting)

			if result.Error != nil {
				if result.Error == gorm.ErrRecordNotFound {
					setting = database.Setting{
						Key:   key,
						Value: value,
					}

					if err := database.DB.Create(&setting).Error; err != nil {
						c.JSON(http.StatusInternalServerError, gin.H{
							"success": false,
							"message": err.Error(),
						})
						return
					}
				} else {
					c.JSON(http.StatusInternalServerError, gin.H{
						"success": false,
						"message": result.Error.Error(),
					})
					return
				}
			} else {
				setting.Value = value

				if err := database.DB.Save(&setting).Error; err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{
						"success": false,
						"message": err.Error(),
					})
					return
				}
			}
		}

		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "Settings updated",
		})
	})

	// =========================
	// SETTINGS - UPDATE USER
	// =========================

	r.POST("/api/settings/updateUser", func(c *gin.Context) {
		var req struct {
			Username string `json:"username"`
			Password string `json:"password"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "Invalid request",
			})
			return
		}

		if req.Username == "" || req.Password == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "Username and password cannot be empty",
			})
			return
		}

		var user database.User

		result := database.DB.First(&user)

		if result.Error != nil {
			user = database.User{
				Username: req.Username,
				Password: req.Password,
			}

			if err := database.DB.Create(&user).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"success": false,
					"message": err.Error(),
				})
				return
			}
		} else {
			user.Username = req.Username
			user.Password = req.Password

			if err := database.DB.Save(&user).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"success": false,
					"message": err.Error(),
				})
				return
			}
		}

		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "User updated",
		})
	})

	// =========================
	// START SERVER
	// =========================

	println("ttcs panel starting...")
	println("Database:", dbPath)
	println("Port:", port)

	if err := r.Run(":" + port); err != nil {
		println("Server error:", err.Error())
		os.Exit(1)
	}
}

// =========================
// LOGIN HANDLER
// =========================

func loginHandler(c *gin.Context) {
	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid request",
		})
		return
	}

	var user database.User

	result := database.DB.Where("username = ?", req.Username).First(&user)

	if result.Error != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "Invalid username or password",
		})
		return
	}

	if req.Password != user.Password {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "Invalid username or password",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Login successful",
	})
}
