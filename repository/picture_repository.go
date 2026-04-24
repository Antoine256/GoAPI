package repository

import (
	"GoAPI/database"
	"GoAPI/ressources"
)

func GetMedias(page int, pageSize int) ([]ressources.Media, error) {
	var medias []ressources.Media
	rows, err := database.DB.Query(`
	SELECT id, media_name, media_type, uuid, full_quality_path, thumb_path, compress_path, size, compress_size, created_at 
	FROM medias 
	ORDER BY created_at DESC
	LIMIT $1 OFFSET $2
	`, pageSize, (page-1)*pageSize)
	if err != nil {
		return medias, err
	}
	defer rows.Close()

	i := 0
	for rows.Next() {
		newRow := ressources.Media{}
		err := rows.Scan(&newRow.ID, &newRow.MediaName, &newRow.MediaType, &newRow.UUID, &newRow.FullQualityPath, &newRow.ThumbPath, &newRow.CompressPath, &newRow.Size, &newRow.CompressSize, &newRow.CreatedAt)
		if err != nil {
			return medias, err
		}
		medias = append(medias, newRow)
		i++
	}

	return medias, nil
}

func GetMediaByID(id int) (ressources.Media, error) {
	var m ressources.Media
	err := database.DB.QueryRow(`
	SELECT id, media_name, media_type, uuid, full_quality_path, thumb_path, compress_path, size, compress_size, created_at 
	FROM medias 
	WHERE id = $1
	`, id).Scan(&m.ID, &m.MediaName, &m.MediaType, &m.UUID, &m.FullQualityPath, &m.ThumbPath, &m.CompressPath, &m.Size, &m.CompressSize, &m.CreatedAt)
	return m, err
}

func CreateMedia(p ressources.Media) (ressources.Media, error) {
	err := database.DB.QueryRow(`
		INSERT INTO medias (media_name, media_type, uuid, full_quality_path, thumb_path, compress_path, size, compress_size)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id, created_at
	`, p.MediaName, p.MediaType, p.UUID, p.FullQualityPath, p.ThumbPath, p.CompressPath, p.Size, p.CompressSize).Scan(&p.ID, &p.CreatedAt)
	return p, err
}

func UpdateMedia(p ressources.Media) (ressources.Media, error) {
	err := database.DB.QueryRow(`
		UPDATE medias SET thumb_path = $1, compress_size = $2, compress_path = $3
		WHERE id = $4
		RETURNING id
	`, p.ThumbPath, p.CompressSize, p.CompressPath, p.ID).Scan(&p.ID)
	return p, err
}
