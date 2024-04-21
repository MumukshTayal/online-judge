package get_leaderboard

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"sort"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	_ "github.com/tursodatabase/libsql-client-go/libsql"
)

type LeaderboardEntry struct {
	UserEmail          string
	Result             int
	LastSubmissionTime time.Time
}

type Submission struct {
	UserID    string
	ProblemID int
	Result    string
	Time      time.Time
}

type UserScore struct {
	UserID            string
	TotalScore        int
	MaxSubmissionTime time.Time
}

func GetLeaderboard(c *fiber.Ctx) error {
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
		fmt.Println("Error Identifying the contestID for Leaderboard!")
	}

	problemIDs, err := fetchContestProblems(db, contestID)
	if err != nil {
		return c.Status(500).SendString(err.Error())
	}
	userEmails, err := fetchContestUsers(db, contestID)
	if err != nil {
		return c.Status(500).SendString(err.Error())
	}

	submissions, err := fetchSubmissions(db, userEmails, problemIDs)
	if err != nil {
		return c.Status(500).SendString(err.Error())
	}

	userScores := calculateUserScores(submissions)

	var leaderboard []LeaderboardEntry
	for _, score := range userScores {
		leaderboard = append(leaderboard, LeaderboardEntry{
			UserEmail:          score.UserID,
			Result:             score.TotalScore,
			LastSubmissionTime: score.MaxSubmissionTime,
		})
	}

	sort.Slice(leaderboard, func(i, j int) bool {
		if leaderboard[i].Result == leaderboard[j].Result {
			return leaderboard[i].LastSubmissionTime.After(leaderboard[j].LastSubmissionTime)
		}
		return leaderboard[i].Result > leaderboard[j].Result
	})

	return c.JSON(leaderboard)
}

func calculateUserScores(submissions map[string][]Submission) map[string]UserScore {
	userScores := make(map[string]UserScore)

	for userID, userSubmissions := range submissions {
		totalScore := 0
		maxSubmissionTime := time.Time{}

		for _, submission := range userSubmissions {
			if submission.Result == "AC" {
				totalScore += 10
			}

			if submission.Time.After(maxSubmissionTime) {
				maxSubmissionTime = submission.Time
			}
		}

		userScores[userID] = UserScore{
			UserID:            userID,
			TotalScore:        totalScore,
			MaxSubmissionTime: maxSubmissionTime,
		}
	}

	return userScores
}

func fetchSubmissions(db *sql.DB, userEmails []string, problemIDs []int) (map[string][]Submission, error) {
	submissions := make(map[string][]Submission)

	for _, userEmail := range userEmails {
		for _, problemID := range problemIDs {
			rows, err := db.Query("SELECT result, submission_date_time FROM submission WHERE user_email = ? AND problem_id = ?", userEmail, problemID)
			if err != nil {
				return nil, err
			}
			// fmt.Println(rows)
			defer rows.Close()

			for rows.Next() {
				var result string
				var submissionTime time.Time
				if err := rows.Scan(&result, &submissionTime); err != nil {
					return nil, err
				}
				submissions[userEmail] = append(submissions[userEmail], Submission{
					UserID:    userEmail,
					ProblemID: problemID,
					Result:    result,
					Time:      submissionTime,
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
