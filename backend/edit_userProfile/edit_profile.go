package edit_userProfile

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	_ "github.com/tursodatabase/libsql-client-go/libsql"
)

func EditUserProfile(c *fiber.Ctx) error {
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

	if err := editUserProfileDetails(c, db); err != nil {
		return err
	}

	return c.SendString("User details modified successfully!")

}

type UserProfile struct {
	UserID    int    `json:"user_id"`
	UserName  string `json:"user_name"`
	UserEmail string `json:"user_email"`
}

func editUserProfileDetails(c *fiber.Ctx, db *sql.DB) error {
	var user UserProfile
	if err := c.BodyParser(&user); err != nil {
		return err
	}

	stmt, err := db.Prepare("UPDATE user_profile SET user_name = ?, user_email = ? WHERE user_id = ?")
	if err != nil {
		return err
	}

	_, err = stmt.Exec(user.UserName, user.UserEmail, user.UserID)
	if err != nil {
		return err
	}

	return nil
}
