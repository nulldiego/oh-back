package utils

import "gorm.io/gorm"

type PaginatedResult struct {
	Data  interface{} `json:"data"`
	Page  int         `json:"page"`
	From  int         `json:"from"`
	To    int         `json:"to"`
	Limit int         `json:"limit"`
	Total int64       `json:"total"`
}

func Paginated(db *gorm.DB, page, limit int, rawFunc func(*gorm.DB) *gorm.DB, output interface{}) (PaginatedResult, error) {
	offset := (page - 1) * limit

	query := db
	if rawFunc != nil {
		query = rawFunc(query)
	}

	var total int64
	query.Model(output).Count(&total)

	err := query.Offset(offset).Limit(limit).Find(output).Error
	if err != nil {
		return PaginatedResult{}, nil
	}

	to := offset + limit
	if to > int(total) {
		to = int(total)
	}

	return PaginatedResult{
		Data:  output,
		Page:  page,
		From:  offset + 1,
		To:    to,
		Limit: limit,
		Total: total,
	}, nil
}
