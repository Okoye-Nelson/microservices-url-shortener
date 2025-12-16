package main

import (
	"context"
	"crypto/rand"
	"database/sql"
	"fmt"
	"log"
	"math/big"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/itsbaivab/url-shortener/internal/adapters/cache"
	"github.com/itsbaivab/url-shortener/internal/adapters/repository/postgres"
	"github.com/itsbaivab/url-shortener/internal/core/domain"
	"github.com/itsbaivab/url-shortener/internal/core/services"
	_ "github.com/lib/pq"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type LinkServiceHandler struct {
	linkService *services.LinkService
}

func main() {
	// Database connection
	dbHost := getEnv("DB_HOST", "localhost")
	dbPort := getEnv("DB_PORT", "5432")
	dbUser := getEnv("DB_USER", "postgres")
	dbPassword := getEnv("DB_PASSWORD", "postgres")
	dbName := getEnv("DB_NAME", "urlshortener")

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbPassword, dbName)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	// Test database connection
	if err := db.Ping(); err != nil {
		log.Fatal("Failed to ping database:", err)
	}

	// Redis connection
	redisHost := getEnv("REDIS_HOST", "localhost")
	redisPort := getEnv("REDIS_PORT", "6379")
	redisCache := cache.NewRedisCache(redisHost+":"+redisPort, "", 0)

	// Initialize repository and service
	linkRepo := postgres.NewPostgresLinkRepository(db)
	linkService := services.NewLinkService(linkRepo, redisCache)

	// Initialize handler
	handler := &LinkServiceHandler{
		linkService: linkService,
	}

	// Setup router
	router := gin.Default()

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "healthy"})
	})

	// Metrics endpoint
	router.GET("/metrics", gin.WrapH(promhttp.Handler()))

	// Link endpoints
	router.PUT("/generate", handler.CreateLink)
	router.GET("/links", handler.GetAllLinks)
	router.DELETE("/delete/:id", handler.DeleteLink) // Use path parameter for ID

	// Start server
	port := getEnv("SERVICE_PORT", "8001")
	srv := &http.Server{
		Addr:    ":" + port,
		Handler: router,
	}

	// Graceful shutdown
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	log.Printf("Link Service started on port %s", port)

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down Link Service...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Link Service stopped")
}

func (h *LinkServiceHandler) CreateLink(c *gin.Context) {
	var req struct {
		Long string `json:"long" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate URL
	parsedURL, err := url.ParseRequestURI(req.Long)
	if err != nil || (parsedURL.Scheme != "http" && parsedURL.Scheme != "https") || parsedURL.Host == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "A valid URL with http or https scheme is required"})
		return
	}
	if len(req.Long) > 2048 { // Add a reasonable length limit
		c.JSON(http.StatusBadRequest, gin.H{"error": "URL cannot exceed 2048 characters"})
		return
	}

	var link domain.Link
	var errCreate error
	maxRetries := 3 // Number of times to retry on ID collision

	for i := 0; i < maxRetries; i++ {
		link = domain.Link{
			Id:          generateShortURLID(8),
			OriginalURL: req.Long,
			CreatedAt:   time.Now(),
		}
		errCreate = h.linkService.Create(c.Request.Context(), link)
		if errCreate == nil {
			break // Success
		}
	}

	if errCreate != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create unique link, please try again"})
		return
	}

	// Return the full link object so frontend can display id, original_url, created_at
	c.JSON(http.StatusOK, link)
}

func (h *LinkServiceHandler) GetAllLinks(c *gin.Context) {
	links, err := h.linkService.GetAll(c.Request.Context())
	if err != nil {
		log.Printf("Error getting all links: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get links"})
		return
	}

	c.JSON(http.StatusOK, links)
}

func (h *LinkServiceHandler) DeleteLink(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID parameter is required"})
		return
	}

	if err := h.linkService.Delete(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func generateShortURLID(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	result := make([]byte, length)
	for i := range result {
		charIndex, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			// Fallback to a less random but safe method if crypto/rand fails
			return "fallback" + fmt.Sprintf("%d", time.Now().UnixNano())
		}
		result[i] = charset[charIndex.Int64()]
	}
	return string(result)
}
