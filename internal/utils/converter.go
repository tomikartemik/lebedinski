package utils

import "lebedinski/internal/model"

func ConvertItemToShortInfo(item model.Item) model.ItemShortInfo {
	return model.ItemShortInfo{
		ID:          item.ID,
		Name:        item.Name,
		Price:       item.Price,
		ActualPrice: item.ActualPrice,
		Discount:    item.Discount,
		SoldOut:     item.SoldOut,
		CategoryID:  item.CategoryID,
		Category:    item.Category,
		Sizes:       item.Sizes,
		Photos:      item.Photos,
	}
}
