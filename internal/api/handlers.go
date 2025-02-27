package api

import (
	"encoding/json"
	"fmt"
	"grok-proxy/config"
	"grok-proxy/internal/client"
	"grok-proxy/internal/cookie"
	"grok-proxy/internal/models"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// Handler contains the handlers for the API
type Handler struct {
	grokClient    *client.GrokClient
	cookieManager *cookie.Manager
	cfg           *config.Config
}

// NewHandler creates a new handler
func NewHandler(grokClient *client.GrokClient, cookieManager *cookie.Manager, cfg *config.Config) *Handler {
	return &Handler{
		grokClient:    grokClient,
		cookieManager: cookieManager,
		cfg:           cfg,
	}
}

// AuthMiddleware authenticates requests
func (h *Handler) AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		auth := c.GetHeader("Authorization")
		if auth == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": map[string]interface{}{
					"message": "Missing API key",
					"type":    "invalid_request_error",
					"param":   nil,
					"code":    "invalid_api_key",
				},
			})
			return
		}

		// Extract bearer token
		parts := strings.SplitN(auth, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": map[string]interface{}{
					"message": "Invalid API key format",
					"type":    "invalid_request_error",
					"param":   nil,
					"code":    "invalid_api_key",
				},
			})
			return
		}

		apiKey := parts[1]
		if !h.cfg.IsValidAPIKey(apiKey) {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": map[string]interface{}{
					"message": "Invalid API key",
					"type":    "invalid_request_error",
					"param":   nil,
					"code":    "invalid_api_key",
				},
			})
			return
		}

		c.Next()
	}
}

// GetModels handles the GET /v1/models endpoint
func (h *Handler) GetModels(c *gin.Context) {
	c.JSON(http.StatusOK, models.NewModelList())
}

// HandleChatCompletions handles the POST /v1/chat/completions endpoint
func (h *Handler) HandleChatCompletions(c *gin.Context) {
	var req models.OpenAIRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Convert messages to string
	messagesStr := fmt.Sprintf("%v", req.Messages)

	if req.Stream {
		// Handle streaming response
		c.Header("Content-Type", "text/event-stream")
		c.Header("Cache-Control", "no-cache")
		c.Header("Connection", "keep-alive")

		// Create context that's canceled when client disconnects
		ctx := c.Request.Context()

		// Start streaming
		tokenCh, errCh := h.grokClient.StreamRequest(ctx, messagesStr, req.Model)

		c.Stream(func(w io.Writer) bool {
			select {
			case token, ok := <-tokenCh:
				if !ok {
					// Send end message
					endData := models.StreamResponseChunk{
						ID:      "grok-proxy-end",
						Object:  "chat.completion.chunk",
						Created: time.Now().Unix(),
						Model:   req.Model,
						Choices: []struct {
							Delta        map[string]string `json:"delta"`
							Index        int               `json:"index"`
							FinishReason interface{}       `json:"finish_reason"`
						}{
							{
								Delta:        map[string]string{},
								Index:        0,
								FinishReason: "stop",
							},
						},
					}

					endBytes, _ := json.Marshal(endData)
					// Fix: Remove the data: prefix from within the string
					fmt.Fprintf(w, "data: %s\n\n", endBytes)
					fmt.Fprintf(w, "data: [DONE]\n\n")
					return false
				}

				// Send token
				data := models.StreamResponseChunk{
					ID:      "grok-proxy",
					Object:  "chat.completion.chunk",
					Created: time.Now().Unix(),
					Model:   req.Model,
					Choices: []struct {
						Delta        map[string]string `json:"delta"`
						Index        int               `json:"index"`
						FinishReason interface{}       `json:"finish_reason"`
					}{
						{
							Delta:        map[string]string{"content": token},
							Index:        0,
							FinishReason: nil,
						},
					},
				}

				bytes, _ := json.Marshal(data)
				// Fix: Remove the data: prefix from within the string
				fmt.Fprintf(w, "data: %s\n\n", bytes)
				return true

			case err, ok := <-errCh:
				if ok && err != nil {
					log.Printf("Stream error: %v", err)
					// Fix: Use direct writing instead of SSEvent
					fmt.Fprintf(w, "data: {\"error\": \"%s\"}\n\n", err.Error())
				}
				return false

			case <-ctx.Done():
				return false
			}
		})
	} else {
		// Handle non-streaming response
		tokens, err := h.grokClient.Execute(c.Request.Context(), messagesStr, req.Model)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		response := gin.H{
			"id":      "grok_proxy",
			"object":  "chat.completion",
			"created": time.Now().Unix(),
			"model":   req.Model,
			"choices": []gin.H{
				{
					"message": gin.H{
						"role":    "assistant",
						"content": tokens,
					},
					"finish_reason": "stop",
					"index":         0,
				},
			},
		}

		c.JSON(http.StatusOK, response)
	}
}

// RegisterRoutes registers API routes
func (h *Handler) RegisterRoutes(router *gin.Engine) {
	v1 := router.Group("/v1")
	v1.Use(h.AuthMiddleware())
	{
		v1.GET("/models", h.GetModels)
		v1.POST("/chat/completions", h.HandleChatCompletions)
	}
}
