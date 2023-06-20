package database

import (
	"fmt"
	"polling-app/model"
	"database/sql"
)

type DB interface {
	GetPolls() (model.PollCollection,error)
	UpdatePoll( index int, name string, upvotes int, downvotes int) (int64, error)
}

type PostgresDB struct {
	db *sql.DB
}

func NewDB(db *sql.DB) DB {
	return PostgresDB{db: db}
}

func (d PostgresDB) GetPolls() (model.PollCollection, error) {
	sql := "SELECT * FROM polls"

	rows, err := d.db.Query(sql)

	if err != nil {
		panic(err)
	}

	defer rows.Close()

	result := model.PollCollection{}

	for rows.Next() {
		poll := model.Poll{}

		err2 := rows.Scan(&poll.ID, &poll.Name, &poll.Topic, &poll.Src, &poll.Upvotes, &poll.Downvotes)

		if err2 != nil {
			panic(err2)
		}

		result.Polls = append(result.Polls, poll)
	}

	return result,nil
}

func (d PostgresDB) UpdatePoll( index int, name string, upvotes int, downvotes int) (int, error) {
	sql := fmt.Sprintf("UPDATE polls SET (upvotes, downvotes) = (%d, %d) WHERE id = %d",upvotes,downvotes,index)


	rows, err := d.db.Query(sql)
	// Exit if we get an error
	if err != nil {
		panic(err)
	}

	// Make sure to cleanup after the program exits
	defer rows.Close()

	return index,nil
}


func migrate(db *sql.DB) {
	sql := `
		CREATE TABLE IF NOT EXISTS polls(
				id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
				name VARCHAR NOT NULL,
				topic VARCHAR NOT NULL,
				src VARCHAR NOT NULL,
				upvotes INTEGER NOT NULL,
				downvotes INTEGER NOT NULL,
				UNIQUE(name)
		);
	
		INSERT OR IGNORE INTO polls(name, topic, src, upvotes, downvotes) VALUES('Angular','Awesome Angular', 'https://cdn.colorlib.com/wp/wp-content/uploads/sites/2/angular-logo.png', 1, 0);
	
		INSERT OR IGNORE INTO polls(name, topic, src, upvotes, downvotes) VALUES('Vue', 'Voguish Vue','https://upload.wikimedia.org/wikipedia/commons/thumb/5/53/Vue.js_Logo.svg/400px-Vue.js_Logo.svg.png', 1, 0);
	
		INSERT OR IGNORE INTO polls(name, topic, src, upvotes, downvotes) VALUES('React','Remarkable React','https://upload.wikimedia.org/wikipedia/commons/thumb/a/a7/React-icon.svg/1200px-React-icon.svg.png', 1, 0);
	
		INSERT OR IGNORE INTO polls(name, topic, src, upvotes, downvotes) VALUES('Ember','Excellent Ember','https://cdn-images-1.medium.com/max/741/1*9oD6P0dEfPYp3Vkk2UTzCg.png', 1, 0);
	
		INSERT OR IGNORE INTO polls(name, topic, src, upvotes, downvotes) VALUES('Knockout','Knightly Knockout','https://images.g2crowd.com/uploads/product/image/social_landscape/social_landscape_1489710848/knockout-js.png', 1, 0);
	`
	_, err := db.Exec(sql)
	
	if err != nil {
			panic(err)
	}
}
	
