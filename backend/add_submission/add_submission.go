package add_submission

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
)

type Submission struct {
	ProblemID          string    `json:"problem_id"`
	UserEmail          string    `json:"user_email"`
	Code               string    `json:"code"`
	Language           string    `json:"language"`
	Result             string    `json:"result"`
	ExecutionTime      int       `json:"execution_time"`
	MemoryUsed         int       `json:"memory_used"`
	SubmissionDateTime time.Time `json:"submission_date_time"`
}

func AddSubmission(c *fiber.Ctx) error {
	authHeader := c.Get("Authorization")
	tokenStr := ""
	if authHeader != "" {
		authValue := strings.Split(authHeader, " ")
		if len(authValue) == 2 && authValue[0] == "Bearer" {
			tokenStr = authValue[1]
		}
	}

	token, err := jwt.Parse(tokenStr, nil)
	if token == nil {
		fmt.Println("Error parsing token:", err)
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Unauthorized",
		})
	}
	claims, _ := token.Claims.(jwt.MapClaims)
	userEmail, ok := claims["email"].(string)
	if !ok {
		fmt.Println("Invalid token claims")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid token claims",
		})
	}

	err = godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file: %s", err)
	}

	db_name := os.Getenv("DATABASE_NAME")
	db_token := os.Getenv("DATABASE_TOKEN")

	url := fmt.Sprintf("libsql://%s.turso.io?authToken=%s", db_name, db_token)

	db, err := sql.Open("libsql", url)
	if err != nil {
		log.Fatalf("Error connecting to the database: %s", err)
	}
	defer db.Close()

	var submission Submission
	if err := c.BodyParser(&submission); err != nil {
		return err
	}

	stmt, err := db.Prepare("INSERT INTO submission (problem_id, user_email, code, language, result, execution_time, memory_used, submission_date_time) VALUES (?, ?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	// Initially when the User submits the code through frontend, it will have only the following fields:
	// problem_id, user_email, code, language, submission_date_time

	// After running the code using judge server, the database of the submission will update to contain:
	// result, execution_time, memory_used
	_, err = stmt.Exec(submission.ProblemID, userEmail, submission.Code, submission.Language, submission.Result, submission.ExecutionTime, submission.MemoryUsed, submission.SubmissionDateTime)

	fmt.Println(submission.ProblemID, userEmail, submission.Code, submission.Language, submission.Result, submission.ExecutionTime, submission.MemoryUsed, submission.SubmissionDateTime)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Submission added successfully!",
	})
}
