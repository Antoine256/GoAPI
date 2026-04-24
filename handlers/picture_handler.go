package handlers

import (
	"GoAPI/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type PictureHandler struct {
	logger *zap.Logger
}

func NewMediaHandler(logger *zap.Logger) *PictureHandler {
	return &PictureHandler{logger: logger}
}

func (h *PictureHandler) GetMedias(c *gin.Context) {
	page, err1 := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, err2 := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	if err1 != nil || err2 != nil {
		c.JSON(400, gin.H{"error": "Invalid page or page_size"})
		return
	}

	medias, err := services.GetMedias(page, pageSize)
	if err != nil {
		h.logger.Error("GetMedias - failed to retrieve medias", zap.Error(err))
		c.JSON(500, gin.H{"error": "Erreur récupération médias"})
		return
	}
	c.JSON(200, medias)
}

func (h *PictureHandler) Upload(c *gin.Context) {

	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, 10<<20) // Limite à 10MB
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(400, gin.H{"error": "file is required"})
		return
	}

	mediaType := c.PostForm("type")
	mediaName := c.PostForm("picture_name")

	openFile, err := file.Open()
	if err != nil {
		h.logger.Error("Upload - failed to open file", zap.Error(err))
		c.JSON(500, gin.H{"error": "Erreur ouverture fichier"})
		return
	}
	defer openFile.Close()

	media, err := services.SaveMedia(openFile, mediaType, mediaName, h.logger)
	if err != nil {
		h.logger.Error("Upload - failed to process upload", zap.Error(err))
		c.JSON(500, gin.H{"error": "Erreur traitement fichier"})
		return
	}

	h.logger.Info("Upload successful", zap.String("filename", mediaName), zap.Int64("size", media.Size))
	c.JSON(201, media)
}
