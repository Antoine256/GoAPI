package ressources

type Media struct {
	ID              int    `json:"id"`
	MediaType       string `json:"type"` // "picture" or "video"
	MediaName       string `json:"media_name"`
	UUID            string `json:"uuid"`
	FullQualityPath string `json:"full_quality_path"`
	ThumbPath       string `json:"thumb_path"`
	CompressPath    string `json:"compress_path"`
	Size            int64  `json:"size"`
	CompressSize    int64  `json:"compress_size"`
	CreatedAt       string `json:"created_at"`
}

type MediaResponse struct {
	ID              int    `json:"id"`
	MediaName       string `json:"media_name"`
	MediaType       string `json:"type"`
	UUID            string `json:"uuid"`
	ThumbPath       string `json:"thumb_path"`
	FullQualityPath string `json:"full_quality_path"`
	CompressPath    string `json:"compress_path"`
	Size            int64  `json:"size"`
	CompressSize    int64  `json:"compress_size"`
	CreatedAt       string `json:"created_at"`
}
