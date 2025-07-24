package initializers

import "github.com/Divyshekhar/7th-sem-project-be/models"

func SeedSubjects() error {
	subjects := []models.Subject{
		{Name: "DSA"},
		{Name: "DBMS"},
		{Name: "OS"},
		{Name: "CN"},
		{Name: "System Design"},
		{Name: "ML"},
		{Name: "Compiler Design"},
		{Name: "OOPs"},
	}

	for _, s := range subjects {
		Db.FirstOrCreate(&s, models.Subject{Name: s.Name})
	}
	return nil
}
