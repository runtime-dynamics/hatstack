package web

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"runtime-dynamics/web/api"
	"runtime-dynamics/web/app"
)

func GetStaticFiles(staticDir string) (map[string]string, error) {
	staticFiles := make(map[string]string)
	err := filepath.Walk(staticDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			relativePath := filepath.ToSlash(strings.TrimPrefix(path, staticDir))
			url := "/" + strings.TrimPrefix(relativePath, "/")
			log.Debug().Msgf("Found static file: %s, setting url: %s", path, url)
			staticFiles[url] = path
			if strings.HasSuffix(url, ".html") {
				staticFiles[strings.TrimSuffix(url, ".html")] = path
			}
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return staticFiles, nil
}

func HandleRoutes(r *gin.Engine, staticDir string) *gin.Engine {
	// Register API routes (JSON endpoints under /api/*)
	api.RegisterRoutes(r)

	// Register web routes (HTML endpoints under /app/* and /)
	// Note: This includes the homepage at /
	app.RegisterWebRoutes(r)

	// In development mode, serve static directories
	// In production, individual files are registered below
	if os.Getenv("IS_DEV") == "true" {
		r.Static("/images", fmt.Sprintf("%s/images", staticDir))
		r.Static("/css", fmt.Sprintf("%s/css", staticDir))
		r.Static("/js", fmt.Sprintf("%s/js", staticDir))
	} else {
		/* Production: handle all the static files individually */
		if len(os.Getenv("NO_STATIC")) < 1 {
			staticFiles, err := GetStaticFiles(staticDir)
			if err != nil {
				log.Fatal().Msgf("Error getting static files: %s", err)
			}
			for url, path := range staticFiles {
				// Skip root path as it's handled by webhandlers
				if url == "/" || url == "/index.html" || url == "/index" {
					continue
				}
				// Skip desktop-login as it's now handled by webhandlers
				if url == "/desktop-login" || url == "/desktop-login.html" {
					continue
				}
				//log.Info().Msgf("Setting static file: %s -> %s", url, path)
				r.StaticFile(url, path)
			}
		}
	}

	return r
}

func Start(r *gin.Engine) *gin.Engine {
	HandleRoutes(r, "static")
	return r
}
