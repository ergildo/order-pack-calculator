package dto

type CalculatePackSizesRequest struct {
	ProductID     int `json:"product_id" binding:"required"`
	OrderQuantity int `json:"order_quantity" binding:"required,min=1"`
}
