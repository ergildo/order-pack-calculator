package dto

type PackSizeResponse struct {
	ID        int64 `json:"id"`
	ProductID int   `json:"product_id"`
	Size      int   `json:"size"`
	Active    bool  `json:"active"`
}
