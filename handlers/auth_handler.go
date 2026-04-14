package handlers

import (
	"GoAPI/ressources"
	"GoAPI/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

func setAuthCookies(c *gin.Context, tokens ressources.TokenResponseDTO) {
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

func Login(c *gin.Context) {
	var dto ressources.LoginDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	tokens, err := services.Login(dto)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	setAuthCookies(c, tokens)
	c.JSON(http.StatusCreated, gin.H{
		"access_token": tokens.AccessToken,
	})
}

func Register(c *gin.Context) {
	var dto ressources.RegisterDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	tokens, err := services.Register(dto)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	setAuthCookies(c, tokens)
	c.JSON(http.StatusCreated, gin.H{
		"access_token": tokens.AccessToken,
	})
}

func RefreshToken(c *gin.Context) {
	var dto ressources.RefreshTokenDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	tokens, err := services.RefreshToken(dto)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	setAuthCookies(c, tokens)
	c.JSON(http.StatusOK, gin.H{
		"access_token": tokens.AccessToken,
	})
}

func Logout(c *gin.Context) {
	var dto ressources.RefreshTokenDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := services.Logout(dto.RefreshToken); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.SetCookie("refresh_token", "", -1, "/", "", false, true)
	c.JSON(http.StatusOK, gin.H{"message": "déconnecté"})
}
