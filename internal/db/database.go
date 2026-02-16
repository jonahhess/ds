package db

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

// Initialize opens a SQLite database file
func InitDB(path string) error {
	var err error
	DB, err = sql.Open("sqlite3", path)
	if err != nil {
		return err
	}

	if _, err := DB.Exec(`
        PRAGMA foreign_keys = ON;
        PRAGMA journal_mode = WAL;
        PRAGMA synchronous = NORMAL;
    `); err != nil {
		return err
	}

	return DB.Ping()
}

func CloseDB() error {
	if DB != nil {
		return DB.Close()
	}
	return nil
}

func execSQL(db *sql.DB, sqlStmt string) error {
    _, err := db.Exec(sqlStmt)
    if err != nil {
        log.Printf("Error executing SQL: %v\nSQL:\n%s", err, sqlStmt)
        return err
    }
    return nil
}


func CreateTables() error {
	users := `
	CREATE TABLE IF NOT EXISTS users (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  name TEXT NOT NULL,
  email TEXT NOT NULL UNIQUE,
  password_hash TEXT NOT NULL,
  created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);
    `
	courses := `
	CREATE TABLE IF NOT EXISTS courses (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  title TEXT NOT NULL,
  description TEXT,
  created_by INTEGER NOT NULL,
  created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY (created_by) REFERENCES users(id)
);
    `

  user_courses := `
	CREATE TABLE IF NOT EXISTS user_courses (
  user_id INTEGER NOT NULL,
  course_id INTEGER NOT NULL,
  current_lesson INTEGER DEFAULT 0,
  PRIMARY KEY (user_id, course_id),
  FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
  FOREIGN KEY (course_id) REFERENCES courses(id)
);
    `

 	lessons := `
	CREATE TABLE IF NOT EXISTS lessons (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  course_id INTEGER NOT NULL,
  lesson_index INTEGER NOT NULL,
  title TEXT NOT NULL,
  text TEXT NOT NULL,
  created_by INTEGER NOT NULL,
  created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY (course_id) REFERENCES courses(id) ON DELETE CASCADE,
  FOREIGN KEY (created_by) REFERENCES users(id)
);
    `

  quizzes := `
	CREATE TABLE IF NOT EXISTS quizzes (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  lesson_id INTEGER NOT NULL UNIQUE,
  FOREIGN KEY (lesson_id) REFERENCES lessons(id) ON DELETE CASCADE
);
    `

  questions := `
	CREATE TABLE IF NOT EXISTS questions (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  quiz_id INTEGER NOT NULL,
  text TEXT NOT NULL,
  created_by INTEGER NOT NULL,
  created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY (quiz_id) REFERENCES quizzes(id) ON DELETE CASCADE,
  FOREIGN KEY (created_by) REFERENCES users(id)
);
    `

  answers := `
	CREATE TABLE IF NOT EXISTS answers (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  question_id INTEGER NOT NULL,
  text TEXT NOT NULL,
  FOREIGN KEY (question_id) REFERENCES questions(id) ON DELETE CASCADE

);
    `

  correct_answers := `
	CREATE TABLE IF NOT EXISTS correct_answers (
  question_id INTEGER PRIMARY KEY,
  answer_id INTEGER NOT NULL UNIQUE,
  FOREIGN KEY (question_id) REFERENCES questions(id) ON DELETE CASCADE,
  FOREIGN KEY (answer_id) REFERENCES answers(id) ON DELETE CASCADE
);
    `

  reviewcards := `
  CREATE TABLE IF NOT EXISTS reviewcards (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    question_id INTEGER NOT NULL,
    review_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    interval INTEGER DEFAULT 1 CHECK(interval > 0),
    easiness REAL DEFAULT 2.5 CHECK(easiness >= 1.3),
    repetitions INTEGER DEFAULT 0 CHECK(repetitions >= 0),
    successes INTEGER DEFAULT 0 CHECK(successes >= 0),
    reviews INTEGER DEFAULT 0 CHECK(reviews >= 0),
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(user_id, question_id),
    FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY(question_id) REFERENCES questions(id) ON DELETE CASCADE
);

    `

    tables := []string{users,courses,user_courses,lessons,quizzes,questions,answers,correct_answers,reviewcards}

    for _, tbl := range tables {
        if err := execSQL(DB, tbl); err != nil {
            return err
        }
    }

        _, err := DB.Exec(`
    CREATE TRIGGER IF NOT EXISTS validate_current_lesson_insert
    BEFORE INSERT ON user_courses
    FOR EACH ROW
    BEGIN
        SELECT CASE
            WHEN NEW.current_lesson > (SELECT COUNT(*) FROM lessons WHERE course_id = NEW.course_id)
            THEN RAISE(ABORT, 'current_lesson cannot exceed total lessons')
        END;
    END;

    CREATE TRIGGER IF NOT EXISTS validate_current_lesson_update
    BEFORE UPDATE ON user_courses
    FOR EACH ROW
    BEGIN
        SELECT CASE
            WHEN NEW.current_lesson > (SELECT COUNT(*) FROM lessons WHERE course_id = NEW.course_id)
            THEN RAISE(ABORT, 'current_lesson cannot exceed total lessons')
        END;
    END;
`)
    if err != nil {
        return err
    }

	return nil
}
