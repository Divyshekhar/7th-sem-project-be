package models

type User struct {
	ID        uint `gorm:"primaryKey"`
	FirstName string
	LastName  string
	Email     string `gorm:"unique"`
	Password  string
	Subjects  []UserSubject `gorm:"foreignKey:UserID"`
}
