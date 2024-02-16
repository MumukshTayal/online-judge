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

func GetProblem(c *fiber.Ctx, problemID int) error {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Some error occured. Err: %s", err)
	}

	db_name := os.Getenv("DATABASE_NAME")
	db_token := os.Getenv("DATABASE_TOKEN")

	url := fmt.Sprintf("libsql://%s.turso.io?authToken=%s", db_name, db_token)

	db, err := sql.Open("libsql", url)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to open db %s: %s", url, err)
		os.Exit(1)
	}

	if err := db.Ping(); err != nil {
		fmt.Fprintf(os.Stderr, "failed to ping db %s: %s\n", url, err)
		os.Exit(1)
	}
	defer db.Close()

	problem_data, err := fetchProblem(db, problemID)
	// problem_data, err := fetAllProblems(db)
	if err != nil {
		return err
	}

	return c.JSON(problem_data)

}

type Problem struct {
	ProblemID   int
	Title       string
	Description string
	Constraints string
	CreatorID   int
	IsPrivate   bool
}

func fetchProblem(db *sql.DB, problemID int) (Problem, error) {
	var problem Problem

	row := db.QueryRow("SELECT * FROM problem WHERE problem_id = ?", problemID)
	err := row.Scan(&problem.ProblemID, &problem.Title, &problem.Description, &problem.Constraints, &problem.CreatorID, &problem.IsPrivate)

	if err != nil {
		if err == sql.ErrNoRows {
			fmt.Println("No problem found with ID:", problemID)
			return problem, nil
		}

		fmt.Fprintf(os.Stderr, "failed to scan problem: %v\n", err)
		return problem, err
	}

	fmt.Println(problem.ProblemID, problem.Title, problem.Description, problem.Constraints, problem.CreatorID, problem.IsPrivate)
	return problem, nil
}

func fetAllProblems(db *sql.DB) ([]Problem, error) {
	rows, err := db.Query("SELECT * FROM problem")
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to execute query: %v\n", err)
		return nil, err
	}
	defer rows.Close()

	var problems []Problem

	for rows.Next() {
		var problem Problem

		if err := rows.Scan(&problem.ProblemID, &problem.Title, &problem.Description, &problem.Constraints, &problem.CreatorID, &problem.IsPrivate); err != nil {
			fmt.Println("Error scanning row:", err)
			return nil, err
		}

		problems = append(problems, problem)
		fmt.Println(problem.ProblemID, problem.Title, problem.Description, problem.Constraints, problem.CreatorID, problem.IsPrivate)
	}

	if err := rows.Err(); err != nil {
		fmt.Println("Error during rows iteration:", err)
	}
	return problems, nil
}
