package repository

import "github.com/glekoz/online_shop_product/internal/models"

func ExtractProductFromDB(r Product) (models.Product, error) {
	p := models.Product{}
	err := p.ID.Scan(r.ID)
	if err != nil {
		return models.Product{}, err
	}
	err = p.UserID.Scan(r.UserID)
	if err != nil {
		return models.Product{}, err
	}
	p.Name = r.Name
	p.Description = r.Description
	p.Price = int(r.Price)
	p.CreatedAt = r.CreatedAt.Time
	return p, nil
}

/*
func ExtractProductImageFromDB(r ProductImage) (models.ProductImage, error) {
	pi := models.ProductImage{}
	err := pi.ProductID.Scan(r.ProductID)
	if err != nil {
		return models.ProductImage{}, err
	}
	pi.MaxCount = int(r.MaxCount)
	return pi, nil
}

func ExtractFullProductInfoFromDB(r GetOneRow) (models.FullProductInfo, error) {
	p := models.FullProductInfo{}
	err := p.ID.Scan(r.ID)
	if err != nil {
		return models.FullProductInfo{}, err
	}
	err = p.UserID.Scan(r.UserID)
	if err != nil {
		return models.FullProductInfo{}, err
	}
	p.Name = r.Name
	p.Description = r.Description
	p.Price = int(r.Price)
	p.CreatedAt = r.CreatedAt.Time
	p.MaxCount = int(r.MaxCount)
	return p, nil
}
*/
