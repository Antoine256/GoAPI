package handlers

import (
	"GoAPI/ressources"
	"GoAPI/services"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type AuthHandler struct {
	logger *zap.Logger
}

func NewAuthHandler(logger *zap.Logger) *AuthHandler {
	return &AuthHandler{logger: logger}
}

func (h *AuthHandler) setAuthCookies(c *gin.Context, tokens ressources.TokenResponse) {
	// Refresh token en cookie httpOnly
	c.SetCookie(
		"refresh_token",     // nom
		tokens.RefreshToken, // valeur
		7*24*3600,           // durée en secondes (7 jours)
		"/",                 // path
		"",                  // domaine
		false,               // secure (true si en production avec HTTPS)
		true,                // httpOnly
	)
}

func (h *AuthHandler) Login(c *gin.Context) {
	var dto ressources.LoginRequest
	// Vérifie que le JSON de la requête correspond au DTO attendu
	// Si c'est le cas, il est automatiquement converti en struct Go
	if err := c.ShouldBindJSON(&dto); err != nil {
		h.logger.Error("Login - invalid request", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// Appelle la logique métier du service d'authentification
	tokens, err := services.Login(dto, h.logger)
	if err != nil {
		h.logger.Error("Login - authentication failed", zap.Error(err))
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	h.setAuthCookies(c, tokens)
	h.logger.Info("Login successful", zap.String("email", dto.Email))
	c.JSON(http.StatusCreated, gin.H{
		"access_token": tokens.AccessToken,
	})
}

func (h *AuthHandler) Register(c *gin.Context) {
	var dto ressources.RegisterRequest
	if err := c.ShouldBindJSON(&dto); err != nil {
		h.logger.Error("Register - invalid request", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	tokens, err := services.Register(dto, h.logger)
	if err != nil {
		h.logger.Error("Register - registration failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	h.setAuthCookies(c, tokens)
	h.logger.Info("Registration successful", zap.String("email", dto.Email))
	c.JSON(http.StatusCreated, gin.H{
		"access_token": tokens.AccessToken,
	})
}

func (h *AuthHandler) RefreshToken(c *gin.Context) {
	var refreshToken, err = c.Cookie("refresh_token")
	if err != nil {
		h.logger.Error("RefreshToken - failed to get refresh token", zap.Error(err))
		c.JSON(http.StatusUnauthorized, gin.H{"error": "token non trouvé"})
		return
	}

	tokens, err := services.RefreshToken(refreshToken, h.logger)
	if err != nil {
		h.logger.Error("RefreshToken - failed to refresh token", zap.String("refresh_token", refreshToken), zap.Error(err))
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	h.setAuthCookies(c, tokens)
	h.logger.Info("Token refreshed successfully", zap.String("refresh_token", refreshToken))
	c.JSON(http.StatusOK, gin.H{
		"access_token": tokens.AccessToken,
	})
}

func (h *AuthHandler) Logout(c *gin.Context) {
	var refreshToken, err = c.Cookie("refresh_token")
	if err != nil {
		h.logger.Error("Logout - failed to get refresh token", zap.Error(err))
		c.JSON(http.StatusUnauthorized, gin.H{"error": "token non trouvé"})
		return
	}
	if err := services.Logout(refreshToken); err != nil {
		h.logger.Error("Logout - failed to logout", zap.String("refresh_token", refreshToken), zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.SetCookie("refresh_token", "", -1, "/", "", false, true)
	h.logger.Info("Logout successful", zap.String("refresh_token", refreshToken))
	c.JSON(http.StatusOK, gin.H{"message": "déconnecté"})
}
