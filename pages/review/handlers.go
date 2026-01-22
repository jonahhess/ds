package review

import (
	"database/sql"
	"encoding/json"
	reviewcard "myapp/components/reviewCard"
	"myapp/layouts"
	"myapp/types"
	"net/http"
	"time"
)

type ReviewCard struct {
    ID       int       `json:"card_id"`
    Text     string    `json:"text"`
    ReviewAt time.Time `json:"-"`
    // SM2 fields
    Repetition int     `json:"-"`
    Interval   int     `json:"-"`
    EaseFactor float64 `json:"-"`
    Answers []types.Answer
}

func Page(DB *sql.DB) http.HandlerFunc {
 return func(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
answers := [2]types.Answer{
        {ID: 1, Text: "Baby don't hurt me"},
        {ID: 2, Text: "I want you to show me"},
    }

	layouts.Base("login", reviewcard.ReviewCard(2,"What is love?", answers)).Render(ctx, w)
 }
}

// GET /review/next
func GetNextCard(DB *sql.DB) http.HandlerFunc {
return func (w http.ResponseWriter, r *http.Request) {
    card, err := fetchNextCard(DB)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    if card == nil {
        json.NewEncoder(w).Encode(map[string]interface{}{
            "message": "No cards due for review",
            "card_id": nil,
        })
        return
    }

    json.NewEncoder(w).Encode(card)
    }
}

// POST /review/submit
func SubmitAnswer(DB *sql.DB) http.HandlerFunc {
return func (w http.ResponseWriter, r *http.Request) {
    var input struct {
        CardID  int `json:"card_id"`
        Quality int `json:"quality"` // 0-5
    }
    if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
        http.Error(w, "Invalid JSON", http.StatusBadRequest)
        return
    }

    // Fetch the card from DB
    card, err := fetchCarDByID(DB, input.CardID)
    if err != nil {
        http.Error(w, "Card not found", http.StatusNotFound)
        return
    }

    // Apply SM2 algorithm (stub)
    card.ReviewAt, card.Repetition, card.Interval, card.EaseFactor = applySM2(card, input.Quality)

    // Update DB
    if err := updateCardReview(DB, card); err != nil {
        http.Error(w, "Failed to update card", http.StatusInternalServerError)
        return
    }

    // Fetch next card
    nextCard, err := fetchNextCard(DB)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    if nextCard == nil {
        json.NewEncoder(w).Encode(map[string]interface{}{
            "message": "Session complete",
            "card_id": nil,
        })
        return
    }

    json.NewEncoder(w).Encode(nextCard)
}}

// --- DB Helpers ---

func fetchNextCard(DB *sql.DB) (*ReviewCard, error) {
    q := ReviewCard{}
row := DB.QueryRow(`
    SELECT id, text, review_at, repetition, interval, ease_factor
    FROM question
    WHERE review_at <= CURRENT_TIMESTAMP
    ORDER BY review_at
    LIMIT 1
`)
err := row.Scan(&q.ID, &q.Text, &q.ReviewAt, &q.Repetition, &q.Interval, &q.EaseFactor)
if err == sql.ErrNoRows {
    return nil, nil
}

// fetch answers
rows, err := DB.Query(`SELECT id, text FROM answers WHERE question_id = ?`, q.ID)
if err != nil {
    return nil, err
}
defer rows.Close()

for rows.Next() {
    var a types.Answer
    if err := rows.Scan(&a.ID, &a.Text); err != nil {
        return nil, err
    }
    q.Answers = append(q.Answers, a)
}
    return &q, nil
}

func fetchCarDByID(DB *sql.DB, id int) (*ReviewCard, error) {
    row := DB.QueryRow(`SELECT id, text, review_at, repetition, interval, ease_factor FROM question WHERE id=$1`, id)
    var q ReviewCard
    err := row.Scan(&q.ID, &q.Text, &q.ReviewAt, &q.Repetition, &q.Interval, &q.EaseFactor)
    if err == sql.ErrNoRows {
        return nil, nil
    }
    return &q, err
}

func updateCardReview(DB *sql.DB, q *ReviewCard) error {
    _, err := DB.Exec(`UPDATE question SET review_at=$1, repetition=$2, interval=$3, ease_factor=$4 WHERE id=$5`,
        q.ReviewAt, q.Repetition, q.Interval, q.EaseFactor, q.ID)
    return err
}

// --- SM2 Algorithm Stub ---
func applySM2(card *ReviewCard, quality int) (reviewAt time.Time, repetition int, interval int, easeFactor float64) {
    // Simplified SM2 algorithm logic
    easeFactor = card.EaseFactor
    repetition = card.Repetition
    interval = card.Interval

    if quality < 3 {
        repetition = 0
        interval = 1
    } else {
        repetition++
        easeFactor = max(1.3, easeFactor + 0.1 - float64(5-quality)*(0.08+float64(5-quality)*0.02))
        switch repetition {
        case 1: interval = 1
        case 2: interval = 6
        default: interval = int(float64(interval) * easeFactor)
        }
    }

    reviewAt = time.Now().Add(time.Duration(interval) * 24 * time.Hour)
    return reviewAt, repetition, interval, easeFactor
}