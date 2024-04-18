// RUN the code which is not submitted for evaluation
// runtime checking of the code on the open testcases of a problem
// Two ways to implement:
// 		- Save in database and then send to judge server like normal submission of code
// 		- Don't save in database and send to judge server
// check for security and ease of implementation between the two

package run_code

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
)

type Runny struct {
	ProblemID string `json:"problem_id"`
	UserEmail string `json:"user_email"`
	Code      string `json:"code"`
}

type Testcase struct {
	TestInput  []byte
	TestOutput []byte
}

type PrepareForJuding struct {
	test_input  string `json:"test_input"`
	test_output string `json:"test_output"`
	test_code   string `json:"test_code"`
}

func RetrieveTextFromBlob(inputBlob []byte, outputBlob []byte) (string, string, error) {
	inputText := string(inputBlob)
	outputText := string(outputBlob)
	return inputText, outputText, nil
}

func RunCode(c *fiber.Ctx) error {
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

	var runny Runny
	if err := c.BodyParser(&runny); err != nil {
		return err
	}

	problemID := runny.ProblemID
	err2 := godotenv.Load(".env")
	if err2 != nil {
		log.Fatalf("Some error occurred. Err: %s", err2)
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
	input, output, err := fetchTestcases(db, problemID)
	if err != nil {
		return err
	}
	// showTestcases(db)
	fmt.Println("User EMAIL:", userEmail)
	fmt.Println("KI KI KI RRR IN:", input)
	fmt.Println("KI KI KI RRR OUT:", output)

	sendToJudge := PrepareForJuding{
		test_input:  input,
		test_output: output,
		test_code:   runny.Code,
	}

	data, err := json.Marshal(sendToJudge)
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	judgeURL := "http://judge-server:3001"
	resp, err := http.Post(judgeURL, "application/json", bytes.NewReader(data))
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("Received non-OK response status code: %d", resp.StatusCode)
		return c.SendStatus(http.StatusInternalServerError)
	}

	return c.SendString("POST request sent successfully")
	// return c.Status(fiber.StatusCreated).JSON(fiber.Map{
	// 	"message": "Code Ran Successfully!",
	// })
}

func fetchTestcases(db *sql.DB, problemID string) (string, string, error) {
	var testcases Testcase

	row := db.QueryRow("SELECT testcase_input, testcase_output FROM testcase WHERE problem_id = ?", problemID)
	err := row.Scan(&testcases.TestInput, &testcases.TestOutput)
	if err != nil {
		if err == sql.ErrNoRows {
			fmt.Println("No testcases found with ID:", problemID)
			return "", "", err
		}

		fmt.Fprintf(os.Stderr, "failed to scan problem: %v\n", err)
		return "", "", err
	}

	inputText, outputText, err := RetrieveTextFromBlob(testcases.TestInput, testcases.TestOutput)
	if err != nil {
		return "", "", err
	}

	return inputText, outputText, nil
}

func showTestcases(db *sql.DB) error {
	rows, err := db.Query("SELECT * FROM testcase")
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to query testcases: %v\n", err)
		return err
	}
	defer rows.Close()
	fmt.Println(rows)
	return nil
}
