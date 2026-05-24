// Package handlers 提供 HTTP 处理器
package handlers

import (
	"github.com/BenedictKing/ccx/internal/config"
	"github.com/gin-gonic/gin"
)

// GetFuzzyMode 获取 Fuzzy 模式状态
func GetFuzzyMode(cfgManager *config.ConfigManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(200, gin.H{
			"fuzzyModeEnabled": cfgManager.GetFuzzyModeEnabled(),
		})
	}
}

// SetFuzzyMode 设置 Fuzzy 模式状态
func SetFuzzyMode(cfgManager *config.ConfigManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			Enabled bool `json:"enabled"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(400, gin.H{"error": "Invalid request body"})
			return
		}

		if err := cfgManager.SetFuzzyModeEnabled(req.Enabled); err != nil {
			c.JSON(500, gin.H{"error": "Failed to save config"})
			return
		}

		c.JSON(200, gin.H{
			"success":          true,
			"fuzzyModeEnabled": req.Enabled,
		})
	}
}

// GetStripBillingHeader 获取移除计费头状态
func GetStripBillingHeader(cfgManager *config.ConfigManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(200, gin.H{
			"stripBillingHeader": cfgManager.GetStripBillingHeader(),
		})
	}
}

// SetStripBillingHeader 设置移除计费头状态
func SetStripBillingHeader(cfgManager *config.ConfigManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			Enabled bool `json:"enabled"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(400, gin.H{"error": "Invalid request body"})
			return
		}

		if err := cfgManager.SetStripBillingHeader(req.Enabled); err != nil {
			c.JSON(500, gin.H{"error": "Failed to save config"})
			return
		}

		c.JSON(200, gin.H{
			"success":            true,
			"stripBillingHeader": req.Enabled,
		})
	}
}

// GetThemePreference 获取主题偏好
func GetThemePreference(cfgManager *config.ConfigManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(200, gin.H{
			"themePreference": cfgManager.GetThemePreference(),
		})
	}
}

// SetThemePreference 设置主题偏好
func SetThemePreference(cfgManager *config.ConfigManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			Theme string `json:"theme"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(400, gin.H{"error": "Invalid request body"})
			return
		}

		if req.Theme != "light" && req.Theme != "dark" && req.Theme != "auto" {
			c.JSON(400, gin.H{"error": "Theme must be 'light', 'dark', or 'auto'"})
			return
		}

		if err := cfgManager.SetThemePreference(req.Theme); err != nil {
			c.JSON(500, gin.H{"error": "Failed to save config"})
			return
		}

		c.JSON(200, gin.H{
			"success":         true,
			"themePreference": req.Theme,
		})
	}
}
