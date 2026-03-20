package middleware

import (
	"agnos-test/model"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func HospitalMiddleware(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		slug := c.GetHeader("X-Hospital-Slug")
		if slug == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "hospital context is required (X-Hospital-Slug header)"})
			c.Abort()
			return
		}

		var hospital model.Hospital
		if err := db.Where("slug = ?", slug).First(&hospital).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid hospital slug: " + slug})
			c.Abort()
			return
		}

		c.Set("hospital_id", hospital.ID)
		c.Set("hospital_slug", hospital.Slug)
		c.Set("hospital_name", hospital.Name)
		c.Next()
	}
}
