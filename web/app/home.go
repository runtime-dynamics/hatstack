package app

import (
	"runtime-dynamics/views/pages"

	"github.com/gin-gonic/gin"
)

// HomePageHandler renders the homepage
func HomePageHandler(c *gin.Context) {
	// Render the homepage using templ
	component := pages.Home()
	if err := component.Render(c.Request.Context(), c.Writer); err != nil {
		c.String(500, "Error rendering page")
	}
}
