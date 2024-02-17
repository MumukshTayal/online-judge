package add_problem

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	_ "github.com/tursodatabase/libsql-client-go/libsql"
)

func AddProblem(c *fiber.Ctx) error {
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

	if err := AddProblemtoProblemTable(c, db); err != nil {
		return err
	}

	return c.SendString("Problem added successfully!")

}

type Problem struct {
	ProblemID          int    `json:"problem_id"`
	ProblemTitle       string `json:"problem_title"`
	ProblemDescription string `json:"problem_description"`
	ProblemConstraints string `json:"problem_constraints"`
	CreatorID          int    `json:"creator_id"`
	IsPrivate          bool   `json:"is_private"`
}

func AddProblemtoProblemTable(c *fiber.Ctx, db *sql.DB) error {
	var problem Problem
	if err := c.BodyParser(&problem); err != nil {
		return err
	}

	stmt, err := db.Prepare("INSERT INTO problem (problem_title, problem_description, constraints_desc, creator_id, is_private) VALUES (?, ?, ?, ?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(problem.ProblemTitle, problem.ProblemDescription, problem.ProblemConstraints, problem.CreatorID, problem.IsPrivate)
	if err != nil {
		return err
	}

	return nil
}
