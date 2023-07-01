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
	sql := `select a.id, a.name, a.topic, a.src,a.upvotes, a.downvotes,b.option from polls a
	left join options b
	on a.id = b.pollid`

	rows, err := d.db.Query(sql)

	if err != nil {
		panic(err)
	}

	defer rows.Close()

	result := model.PollCollection{}

	log.Println("Fetched rows from the db")

	for rows.Next() {
		poll := model.Poll{}

		err2 := rows.Scan(&poll.ID, &poll.Name, &poll.Topic, &poll.Src, &poll.Upvotes, &poll.Downvotes, &poll.Options)

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

func (d *PostgresDB) CreateOption (option string, pollId int) (int, error) {
	// query := fmt.Sprintf("INSERT INTO options( option, pollId) VALUES('%s',%d)", option, pollId)
	sql := `INSERT INTO options(option, pollId) 
	VALUES ($1, $2) returning id`
	id := 0
	err := d.db.QueryRow(sql, option, pollId).Scan(&id)
  	fmt.Println("New Option ID is:", id)
	  if err != nil {
		panic(err)
	}

	return id,nil
}

func (d *PostgresDB) CreatePoll (name string, topic string, src string ) (int, error ){
	fmt.Println("create post db func",name,topic,src)
	
	sqlStatement := `
	INSERT INTO polls(name, topic, src, upvotes, downvotes) 
	VALUES ($1, $2, $3, $4, $5) returning id`

	fmt.Println("Sql Statement", sqlStatement)
	id := 0
 	err := d.db.QueryRow(sqlStatement, name, topic, src, 0, 0).Scan(&id)
  	fmt.Println("New Poll ID is:", id)

	if err != nil {
		panic(err)
	}

	return  id,nil
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
func (d *PostgresDB) CreateUser (name string, email string, password string) (model.User, error) {
	fmt.Println("creating user ",name,email,password)
 
	user, userErr := d.FetchUser(email)
	log.Println("User fetched from DB", user)
	if userErr != nil {
		log.Fatal( "Error fetching user with email")
		panic(userErr)
	}
	if user.Email == email {
		log.Println("User already exists with this email")
		return user, errors.New("user already exists with this email")

	}
	query := fmt.Sprintf("INSERT INTO users(name, email, password) VALUES('%s','%s', '%s')", name, email, password)
	_,err := d.db.Query(query) 
	fetchedUser, userErr := d.FetchUser(email)

	if userErr != nil {
		panic(err)
	}
	return fetchedUser,nil
}

func (d *PostgresDB) SaveToken(token string , userId int) (error) {
	fmt.Println("Saving token to DB ", token, userId )
 
	sql:= fmt.Sprintf("INSERT into tokens (token, userId) values( '%s' , %d)", token, userId);
	_, err := d.db.Query(sql)
	if err != nil {
		panic(err)
	}
	return nil

}




func (d *PostgresDB) Migrate (){
	sql := `
	DROP TABLE polls;
	DROP TABLE users;
	DROP TABLE tokens;
	drop table options;


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
		
		INSERT INTO users (name, email, password) VALUES(
			'Aman Kumar','aman@gmail.com', 'https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcT7F7Ca089qQJSIBKJuWNC2Wnb9nmtsMhvgYtXDa7-9jA&s');
		  
			CREATE TABLE IF NOT EXISTS tokens (
				id SERIAL NOT NULL PRIMARY KEY,
				userId int UNIQUE NOT NULL,
				token VARCHAR(255) UNIQUE NOT NULL
			  ); 


	INSERT INTO tokens (token,userId) VALUES(
				'346ae688-e640-45e5-9e84-f526df595f0d',1);

	CREATE TABLE IF NOT EXISTS options (
				id SERIAL NOT NULL PRIMARY KEY,
				pollId int NOT NULL,
				option VARCHAR(255) NOT NULL
			  ); 
			  

	`
	  rows, err := d.db.Query(sql)

	  if err != nil {
		  panic(err)
	  }
  
	  defer rows.Close()
}
	

func (d *PostgresDB) FetchUserToken(userId int) (string, error) {
	sql:= fmt.Sprintf("Select token from tokens where id = %d", userId);
	tokenToUserId := model.TokenToUserId{}

	err:= d.db.QueryRow(sql).Scan(&tokenToUserId.Token)
	if err != nil {
		log.Fatal( "Error fetching token with userId", userId)	
	}

	return tokenToUserId.Token, nil
}
