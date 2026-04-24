package services

import (
	"GoAPI/repository"
	"GoAPI/ressources"
	"errors"
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"math/rand"
	"mime/multipart"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/HugoSmits86/nativewebp"
	"go.breu.io/ulid"
	"go.uber.org/zap"
)

func SaveMedia(file multipart.File, mediaType string, mediaName string, logger *zap.Logger) (*ressources.Media, error) {
	buf := make([]byte, 512)
	_, err := file.Read(buf)
	if err != nil {
		logger.Error("Erreur lecture fichier", zap.Error(err))
		return nil, err
	}

	// détecter MIME sans casser le flux
	format := strings.ToLower(http.DetectContentType(buf))

	// reset curseur
	_, err = file.Seek(0, 0)
	if err != nil {
		return nil, err
	}

	if format != "image/jpeg" && format != "image/png" && format != "video/mp4" && format != "video/webm" && format != "video/quicktime" {
		logger.Warn("Format d'image non supporté", zap.String("format", format))
		return nil, errors.New("format non supporté")
	}

	var ext string
	switch format {
	case "video/mp4":
		ext = ".mp4"
	case "video/webm":
		ext = ".webm"
	case "video/quicktime": // .mov
		ext = ".mov"
	case "image/jpeg":
		ext = ".jpg"
	case "image/png":
		ext = ".png"
	default:
		logger.Error("Format d'image non supporté", zap.String("format", format))
		return nil, errors.New("format non supporté")
	}

	// Génération d'un nom sécurisé (UUID)
	path := os.Getenv("PICTURES_DIRECTORIES")
	if path == "" {
		logger.Error("Variable d'environnement non définie: PICTURES_DIRECTORIES")
		return nil, errors.New("configuration du serveur incorrecte")
	}

	t := time.Now()
	entropy := ulid.Monotonic(rand.New(rand.NewSource(t.UnixNano())), 0)
	newID, err := ulid.New(ulid.Timestamp(t), entropy)
	if err != nil {
		logger.Error("Erreur génération nom fichier", zap.Error(err))
		return nil, err
	}
	newName := newID.String()

	// Création du fichier physique de full quality
	full_quality_path := filepath.Join(path, "full_quality", newName+ext)
	dst, err := os.Create(full_quality_path)
	if err != nil {
		logger.Error("Erreur création fichier", zap.Error(err))
		return nil, err
	}
	defer dst.Close()

	// Remise du curseur au début du fichier pour la copie
	file.Seek(0, io.SeekStart)

	// Copie des bits
	written, err := io.Copy(dst, file)
	if err != nil {
		logger.Error("Erreur copie fichier", zap.Error(err))
		return nil, err
	}

	//Enregistrement en base de données
	mediadto := ressources.Media{
		MediaName:       filepath.Base(mediaName),
		MediaType:       mediaType,
		UUID:            newName,
		FullQualityPath: full_quality_path,
		ThumbPath:       "",
		Size:            written,
		CompressSize:    0,
	}

	media, err := repository.CreateMedia(mediadto)
	if err != nil {
		logger.Error("Erreur enregistrement base de données", zap.Error(err))
		return nil, err
	}

	// Compression de l'image et création du thumbnail dabs le dossier webp
	go compressWithRetry(media, file, logger, 3)

	return &media, nil

}

func GetMedias(page int, pageSize int) ([]ressources.Media, error) {
	medias, err := repository.GetMedias(page, pageSize)
	if err != nil {
		return []ressources.Media{}, err
	}
	return medias, nil
}

///************************************************************///

func compressWithRetry(media ressources.Media, file multipart.File, logger *zap.Logger, retries int) error {
	for i := 0; i < retries; i++ {
		err := error(nil)
		if media.MediaType != "picture" {
			err = compressVideoToFFMPEG(media, logger)
		} else {
			err = compressPictureToWebP(media, file, logger)
		}
		if err == nil {
			return nil
		}
		logger.Warn("Compression échouée, tentative de retry", zap.Int("attempt", i+1), zap.Error(err))
		time.Sleep(2 * time.Second) // Attente avant le retry
	}
	return errors.New("échec de la compression après plusieurs tentatives")
}

func compressPictureToWebP(media ressources.Media, file multipart.File, logger *zap.Logger) (err error) {

	path := os.Getenv("PICTURES_DIRECTORIES")
	if path == "" {
		logger.Error("Variable d'environnement non définie: PICTURES_DIRECTORIES")
		return
	}

	thumbPath := filepath.Join(path, "thumbnail", media.UUID+".webp")

	buf := make([]byte, 512)
	file.Read(buf)
	fmt.Println(http.DetectContentType(buf))
	file.Seek(0, 0)

	img, _, err := image.Decode(file)
	if img == nil {
		logger.Error("Erreur décodage image", zap.Error(err))
		return
	}

	out, _ := os.Create(thumbPath)
	defer out.Close()

	err = nativewebp.Encode(out, img, nil)
	if err != nil {
		logger.Error("Erreur création options WebP", zap.Error(err))
		return
	}

	info, err := os.Stat(thumbPath)
	if err != nil {
		logger.Error("Erreur stat webp", zap.Error(err))
		return
	}

	_, err = repository.UpdateMedia(ressources.Media{
		ID:           media.ID,
		ThumbPath:    thumbPath,
		CompressSize: info.Size(), // Taille du fichier compressé
	})
	if err != nil {
		logger.Error("Erreur mise à jour base de données", zap.Error(err))
		return
	}

	logger.Info("Thumbnail WebP généré", zap.String("path", thumbPath))

	return nil
}

func compressVideoToFFMPEG(media ressources.Media, logger *zap.Logger) error {
	path := os.Getenv("PICTURES_DIRECTORIES")
	if path == "" {
		logger.Error("Variable d'environnement non définie: PICTURES_DIRECTORIES")
		return errors.New("configuration du serveur incorrecte")
	}

	// Implémentation de la compression vidéo avec FFMPEG
	inputFile, err := os.Open(media.FullQualityPath)
	if err != nil {
		logger.Error("Erreur ouverture fichier vidéo", zap.Error(err))
		return err
	}

	err = inputFile.Close()
	if err != nil {
		return err
	}

	thumbPath := filepath.Join(path, "thumbnail", media.UUID+".jpg")
	compressedPath := filepath.Join(path, "compressed", media.UUID+"_compressed.mp4")

	cmd1 := exec.Command(
		"ffmpeg",
		"-i", media.FullQualityPath,
		"-ss", "00:00:01",
		"-vframes", "1",
		"-q:v", "2",
		thumbPath,
	)

	cmd2 := exec.Command(
		"ffmpeg",
		"-i", media.FullQualityPath,
		"-vcodec", "libx264",
		"-crf", "28",
		"-preset", "fast",
		"-acodec", "aac",
		"-movflags", "+faststart",
		compressedPath,
	)

	err = cmd1.Run()
	if err != nil {
		logger.Error("Erreur exécution commande FFMPEG", zap.Error(err))
		return errors.New("échec de l'exécution de la commande FFMPEG 1")
	}

	err = cmd2.Run()
	if err != nil {
		logger.Error("Erreur exécution commande FFMPEG", zap.Error(err))
		return errors.New("échec de l'exécution de la commande FFMPEG 2")
	}

	info, err := os.Stat(compressedPath)
	if err != nil {
		return err
	}

	_, err = repository.UpdateMedia(ressources.Media{
		ID:           media.ID,
		ThumbPath:    thumbPath,
		CompressPath: compressedPath,
		CompressSize: info.Size(), // Taille du fichier compressé
	})
	if err != nil {
		logger.Error("Erreur mise à jour base de données", zap.Error(err))
		return err
	}

	return nil
}
