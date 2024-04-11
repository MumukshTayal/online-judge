package add_contest

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	_ "github.com/tursodatabase/libsql-client-go/libsql"
)

type ContestData struct {
	Name           string          `json:"name"`
	StartTime      string          `json:"startTime"`
	EndTime        string          `json:"endTime"`
	Description    string          `json:"description"`
	Problems       []string        `json:"problems"`
	Emails         []string        `json:"emails"`
	LanguageLimits []LanguageLimit `json:"languageLimits"`
}

type LanguageLimit struct {
	ID          int    `json:"id"`
	Language    string `json:"language"`
	TimeLimit   string `json:"timeLimit"`
	MemoryLimit string `json:"memoryLimit"`
}

func AddContest(c *fiber.Ctx) error {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Some error occurred. Err: %s", err)
	}

	dbName := os.Getenv("DATABASE_NAME")
	dbToken := os.Getenv("DATABASE_TOKEN")

	url := fmt.Sprintf("libsql://%s.turso.io?authToken=%s", dbName, dbToken)

	db, err := sql.Open("libsql", url)
	if err != nil {
		return err
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		return err
	}

	var contestData ContestData
	if err := c.BodyParser(&contestData); err != nil {
		return err
	}

	// Insert contest data into the contest table
	if err := AddContestToTable(db, contestData); err != nil {
		return err
	}

	// Insert problem IDs and contest ID into the distribute_problems_to_contest table
	if err := AddProblemIDsToContest(db, contestData); err != nil {
		return err
	}

	// Insert language limits into the allowed_list table
	if err := AddLanguageLimits(db, contestData); err != nil {
		return err
	}

	return c.SendString("Contest added successfully!")
}

func AddContestToTable(db *sql.DB, contestData ContestData) error {
	// Parse the start and end time strings into time.Time format
	startTime, err := time.Parse("2006-01-02T15:04", contestData.StartTime)
	if err != nil {
		return err
	}

	endTime, err := time.Parse("2006-01-02T15:04", contestData.EndTime)
	if err != nil {
		return err
	}

	stmt, err := db.Prepare("INSERT INTO contest (contest_title, contest_description, contest_start_time, contest_end_time) VALUES (?, ?, ?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	result, err := stmt.Exec(contestData.Name, contestData.Description, startTime, endTime)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("failed to insert contest data into contest table")
	}

	return nil
}

func AddProblemIDsToContest(db *sql.DB, contestData ContestData) error {
	// Fetch the contest ID of the recently added contest
	var contestID int
	err := db.QueryRow("SELECT contest_id FROM contest WHERE contest_title = ?", contestData.Name).Scan(&contestID)
	if err != nil {
		return err
	}

	// Insert problem IDs and contest ID into the distribute_problems_to_contest table
	stmt, err := db.Prepare("INSERT INTO distribute_problems_to_contest (contest_id, problem_id) VALUES (?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, problem := range contestData.Problems {
		// You need to fetch the problem ID from the problem table using the problem title
		problemID, err := getProblemIDFromTitle(db, problem)
		if err != nil {
			return err
		}

		_, err = stmt.Exec(contestID, problemID)
		if err != nil {
			return err
		}
	}

	return nil
}

func AddLanguageLimits(db *sql.DB, contestData ContestData) error {
	stmt, err := db.Prepare("INSERT INTO allowed_list (contest_id, language, time_limit, memory_limit) VALUES (?, ?, ?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, limit := range contestData.LanguageLimits {
		// You may need to parse the timeLimit and memoryLimit strings to integers
		timeLimit, err := strconv.Atoi(limit.TimeLimit)
		if err != nil {
			return err
		}
		memoryLimit, err := strconv.Atoi(limit.MemoryLimit)
		if err != nil {
			return err
		}

		// Fetch the contest ID of the recently added contest
		var contestID int
		err = db.QueryRow("SELECT contest_id FROM contest WHERE contest_title = ?", contestData.Name).Scan(&contestID)
		if err != nil {
			return err
		}

		_, err = stmt.Exec(contestID, limit.Language, timeLimit, memoryLimit)
		if err != nil {
			return err
		}
	}

	return nil
}

func getProblemIDFromTitle(db *sql.DB, problemTitle string) (int, error) {
	var problemID int
	err := db.QueryRow("SELECT problem_id FROM problem WHERE problem_title = ?", problemTitle).Scan(&problemID)
	if err != nil {
		return 0, err
	}
	return problemID, nil
}
