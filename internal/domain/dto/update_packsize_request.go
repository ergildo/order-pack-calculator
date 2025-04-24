package dto

type UpdatePackSizeRequest struct {
	ID     int64 `json:"id" binding:"required"`
	Size   *int  `json:"size" binding:"min=1"`
	Active *bool `json:"active"`
}
