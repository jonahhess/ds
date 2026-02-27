package root

import (
	"database/sql"
	"net/http"

	"github.com/jonahhess/ds/internal/auth"
	"github.com/jonahhess/ds/internal/errors"
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
