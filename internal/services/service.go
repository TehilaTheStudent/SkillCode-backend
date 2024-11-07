package services

import (
	"github.com/TehilaTheStudent/SkillCode-backend/internal/models"
	"github.com/TehilaTheStudent/SkillCode-backend/internal/repository"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func CreateQuestion(question models.Question) (*models.Question, error) {
	result, err := repository.CreateQuestion(question)
	if err != nil {
		return nil, err
	}
	question.ID = result.InsertedID.(primitive.ObjectID)
	return &question, nil
}

func GetQuestionByID(id string) (models.Question, error) {
	objID, _ := primitive.ObjectIDFromHex(id)
	return repository.GetQuestionByID(objID)
}

func GetAllQuestions() ([]models.Question, error) {
	return repository.GetAllQuestions()
}

func UpdateQuestion(id string, question models.Question) (*models.Question, error) {
	objID, _ := primitive.ObjectIDFromHex(id)
	_, err := repository.UpdateQuestion(objID, question)
	if err != nil {
		return nil, err
	}
	question.ID = objID
	return &question, nil
}

func DeleteQuestion(id string) error {
	objID, _ := primitive.ObjectIDFromHex(id)
	_, err := repository.DeleteQuestion(objID)
	return err
}

func TestQuestion(id primitive.ObjectID, userFunction string) (bool, error) {
	question, err := repository.GetQuestionByID(id)
	if err != nil {
		return false, err
	}

	// Here you would need to execute the userFunction with the test cases
	// This is a complex task and typically involves running the code in a sandboxed environment
	// For simplicity, let's assume we have a function `executeFunction` that does this

	for _, testCase := range question.TestCases {
		result, err := executeFunction(userFunction, testCase.Input)
		if err != nil || result != testCase.ExpectedOutput {
			return false, nil
		}
	}

	return true, nil
}

// Placeholder for the function execution logic
func executeFunction(userFunction string, input interface{}) (interface{}, error) {
	// Implement the logic to execute the userFunction with the given input
	// This is a complex task and typically involves running the code in a sandboxed environment
	return nil, nil
}
