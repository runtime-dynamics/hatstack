package app

import (
	"runtime-dynamics/web/app/seo"

	"github.com/gin-gonic/gin"
)

// RegisterWebRoutes registers all web-facing routes (non-API)
func RegisterWebRoutes(r *gin.Engine) {
	// SEO and crawler files
	r.GET("/robots.txt", seo.RobotsTxtHandler)
	r.GET("/sitemap.xml", seo.SitemapXMLHandler)

	// Homepage (public)
	r.GET("/", HomePageHandler)

}
