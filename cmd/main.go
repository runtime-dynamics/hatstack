package main

import (
	"os"
	"runtime-dynamics/config"
	"time"

	"runtime-dynamics/web"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var (
	version = "source"
)

func setLogger() {
	zerolog.TimeFieldFormat = time.RFC3339Nano
	//if len(os.Getenv("CONSOLE_LOG")) > 0 {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout})
	//}
	if os.Getenv("DEBUG") != "" {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
		return
	}
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
}

func main() {
	setLogger()
	err := config.LoadConfig()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to load config")
	}
	log.Info().Msgf("Starting StarXAPI (Version: %s)", version)
	if os.Getenv("DEBUG") == "" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Desktop token cleanup not required; tokens are stored in datastore and removed on connect
	router := gin.New()
	router.Use(func(c *gin.Context) {
		start := time.Now()
		c.Next()
		end := time.Now()
		latency := end.Sub(start)
		log.Info().
			Str("host", c.Request.Host).
			Str("path", c.Request.URL.Path).
			Str("action", c.Request.Method).
			Int("status", c.Writer.Status()).
			Dur("latency", latency)
	})

	web.Start(router)
	port := os.Getenv("LISTEN_PORT")
	listenPort := ":8080"
	if port != "" {
		listenPort = ":" + port
	}
	log.Info().Msgf("Listening on %s", listenPort)
	log.Fatal().Err(router.Run(listenPort))
}
