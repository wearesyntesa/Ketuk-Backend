package handlers

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"time"

	"ketukApps/internal/models"
	"ketukApps/internal/services"
	"ketukApps/internal/utils"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// AuthHandler handles authentication-related requests
type AuthHandler struct {
	db          *gorm.DB
	googleOAuth *services.GoogleOAuthService
	stateStore  map[string]time.Time // Simple in-memory store for OAuth state (use Redis in production)
}

// NewAuthHandler creates a new AuthHandler
func NewAuthHandler(db *gorm.DB, googleOAuth *services.GoogleOAuthService) *AuthHandler {
	return &AuthHandler{
		db:          db,
		googleOAuth: googleOAuth,
		stateStore:  make(map[string]time.Time),
	}
}

// LoginRequest represents the login request body
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email" example:"user@example.com"`
	Password string `json:"password" binding:"required" example:"password123"`
}

// RegisterRequest represents the registration request body
type RegisterRequest struct {
	Email    string `json:"email" binding:"required,email" example:"user@example.com"`
	Password string `json:"password" binding:"required,min=6" example:"password123"`
	Name     string `json:"name" binding:"required" example:"John Doe"`
}

// LoginResponse represents the login response
type LoginResponse struct {
	Token        string      `json:"token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
	RefreshToken string      `json:"refresh_token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
	User         models.User `json:"user"`
}

// Login godoc
// @Summary Login user
// @Description Authenticate user and return JWT token
// @Tags auth
// @Accept json
// @Produce json
// @Param request body LoginRequest true "Login credentials"
// @Success 200 {object} models.APIResponse{data=LoginResponse}
// @Failure 400 {object} models.APIResponse
// @Failure 401 {object} models.APIResponse
// @Router /api/auth/v1/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "Invalid request",
			Error:   err.Error(),
		})
		return
	}

	// Find user by email
	var user models.User
	if err := h.db.Where("email = ?", req.Email).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusUnauthorized, models.APIResponse{
				Success: false,
				Message: "Invalid credentials",
				Error:   "Email or password is incorrect",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Message: "Failed to authenticate",
			Error:   err.Error(),
		})
		return
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, models.APIResponse{
			Success: false,
			Message: "Invalid credentials",
			Error:   "Email or password is incorrect",
		})
		return
	}

	// Generate JWT token
	token, err := utils.GenerateToken(user.ID, user.Email, user.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Message: "Failed to generate token",
			Error:   err.Error(),
		})
		return
	}

	// Generate refresh token
	refreshToken, err := utils.GenerateRefreshToken(user.ID, user.Email, user.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Message: "Failed to generate refresh token",
			Error:   err.Error(),
		})
		return
	}

	// Remove password from response
	user.Password = ""

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "Login successful",
		Data: LoginResponse{
			Token:        token,
			RefreshToken: refreshToken,
			User:         user,
		},
	})
}

// Register godoc
// @Summary Register new user
// @Description Create a new user account
// @Tags auth
// @Accept json
// @Produce json
// @Param request body RegisterRequest true "Registration details"
// @Success 201 {object} models.APIResponse{data=LoginResponse}
// @Failure 400 {object} models.APIResponse
// @Failure 409 {object} models.APIResponse
// @Router /api/auth/v1/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "Invalid request",
			Error:   err.Error(),
		})
		return
	}

	// Check if user already exists
	var existingUser models.User
	if err := h.db.Where("email = ?", req.Email).First(&existingUser).Error; err == nil {
		c.JSON(http.StatusConflict, models.APIResponse{
			Success: false,
			Message: "User already exists",
			Error:   "An account with this email already exists",
		})
		return
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Message: "Failed to create user",
			Error:   err.Error(),
		})
		return
	}

	// Check if this is the first user (make them admin)
	var userCount int64
	h.db.Model(&models.User{}).Count(&userCount)

	// Determine role - first user becomes admin
	userRole := "user"
	if userCount == 0 {
		userRole = "admin"
	}

	// Create user
	user := models.User{
		Email:    req.Email,
		Name:     req.Name,
		Password: string(hashedPassword),
		Role:     userRole,
	}

	if err := h.db.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Message: "Failed to create user",
			Error:   err.Error(),
		})
		return
	}

	// Generate JWT token
	token, err := utils.GenerateToken(user.ID, user.Email, user.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Message: "Failed to generate token",
			Error:   err.Error(),
		})
		return
	}

	// Generate refresh token
	refreshToken, err := utils.GenerateRefreshToken(user.ID, user.Email, user.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Message: "Failed to generate refresh token",
			Error:   err.Error(),
		})
		return
	}

	// Remove password from response
	user.Password = ""

	c.JSON(http.StatusCreated, models.APIResponse{
		Success: true,
		Message: "Registration successful",
		Data: LoginResponse{
			Token:        token,
			RefreshToken: refreshToken,
			User:         user,
		},
	})
}

// Me godoc
// @Summary Get current user
// @Description Get the currently authenticated user's information
// @Tags auth
// @Produce json
// @Security BearerAuth
// @Success 200 {object} models.APIResponse{data=models.User}
// @Failure 401 {object} models.APIResponse
// @Router /api/auth/v1/me [get]
func (h *AuthHandler) Me(c *gin.Context) {
	// Get user from context (set by AuthRequired middleware)
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, models.APIResponse{
			Success: false,
			Message: "Unauthorized",
			Error:   "User not authenticated",
		})
		return
	}

	userData := user.(models.User)
	userData.Password = "" // Remove password from response

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "User retrieved successfully",
		Data:    userData,
	})
}

// RefreshToken godoc
// @Summary Refresh access token
// @Description Get a new access token using refresh token
// @Tags auth
// @Accept json
// @Produce json
// @Param request body models.RefreshTokenRequest true "Refresh token"
// @Success 200 {object} models.APIResponse{data=LoginResponse}
// @Failure 400 {object} models.APIResponse
// @Failure 401 {object} models.APIResponse
// @Router /api/auth/v1/refresh [post]
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	var req models.RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "Invalid request",
			Error:   err.Error(),
		})
		return
	}

	// Validate refresh token
	claims, err := utils.ValidateToken(req.RefreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, models.APIResponse{
			Success: false,
			Message: "Invalid refresh token",
			Error:   err.Error(),
		})
		return
	}

	// Get user from database
	var user models.User
	if err := h.db.First(&user, claims.UserID).Error; err != nil {
		c.JSON(http.StatusUnauthorized, models.APIResponse{
			Success: false,
			Message: "User not found",
			Error:   err.Error(),
		})
		return
	}

	// Generate new tokens
	token, err := utils.GenerateToken(user.ID, user.Email, user.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Message: "Failed to generate token",
			Error:   err.Error(),
		})
		return
	}

	refreshToken, err := utils.GenerateRefreshToken(user.ID, user.Email, user.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Message: "Failed to generate refresh token",
			Error:   err.Error(),
		})
		return
	}

	// Remove password from response
	user.Password = ""

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "Token refreshed successfully",
		Data: LoginResponse{
			Token:        token,
			RefreshToken: refreshToken,
			User:         user,
		},
	})
}

// generateState generates a random state for OAuth
func (h *AuthHandler) generateState() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	state := base64.URLEncoding.EncodeToString(b)

	// Store state with expiration (10 minutes)
	h.stateStore[state] = time.Now().Add(10 * time.Minute)

	// Clean up expired states
	go h.cleanupExpiredStates()

	return state, nil
}

// validateState validates OAuth state and removes it
func (h *AuthHandler) validateState(state string) bool {
	expiry, exists := h.stateStore[state]
	if !exists {
		return false
	}

	// Check if expired
	if time.Now().After(expiry) {
		delete(h.stateStore, state)
		return false
	}

	// Remove state after use
	delete(h.stateStore, state)
	return true
}

// cleanupExpiredStates removes expired states from store
func (h *AuthHandler) cleanupExpiredStates() {
	now := time.Now()
	for state, expiry := range h.stateStore {
		if now.After(expiry) {
			delete(h.stateStore, state)
		}
	}
}

// GoogleLogin godoc
// @Summary Initiate Google OAuth login
// @Description Redirects user to Google OAuth consent screen
// @Tags auth
// @Produce json
// @Success 200 {object} models.APIResponse{data=map[string]string}
// @Failure 500 {object} models.APIResponse
// @Router /api/auth/v1/google/login [get]
func (h *AuthHandler) GoogleLogin(c *gin.Context) {
	// Generate random state for CSRF protection
	state, err := h.generateState()
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Message: "Failed to generate state",
			Error:   err.Error(),
		})
		return
	}

	// Get authorization URL
	authURL := h.googleOAuth.GetAuthURL(state)

	// Debug log
	log.Printf("Generated Google OAuth URL: %s", authURL)

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "Redirect to Google OAuth",
		Data: map[string]string{
			"auth_url": authURL,
			"state":    state,
		},
	})
}

// GoogleCallback godoc
// @Summary Google OAuth callback
// @Description Handle callback from Google OAuth and authenticate user
// @Tags auth
// @Produce json
// @Param code query string true "Authorization code from Google"
// @Param state query string true "State for CSRF protection"
// @Success 200 {object} models.APIResponse{data=LoginResponse}
// @Failure 400 {object} models.APIResponse
// @Failure 401 {object} models.APIResponse
// @Failure 500 {object} models.APIResponse
// @Router /api/auth/v1/google/callback [get]
func (h *AuthHandler) GoogleCallback(c *gin.Context) {
	// Get code and state from query params
	code := c.Query("code")
	state := c.Query("state")

	if code == "" {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "Invalid request",
			Error:   "Authorization code is required",
		})
		return
	}

	// Validate state for CSRF protection
	if !h.validateState(state) {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "Invalid request",
			Error:   "Invalid or expired state",
		})
		return
	}

	// Exchange code for token
	ctx := context.Background()
	token, err := h.googleOAuth.ExchangeCode(ctx, code)
	if err != nil {
		c.JSON(http.StatusUnauthorized, models.APIResponse{
			Success: false,
			Message: "Failed to exchange code",
			Error:   err.Error(),
		})
		return
	}

	// Get user info from Google
	googleUser, err := h.googleOAuth.GetUserInfo(ctx, token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, models.APIResponse{
			Success: false,
			Message: "Failed to get user info",
			Error:   err.Error(),
		})
		return
	}

	// Find or create user in database
	var user models.User
	result := h.db.Where("email = ?", googleUser.Email).First(&user)

	if result.Error == gorm.ErrRecordNotFound {
		// Check if this is the first user (make them admin)
		var userCount int64
		h.db.Model(&models.User{}).Count(&userCount)

		// Determine role - first user becomes admin
		userRole := "user"
		if userCount == 0 {
			userRole = "admin"
		}

		// Create new user
		user = models.User{
			Email:     googleUser.Email,
			Name:      googleUser.Name,
			GoogleSub: googleUser.ID,
			Role:      userRole,
		}

		if err := h.db.Create(&user).Error; err != nil {
			c.JSON(http.StatusInternalServerError, models.APIResponse{
				Success: false,
				Message: "Failed to create user",
				Error:   err.Error(),
			})
			return
		}
	} else if result.Error != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Message: "Database error",
			Error:   result.Error.Error(),
		})
		return
	} else {
		// Update existing user info
		user.Name = googleUser.Name
		user.GoogleSub = googleUser.ID

		if err := h.db.Save(&user).Error; err != nil {
			c.JSON(http.StatusInternalServerError, models.APIResponse{
				Success: false,
				Message: "Failed to update user",
				Error:   err.Error(),
			})
			return
		}
	}

	// Generate JWT tokens
	jwtToken, err := utils.GenerateToken(user.ID, user.Email, user.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Message: "Failed to generate token",
			Error:   err.Error(),
		})
		return
	}

	refreshToken, err := utils.GenerateRefreshToken(user.ID, user.Email, user.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Message: "Failed to generate refresh token",
			Error:   err.Error(),
		})
		return
	}

	// Remove password from response
	user.Password = ""

	// Check if redirect_uri is provided (for frontend redirect)
	redirectURI := c.Query("redirect_uri")
	if redirectURI != "" {
		// Redirect to frontend with tokens in URL
		params := url.Values{}
		params.Add("token", jwtToken)
		params.Add("refresh_token", refreshToken)

		redirectURL := fmt.Sprintf("%s?%s", redirectURI, params.Encode())
		c.Redirect(http.StatusFound, redirectURL)
		return
	}

	// Default: return JSON response
	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "Login successful",
		Data: LoginResponse{
			Token:        jwtToken,
			RefreshToken: refreshToken,
			User:         user,
		},
	})
}
