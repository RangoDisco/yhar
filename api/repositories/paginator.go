package repositories

import "gorm.io/gorm"

func Paginate(page, limit int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if page < 1 {
			page = 1
		}
		if limit < 1 {
			limit = 1
		}
		return db.Offset((page - 1) * limit).Limit(limit)
	}
}
