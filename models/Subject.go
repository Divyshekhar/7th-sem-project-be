package models

type Subject struct {
	ID   uint `gorm:"primaryKey"`
	Name string
}
