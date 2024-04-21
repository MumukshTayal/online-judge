package get_submissions

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

type Submission struct {
	UserID    string
	ProblemID int
	Result    string
	Time      time.Time
	Language  string
}

func GetSubmissions(c *fiber.Ctx) error {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Some error occurred. Err: %s", err)
	}

	dbName := os.Getenv("DATABASE_NAME")
	dbToken := os.Getenv("DATABASE_TOKEN")

	url := fmt.Sprintf("libsql://%s.turso.io?authToken=%s", dbName, dbToken)

	db, err := sql.Open("libsql", url)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to open db %s: %s", url, err)
		return c.Status(500).SendString(err.Error())
	}

	if err := db.Ping(); err != nil {
		fmt.Fprintf(os.Stderr, "failed to ping db %s: %s\n", url, err)
		return c.Status(500).SendString(err.Error())
	}
	defer db.Close()

	contestIDStr := c.Query("contestId")
	contestID, err := strconv.Atoi(contestIDStr)
	if err != nil {
		fmt.Println("Error Identifying the contestID for Submissions!")
	}

	problemIDs, err := fetchContestProblems(db, contestID)
	if err != nil {
		return c.Status(500).SendString(err.Error())
	}
	userEmails, err := fetchContestUsers(db, contestID)
	if err != nil {
		return c.Status(500).SendString(err.Error())
	}

	// fmt.Println(problemIDs, userEmails)

	submissions, err := fetchSubmissions(db, userEmails, problemIDs)
	if err != nil {
		return c.Status(500).SendString(err.Error())
	}

	// fmt.Println(submissions)

	return c.JSON(submissions)
}

func fetchSubmissions(db *sql.DB, userEmails []string, problemIDs []int) (map[string][]Submission, error) {
	submissions := make(map[string][]Submission)

	for _, userEmail := range userEmails {
		for _, problemID := range problemIDs {
			rows, err := db.Query("SELECT result, submission_date_time, language FROM submission WHERE user_email = ? AND problem_id = ?", userEmail, problemID)
			if err != nil {
				return nil, err
			}

			defer rows.Close()

			for rows.Next() {
				var result string
				var submissionTime time.Time
				var lang string
				if err := rows.Scan(&result, &submissionTime, &lang); err != nil {
					return nil, err
				}
				submissions[userEmail] = append(submissions[userEmail], Submission{
					UserID:    userEmail,
					ProblemID: problemID,
					Result:    result,
					Time:      submissionTime,
					Language:  lang,
				})
			}
			if err := rows.Err(); err != nil {
				return nil, err
			}
		}
	}

	return submissions, nil
}

func fetchContestProblems(db *sql.DB, contestID int) ([]int, error) {
	var problemIDs []int

	rows, err := db.Query("SELECT problem_id FROM distribute_problems_to_contest WHERE contest_id = ?", contestID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var problemID int
		if err := rows.Scan(&problemID); err != nil {
			return nil, err
		}
		problemIDs = append(problemIDs, problemID)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return problemIDs, nil
}

func fetchContestUsers(db *sql.DB, contestID int) ([]string, error) {
	var userEmails []string

	rows, err := db.Query("SELECT user_email FROM contest_user WHERE contest_id = ?", contestID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var userEmail string
		if err := rows.Scan(&userEmail); err != nil {
			return nil, err
		}
		userEmails = append(userEmails, userEmail)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return userEmails, nil
}

func GetAllSubmissions(c *fiber.Ctx) error {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Some error occurred. Err: %s", err)
	}

	dbName := os.Getenv("DATABASE_NAME")
	dbToken := os.Getenv("DATABASE_TOKEN")

	url := fmt.Sprintf("libsql://%s.turso.io?authToken=%s", dbName, dbToken)

	db, err := sql.Open("libsql", url)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to open db %s: %s", url, err)
		return c.Status(500).SendString(err.Error())
	}

	if err := db.Ping(); err != nil {
		fmt.Fprintf(os.Stderr, "failed to ping db %s: %s\n", url, err)
		return c.Status(500).SendString(err.Error())
	}
	defer db.Close()

	if err != nil {
		fmt.Println("Error Identifying the contestID for Submissions!")
	}

	submissions, err := fetchAllSubmissions(db)
	if err != nil {
		return c.Status(500).SendString(err.Error())
	}

	// fmt.Println(submissions)

	return c.JSON(submissions)
}

func fetchAllSubmissions(db *sql.DB) ([]Submission, error) {
	var submissions []Submission
	rows, err := db.Query("SELECT problem_id, user_email, result, submission_date_time, language FROM submission")
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var result string
		var submissionTime time.Time
		var problemId int
		var userEmail string
		var lang string
		if err := rows.Scan(&problemId, &userEmail, &result, &submissionTime, &lang); err != nil {
			return nil, err
		}
		submissions = append(submissions, Submission{
			UserID:    userEmail,
			ProblemID: problemId,
			Result:    result,
			Time:      submissionTime,
			Language:  lang,
		})
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return submissions, nil
}
