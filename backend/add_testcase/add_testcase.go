package add_testcase

import (
	"database/sql"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"os"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
)

type TestCase struct {
	ProblemName string                `form:"problem_name"`
	Input       *multipart.FileHeader `form:"input"`
	Output      *multipart.FileHeader `form:"output"`
}

func AddTestCase(c *fiber.Ctx) error {
	authHeader := c.Get("Authorization")
	tokenStr := ""
	if authHeader != "" {
		authValue := strings.Split(authHeader, " ")
		if len(authValue) == 2 && authValue[0] == "Bearer" {
			tokenStr = authValue[1]
		}
	}

	// fmt.Println("token str --> ", tokenStr)

	token, err := jwt.Parse(tokenStr, nil)
	if token == nil {
		fmt.Println("err --> ", err)
		return nil
	}
	claims := token.Claims.(jwt.MapClaims)
	fmt.Println("email --> ", claims["email"])

	err = godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Some error occured. Err: %s", err)
	}

	dbName := os.Getenv("DATABASE_NAME")
	dbToken := os.Getenv("DATABASE_TOKEN")

	url := fmt.Sprintf("libsql://%s.turso.io?authToken=%s", dbName, dbToken)

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

	var creatorEmail = claims["email"].(string)

	if err := AddTestCaseToDB(c, db, creatorEmail); err != nil {
		return err
	}

	return c.SendString("Test case added successfully!")
}

func AddTestCaseToDB(c *fiber.Ctx, db *sql.DB, creatorEmail string) error {
	form, err := c.MultipartForm()
	if err != nil {
		return err
	}

	var testCase TestCase
	testCase.ProblemName = form.Value["problem_name"][0]
	testCase.Input = form.File["input"][0]
	testCase.Output = form.File["output"][0]

	// fmt.Println("test case --> ", testCase)

	// Query the problem table to get problem_id and creator_email
	// fmt.Println("before select query! --> ", testCase.ProblemName)
	row := db.QueryRow("SELECT problem_id, creator_email FROM problem WHERE problem_title = ?", testCase.ProblemName)
	var problemID int
	var problemCreatorEmail string
	errr := row.Scan(&problemID, &problemCreatorEmail)
	if errr != nil {
		return errr
	}

	// Check if the creator email from the problem table matches the one from JWT claims
	if creatorEmail != problemCreatorEmail {
		return fmt.Errorf("Unauthorized access: creator email mismatch")
	}

	// Insert test case into the database
	stmt, err := db.Prepare("INSERT INTO testcase (testcase_input, testcase_output, problem_id) VALUES (?, ?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	inputFile, err := testCase.Input.Open()
	if err != nil {
		return err
	}
	defer inputFile.Close()

	outputFile, err := testCase.Output.Open()
	if err != nil {
		return err
	}
	defer outputFile.Close()

	// Read input and output files
	inputData, err := io.ReadAll(inputFile)
	if err != nil {
		return err
	}

	outputData, err := io.ReadAll(outputFile)
	if err != nil {
		return err
	}

	_, err = stmt.Exec(inputData, outputData, problemID)
	if err != nil {
		return err
	}

	return nil
}
