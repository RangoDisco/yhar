package dto

type UpdateAlbumInput struct {
	Title   *string `json:"title" binding:"omitempty,min=2,max=150"`
	ImageID *int64  `json:"image_id" binding:"omitempty,gt=0"`
}
