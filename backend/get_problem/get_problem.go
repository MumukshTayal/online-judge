package get_problem

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	_ "github.com/tursodatabase/libsql-client-go/libsql"
)

func GetProblemByProblemId(c *fiber.Ctx) error {
	problemID := c.Params("problemId")

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
		return err
	}
	defer db.Close()

	problemData, err := fetchProblem(db, problemID)
	if err != nil {
		return err
	}

	return c.JSON(problemData)
}

type Problem struct {
	ProblemID          int
	ProblemTitle       string
	ProblemDescription string
	ConstraintsDesc    string
	InputFormat        string
	OutputFormat       string
	SampleInput        string
	SampleOutput       string
	CreatorEmail       string
	IsPrivate          bool
}

type ProblemDisplay struct {
	ProblemID    int
	ProblemTitle string
	CreatorEmail string
	IsPrivate    bool
}

func fetchProblem(db *sql.DB, problemID string) (Problem, error) {
	var problem Problem

	row := db.QueryRow("SELECT * FROM problem WHERE problem_id = ?", problemID)
	err := row.Scan(&problem.ProblemID, &problem.ProblemTitle, &problem.ProblemDescription, &problem.ConstraintsDesc, &problem.InputFormat, &problem.OutputFormat, &problem.SampleInput, &problem.SampleOutput, &problem.CreatorEmail, &problem.IsPrivate)

	if err != nil {
		if err == sql.ErrNoRows {
			fmt.Println("No problem found with ID:", problemID)
			return problem, nil
		}

		fmt.Fprintf(os.Stderr, "failed to scan problem: %v\n", err)
		return problem, err
	}

	// fmt.Println(problem.ProblemID, problem.ProblemTitle, problem.ProblemDescription, problem.ConstraintsDesc, problem.InputFormat, problem.OutputFormat, problem.SampleInput, problem.SampleOutput, problem.CreatorEmail, problem.IsPrivate)
	return problem, nil
}

func GetAllProblems(c *fiber.Ctx) error {
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
		return err
	}
	defer db.Close()

	problemData, err := fetchAllProblems(db)
	if err != nil {
		return err
	}

	return c.JSON(problemData)
}

func fetchAllProblems(db *sql.DB) ([]ProblemDisplay, error) {
	var problems []ProblemDisplay
	rows, err := db.Query("SELECT problem_id, problem_title, creator_email, is_private FROM problem")
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to query problems: %v\n", err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var problem ProblemDisplay
		err := rows.Scan(&problem.ProblemID, &problem.ProblemTitle, &problem.CreatorEmail, &problem.IsPrivate)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to scan problem: %v\n", err)
			return nil, err
		}
		problems = append(problems, problem)
	}

	if err := rows.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "error iterating over rows: %v\n", err)
		return nil, err
	}

	return problems, nil
}
