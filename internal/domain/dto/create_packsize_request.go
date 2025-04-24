package dto

type CreatePackSizeRequest struct {
	ProductID int `json:"product_id" binding:"required"`
	Size      int `json:"size" binding:"required,min=1"`
}
