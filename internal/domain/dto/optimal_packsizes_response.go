package dto

type OptimalPackSizesResponse struct {
	PackCombination []PackDetail `json:"pack_combination"`
	TotalItems      int          `json:"total_items"`
	TotalPacks      int          `json:"total_packs"`
}

type PackDetail struct {
	Size  int `json:"size"`
	Count int `json:"count"`
}
