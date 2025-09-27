package models

import "time"

type FullProduct struct {
	ID          string
	Name        string
	Price       int
	Description string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type Product struct {
	Name        string
	Price       int
	Description string
}

type ProductDigest struct {
	ID    string
	Name  string
	Price int
}
