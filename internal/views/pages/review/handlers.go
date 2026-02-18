package review

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"github.com/jonahhess/ds/internal/auth"
	"github.com/jonahhess/ds/internal/errors"
	"github.com/jonahhess/ds/internal/params"
	"github.com/jonahhess/ds/internal/types"
	"github.com/jonahhess/ds/internal/views/layouts"
)

// Page - GET /review - Main review page with stats
func Page(DB *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		userID, ok := auth.UserIDFromContext(ctx)
		if !ok {
			errors.HandleUnauthorized(w, r)
			return
		}

		stats, err := GetReviewStats(DB, userID)
		if err != nil {
			errors.HandleInternalError(w, r, err)
			return
		}

		if err := layouts.
			Base("Review", Review(*stats)).
			Render(ctx, w); err != nil {
			errors.HandleInternalError(w, r, err)
		}
	}
}

// NextCard - GET /review/next - Get next card due for review
func NextCard(DB *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		userID, ok := auth.UserIDFromContext(ctx)
		if !ok {
			errors.HandleUnauthorized(w, r)
			return
		}

		card, err := fetchNextCard(DB, userID)
		if err != nil {
			errors.HandleInternalError(w, r, err)
			return
		}

		if card == nil {
			// No cards due, redirect to complete page
			http.Redirect(w, r, "/review/complete", http.StatusSeeOther)
			return
		}

		if err := layouts.
			Base("Review Card", ReviewCardQuestion(*card)).
			Render(ctx, w); err != nil {
			errors.HandleInternalError(w, r, err)
		}
	}
}

// ShowAnswer - GET /review/card/{questionID}/answer - Show answer side of card
func ShowAnswer(DB *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		userID, ok := auth.UserIDFromContext(ctx)
		if !ok {
			errors.HandleUnauthorized(w, r)
			return
		}
		
		questionID, ok := params.IntFrom(ctx, "questionID")
		if !ok {
			errors.HandleBadRequest(w, r, "question id not found")
			return
		}

		card, correctAnswer, err := fetchCardByID(DB, userID, questionID)
		if err != nil {
			errors.HandleInternalError(w, r, err)
			return
		}
		if card == nil {
			errors.HandleNotFound(w, r, "Card")
			return
		}

		csrfToken := auth.CSRFToken(r)
		if err := layouts.
			Base("Review Card Answer", ReviewCardAnswer(*card, correctAnswer, csrfToken)).
			Render(ctx, w); err != nil {
			errors.HandleInternalError(w, r, err)
		}
	}
}

// RateCard - POST /review/card/{questionID}/rate - Rate card quality and update SM2
func RateCard(DB *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		userID, ok := auth.UserIDFromContext(ctx)
		if !ok {
			errors.HandleUnauthorized(w, r)
			return
		}

		questionID, ok := params.IntFrom(ctx, "questionID")
		if !ok {
			errors.HandleBadRequest(w, r, "question id not found")
			return
		}

		if err := r.ParseForm(); err != nil {
			errors.HandleBadRequest(w, r, "Error parsing form")
			return
		}

		qualityStr := r.FormValue("quality")
		if qualityStr == "" {
			errors.HandleBadRequest(w, r, "quality rating required")
			return
		}

		var quality int
		if _, err := fmt.Sscanf(qualityStr, "%d", &quality); err != nil || quality < 0 || quality > 5 {
			errors.HandleBadRequest(w, r, "quality must be 0-5")
			return
		}

		// Fetch card
		card, correctAnswer, err := fetchCardByID(DB, userID, questionID)
		if err != nil || card == nil || correctAnswer == "" {
			errors.HandleInternalError(w, r, err)
			return
		}

		// Apply SM2 algorithm
		newReviewAt, newRepetitions, newInterval, newEasiness := applySM2(card, quality)

		// Update card in database
		if err := updateCardReview(DB, userID, questionID, newReviewAt, newRepetitions, newInterval, newEasiness, quality); err != nil {
			errors.HandleInternalError(w, r, err)
			return
		}

		// Redirect to next card
		http.Redirect(w, r, "/review/next", http.StatusSeeOther)
	}
}

// Complete - GET /review/complete - Show completion page
func Complete(DB *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		userID, ok := auth.UserIDFromContext(ctx)
		if !ok {
			errors.HandleUnauthorized(w, r)
			return
		}

		stats, err := GetReviewStats(DB, userID)
		if err != nil {
			errors.HandleInternalError(w, r, err)
			return
		}

		if err := layouts.
			Base("Review Complete", ReviewComplete(*stats)).
			Render(ctx, w); err != nil {
			errors.HandleInternalError(w, r, err)
		}
	}
}

// --- Database Helpers ---

func GetReviewStats(DB *sql.DB, userID int) (*types.ReviewStats, error) {
	var stats types.ReviewStats

	// Total cards
	err := DB.QueryRow(`
		SELECT COUNT(*) 
		FROM reviewcards 
		WHERE user_id = ?
	`, userID).Scan(&stats.TotalCards)
	if err != nil {
		return nil, err
	}

	// Cards due today
	err = DB.QueryRow(`
		SELECT COUNT(*) 
		FROM reviewcards 
		WHERE user_id = ? AND review_at <= CURRENT_TIMESTAMP
	`, userID).Scan(&stats.CardsDueToday)
	if err != nil {
		return nil, err
	}

	// Cards due in next 7 days
	err = DB.QueryRow(`
		SELECT COUNT(*) 
		FROM reviewcards 
		WHERE user_id = ? AND review_at <= datetime('now', '+7 days')
	`, userID).Scan(&stats.CardsDueSoon)
	if err != nil {
		return nil, err
	}

	// Total reviews
	err = DB.QueryRow(`
		SELECT COALESCE(SUM(reviews), 0)
		FROM reviewcards 
		WHERE user_id = ?
	`, userID).Scan(&stats.TotalReviews)
	if err != nil {
		return nil, err
	}

	// Success rate
	var totalSuccesses, totalReviews int
	err = DB.QueryRow(`
		SELECT 
			COALESCE(SUM(successes), 0),
			COALESCE(SUM(reviews), 0)
		FROM reviewcards 
		WHERE user_id = ?
	`, userID).Scan(&totalSuccesses, &totalReviews)
	if err != nil {
		return nil, err
	}

	if totalReviews > 0 {
		stats.SuccessRate = float64(totalSuccesses) / float64(totalReviews) * 100
	}

	return &stats, nil
}

func fetchNextCard(DB *sql.DB, userID int) (*types.ReviewCard, error) {
	var card types.ReviewCard
	var questionID int

	row := DB.QueryRow(`
		SELECT 
			rc.id,
			rc.question_id,
			q.text,
			rc.review_at,
			rc.repetitions,
			rc.interval,
			rc.easiness,
			rc.successes,
			rc.reviews,
			rc.created_at
		FROM reviewcards rc
		JOIN questions q ON q.id = rc.question_id
		WHERE rc.user_id = ? AND rc.review_at <= CURRENT_TIMESTAMP
		ORDER BY rc.review_at ASC
		LIMIT 1
	`, userID)

	err := row.Scan(
		&card.ID,
		&questionID,
		&card.QuestionText,
		&card.ReviewAt,
		&card.Repetitions,
		&card.Interval,
		&card.Easiness,
		&card.Successes,
		&card.Reviews,
		&card.CreatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	card.UserID = userID
	card.QuestionID = questionID

	// Fetch answers for this question
	rows, err := DB.Query(`
		SELECT id, text 
		FROM answers 
		WHERE question_id = ?
		ORDER BY id
	`, questionID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var answer types.Answer
		if err := rows.Scan(&answer.ID, &answer.Text); err != nil {
			return nil, err
		}
		answer.QuestionID = questionID
		card.Answers = append(card.Answers, answer)
	}

	return &card, nil
}

func fetchCardByID(DB *sql.DB, userID int, questionID int) (*types.ReviewCard, string, error) {
	var card types.ReviewCard

	row := DB.QueryRow(`
		SELECT 
			rc.id,
			rc.question_id,
			q.text,
			rc.review_at,
			rc.repetitions,
			rc.interval,
			rc.easiness,
			rc.successes,
			rc.reviews,
			rc.created_at
		FROM reviewcards rc
		JOIN questions q ON q.id = rc.question_id
		WHERE rc.question_id = ? AND rc.user_id = ?
	`, questionID, userID)

	err := row.Scan(
		&card.ID,
		&questionID,
		&card.QuestionText,
		&card.ReviewAt,
		&card.Repetitions,
		&card.Interval,
		&card.Easiness,
		&card.Successes,
		&card.Reviews,
		&card.CreatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, "", nil
	}
	if err != nil {
		return nil, "", err
	}

	card.UserID = userID
	card.QuestionID = questionID

	// Fetch answers
	rows, err := DB.Query(`
		SELECT id, text 
		FROM answers 
		WHERE question_id = ?
		ORDER BY id
	`, questionID)
	if err != nil {
		return nil, "", err
	}
	defer rows.Close()

	for rows.Next() {
		var answer types.Answer
		if err := rows.Scan(&answer.ID, &answer.Text); err != nil {
			return nil, "", err
		}
		answer.QuestionID = questionID
		card.Answers = append(card.Answers, answer)
	}

	var correctAnswer string
	err = DB.QueryRow(`
		SELECT a.text 
		FROM correct_answers ca
		JOIN answers a ON a.id = ca.answer_id
		WHERE ca.question_id = ?
	`, questionID).Scan(&correctAnswer)
	if err != nil {
		return nil, "", err
	}

	return &card, correctAnswer, nil
}

func updateCardReview(DB *sql.DB, userID int, questionID int, reviewAt time.Time, repetitions int, interval int, easiness float64, quality int) error {
	// Update SM2 fields and increment counters
	_, err := DB.Exec(`
		UPDATE reviewcards 
		SET 
			review_at = ?,
			repetitions = ?,
			interval = ?,
			easiness = ?,
			reviews = reviews + 1,
			successes = successes + ?
		WHERE id = ? AND user_id = ?
	`, reviewAt, repetitions, interval, easiness, 
		// If quality >= 3, count as success (1), otherwise 0
		map[bool]int{true: 1, false: 0}[quality >= 3],
		questionID, userID)
	return err
}

// --- SM2 Algorithm ---
func applySM2(card *types.ReviewCard, quality int) (reviewAt time.Time, repetitions int, interval int, easiness float64) {
	easiness = card.Easiness
	repetitions = card.Repetitions
	interval = card.Interval

	if quality < 3 {
		// Failed - reset
		repetitions = 0
		interval = 1
	} else {
		// Passed
		repetitions++
		// Update easiness factor
		easiness = easiness + (0.1 - float64(5-quality)*(0.08+float64(5-quality)*0.02))
		if easiness < 1.3 {
			easiness = 1.3
		}

		// Calculate interval
		switch repetitions {
		case 1:
			interval = 1
		case 2:
			interval = 6
		default:
			interval = int(float64(interval) * easiness)
		}
	}

	reviewAt = time.Now().Add(time.Duration(interval) * 24 * time.Hour)
	return reviewAt, repetitions, interval, easiness
}
