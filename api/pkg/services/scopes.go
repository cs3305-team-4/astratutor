package services

import (
	"fmt"
	"strings"

	"gorm.io/gorm"
)

func Paginate(pageSize int, page int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if page <= 0 {
			page = 1
		}

		if pageSize <= 0 {
			pageSize = 1
		}

		offset := (page - 1) * pageSize
		return db.Offset(offset).Limit(pageSize)
	}
}

type SearchQuery struct {
	field,
	query string
}

func Search(queries ...SearchQuery) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		q := ""
		conds := []interface{}{}
		for i, query := range queries {
			if query.query != "" {
				if i > 0 {
					q = fmt.Sprintf("%s OR LOWER( %s ) LIKE ?", q, query.field)
				} else {
					q = fmt.Sprintf("LOWER( %s ) LIKE ?", query.field)
				}
				conds = append(conds, "%"+strings.ToLower(query.query)+"%")
			}
		}
		if q != "" {
			db = db.Where(q, conds...)
		}
		return db
	}
}
func SortTutors(sort string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		switch sort {
		case "low":
			db = db.Order("AVG( subject_taughts.price ) asc")
		case "high":
			db = db.Order("AVG( subject_taughts.price ) desc")
		}
		return db
	}
}
