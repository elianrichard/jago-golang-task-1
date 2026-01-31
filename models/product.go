package models

type Product struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Stock       int       `json:"stock"`
	Price       int       `json:"price"`
	Category_ID *string   `json:"category_id,omitempty"`
	Category    *Category `json:"category,omitempty"`
}
