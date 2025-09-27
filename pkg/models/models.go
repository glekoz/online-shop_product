package models

import "time"

type ProductCreation struct {
	ID          string
	Name        string
	Price       int
	Description string
}

type Product struct {
	ID          string
	Name        string
	Price       int
	Description string
	CreatedAt   time.Time
}

type ProductDigest struct {
	ID    string
	Name  string
	Price int
}
