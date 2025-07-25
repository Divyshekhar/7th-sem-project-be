package models

type User struct {
	ID       uint `gorm:"primaryKey"`
	Name     string
	Email    string `gorm:"unique"`
	Password string
	Subjects []UserSubject `gorm:"foreignKey:UserID"`
}
