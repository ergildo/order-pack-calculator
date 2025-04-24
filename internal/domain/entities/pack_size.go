package entities

type PackSize struct {
	ID        int64 `db:"id"`
	ProductID int   `db:"product_id"`
	Size      int   `db:"size"`
	Active    bool  `db:"active"`
}
