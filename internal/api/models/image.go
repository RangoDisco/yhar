package models

import (
	"fmt"
	"os"

	"gorm.io/gorm"
)

type Image struct {
	Timestamps
	ID     int64  `json:"id" gorm:"primary_key;autoIncrement"`
	Path   string `json:"path" gorm:"not null"`
	Type   string `json:"type" gorm:"type:varchar(50);not null"`
	Domain string `json:"domain"`
	// Calculated field based on type, domain and path
	ContentURL string `json:"content_url" gorm:"-"`
}

func (i *Image) AfterFind(_ *gorm.DB) error {
	switch i.Type {
	case "distant":
		i.ContentURL = fmt.Sprintf("%s/%s", i.Domain, i.Path)
	default:
		i.ContentURL = fmt.Sprintf("%s/public/img/%s", os.Getenv("BASE_URL"), i.Path)
	}
	return nil
}
