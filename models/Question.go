package models

type Question struct {
	ID            uint `gorm:"primaryKey"`
	UserID        uint
	UserSubjectID uint
	QuestionText  string
	AnswerText    string
}
