package dto

type UpdateArtistInput struct {
	Name    *string `json:"name" binding:"omitempty,min=2,max=150"`
	ImageID *int64  `json:"image_id" binding:"omitempty,gt=0"`
}
