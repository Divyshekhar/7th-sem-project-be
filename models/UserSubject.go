package models

type UserSubject struct {
	ID        uint `gorm:"primaryKey"`
	UserID    uint
	SubjectID uint
	Subject   Subject
	Questions []Question `gorm:"foreignKey:UserSubjectID"`
}
