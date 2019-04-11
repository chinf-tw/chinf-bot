package userinfo

import (
	"database/sql"
	"log"
)

func getUserID(db *sql.DB) []string {
	query := `select line_id from spotify_user ;`
	userIDList := make([]string, 0)
	rows, err := db.Query(query)
	if err != nil {
		log.Println(err)
	}
	defer rows.Close()
	for rows.Next() {
		var userID string
		if err := rows.Scan(&userID); err != nil {
			log.Println(err)
		}
		userIDList = append(userIDList, userID)
	}
	return userIDList
}
