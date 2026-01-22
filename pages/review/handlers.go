package review

import (
	"database/sql"
	"myapp/auth"
	pleaseLogin "myapp/components/pleaseLogin"
	reviewcard "myapp/components/reviewCard"
	"myapp/layouts"
	"myapp/types"
	"net/http"
	"strconv"
)

func Page(DB *sql.DB) http.HandlerFunc {
 return func(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID, ok := auth.UserIDFromContext(ctx)
	if !ok {
		layouts.Base("login", pleaseLogin.PleaseLogin()).Render(ctx, w)
		return
	}

	card, err := GetNextReviewCard(DB, userID)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	
	if card == nil {
		layouts.Base("login", reviewcard.ReviewCard(0, "", nil)).Render(ctx, w)
		return
	}
	
	answers, err := GetAnswersForQuestion(DB, card.QuestionID)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	
	layouts.Base("login", reviewcard.ReviewCard(
		card.QuestionID,
		card.Question,
		answers,
	)).Render(ctx, w)
 }
}

func SubmitReviewAnswer(DB *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, ok := auth.UserIDFromContext(r.Context())
        if !ok {
            http.Error(w, "user not logged in", 400)
            http.Redirect(w,r,"/login",401)
            return
        }

		err := r.ParseForm()
		if err != nil {
			http.Error(w, "invalid form", 400)
			return
		}

		questionID, _ := strconv.Atoi(
			r.FormValue("question_id"),
		)
		answerID, _ := strconv.Atoi(
			r.FormValue("answer_id"),
		)

		correct, err := IsAnswerCorrect(DB, questionID, answerID)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		err = UpdateReviewCard(DB, userID, questionID, correct)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		http.Redirect(w, r, "/review", http.StatusSeeOther)
	}
}

func GetNextReviewCard(db *sql.DB, userID int) (*types.ReviewCard, error) {
    row := db.QueryRow(`
        SELECT
            rc.user_id,
            rc.question_id,
            q.text
        FROM reviewcards rc
        JOIN questions q ON q.id = rc.question_id
        WHERE rc.user_id = ?
        ORDER BY rc.consecutive_successes ASC, rc.reviews ASC
        LIMIT 1;
    `, userID)

    var card types.ReviewCard
    err := row.Scan(
        &card.UserID,
        &card.QuestionID,
        &card.Question,
    )

    if err == sql.ErrNoRows {
        return nil, nil // no more cards
    }
    if err != nil {
        return nil, err
    }

    return &card, nil
}

func GetAnswersForQuestion(db *sql.DB, questionID int) ([]types.Answer, error) {
    rows, err := db.Query(`
        SELECT id, text
        FROM answers
        WHERE question_id = ?
        ORDER BY id;
    `, questionID)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var answers []types.Answer

    for rows.Next() {
        var a types.Answer
        if err := rows.Scan(&a.ID, &a.Text); err != nil {
            return nil, err
        }
        answers = append(answers, a)
    }

    return answers, rows.Err()
}

func IsAnswerCorrect(db *sql.DB, questionID, answerID int) (bool, error) {
    row := db.QueryRow(`
        SELECT 1
        FROM correct_answers
        WHERE question_id = ? AND answer_id = ?;
    `, questionID, answerID)

    var dummy int
    err := row.Scan(&dummy)

    if err == sql.ErrNoRows {
        return false, nil
    }
    if err != nil {
        return false, err
    }

    return true, nil
}

func UpdateReviewCard(
    db *sql.DB,
    userID int,
    questionID int,
    correct bool,
) error {

    _, err := db.Exec(`
        UPDATE reviewcards
        SET
            reviews = reviews + 1,
            successes = successes + CASE WHEN ? THEN 1 ELSE 0 END,
            consecutive_successes = CASE
                WHEN ? THEN consecutive_successes + 1
                ELSE 0
            END
        WHERE user_id = ? AND question_id = ?;
    `, correct, correct, userID, questionID)

    return err
}


