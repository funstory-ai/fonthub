package handlers

import (
	"github.com/funstory-ai/fonthub/internal/fonts"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// GetFontsHandler handles the request to get all fonts
func GetFontsHandler(c *gin.Context) {
	fonts := fonts.GetAllFonts()
	c.JSON(200, gin.H{
		"success": true,
		"data":    fonts,
	})
}

// GetFontHandler handles the request to get a specific font
func GetFontsBySelectorHandler(c *gin.Context) {
	selector := fonts.FontSelector{
		Widths:    c.QueryArray("width"),
		Weights:   c.QueryArray("weight"),
		Styles:    c.QueryArray("style"),
		Languages: c.QueryArray("language"),
	}

	// Get matching fonts
	logrus.Info("selector: ", selector)
	matchingFonts := fonts.GetFontsBySelector(selector)

	c.JSON(200, gin.H{
		"success": true,
		"data":    matchingFonts,
	})
}
