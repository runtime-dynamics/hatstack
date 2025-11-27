package api

import (
	"net/http"
	"runtime-dynamics/data"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

func RegisterRoutes(r *gin.Engine) {

	//r.GET("/api/exec-data-transit", execDataTransit)
}

func renderError(c *gin.Context, e error, statusCode int, message string) {
	if data.IsNotFound(e) {
		c.JSON(http.StatusNotFound, gin.H{
			"error": message,
		})
		log.Error().Err(e).Msg(message)
		return
	}
	log.Error().Err(e).Msg(message)
	c.JSON(statusCode, gin.H{
		"error": message,
	})
}

func renderSuccess(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "success",
	})
}

func renderFinal(c *gin.Context, err error, message string) {
	if err != nil {
		renderError(c, err, http.StatusInternalServerError, message)
		return
	}
	renderSuccess(c)
}

func renderFinalContent(c *gin.Context, content interface{}, key string, err error) {
	if err != nil {
		renderError(c, err, http.StatusInternalServerError, key)
		return
	}
	if key == "" {
		c.JSON(http.StatusOK, content)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		key: content,
	})
}
