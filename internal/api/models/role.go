package models

type Role struct {
	Timestamps
	ID          int64        `json:"id" gorm:"primary_key;auto_increment"`
	Name        string       `json:"name" gorm:"unique;type:varchar(255);not null"`
	Permissions []Permission `json:"permissions" gorm:"many2many:role_permissions"`
}
