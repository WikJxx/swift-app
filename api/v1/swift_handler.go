package v1

import (
	"net/http"

	"swift-app/internal/services" // Importujemy logikę biznesową

	"github.com/gin-gonic/gin"
)

// Handler pobierający SWIFT Code
func GetSwiftCode(c *gin.Context) {
	code := c.Param("swift-code")
	swift, err := services.FindSwiftByCode(code) // Wywołanie funkcji z services
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "SWIFT Code not found"})
		return
	}
	c.JSON(http.StatusOK, swift)
}
