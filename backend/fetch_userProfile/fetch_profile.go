package fetch_userProfile

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	_ "github.com/tursodatabase/libsql-client-go/libsql"
)

type UserProfile struct {
	UserID    int    `json:"user_id"`
	UserName  string `json:"user_name"`
	UserEmail string `json:"user_email"`
}

func GetUserProfileByUserId(c *fiber.Ctx, userId int) (UserProfile, error) {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Some error occurred. Err: %s", err)
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
	row, err := db.Query("SELECT * FROM user_profile WHERE user_id = ?", userId)
	if err != nil {
		return UserProfile{}, err
	}
	defer row.Close()

	var userProfile UserProfile

	if row.Next() {
		err := row.Scan(&userProfile.UserID, &userProfile.UserName, &userProfile.UserEmail)
		if err != nil {
			return UserProfile{}, err
		}
	}

	if err := row.Err(); err != nil {
		return UserProfile{}, err
	}

	return userProfile, nil
}
