package database

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"polling-app/model"
)

type PostgresDB struct {
	db *sql.DB
}

func NewDB(db *sql.DB) PostgresDB {
	return PostgresDB{db: db}
}

func (d *PostgresDB) GetPolls() (model.PollCollection, error) {
	sql := `SELECT a.id, a.name, a.topic, a.src, a.upvotes, a.downvotes
			FROM polls a`
	rows, err := d.db.Query(sql)

	if err != nil {
		panic(err)
	}

	defer rows.Close()

	result := model.PollCollection{}

	log.Println("Fetched rows from the db")

	for rows.Next() {
		poll := model.Poll{}

		err2 := rows.Scan(&poll.ID, &poll.Name, &poll.Topic, &poll.Src, &poll.Upvotes, &poll.Downvotes)

		if err2 != nil {
			panic(err2)
		}

		sql = `SELECT id, option, votes, pollId FROM options WHERE pollId = $1`

		rows2, err3 := d.db.Query(sql, poll.ID)

		if err3 != nil {
			panic(err3)
		}

		defer rows2.Close()

		for rows2.Next() {
			var option model.Option
			err4 := rows2.Scan(&option.ID, &option.Option, &option.Votes, &option.PollId)

			if err4 != nil {
				panic(err4)
			}

			poll.Options = append(poll.Options, option)
		}

		result.Polls = append(result.Polls, poll)
	}
	if result.Polls == nil {
		fmt.Println("result polls", result.Polls)
		result.Polls = make([]model.Poll, 0)
	}

	return result, nil
}

func (d *PostgresDB) GetPoll(id int) (model.Poll, error) {
	query := fmt.Sprintf("Select * from polls where id = %d", id)
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
		fmt.Println("What's inside the rows", poll)
		return poll, nil

	}
	return poll, nil
}

func (d *PostgresDB) UpdatePoll(id int, name string, upvotes int, downvotes int) error {
	sql := fmt.Sprintf("UPDATE polls SET (upvotes, downvotes) = (%d, %d) WHERE id = %d", upvotes, downvotes, id)

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

func (d *PostgresDB) CreateOption(option string, pollId int) (int, error) {
	// query := fmt.Sprintf("INSERT INTO options( option, pollId) VALUES('%s',%d)", option, pollId)
	sql := `INSERT INTO options(option, pollId, votes) 
	VALUES ($1, $2, $3) returning id`
	id := 0
	err := d.db.QueryRow(sql, option, pollId, 0).Scan(&id)
	fmt.Println("New Option ID is:", id)
	if err != nil {
		panic(err)
	}

	return id, nil
}

func (d *PostgresDB) UpdateOption(id int,votes int) (int, error) {
	sql := fmt.Sprintf("UPDATE options SET votes = (%d) WHERE id = %d", votes, id)

	rows, err := d.db.Query(sql)
	// Exit if we get an error
	if err != nil {
		panic(err)
	}
	var option model.Option
	rows.Scan(&option)
	defer rows.Close()
	fmt.Println(option)
	return id,nil
}


func (d *PostgresDB) CreatePoll(name string, topic string, src string) (int, error) {
	fmt.Println("create post db func", name, topic, src)

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

	return id, nil
}

func (d *PostgresDB) DeletePoll(id int) error {
	query := fmt.Sprintf("DELETE FROM polls where id = %d", id)
	_, err := d.db.Query(query)
	if err != nil {
		panic(err)
	}
	return nil

}

func (d *PostgresDB) FetchUser(email string) (model.User, error) {
	sql := fmt.Sprintf("Select id, email, password from users where email = '%s'", email)
	user := model.User{}

	// row:= d.db.QueryRow(sql)
	err := d.db.QueryRow(sql).Scan(&user.ID, &user.Email, &user.Password)
	log.Println("UserInfo in DB", user.ID, user.Email, user.Password)
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
func (d *PostgresDB) CreateUser(name string, email string, password string) (model.User, error) {
	fmt.Println("creating user ", name, email, password)

	user, userErr := d.FetchUser(email)
	log.Println("User fetched from DB", user)
	if userErr != nil {
		log.Fatal("Error fetching user with email")
		panic(userErr)
	}
	if user.Email == email {
		log.Println("User already exists with this email")
		return user, errors.New("user already exists with this email")

	}
	query := fmt.Sprintf("INSERT INTO users(name, email, password) VALUES('%s','%s', '%s')", name, email, password)
	_, err := d.db.Query(query)
	fetchedUser, userErr := d.FetchUser(email)

	if userErr != nil {
		panic(err)
	}
	return fetchedUser, nil
}

func (d *PostgresDB) SaveToken(token string, userId int) error {
	fmt.Println("Saving token to DB ", token, userId)

	sql := fmt.Sprintf("INSERT into tokens (token, userId) values( '%s' , %d)", token, userId)
	_, err := d.db.Query(sql)
	if err != nil {
		panic(err)
	}
	return nil

}

func (d *PostgresDB) Migrate() {
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
		'Who is the best Cricketer ','Cricket', 'https://images.cnbctv18.com/wp-content/uploads/2022/07/Cricket-Shutterstock-1019x573.jpg?im=FitAndFill,width=1200,height=900', 1, 0
	  );
	  
	  INSERT INTO polls(name, topic, src, upvotes, downvotes) VALUES(
		'What is the best Frontend framework', 'FrontEnd','https://www.vshsolutions.com/wp-content/uploads/2020/04/blog-featured-choosing-front-end-framework.png	', 1, 0
	  );
	  
	  INSERT  INTO polls(name, topic, src, upvotes, downvotes) VALUES(
		'What is the best Backend framework','Backend Framework','data:image/jpeg;base64,/9j/4AAQSkZJRgABAQAAAQABAAD/2wCEAAkGBxASEhUSEBIVFRUVGBUWFxgXFRUWFRgVFxgaFhgXGxYYHSggGBolGxgVLTEhJSsrLi4uFyAzODMtNygtLisBCgoKDg0OGxAQGy0lICYtLS0rLTAtLS0tLS0tLS0tLS0tLTUtLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLf/AABEIAJoBSAMBIgACEQEDEQH/xAAcAAABBQEBAQAAAAAAAAAAAAAAAQQFBgcCAwj/xABVEAACAQMCAwQFBA0IBwUJAAABAgMABBESIQUGMRMiQVEHMmFxkRRzgaEVIzM0QlJTkpOxstHTF2Jys8HS4fAWJDVDVHWUdKK0w/ElNkVVY4KDhKP/xAAaAQACAwEBAAAAAAAAAAAAAAAAAQIDBAUG/8QAOREAAQMCBAQEBQIEBgMAAAAAAQACEQMhBBIxURNBYXEFIpGhMoGxwfBS0RQzQnIGNGLS4fEjksL/2gAMAwEAAhEDEQA/AMjpKWuaisSWkpaShCSjFFFCaDSUtJQhFBoooTXNdVzS0JlJSUtJQhFBopTQmkooooQuaWkpaE0lJS0UIXJNJqHmKtno3vBDeGYxrJ2cEzaG9UkKPYfOtET0gg//AA+3+lgB8SlImFfTol+nssP1DzFGoeYrbX9IY3B4fAPpGR/3K8rfnQSGOL5Ki5eMalcavW8fte/WlLp0t3UxREwTH5tqsYoqw8I5cmv76W2tzGr6pm75ZV0q+47qk+I8KtP8inFPytp+km/hUy4DVUhpKzQ0laZ/InxT8rafpJv4VJ/IlxT8rafpJv4VGdu6eQrNKK0FfRBxEsyiW1ypAPfm6nOP917DVa5l5YmsghleNtbSKNBc4MWjVnWi/jrjGfGosqsfOUzH59QVLhPy5osoOirBDyTxR1V0sp2VgGUhNipGQR7CK6/0E4t/wFx+ZU8wUMpVcNJVjfkbio62M48N1xv5daacR5ZvoELz28kargksAMAnSD16aiBnzozCYlSynZQ9FLRTUUlFFFCEUUUUIT+iikxSWVFFdaaMUIlcUuK6xRQiVzpoxXVJQnKTFBFLRQiVzppMV1XJoTSYo00tFCcrnFGK6FJQiVzppMV3SUJyua5rukIoTXNFdYoxQiVOcmfdpf8As0/6hU2jbINiNDnfzBc/2CoXk0fbpf8As0/6hUyn4JA1aVYMPHct/Yag/T82K6mAME/LrbM2bc7TNjabLtMkqxGrSv15OP8A0pxYHMkecthoSrHwy6bfA/4U2YasD19sqPVQD2/5+mvW2kJuIstq78PuzldWPpqAEn839fp87FbqlTLSM3JtyvLbknS0Ai7jp8BMKK5dci7vSCQRHc7g4P3ZPGpa2S7kGqMysB46z/afaK8uG8W7Sa9j+T28eI5+/HGVkOmVBu2o5z41Z+W+ZoYYVRyVKk7DOHGVbvYB2yg+LDxpvc9rZaPaeWwXFc5zW+USqobyb8pJ+e376Pls35WT89v317y8QUszdlGcsSMr546/D4k15m9BUgxx5xgELg523qYc7m33CtzO29wuPls35WT89v31zxniES2sBuYPlJM1zp1zSqVAS2yBoO+cjr5U4+yAznsos5z6nj1/z/nMLzZ962/z13/V2tSBc7UR81EkxcK533Nt9EwjhnaONUhCIAhCr2SEKCyknGepNEHNfGHGY5pnA2JWJWAPlkJUJxb7p/8Ajg/qY6vXo1mVLWZmIAE3UkAeogG58zj41ix2J/haHFDZNhHcxsVcxuYwq9dcy8WABlklABGC8KgavDBZOvWorjXGri4tboTya8RJjKRg7XUGO8qg47x26b1aebuNm4tcMgQhreQAOXJEkUhOe6AMNkYyTtuBtmmx24khukaVIgYV78mvQMXNud9Csd/dRgcQcRS4rmBpmLQdjrCHWBAPJUTUK5LCtg4ZfKhgs4bDh9w3Z26rIbfUZC0KOXLMASNyckDbrSXPMQjJWThXD1byNqo+kb7j2itku291n4ayDNLWjcf4nHPaXC/JLaIpGjBoolVs/KIFzq9xPxPnWc02kkXEJPblRRRRTUFI4qwcp8pT8Q7XsZYoxCoZzKxVQpzvkKcYwetQNaN6JkhNvxUTsyxG1IkZRlhHpk1FR4kDNILPTEugqt8x8nSWcXbPdWko1BdMM2t8nO+nSNtqrX01ofKvC+CNxKyjs3luQ7yiVLiJdGkROVOCoBOoe3pTqC9F1xEwWfCLV0tDcIFLdnkKwRZZpCO+F0k6Tn1+ucGhWOp7LMzXFa9zFwiEx8MuXhslme/it5vkWDbOjMTpI3BYBQDnPU+4SQaxe74lbNw217OwjNxGVQrIzoO0Ks4/AJPqjYDbBFEJcE7rEa5IrVfsZbcSg4XdNbw28lxdPBMLdOzjdBrPTJw2Ixg7+sfdUlx9OEk3VtO/C4o41kSHsQ63sUsey63x3zkHK/R3vFpikVi9FbDwzgtpObHirQRi2Sznku0WNRGZrbutlMYLMzNt4iOm/MPLttYQcVuTFE0c7RR2WpV0r8oUOzJkd3RqOCOnZ0I4RWTUhrYuKzWdpxG34QvDraaEiFJZJI9Vwxm2MiSZ7mMg/QQMbY5tuH2nDbXikj20V0bO6RYe1AOFk7MIrNjJC9puPHT8BHCO6zmXlyRbBeIa17Npuw076w2ktnpjG3nUL9NbXyxwuLivDY1m7K3STiLy9lGAgYLEx7GIHbJAz7gxqC5e4tw+biF18ot7W0ZYzDaR3EYWBZUZgflGMapM4ySfAgHYUQpcNZhVr5a5DuL23a6jnt4o1kMRM0hj7wCt10kdGFL6SLGeK5Tt7W1ty0akG0z8nmGT9sUeHUDGAeh3yKs/KMNm/L9wt/JJHB8tGWiUM+rRDpGCDtn2UJBvmgqj80csvYmMPcW03aa8dhL2mnRpzq7o051DHng+VQdaZyrwfg4lu7q2El3DZ2hn0TqE1Td86SAoygVB1B9bxxim3EGgv+E3F+bSC2uLSWIBrZOyjlSR0TSyZOSC2c9dh0yQSEyzZZ3SCtQ9MN/a280nD7awto9SxSNMEAlViQ2I8eoCqgHHXU3nWYChRcIRijFFFJJSfL1/HBI7S69LRSR9xVZgXAwcFlB6edTlpcWrpM6yTgQoshzBHkgyxwgDE/XMgPuBqP5KaNbhnlhjmVIpG7OVdSE7AZHszVyj5ltlDBeFWIDjSwEeAyhg+D5jUqn3qKRnktNLNFlV5OOWrbGS4/6eL+PvS2vG7NHR9VwdLK2Owi30kHGe39lWeLjNqfV4RYY8zF/nNOo7qFkLrwjhpCkA91cjOSOpA3wfHwNQzBvlt6ha+DiqjeLlcRvBPus74PxeKOeeSQOFmSVRpVWYF3VxkFlGMDzp39l7P8ef9BH/ABqtM3GbXIzwmxXHgIhg5/XVO5y7Pto2iiSIPCrFVAVc65BnAA8APhUwXTcLNBa2Z9k5+y9n+PP+gj/jUfZiz/Hn/QR/xqiOXlBuoNSqw7RcqyhlODnDKdmHmDWgWV4XBJsrILoY5+QW+Mjwzppqs1srczjAVW+zFn+PP+gj/jUy5h4pDLDFHD2hMbzuxdFQYkWFQBpds/cmz06irHx+1S4i0qtpbEGJ9bRR2uQyyZTUq97fBA/m+yq4eWj/AMbY/wDUj+7Uc4/qUsxcLXVn4gilyWPRLfG+OsK/zTnp7K9raW4hBRJ2QNkso1aTgDJO2+2PhiueK28qykIGI0QjUmoq2IkGQR1HkaZdhNjGiTHlpbHwp5Q4Q4T3urgpC8muJQVlnZxq3BOcuoPl1wM+3fpioq/j0290M5+0xnP/AO1b+W1evYz/AIsnTHR+nl7q44hFILa7Lqw+1R7sD/xNv4mm1jWCGiB0sk7QqT4FxFbee3kf1OwtlfYHuPaxqSAepGcj2qKd8z8DnV1dA0sbjKOuX1AnORjc9QCfPyyMxllJCHjM4ygtoNu9kt8kj0Du/wA7HXbzqRtLyBAGju5IjhCY4zMqFwi5BwPF9eTk9akEQq/fIRbXQIIPZJsRg/fNvVJrQOYpkdL1o5GkUxoVZixYj5Tb7HVvt5b+89az+kqqmqKKKKFWpQVMcG5jltYbqGNEYXcRhctnKqQwyuD17x61DVq3oz5TsLuzMlzAHcSumdci7AKQO6w8zWLGY2ng6XFqzEgWAJv3I+qoosc90N1Wb8tcYeyuIrmNVZ4iSA2dJypXfG/Q0+4FzXNbXFxOI4pBciVZonBMbpKxZl65HX27E9a0az4RyxcOIYTEZGOABNOrE+Q1Ngn2VR/SHyiOHSp2bl4ZgxQtjUrLjUjEDB9YYO3j5ZOfD+LUa1YUcr2uIkZm5Z9z19N7K51GoxsyD2XF/wA9zyRwRLbW0MVtcpcxJEjIAyZwh73eBJJJ6kmuE53uRPfXHZxar6NonHe0qGULld85wPGq1pPlSV1YKp4rlMf6TXAtLa0QKgtZjcRyLntBJliOpxsW8vAVL8T9IMsyyarGxE8yGOS4EGZGUjSThiQGxjc591U/FIaIKOKVO2vNd1Hw+XhqaeylfUzHOsDukopzgKSozt4t50nHOa7m6tLWyl09najCkZ1PhdCasnHdXI2860HnHgdonBFnS3iWXs7Q6wgD5dowx1ddwT8ayEVhwGNbjGOewEQ4tv0j91bUDqcAnkrta+ku6URs9raTXESCOO5kjJnVR0yQe8Rk77dffmFHNVybS6tXCv8AK5VnlkbPaa1ZXOMHGCVHh4moQCkrddV8Uqaj5nuFso7JAqiO4F0ki6hKsijAwc42O+cVJz8+vJM089hZTNLGiTB42IkKerJ63cfpkjqFHkMVE0hFAB5IFQqY5m5imvpEaRI40iQRRRRLpjjjHQAEn2b+wdKWPmOVeHycOCIY5JhOWOdYYBBpG+MdweHiahqUCgAkwjOZlSnLXMNxYTdtb6TlSjo66o5EbqjrkZGw6Y6e/L/mLnSa6txaR28Fpb6tbRwIVDv4FiTuBttjwHXAxe/SJwG0i4Ss0VvEkn+r99UAfvY1d4bnOax+sOAxzcZTNRoIglt+kfurXh1OylubOYZeIXTXUyIjMqLhM6cIMD1iTUPXVJW1VkykpKWihJTXKn3SX5iT9a1YuEcMefWUUv2a6mCgknJC9Bv4k/RVd5U+6S/MSfrWrBwa8EUmps6Qr7DxONtvM4A+movJDDGq6nhrWOqNzxE89D05a/RScnCpgczJ2QO+XBUAeYBGW9yg0x4nfQaewXWFGC74HefGQ2hScBASAOvrE9dIaRPczysygsW3K5AUDwG+MAClns2S4xKBgDXgHIIxgAHG/eHxFc4MuvdOxOZoJ328sQb68hsfls3ZSAATnA+AJLDHxz9JqL5s9eH5hf6yWpRmJJJ6k5qJ5sPfh+YX+slrphsABeBxVTiFz9zP58k25c++of6a1c4L5MLntNQj7IKMaDnYN61POFeiK8xFMLuCN8LIBpZ9JIyAc7N8Me+pT+TPiP8Ax9r/ANLF/cqHEZusrsOHiHA/ghUznH7hj8VrZD/SWOaqWa2G79Fd9Kuh7+205DYECp3gCAcoATszfGmX8itx/wAbB+Y/76OKzdNtEtaG7CFnKcauwAFuZwAAABNIAAOgA1bCva04reySJGLucF2VATNLgFiFz63tq9XfoYuljd0uoHKqWC6XGrAzjO+Kz7gP3zb/AD0P7a0w5p0TIIN1qx9FPGf/AJp//S5ri49EXFZF0ycRV1OMqz3DLtuNjtVX4gB2sm34b/tGvIR/0R78UiQ3VwHyWulhKlWeG0mNeiuEfoo4woCrxPSqgAASXIAA2AAHQAeFOL7kDinZQQDiBEgad2kD3GGVgmlS3XbQ/X8b2mqMUHlSaR5UnMLhZw9P+dlA08phwupHmzlDiNtBK8988qIqMyFpijBpUQDvHBIZlOD+LWdVdpx/qt58yn/i7aqTirGggQf2VVXLPlEfOUUUUU1UpOtx9DH3gfn3/YjrDq3L0MfeB+ff9mOvPf4l/wAgf7m/dRwf835JrwL0XW8Fwk7XDSlHEirpCjUp1DO52BxVf9M9/KbmCKWBhBHqcNlftxOnXpIzp0gAb797OMYrz4J6Lb9LlJnaOILIH1K5LABs7YHXFWr0myW003D7KVhqkuY3YeIiw0e/lqZlA88HyrmfxOTxClVfVFazpgAFogyYFjAk+ulitWWWOAGX7pjDz5wmAJHZWcsuw1aIQHUeRL7u3uyPbTL0wctwJHFeQII2aRY5Ao0hg4LKxUfhAjHmdW/QVcOPSXtqIoeE2MThs6nYhY4yMYyoILZ3Oc+FRHpoZhw6MnGoXEJ26agrn4ZrPha7BjKFSlIzugl1QPc4EgHMABl6TrreJUntljgdtoAXotlYcFsRPNEJJO4CdKtI0jfgqW2AGD9Cn3VSucucOF3tsdNmyXWQEOlVKjxYyJ664z3T4+HjV4534c3F+GxvZkM2pLhFyBq7rq0eTsrDU3XxXBxWbcU9Hl7b2jXU2hdG7x6gXCdNWR3T4bZz+o6/Czh6zhiMVUIq8T9UbQI/TMgxEaTChVzNGVjfLC0rjvDJrrgkMEK6pHisQPId6Ikk+AABJ91MOPy2XBLFYkjjluHBCBlBLP8AhSNncIvl47DxJqW4jxuSz4NBcxAMyQ2Wx6MrdmrDPhlSd/Cofn7gkXFbKO/shqkRNQwO88QzriIH4anVgeYYeNc3COPlp1zFA1XSQdXQLOPJum3M8rWu3GsL15J5ehtbH5bNAbiaRe20rEHbDbqscYGATkeWM+AFPoLOHi0Mkd3w6S1dfUZ4wrjOcMkgA3GN1O3TOQa9OVeIyXXCI/kMipPHEsILAMFliAGGBzswA38nBqIsIeaZAxeW3hIBwHSJtR8B9r1Y95qVU1ataq+pUax7XQJcWlsHRouCPUj3I0ANAixGy8PRBwxNF5DPGjtHN2TalDDKgqeo6ZBqS5f5l4ZLctwyC00qvaIC0aFGMedQIJJPQ7t18aaehdpCt40xzIbjMnT7oc6/V29bPTaqn6Pv9uP87d/+ZWyvh/4jE4rim7WBwymBmDBHcbdFAHI1kDUx7pzxTkqH7NLaRjTDIouNI/AXvalHsyhx5agKt/NPNHDuFGO1+Sa9SBiqImBGSVBYt6xOltvZv4U041fpDzHbazgSWoiyegZ3l0/FgB/91NvSXyNeXl7HNbqrI0ao5LqugqzHJB3IIYdMnY1FlVuKqYaljHxT4ZN3QCQXAEnew+Y3JRlyhxYLypT0tOrcIJQYUtblRjGFLDAx4bVhFbv6WYSnCCh3Ktbr8GArCBXV/wAM/wCUdH6z9GqnF/EOyKKKK9Csq5opaShNSnL95HFI5lJCtG6ZC6sFsEbZG21Sn2QtPyz/AKE/3qq5qSs+CTSwmddOgSCEZJ1GQgHAUA+BFBVrHuAgKXiu7VmVRM+SQB9qPUnA/Cp3xPiUSO8Mtw7NEzxn7UW3VsNhi2SMqPhUPa8u3aujGJsKyk9x84BBP4Ne3GuB3MtxPKkT6ZJZZFzG4Ol3ZhkY64NMseTcH0P7KxuJLWkBwE9eXrC9fsjafln/AEJ/vVEcx3kcrIYiSqRBMldOSGdjtk7d4U6blO+69g/5rfuqO4lwqaEZlUrnzBHn5inlcNQfRQ4ma0j2Wzc/8Sa3jikVVY9nAuGzjBDeXuqrjjN/pDfJ7fScYPyiHO50ju9rkbkdRt1O1TfpW+4Rf0bf9hqiOTGkaDUkt1I64iCI12yQqz4BUQDunSM4JIwvQ4xU+PUp0mZTbKPqVW3D0qjnlzZOY/ZR91zXcxsUeGEMMZAYuNxkd5GKnY+BrxPOs35KP/vfvrvnGS/VRHPHIkJZSm18I8gNpTFycawNyAPDY4zVVfoaYxdUtnMVMYOhPwrdeXPXf5mT9QrA+A/fFv8AOw/trW+cueu/zMn6hXz9wy4EcsUhBIR43IHUhWDED27U8X/PPYfdUYP+QO5+yud8v2yTww7+OOrHevEKM742G2TkeHl9NNrjjlozs2ZxqZmx2Me2Tn8tXl9mLT8af9DH/HrIaUuLp5R2XoKfiTWUW0coIBk/6u/a0dgU/A6+r4+Jz4YHl50rAZ6D8LbPtOPEbYxUf9mbT8af9DH/ABqPsxafjT/oY/41HCP6j+fNB8QpGxpN5aRyM7RfmIi5T65+9rz5mP8A8Vbe01SM1Z7rjVt2E8adsWlRUGqNFUYmilJJErHpGR08aq9WNEBcyu4PeXC0paKSimqlJ1McJ5pv7VDHa3BjQksV0RMNRABOXUnoBUQa1bk/gfDvsX8su4A+jtGdgX1aVY+AYZ2rn+IV6FGkHV2ZhIEQHXMxYwN1RQa5zvKYVLfn/jBGDeN9EcAPxCZFV2aWR3Mkjs8jHJdmJckdDqO+enwrU7DhnL/EyYbXtIJ8Er66kgbnCuSrjzAwcZ6dazvivBZ7e5e1ddUiNgaQTqBwVZfYQR8cVVgK2FdUdTpU+G8C4LA0kb2mR81bVbUAkmQpI8/cW0dmLx8YxnRFrx/T06vpzn21HX/MV9PCtvPcNJEhBCsEJyuQCXxqbqepNHEOAXcCdpPbzRpt3micKCegJxgH314XXDZ4kWSWGREfGh2RlRsjUNLEYbI328KvpUME0h1NrNbQG6jaOY9RqoGpVNjKdcE5kvbPItZ2RTuVwroT56XBAPtGDXPG+Y768A+VXDOoOQmFRM+ehAASPM5NNzwy47Lt+xl7L8oEYxjfG7gYG+OvnTmHl2+ePtUtZ2TrqWJmBHmMDeplmGa81iGBwN3WkHqdZ90g6pGUTCS75lvZLcWklwWgARQmiMd2PGgagurbSPHwrrg3M99aIY7W4aNCdRXTGw1dMgOpx0HSrJ6I+GQXF1MlxEkiiEnDrkBhIg6Hx3NVe/4e73lxDbxM2maZVRFJIVZGAAUeQH1VnacK6s/B8MQAHGzcpnnEaxF4057zJqBofPRefDOOXdvK00E7RyOSXK6cMSS26EaTuTtjbO1SHEueOJ3CGOS7YIdiEVI8jyLIobHszio/iXBbq3ANxBNGDsC6Mqk+QYjGfZRDwW6dVdLeZkc4VhGxVjnGAwGCcg7Ve6jhahFVzWHkDDeXIHoPQdAo56gtdenAuY72yVltJzEHILAJG2SNh66nH0U24fxS4hm+UxSFZss2vCk5fOo4IK75Ph40vE+EXNsQLiGSLV6utSob3E7GnEXLV+xCraT5Yal+0ybrt3gcYI3G48xUow0F5y+aZNvNFjJ59dbImppeybcZ4tc3cglupDI4UIG0qvdBJAwgA6sfjUunPvFRH2QvH0gYyVjMmMY+6FdWfbnPtpha8vXspcR20zGM4cCJ8q3XSRjIOCNqjmhcMUKsHBwVIIYN4gqdwfZRwMNVAp5WnLoIaYnYcp3GvVPO8ea91I3vM1/NALaa4Z4Rp7rBCe6cjL6dR+k1EipS65dvoo+1ktZ1j6ljE+kDzJxsPaa8LbhdxIhkihkdAcFlRmUHyJAwOo6+dSpcBjTwsoE8oiT2tJt1NuiTsx1TI0lSPEeCXduA1xbzRq2wZ0dVz5ZIxn2VaLfkF24ZJdlLj5SHGiJUJ1IXRc6ApY7FjkeVRq4yhSYHucIJDZkRJ6zFuew1TbTcTCo1c16zQujFJFZGU4ZWBVgfIqdwa860AgiQo6JKvXK4/wDZ6/8AMYv2I6otX3lf/Zi/8zj/AGEo/rZ/cPupj4Xdivfhd0ZmdGvZUlEkgVNPd0qSQ2RE22MDGc7dMb004/c3VuQUuyyse6CD2gGkHLZjVfHwJ6j217cqXREtzEHcByxZEKLldegsXfIwNeNOPws9Aaa802CtcN2SxIdixkvLFS+oBslMpobLHOSS2c7VM1Xh5GY+qtbSp5R5R6D9lM8l8XnmE4mkL47PTnG2deeg9gqH9JfRPd/epzyAMGcbbdmNiCNtfQjYj2im3pL6J7v7WrY4k4a/5dYQ0DFkAcv/AJVs9K33CL+jb/sNVK4JckQOMRko8ejVFAxAbWz+upZt1TzA6bZ3unpYP+rxf0bf9lqr3KvD7iW3miKy4+1tCH7RYwdRLlCw0DO2ehO3WsFY5aLD/pH1K6uBDTXObTPf21UVIkbOZZIlfxYLiLUAOmYx3eg6CojiEkbMTFF2S49TW0mDjfvMATmr1Z8Burd1nkhZliOvTE8byEjppCtnIOD57VTuNtGXypuS51GQ3OntCdsHI3Pjkn2VRQfmkaroeIimKjckdY7+krZuXPXf5mT9QrGuV+XRdITqCldPXJznUPA/zfrrZeXPXf5mT9QrLuRbwQ280hUtjsRhcZJZ3UYz7SK6lUNOK82kf7l5ugXDCy3WberQvf8A0BH5RPg376W65BWOPtHlGMqMLHI7ZY4ACq2T1FS9xzIY95bS5UZwSVUe3z8q64jzCl1aO6K8YjmtwctpO8inIZDldvEbim5tCDlue5Sa/EZhnFpA0bzPZU654NZxnDzMG66Wt5lOD7zUxwblO2druKVctFGGV1Z1IJUnOM4Ph1B6V486RSOwkVmeJVB3CAIXPRcHJBIHXxHXHS0cAjzc3/siT9isVQ/+IuAg235yt7Ww8CZsduUbAbrHBSGkXoK6pqC5opaKEKUrWuD/APu3P83c/rNZLV75Z9IMVraC0ltDMMvnLrpYMc4Ksprj+MUKtag0Um5iHtMSBYTuqsK5rXy4xZRnousZJOIQMgOIyXY42C6SDv7c4+mtTsoYJeNXD7GSC2t09zO8jE+/To/OqkXXpW0RlLGxjgY/hEgge3s0VQT7z9Bqo8vczXNpdG7Ddo76u1Dk4kDHJBI6HIGCOmPLasGMwOLxzn1S3hkMytGYEkkyZItBBIHWFoZUp0gGzN5WsWnMlmrTrecVhuI5NSmFoVjCAkgptuRpyCGydvfmB4FjiXBbizQl5LRiIc+uyIS8HXpqUMn0UwuvSPaAO8HColnkB1u4jKkt1LaV1SDzG2arfI3NDcOnaXRrR1KsgOnxypGx6H6iaz0/Ca/DfUawtqNLHNnJctmbNAAtvqYVjqzZAmReVo3Mwii+xnB8jEkkLTAYwUQ50n2PJn801M8xcXjguow/E47VEVCbcwqe0XJyS53AIGBpxjFYtzLzBLeXrXgzGQU7IZyYxHjTg466sn3mrgPSXbyqjX3DY5poxhXGgjPmNakpv4DNRreEYnLTcQXk5i8S2z3XJ80tmIBIv5ZlDa7JImNuysfLd5azccnks3V0e0DOyZ0mXtEB950hM48c+Oa8eT3EEHGLyNQZY7i+xkeEY7RV88Fjv7hVP4T6QOy4hNfPbKBLF2fZxtpxgqQxYjvNhcE4Hh5U25b56ls7i4lEQeG5keR4WbGCzEghsdQGwdsHApP8Ixjg4NB+BgguBnLEtkRbYkAQInRArMkGeZ/7V49G3Frjilvdw8RYTJhFDaEQ4cNqHcAGxVSDjIPj0x58B4tNa8udvCwEidoFYgMATcFNWDsdicZ26dar3E/SSgt3t+G2S2okBDMNIIyMHSsYxqx4k7eVWTlq8ituXxLLCs8a6tUbYwwa40nqCNs594p4rCvaOIaWVr61KKdrw1wOlhmNo9RoVJrwTEyQDded5dyX3L0k14Q0ih3V9KqS0b91gAMAkZU4x1PnT7nbmC6s7Kwa2cK0jQqxKq2V7LOnvDbPxqh848+/LLdbO1txb240lhkZYKdSrhRhVDYPiSQOm+fLmznUXlvawrCUNuVYsXDBiqaNhgYq+j4XUNWmX0op8V7iw5Ya0gR9NL6KLqrQDBvGvVaXz1zDc215w+GBlVLibEuVUll7SKPGSNtmO436V5zcPibmIMyjIsRN06yrMYgx8yFx+aPKs95r56W8urO4W3ZPkr6ypcNr+2I+AQNvU+uuOOc+yS8Rjv7ePsmijEWlm1BxqcsGxjukP9BANZKXguLaymGNyOyPDjI1JdAMGbiBN4+SkazJMmRIWljmS3iu5TccWiKZdDbNEqiMg4xr9YkYOc7HJ26YieRLyOGx4rNZ6THFPdyQZB0aViDoMbHTsPLaoSX0m2eozjhUfyojdyY+uMZ16NR29nszUDwPnUW9ne2zW4Ju2nYFCEjjM0ejATB7o8BnptVZ8IxTqRApEHyiJYAYMkwImP1Ez3TNVsi+++yunAuLz8Q4FfvesJWUXGCUVfUhWVNlAGQ+4PsFe0HMd2vLxvBKe3XAEmiMkAXCxgaSuk9zbJHt61QeXecRa8PubIwlzcCUBw4AXtIhFuuN8YzT/lPn6O2tDZ3doLiIEsoyuN216WVwQcNuD9W1ba3hdUOqOZSBaKrXNb5QC2+YDkJ8si0gDWFBtUWl3Iql317LPI80za5JDqZsAZPnhQAPoFeNOeJ3SSzSSxxCFHYsI1OVTPgDgbZz4ADOAAKbV6lghgERYWtbpa1tLWWN2q5q+8r/AOzF/wCZx/sJVDrQOUYHfhoCIzkcSjJCqWIURx5JAGwHnT/rZ/cPum34XdivHlIg3U4CAtplbUU7XChwCvZnbclTq6jTj8I1NccsLQjXcJHG3e0akNurHAwGdCC3eBABDYCtjpio6DlriyPIYF0BnLfg77tg95T4E/GvC95Q4rIF7VA2jOnddtRGeg9g+FTdSeXSrhWYALr15Fh0PcLrV8dl3kOUOzHY4GfhTrmbgUt7LHFHsMZdz6qLk7+0+Q8fdkj35S5durftO1jI16MeI21ZycbdRV2hiCjA+n31X4hjhhsKGD4zMdBOp9oHM9AU8DgziMW55+Aa9SRED3k8u5C9H72ksASqqoOPBRgf2/GoW95lgjcp3mI2JGMZ+NO+KX4ijkbxVdvedh9ZFZoxrz2EwwxEvqknl+HovaUKDQIiBtorz/pbB5P9X76dsLW/iKuodfEHZlPmCOh9oNZ3U9ybd6J9JOzqR9I3H6jWivgKbKZfTkEX1VlSi0tMK88Jtikj+I7GXB+gdaxjgP3hc+6D+sNbUreIOMg/AjH6jWbcO5NvEhurdYXfJjVH0nQ4Vy2Qenq4yPA7Vv8AD8ecS88U+YNI72N9gbwevoPIY3AjDNHDHlLge3maY35EpxxpxJZqyuYwwTc5VCjqGyxwx0ksfHqRvXPI3DYpoJ4ZcSIZY86WYA4GQQwwetS8fDOLqiKtopKqAxdydRAXvYAGCSD9VOOU+X7y3WczQsGkfWAupvA5395rpYZsP82l1gxT5p+XW2ndRsnKVjk4iOPnJP71OuW/vniPzK/sGpn7Gz/kZPzW/dUfwfhs8U1+8sTojxAKzKQGwhzgnrVmP4Yo+WNRpCowJqGoS+dDrO4/ZYevQUEUJ0FLWZaVwaKU0UJqUpKWkNQWJc0poNJTQikpaDQmuaUCkopppdNIwoopQhc4rpnYro1tp/FydPn6vSikpwpAkLnFJXRooQuaXFLRQhcYpMV0aShNJiilooQkIopTSouSBnGTjJ6UJi64IrQvRrz5bcOhkinjmYtIXBjEZGCiLg6nXfufXVMHDs7CWP40Hhv/ANSP87/CouaHCCtDKdVpkBbL/LPw78jd/mQ/xqslxzdbC2+VsswTsxLpwofSRqAxrxq386+dG4acfdU/O/wrUOKc28OnszZsZ41KRpqQQk4QqfGTx04+mqH0QIyjur2cS+ZWXhXPVtxBXjt0nXRpLGUIMhs4AKu2d1+qvWeYKPb4Css4Vfw2MxazkZ4nVRIJiiPqBOCugsMAHx8z7DVhbnG3O51Z96f3q5ONwdV1XM1trR+d5Xb8PqUxSAqOg8/ztCcc0SnssfjMM/Wf3VVKmuN8QSeOJ09Ul/LquB4E1C1uwbCykARBk/Uj7LuscHNBboUU44bJplQ/zh9ZxTeukOCD5VpIkQpLT7ObI0nqPrFLxLnGHhsYaeOV1kbC9mEJDYyc6mHUD6qq15zLBFIY2zqXGcFfEBvFs+Iqq838Wa90IrosaEsNTd4sRjJxsMDPietcbB4SrxWvItz+Y+64niFSkabmtMnbnIKv7emjh5/3F39Cwj/zasnLnOVtdxdtCswXUUIfTkEYPQMR0I+NfPEfCxkapU05GcNvjO+MjGcVpvKXMfC+HxvHC106u2vv/J9mwF/BceAHwrs1KIA8ouuCwVJ8wVv5k9I9pZMizQztrBYFFjI2OCO8436fGqvxz0v2csTpFb3Goq4GsRKMsMDJDkgfRVf9IvGbbiKw9hlHiL5LlACrgfisd8qPrqkfYo/lY/zv8KbKIy+YXScKk2CjgKWpB+FkbGWIe9v8KZTx6WK5DY8Qcg7Z61eqHMc3VeRorqimoypKiiioLEkNIKU0hppoIopKKE0lJXRrmmE0lLSUUIXeKSkoohEJK5paDQmiiiihCK5rquTQmkooooTRQaKDQhetn6616y+sfea8rP11r1l9Y+80LZh/hXNFAFdiPzoV64oJr2CCvO48KaFNcEui0XZn/dtsfZIBt8VPxp9FGzMFUZZiFA82JwB8TVZ4deSIxEadpqxtgk7Z3295qTj4jdqQy2zAqQwIVsgg5B6eBqpzL2/Oa9BgsawUGtMki1muI1tcCNIVj4jwdolLCSKQK3Zv2bMdD4JCnUq5zpbDLkHB3qNFS3OF7dLHAEtlxIoknEQH32Bh1fSTpZRJ6uwHaHaqq/EbhQS1uQB4kNgfVVbJIlX0sezJL5/9Tp3Ay+nKOaa8Xu+0nkcjGpj8F7o+oCmpFeKvnJPU5J+mnkR2FaA0NAA5LzRdncXbkn1K8aKcFRXmY/KhRhedFKRikoQvLivr/QKaU84r6/0D+2mdCyv+IoooooUFJClI99dGuDUVkSUH3UtIaEJM+ykJ9lKaKaa5z7KSuqQ0JrnNGaBS00JM0E0UpoQuc0ZoooTXJpc0pooTlJmjNLRQhcZpK6NJQmkpaKKEL1tPXWnDrlj7zTe09dffTzxPvNC24f4Vyq4paKKFoRXhedBXvXjdeqfoppHRc8LvOyfXgnYjbrvj91WFOMKRkM/11UqkLf1RUXMB1WrDY6vRbkYba6Kz3nHWfT2kjNgZX6cDJwNycDJO+wqN4jxgaCveOoEdcD66a3P4H9Af20xvOg99IMCud4jiA0tBEdgmdP7b1R/nxqPqRh9Ue4VYua1d0UUVFTQRXk6Yr1obpQkmfFfX+gf20zp5xX1/oH9tM6Fkf8RRRRRQoL//2Q==', 1, 0);
		  
	  INSERT INTO polls(name, topic, src, upvotes, downvotes) VALUES(
		'Who should be the President of USA','Politics','https://media.istockphoto.com/id/1021174468/photo/briefing-of-president-of-us-united-states-in-white-house-podium-speaker-tribune-with-usa.jpg?s=612x612&w=0&k=20&c=bjgAmJTLdeX1LpwMLINaFxC8fN9Tyw6PZgfiTEa5glw=', 1, 0);
		  
	  INSERT INTO polls(name, topic, src, upvotes, downvotes) VALUES(
		'Best Football Player','Football','https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcQDn0dSN-H8iZAlsegO7dsiUV79s3BehwLMDEoNd7PT9WFc4wdGmPq4WlW6rf-xRcjNaRk&usqp=CAU', 1, 0);
		 
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
				option VARCHAR(255) NOT NULL,
				votes int 
			  ); 
	
	INSERT INTO options (pollId,option, votes) VALUES(
				1, 'Sachin',0);
	INSERT INTO options (pollId,option, votes) VALUES(
				1, 'Dhoni', 0);
	INSERT INTO options (pollId,option, votes) VALUES(
				1, 'Virat Kohli', 0);

	INSERT INTO options (pollId,option, votes) VALUES(
					2, 'VueJs',0);
	INSERT INTO options (pollId,option, votes) VALUES(
					2, 'React Js',0);
	INSERT INTO options (pollId,option, votes) VALUES(
						2, 'Angular Js', 0);
	INSERT INTO options (pollId,option, votes) VALUES(
					2, 'Ember Js', 0);
	INSERT INTO options (pollId,option, votes) VALUES(
					3, 'Java', 0);
	INSERT INTO options (pollId,option, votes) VALUES(
						3, 'Golang', 0);
	INSERT INTO options (pollId,option, votes) VALUES(
						3, 'Python', 0);
					
	INSERT INTO options (pollId,option, votes) VALUES(
							4, 'Donald Trump', 0);

	INSERT INTO options (pollId,option, votes) VALUES(
					4, 'Barrack Obama', 0);	

	INSERT INTO options (pollId,option, votes) VALUES(
				4, 'Narendra Modi', 0);

	INSERT INTO options (pollId,option, votes) VALUES(
							5, 'Messi', 0);
		
	INSERT INTO options (pollId,option, votes) VALUES(
							5, 'Ronaldo', 0 );
	INSERT INTO options (pollId,option, votes) VALUES(
								5, 'Neymar', 0 );

	`
	rows, err := d.db.Query(sql)

	if err != nil {
		panic(err)
	}

	defer rows.Close()
}

func (d *PostgresDB) FetchUserToken(userId int) (string, error) {
	sql := fmt.Sprintf("Select token from tokens where id = %d", userId)
	tokenToUserId := model.TokenToUserId{}

	err := d.db.QueryRow(sql).Scan(&tokenToUserId.Token)
	if err != nil {
		log.Fatal("Error fetching token with userId", userId)
	}

	return tokenToUserId.Token, nil
}
