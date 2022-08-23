package core

import (
	"errors"
	"github.com/dominik24c/quizzes_api/internal"
)

// ValidateAndCalculateSolvedQuiz you have to refill only user id and quiz id value
func ValidateAndCalculateSolvedQuiz(questionsBody []internal.QuestionOut, questions []internal.Question) (*QuizSolution, error) {
	var solution QuizSolution
	var isValidQuestion = false
	var totalPoints float32
	var numOfCorrectAnswers int
	var numOfCorrectAnswersSelectedByUser int
	var numOfIncorrectAnswersSelectedByUser int

	var maxPoints float32
	for _, q := range questions {
		maxPoints += float32(q.Points)
	}
	solution.MaxPoints = maxPoints

	for _, q1 := range questionsBody {
		isValidQuestion = false
		for _, q2 := range questions {
			if q2.Question == q1.Question {
				isValidQuestion = true
				for _, a2 := range q2.Answers {
					if a2.IsCorrect {
						numOfCorrectAnswers += 1
					}
				}
				for _, a1 := range q1.Answers {
					for _, a2 := range q2.Answers {
						if a1.Answer == a2.Answer {
							if a2.IsCorrect {
								numOfCorrectAnswersSelectedByUser += 1
							} else {
								numOfIncorrectAnswersSelectedByUser += 1
							}
						}
					}

				}
				diff := numOfCorrectAnswersSelectedByUser - numOfIncorrectAnswersSelectedByUser
				if diff < 0 {
					diff = 0
				}
				totalPoints += float32(diff) / float32(numOfCorrectAnswers) * float32(q2.Points)
				numOfCorrectAnswersSelectedByUser = 0
				numOfIncorrectAnswersSelectedByUser = 0
				numOfCorrectAnswers = 0
			}
		}

		if !isValidQuestion {
			return &solution, errors.New("invalid question or answers")
		}
	}

	solution.TotalPoints = float32(int32(totalPoints*100)) / 100
	solution.Questions = questionsBody

	return &solution, nil
}
