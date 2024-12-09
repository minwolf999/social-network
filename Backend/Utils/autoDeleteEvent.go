package utils

import (
	"database/sql"
	"fmt"
	"time"

	model "social-network/Model"
)

func AutoDeleteEvent(db *sql.DB) {
	for range time.Tick(time.Second * 1) {
		var event model.Event
		eventDetail, err := event.SelectFromDb(db, map[string]any{})
		if err != nil {
			fmt.Println("1", err)
			return
		}

		for i := 0; i < len(eventDetail); i++ {
			t, err := time.Parse("2006-01-02 15:04", eventDetail[i].DateOfTheEvent)
			if err != nil {
				fmt.Println("2", err)
				return
			}

			if !t.After(time.Now()) {
				stmt, err := db.Prepare("DELETE FROM Event WHERE Id = ?")
				if err != nil {
					fmt.Println("3", err)
					return
				}

				_, err = stmt.Exec(eventDetail[i].Id)
				if err != nil {
					fmt.Println("4", err)
					return
				}
			}
		}
	}
}
