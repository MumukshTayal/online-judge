package add_problem

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
	_ "github.com/tursodatabase/libsql-client-go/libsql"
)

func AddProblem(c *fiber.Ctx) error {
	authHeader := c.Get("Authorization")
	tokenStr := ""
	if authHeader != "" {
		authValue := strings.Split(authHeader, " ")
		if len(authValue) == 2 && authValue[0] == "Bearer" {
			tokenStr = authValue[1]
		}
	}

	token, errr := jwt.Parse(tokenStr, nil)
	if token == nil {
		fmt.Println("err --> ", errr)
		return nil
	}
	claims, _ := token.Claims.(jwt.MapClaims)
	fmt.Println("email --> ", claims["email"])

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

	var creator_email = claims["email"]

	if err := AddProblemtoProblemTable(c, db, creator_email.(string)); err != nil {
		return err
	}

	return c.SendString("Problem added successfully!")

}

type Problem struct {
	ProblemID          int    `json:"problem_id"`
	ProblemTitle       string `json:"title"`
	ProblemDescription string `json:"statement"`
	ProblemConstraints string `json:"constraints"`
	SampleInput        string `json:"sampleInput"`
	SampleOutput       string `json:"sampleOutput"`
	InputFormat        string `json:"input"`
	OutputFormat       string `json:"output"`
	IsPrivate          string `json:"isPrivate"`
}

func AddProblemtoProblemTable(c *fiber.Ctx, db *sql.DB, email string) error {
	var problem Problem
	if err := c.BodyParser(&problem); err != nil {
		return err
	}

	b, err := strconv.ParseBool(problem.IsPrivate)
	if err != nil {
		fmt.Println(err)
		return err
	}

	stmt, err := db.Prepare("INSERT INTO problem (problem_title, problem_description, constraints_desc, input_format, output_format, sample_input, sample_output, creator_email, is_private) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(problem.ProblemTitle, problem.ProblemDescription, problem.ProblemConstraints, problem.InputFormat, problem.OutputFormat, problem.SampleInput, problem.SampleOutput, email, b)
	if err != nil {
		return err
	}

	return nil
}
