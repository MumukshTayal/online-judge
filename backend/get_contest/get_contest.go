package get_contest

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	_ "github.com/tursodatabase/libsql-client-go/libsql"
)

type Contest struct {
	ContestID    int       `json:"contest_id"`
	ContestTitle string    `json:"contest_title"`
	ContestDesc  string    `json:"contest_description"`
	StartTime    time.Time `json:"start_time"`
	EndTime      time.Time `json:"end_time"`
	IsPublic     bool      `json:"is_public"`
	CreatorEmail string    `json:"creator_email"`
}

func GetAllContests(c *fiber.Ctx) error {
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

	contests, err := fetchAllContests(db)
	if err != nil {
		return c.Status(500).SendString(err.Error())
	}

	return c.JSON(contests)
}

func fetchAllContests(db *sql.DB) ([]Contest, error) {
	rows, err := db.Query("SELECT contest_id, contest_title, contest_description, contest_start_time, contest_end_time, COALESCE(is_public, 0), creator_email FROM contest")
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to execute query: %v\n", err)
		return nil, err
	}
	defer rows.Close()

	var contests []Contest

	for rows.Next() {
		var contest Contest
		var isPublicInt sql.NullInt64
		var creatorEmail sql.NullString

		err := rows.Scan(&contest.ContestID, &contest.ContestTitle, &contest.ContestDesc, &contest.StartTime, &contest.EndTime, &isPublicInt, &creatorEmail)
		if err != nil {
			fmt.Println("Error scanning row:", err)
			return nil, err
		}

		contest.IsPublic = isPublicInt.Valid && isPublicInt.Int64 != 0
		if creatorEmail.Valid {
			contest.CreatorEmail = creatorEmail.String
		} else {
			contest.CreatorEmail = ""
		}

		contests = append(contests, contest)
	}

	if err := rows.Err(); err != nil {
		fmt.Println("Error during rows iteration:", err)
		return nil, err
	}

	return contests, nil
}

func ContestList(c *fiber.Ctx) error {
	return nil
}
