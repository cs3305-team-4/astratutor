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

type Table struct {
	table  string
	column string
}
type Join struct {
	current Table
	new     Table
}

func Sort(sort string, asc string, joins ...Join) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		for _, join := range joins {
			db = db.Joins(fmt.Sprintf("JOIN %s ON %s.%s = %s.%s", join.new.table, join.new.table, join.new.column, join.current.table, join.current.column))
		}
		if len(joins) > 0 {
			db.Group(fmt.Sprintf("%s.%s", joins[0].current.table, joins[0].current.column))
		}
		if asc != "" {
			db = db.Order(fmt.Sprintf("%s %s", sort, asc))
		}
		return db
	}
}
