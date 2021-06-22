package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/jackc/pgx/v4/stdlib"
)

func ref() {
	conn, err := sql.Open("pgx", "host=localhost port=5432 dbname=bookings user=postgres password=password")
	if err != nil {
		log.Fatal(fmt.Sprintf("Unvable to connect: %v\n", err))
	}
	defer conn.Close()

	log.Println("Connected")

	err = conn.Ping()
	if err != nil {
		log.Fatal("ping failed")
	}

	log.Println("Ping DB success")

	err = getAllRows(conn)
	if err != nil {
		log.Fatal(err)
	}

	query := `insert into users (first_name, last_name) values ($1, $2)`
	_, err = conn.Exec(query, "Uzumaki", "Naruto")
	if err != nil {
		log.Fatal("ping failed")
	}

	stmt := ` update users set first_name = $1 where first_name = $2`
	_, err = conn.Exec(stmt, "Saske", "Sakura")
	if err != nil {
		log.Fatal(err)
	}

	var firstName, lastName string
	var id int
	query = `select id, first_name, last_name from users where id = $1`
	row := conn.QueryRow(query, 1)
	err = row.Scan(&id, &firstName, &lastName)
	if err != nil {
		log.Fatal(err)
	}

	query = `delete from users where id = $1`
	_, err = conn.Exec(query, 1)
	if err != nil {
		log.Fatal(err)
	}
}

func getAllRows(conn *sql.DB) error {
	rows, err := conn.Query("select id, first_name, last_name from users")
	if err != nil {
		log.Println(err)
		return err
	}
	defer rows.Close()

	var firstName, lastName string
	var id int

	for rows.Next() {
		err := rows.Scan(&id, &firstName, &lastName)
		if err != nil {
			log.Println(err)
			return err
		}
		fmt.Println("record is", id, firstName, lastName)
	}

	if err = rows.Err(); err != nil {
		log.Fatal("Errpr scanning rows", err)
	}
	return nil
}
