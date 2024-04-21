package get_contest_details

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

type ContestDetails struct {
	Contest         Contest          `json:"contest"`
	ContestProblems []ContestProblem `json:"contest_problems"`
}

type Contest struct {
	ContestID    int       `json:"contest_id"`
	ContestTitle string    `json:"contest_title"`
	ContestDesc  string    `json:"contest_description"`
	StartTime    time.Time `json:"start_time"`
	EndTime      time.Time `json:"end_time"`
	IsPublic     string    `json:"is_public"`
	CreatorEmail string    `json:"creator_email"`
}

type ContestProblem struct {
	ProblemID       int    `json:"problem_id"`
	ProblemTitle    string `json:"problem_title"`
	ProblemDesc     string `json:"problem_description"`
	ConstraintsDesc string `json:"constraints_desc"`
	InputFormat     string `json:"input_format"`
	OutputFormat    string `json:"output_format"`
	SampleInput     string `json:"sample_input"`
	SampleOutput    string `json:"sample_output"`
	CreatorEmail    string `json:"creator_email"`
	IsPrivate       bool   `json:"is_private"`
}

func GetContestDetails(c *fiber.Ctx) error {
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

	contestID := c.Params("contestId")

	contest, err := fetchContest(db, contestID)
	if err != nil {
		return c.Status(500).SendString(err.Error())
	}

	contestProblems, err := fetchContestProblems(db, contestID)
	if err != nil {
		return c.Status(500).SendString(err.Error())
	}

	contestDetails := ContestDetails{
		Contest:         contest,
		ContestProblems: contestProblems,
	}

	return c.JSON(contestDetails)
}

func fetchContest(db *sql.DB, contestID string) (Contest, error) {
	var contest Contest

	row := db.QueryRow("SELECT contest_id, contest_title, contest_description, contest_start_time, contest_end_time, is_public, creator_email FROM contest WHERE contest_id = ?", contestID)

	var creatorEmail sql.NullString // Use sql.NullString for handling NULL values

	err := row.Scan(&contest.ContestID, &contest.ContestTitle, &contest.ContestDesc, &contest.StartTime, &contest.EndTime, &contest.IsPublic, &creatorEmail)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to scan contest row: %v\n", err)
		return contest, err
	}

	// Check if the creator_email value is valid before assigning it to contest.CreatorEmail
	if creatorEmail.Valid {
		contest.CreatorEmail = creatorEmail.String
	} else {
		// Handle NULL value appropriately, for example, set an empty string
		contest.CreatorEmail = ""
	}

	return contest, nil
}

func fetchContestProblems(db *sql.DB, contestID string) ([]ContestProblem, error) {
	rows, err := db.Query("SELECT p.problem_id, p.problem_title, p.problem_description, p.constraints_desc, p.input_format, p.output_format, p.sample_input, p.sample_output, p.creator_email, p.is_private FROM distribute_problems_to_contest d JOIN problem p ON d.problem_id = p.problem_id WHERE d.contest_id = ?", contestID)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to execute query: %v\n", err)
		return nil, err
	}
	defer rows.Close()

	var contestProblems []ContestProblem

	for rows.Next() {
		var problem ContestProblem

		err := rows.Scan(&problem.ProblemID, &problem.ProblemTitle, &problem.ProblemDesc, &problem.ConstraintsDesc, &problem.InputFormat, &problem.OutputFormat, &problem.SampleInput, &problem.SampleOutput, &problem.CreatorEmail, &problem.IsPrivate)
		if err != nil {
			fmt.Println("Error scanning row:", err)
			return nil, err
		}

		contestProblems = append(contestProblems, problem)
	}

	if err := rows.Err(); err != nil {
		fmt.Println("Error during rows iteration:", err)
		return nil, err
	}

	return contestProblems, nil
}
