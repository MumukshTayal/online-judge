package add_contest

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

func AddContests(c *fiber.Ctx) error {
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

	contest_data := AddContestsToContestTable(c, db)
	if err != nil {
		return err
	}

	return c.JSON(contest_data)

}

type Contest struct {
	ContestID          int       `json:"contest_id"`
	ContestTitle       string    `json:"contest_title"`
	ContestDescription string    `json:"contest_description"`
	StartTime          time.Time `json:"start_time"`
	EndTime            time.Time `json:"end_time"`
	IsPublic           bool      `json:"is_public"`
	CreatorID          int       `json:"creator_id"`
}

func AddContestsToContestTable(c *fiber.Ctx, db *sql.DB) error {
	var contest Contest
	if err := c.BodyParser(&contest); err != nil {
		return err
	}

	stmt, err := db.Prepare("INSERT INTO contest (contest_title, contest_description, contest_start_time, contest_end_time, is_public, creator_id) VALUES (?, ?, ?, ?, ?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(contest.ContestTitle, contest.ContestDescription, contest.StartTime, contest.EndTime, contest.IsPublic, contest.CreatorID)
	if err != nil {
		return err
	}

	return c.JSON(contest)
}
