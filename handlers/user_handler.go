package handlers

import (
	"GoAPI/ressources"
	"GoAPI/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetUsers(c *gin.Context) {
	// Nécsessite le rôle admin pour voir les autres utilisateurs
	if role := c.GetString("user_role"); role != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "accès refusé"})
		return
	}
	users, err := services.GetAllUsers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, users)
}

func GetCurrentUser(c *gin.Context) {
	userID := c.GetInt("user_id")
	user, err := services.GetUserByID(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, ressources.ToUserPublicDTO(user))
}

func GetUser(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id invalide"})
		return
	}
	user, err := services.GetUserByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "utilisateur introuvable"})
		return
	}
	c.JSON(http.StatusOK, ressources.ToUserPublicDTO(user))
}

func CreateUser(c *gin.Context) {
	if role := c.GetString("user_role"); role != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "accès refusé"})
		return
	}
	var dto ressources.UserCreateDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	user, err := services.CreateUser(dto)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, ressources.ToUserPublicDTO(user))
}

func UpdateUser(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id invalide"})
		return
	}
	// Nécsessite le rôle admin pour modifier un utilisateur ou soi même
	if role := c.GetString("user_role"); role != "admin" {
		userID := c.GetInt("user_id")
		if userID != id {
			c.JSON(http.StatusForbidden, gin.H{"error": "accès refusé"})
			return
		}
	}
	var dto ressources.UserUpdateDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := services.UpdateUser(id, dto)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, ressources.ToUserPublicDTO(user))
}

func DeleteUser(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id invalide"})
		return
	}
	// Nécsessite le rôle admin pour supprimer un utilisateur ou soi même
	if role := c.GetString("user_role"); role != "admin" {
		userID := c.GetInt("user_id")
		if userID != id {
			c.JSON(http.StatusForbidden, gin.H{"error": "accès refusé"})
			return
		}
	}

	if err := services.DeleteUser(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}
