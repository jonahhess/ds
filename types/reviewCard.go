package types

type ReviewCard struct {
    UserID     int
    QuestionID int
    Question   string
    Answers    []Answer
}