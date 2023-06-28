package database

import (
	"log"
	"fmt"
	"polling-app/model"
	"database/sql"

	"errors"
)

type PostgresDB struct {
	db *sql.DB
}

func NewDB(db *sql.DB) PostgresDB {
	return PostgresDB{db: db}
}

func (d *PostgresDB) GetPolls() (model.PollCollection, error) {
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
	if result.Polls == nil {
		fmt.Println("result polls", result.Polls)
		result.Polls = make([]model.Poll, 0)
	}

	return result,nil
}
func (d *PostgresDB) GetPoll(id int) (model.Poll,error) {
	query:= fmt.Sprintf("Select * from polls where id = %d", id)
	rows, err := d.db.Query(query)
	if err != nil {
		panic(err) 
	}
	defer rows.Close()
	poll := model.Poll{}

	for rows.Next() {

		err2 := rows.Scan(&poll.ID, &poll.Name, &poll.Topic, &poll.Src, &poll.Upvotes, &poll.Downvotes)

		if err2 != nil {
			panic(err2)
		}
		fmt.Println("What's inside the rows",poll)
		return poll,nil

	}
	return poll,nil
}
func (d *PostgresDB) UpdatePoll( id int, name string, upvotes int, downvotes int) error {
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

func (d *PostgresDB) CreatePoll (name string, topic string, src string) error {
	fmt.Println("create post db func",name,topic,src)
	query := fmt.Sprintf("INSERT INTO polls(name, topic, src, upvotes, downvotes) VALUES('%s','%s', '%s', 0, 0)", name,topic, src)
	fmt.Print(query)
	_, err := d.db.Query(query)
	
	if err != nil {
		panic(err)
	}
	return nil
}

func (d *PostgresDB) DeletePoll(id int) error{
	query := fmt.Sprintf("DELETE FROM polls where id = %d", id);
	_, err := d.db.Query(query);
	if err != nil {
		panic(err)
	}
	return nil

}

func (d *PostgresDB) FetchUser(email string) (model.User, error){
	sql := fmt.Sprintf("Select id, email, password from users where email = '%s'", email)
	user := model.User{}

	// row:= d.db.QueryRow(sql)
	err:= d.db.QueryRow(sql).Scan(&user.ID, &user.Email, &user.Password)
	log.Println("UserInfo in DB",user.ID, user.Email, user.Password)
	if user.ID == 0 {
		return model.User{}, nil
	}

	if err != nil {
		panic(err)
	} else {
		fmt.Println(user.ID, user.Email)
	}
	return user, nil

}
func (d *PostgresDB) CreateUser (name string, email string, password string) error {
	fmt.Println("creating user ",name,email,password)
 
	user, userErr := d.FetchUser(email)
	log.Println("User fetched from DB", user)
	if userErr != nil {
		log.Fatal( "Error fetching user with email")
		panic(userErr)
	}
	if user.Email == email {
		log.Println("User already exists with this email")
		return errors.New("User already exists with this email")

	}
	query := fmt.Sprintf("INSERT INTO users(name, email, password) VALUES('%s','%s', '%s')", name, email, password)
	fmt.Print(query)
	_, err := d.db.Query(query)
	
	if err != nil {
		panic(err)
	}
	return nil
}





func (d *PostgresDB) Migrate (){
	sql := `
	DROP TABLE polls;
	DROP TABLE users;


	CREATE TABLE IF NOT EXISTS polls (
		id SERIAL NOT NULL PRIMARY KEY,
					  name VARCHAR(255) UNIQUE NOT NULL,
					  topic VARCHAR(255),
					  src VARCHAR NOT NULL,
					  upvotes INT NOT NULL,
					  downvotes INT NOT NULL
	  );

	  INSERT INTO polls (name, topic, src, upvotes, downvotes) VALUES(
		'Angular','Awesome Angular', 'https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcT7F7Ca089qQJSIBKJuWNC2Wnb9nmtsMhvgYtXDa7-9jA&s', 1, 0
	  );
	  
	  INSERT INTO polls(name, topic, src, upvotes, downvotes) VALUES(
		'Vue', 'Voguish Vue','https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcTrRNrlbypiGF07Uxf_PyVpJv4ipjKCc6s13EKqsLxnHg&s	', 1, 0
	  );
	  
	  INSERT  INTO polls(name, topic, src, upvotes, downvotes) VALUES(
		'React','Remarkable React','https://upload.wikimedia.org/wikipedia/commons/thumb/a/a7/React-icon.svg/1200px-React-icon.svg.png', 1, 0);
		  
	  INSERT INTO polls(name, topic, src, upvotes, downvotes) VALUES(
		'Ember','Excellent Ember','https://cdn-images-1.medium.com/max/741/1*9oD6P0dEfPYp3Vkk2UTzCg.png', 1, 0);
		  
	  INSERT INTO polls(name, topic, src, upvotes, downvotes) VALUES(
		'Knockout','Knightly Knockout','https://images.g2crowd.com/uploads/product/image/social_landscape/social_landscape_1489710848/knockout-js.png', 1, 0);
		 
		CREATE TABLE IF NOT EXISTS users (
			id SERIAL NOT NULL PRIMARY KEY,
						  name VARCHAR(255) NOT NULL,
						  email VARCHAR(255) unique NOT NULL,
						  password VARCHAR  NOT NULL	
		);
		INSERT INTO users (name,email, password) VALUES(
			'Aman Kumar','aman@gmail.com', 'https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcT7F7Ca089qQJSIBKJuWNC2Wnb9nmtsMhvgYtXDa7-9jA&s');
		  
	`
	  rows, err := d.db.Query(sql)

	  if err != nil {
		  panic(err)
	  }
  
	  defer rows.Close()
}
	
