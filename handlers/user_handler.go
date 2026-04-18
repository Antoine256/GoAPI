package handlers

import (
	"GoAPI/ressources"
	"GoAPI/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type UserHandler struct {
	logger *zap.Logger
}

func NewUserHandler(logger *zap.Logger) *UserHandler {
	return &UserHandler{logger: logger}
}

func (h *UserHandler) GetUsers(c *gin.Context) {
	// Nécsessite le rôle admin pour voir les autres utilisateurs
	if role := c.GetString("user_role"); role != "admin" {
		h.logger.Error("GetUsers - access denied", zap.String("user_role", role))
		c.JSON(http.StatusForbidden, gin.H{"error": "accès refusé"})
		return
	}
	users, err := services.GetAllUsers()
	if err != nil {
		h.logger.Error("GetUsers - failed to retrieve users", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	h.logger.Info("GetUsers - users retrieved successfully", zap.Int("count", len(users)))
	c.JSON(http.StatusOK, users)
}

func (h *UserHandler) GetCurrentUser(c *gin.Context) {
	userID := c.GetInt("user_id")
	user, err := services.GetUserByID(userID, h.logger)
	if err != nil {
		h.logger.Error("GetCurrentUser - failed to retrieve user", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	h.logger.Info("GetCurrentUser - user retrieved successfully", zap.Int("user_id", user.ID))
	c.JSON(http.StatusOK, user.ToUserResponse())
}

func (h *UserHandler) GetUser(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		h.logger.Error("GetUser - invalid id", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "id invalide"})
		return
	}
	user, err := services.GetUserByID(id, h.logger)
	if err != nil {
		h.logger.Error("GetUser - failed to retrieve user", zap.Error(err))
		c.JSON(http.StatusNotFound, gin.H{"error": "utilisateur introuvable"})
		return
	}

	h.logger.Info("GetUser - user retrieved successfully", zap.Int("user_id", user.ID))
	c.JSON(http.StatusOK, user.ToUserResponse())
}

func (h *UserHandler) CreateUser(c *gin.Context) {
	if role := c.GetString("user_role"); role != "admin" {
		h.logger.Error("CreateUser - access denied", zap.String("user_role", role))
		c.JSON(http.StatusForbidden, gin.H{"error": "accès refusé"})
		return
	}
	var dto ressources.UserCreateRequest
	if err := c.ShouldBindJSON(&dto); err != nil {
		h.logger.Error("CreateUser - invalid request", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	user, err := services.CreateUser(dto, h.logger)
	if err != nil {
		h.logger.Error("CreateUser - failed to create user", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	h.logger.Info("CreateUser - user created successfully", zap.Int("user_id", user.ID))
	c.JSON(http.StatusCreated, user.ToUserResponse())
}

func (h *UserHandler) UpdateUser(c *gin.Context) {

	//Récupération de l'id de l'utilisateur à modifier
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		h.logger.Error("UpdateUser - invalid id", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "id invalide"})
		return
	}

	// Nécsessite le rôle admin pour modifier un utilisateur ou soi même
	if role := c.GetString("user_role"); role != "admin" {
		userID := c.GetInt("user_id")
		if userID != id {
			h.logger.Error("UpdateUser - access denied", zap.String("user_role", role))
			c.JSON(http.StatusForbidden, gin.H{"error": "accès refusé"})
			return
		}
	}

	// Bind de la requête JSON dans un DTO
	var dto ressources.UserUpdateRequest
	h.logger.Info("UpdateUser - received update request", zap.Int("user_id", id), zap.Any("update_fields", c.Request.Body))
	if err := c.ShouldBindJSON(&dto); err != nil {
		h.logger.Error("UpdateUser - invalid request", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	h.logger.Info("UpdateUser - received update request", zap.Int("user_id", id), zap.Any("update_fields", dto))

	// Modification de l'utilisateur via le service
	user, err := services.UpdateUser(id, dto, h.logger)

	if err != nil {
		h.logger.Error("UpdateUser - failed to update user", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	h.logger.Info("UpdateUser - user updated successfully", zap.Int("user_id", user.ID))
	c.JSON(http.StatusOK, user.ToUserResponse())
}

func (h *UserHandler) DeleteUser(c *gin.Context) {

	// Récupération de l'id de l'utilisateur à supprimer
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		h.logger.Error("DeleteUser - invalid id", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "id invalide"})
		return
	}
	// Nécsessite le rôle admin pour supprimer un utilisateur ou soi même
	if role := c.GetString("user_role"); role != "admin" {
		userID := c.GetInt("user_id")
		if userID != id {
			h.logger.Error("DeleteUser - access denied", zap.String("user_role", role))
			c.JSON(http.StatusForbidden, gin.H{"error": "accès refusé"})
			return
		}
	}

	if err := services.DeleteUser(id); err != nil {
		h.logger.Error("DeleteUser - failed to delete user", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	h.logger.Info("DeleteUser - user deleted successfully", zap.Int("user_id", id))
	c.JSON(http.StatusNoContent, nil)
}
