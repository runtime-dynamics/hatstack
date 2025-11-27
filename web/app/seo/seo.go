package seo

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// RobotsTxtHandler serves the robots.txt file for web crawlers
func RobotsTxtHandler(c *gin.Context) {
	// Get the base URL from the request
	scheme := "https"
	if c.Request.TLS == nil {
		scheme = "http"
	}
	baseURL := scheme + "://" + c.Request.Host

	robotsTxt := `User-agent: *
Allow: /
Disallow: /app/
Disallow: /api/
Disallow: /desktop-login

Sitemap: ` + baseURL + `/sitemap.xml
`
	c.Header("Content-Type", "text/plain")
	c.String(http.StatusOK, robotsTxt)
}

// SitemapXMLHandler serves the sitemap.xml file for search engines
func SitemapXMLHandler(c *gin.Context) {
	// Get the base URL from the request
	scheme := "https"
	if c.Request.TLS == nil {
		scheme = "http"
	}
	baseURL := scheme + "://" + c.Request.Host

	// Get the current time for lastmod
	now := time.Now().Format("2006-01-02")

	// Build sitemap XML
	sitemap := `<?xml version="1.0" encoding="UTF-8"?>
<urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">
  <url>
    <loc>` + baseURL + `/</loc>
    <lastmod>` + now + `</lastmod>
    <changefreq>weekly</changefreq>
    <priority>1.0</priority>
  </url>
  <url>
    <loc>` + baseURL + `/features</loc>
    <lastmod>` + now + `</lastmod>
    <changefreq>weekly</changefreq>
    <priority>0.8</priority>
  </url>
  <url>
    <loc>` + baseURL + `/about</loc>
    <lastmod>` + now + `</lastmod>
    <changefreq>monthly</changefreq>
    <priority>0.7</priority>
  </url>
  <url>
    <loc>` + baseURL + `/help</loc>
    <lastmod>` + now + `</lastmod>
    <changefreq>weekly</changefreq>
    <priority>0.6</priority>
  </url>
  <url>
    <loc>` + baseURL + `/contact</loc>
    <lastmod>` + now + `</lastmod>
    <changefreq>monthly</changefreq>
    <priority>0.6</priority>
  </url>
  <url>
    <loc>` + baseURL + `/privacy</loc>
    <lastmod>` + now + `</lastmod>
    <changefreq>yearly</changefreq>
    <priority>0.3</priority>
  </url>
  <url>
    <loc>` + baseURL + `/terms</loc>
    <lastmod>` + now + `</lastmod>
    <changefreq>yearly</changefreq>
    <priority>0.3</priority>
  </url>
  <url>
    <loc>` + baseURL + `/cookies</loc>
    <lastmod>` + now + `</lastmod>
    <changefreq>yearly</changefreq>
    <priority>0.3</priority>
  </url>
</urlset>`

	c.Header("Content-Type", "application/xml")
	c.String(http.StatusOK, sitemap)
}
