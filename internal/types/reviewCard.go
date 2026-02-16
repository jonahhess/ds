package types

import "time"

type ReviewCard struct {
	ID           int
	UserID       int
	QuestionID   int
	QuestionText string
	Answers      []Answer
	ReviewAt     time.Time
	Interval     int
	Easiness     float64
	Repetitions  int
	Successes    int
	Reviews      int
	CreatedAt    time.Time
}

type ReviewStats struct {
	TotalCards     int
	CardsDueToday  int
	CardsDueSoon   int
	TotalReviews   int
	SuccessRate    float64
}
