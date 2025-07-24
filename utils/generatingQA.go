package utils

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/Divyshekhar/7th-sem-project-be/initializers"
	"github.com/Divyshekhar/7th-sem-project-be/models"
	"github.com/tmc/langchaingo/llms/googleai"
	"gorm.io/gorm"
)

func buildPrompt(subject string, past []models.Question) string {
	pastBlock := ""
	for i, q := range past {
		pastBlock += fmt.Sprintf("Q%d: %s\n", i+1, q.QuestionText)
	}
	return fmt.Sprintf(`
		You are an interviewer for a technical company.
		Generate 10 new and non-repeating questions and answer pairs for the subject %s.
		No other texts given start with the questions 

		Do NOT repeat any of the following questions:
		%s

		FORMAT for the question and answer is Q<number> representing the question number and A<number> representing the answer number for the respected Q<number>:
		Q1: ...
		A1: ...
		Q2: ...
		A2: ...
	`, subject, pastBlock)

}

type QAPair struct {
	Question string
	Answer   string
}

func ParseQABlock(output string) []QAPair {
	lines := strings.Split(output, "\n")
	var result []QAPair
	var q, a string

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "Q") {
			q = strings.Join(strings.SplitN(line, ": ", 2)[1:], "")
		} else if strings.HasPrefix(line, "A") {
			a = strings.Join(strings.SplitN(line, ": ", 2)[1:], "")
			result = append(result, QAPair{Question: q, Answer: a})
		}
	}
	return result
}

func SaveQuestion(user models.User, subjectName string, qa []QAPair) error {
	var subject models.Subject
	if err := initializers.Db.FirstOrCreate(&subject, models.Subject{Name: subjectName}).Error; err != nil {
		return err
	}

	var userSubject models.UserSubject
	err := initializers.Db.Where(&models.UserSubject{
		UserID:    user.ID,
		SubjectID: subject.ID,
	}).First(&userSubject).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		userSubject = models.UserSubject{
			UserID:    user.ID,
			SubjectID: subject.ID,
		}
		if err := initializers.Db.Create(&userSubject).Error; err != nil {
			return err
		}
	} else if err != nil {
		return err
	}

	var questions []models.Question
	for _, pair := range qa {
		questions = append(questions, models.Question{
			UserID:        user.ID,
			UserSubjectID: userSubject.ID,
			QuestionText:  pair.Question,
			AnswerText:    pair.Answer,
		})
	}

	if err := initializers.Db.Create(&questions).Error; err != nil {
		return err
	}
	return nil
}

func GetUserSubjectQuestions(user models.User, subjectName string) ([]models.Question, error) {
	var subject models.Subject
	if err := initializers.Db.Where("name = ?", subjectName).First(&subject).Error; err != nil {
		return nil, err // subject not found
	}
	var userSubject models.UserSubject
	err := initializers.Db.
		Where("user_id = ? AND subject_id = ?", user.ID, subject.ID).
		First(&userSubject).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		// Create the UserSubject link if not found
		userSubject = models.UserSubject{
			UserID:    user.ID,
			SubjectID: subject.ID,
		}
		if err := initializers.Db.Create(&userSubject).Error; err != nil {
			return nil, fmt.Errorf("failed to create UserSubject: %w", err)
		}
	} else if err != nil {
		return nil, err // some other DB error
	}

	var questions []models.Question
	if err := initializers.Db.
		Where("user_subject_id = ?", userSubject.ID).
		Find(&questions).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch questions: %w", err)
	}

	return questions, nil
}

func Generate(user models.User, subject string) (string, error) {
	ctx := context.Background()
	pastQuestions, _ := GetUserSubjectQuestions(user, subject)
	fmt.Println("GetUserSubjectQuestions is DONEEEEEEEEE")
	prompt := buildPrompt(subject, pastQuestions)
	fmt.Println("Build Prompt is DONEEEEEEE")
	llm, err := googleai.New(ctx, googleai.WithAPIKey(os.Getenv("GEMINI_API")), googleai.WithDefaultModel("gemini-2.0-flash"))
	if err != nil {
		return "", err
	}
	fmt.Println("GOOGLEAI.NEW DONEEEEEEE")
	output, err := llm.Call(ctx, prompt)
	if err != nil {
		return "", err
	}
	fmt.Println("LLM.CALL DONEEEEEEE")
	qaPairs := ParseQABlock(output)
	fmt.Println("Parse DONEEEEE")
	SaveQuestion(user, subject, qaPairs)
	fmt.Println("SaveQuestion Doneeeee")
	return output, nil
}
