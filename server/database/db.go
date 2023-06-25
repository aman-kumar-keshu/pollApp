package database

import (
	"fmt"
	"polling-app/model"
	"database/sql"
)

type DB interface {
	GetPolls() (model.PollCollection,error)
	UpdatePoll( index int, name string, upvotes int, downvotes int) error
	Migrate()
}

type PostgresDB struct {
	db *sql.DB
}

func NewDB(db *sql.DB) DB {
	return PostgresDB{db: db}
}

func (d PostgresDB) GetPolls() (model.PollCollection, error) {
	sql := "SELECT * FROM polls order by id"

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

func (d PostgresDB) UpdatePoll( id int, name string, upvotes int, downvotes int) error {
	sql := fmt.Sprintf("UPDATE polls SET (upvotes, downvotes) = (%d, %d) WHERE id = %d",upvotes,downvotes,id)


	rows, err := d.db.Query(sql)
	// Exit if we get an error
	if err != nil {
		panic(err)
	}
	var poll model.Poll
	rows.Scan(&poll)
	defer rows.Close()
	fmt.Println(poll)
	return nil
}


func (d PostgresDB) Migrate (){
	sql := `
	DROP TABLE polls;

	CREATE TABLE IF NOT EXISTS polls (
		id SERIAL NOT NULL PRIMARY KEY,
					  name VARCHAR(255) UNIQUE NOT NULL,
					  topic VARCHAR(255),
					  src VARCHAR NOT NULL,
					  upvotes INT NOT NULL,
					  downvotes INT NOT NULL
	  );
	
	  INSERT INTO polls (name, topic, src, upvotes, downvotes) VALUES(
		'Angular','Awesome Angular', 'https://cdn.colorlib.com/wp/wp-content/uploads/sites/2/angular-logo.png', 1, 0
	  );
	  
	  INSERT INTO polls(name, topic, src, upvotes, downvotes) VALUES(
		'Vue', 'Voguish Vue','https://upload.wikimedia.org/wikipedia/commons/thumb/5/53/Vue.js_Logo.svg/400px-Vue.js_Logo.svg.png', 1, 0
	  );
	  
	  INSERT  INTO polls(name, topic, src, upvotes, downvotes) VALUES(
		'React','Remarkable React','https://upload.wikimedia.org/wikipedia/commons/thumb/a/a7/React-icon.svg/1200px-React-icon.svg.png', 1, 0);
		  
	  INSERT INTO polls(name, topic, src, upvotes, downvotes) VALUES(
		'Ember','Excellent Ember','https://cdn-images-1.medium.com/max/741/1*9oD6P0dEfPYp3Vkk2UTzCg.png', 1, 0);
		  
	  INSERT INTO polls(name, topic, src, upvotes, downvotes) VALUES(
		'Knockout','Knightly Knockout','https://images.g2crowd.com/uploads/product/image/social_landscape/social_landscape_1489710848/knockout-js.png', 1, 0);
		  
	  `
	  rows, err := d.db.Query(sql)

	  if err != nil {
		  panic(err)
	  }
  
	  defer rows.Close()
}
	
