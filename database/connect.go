package database

import (
	"database/sql"
	"errors"
	"log"
	"strings"
)

func Connect() (*sql.DB, error) {
	dbConn, err := sql.Open(
		"mysql",
		"Khoroshylov:Valentyn!1@/routines_db?parseTime=true",
	)
	if err != nil {
		return nil, err
	}

	err = dbConn.Ping()
	if err != nil {
		return nil, err
	}
	return dbConn, err
}
func DeleteTables(db *sql.DB) error {
	query := `DELETE FROM restaurant`
	_, err := db.Exec(query)
	q2 := `DELETE FROM product`
	_, err = db.Exec(q2)
	q3 := `DELETE FROM ingredient`
	_, err = db.Exec(q3)
	q4 := `DELETE FROM product_ingredient`
	_, err = db.Exec(q4)
	q5 := `DELETE FROM menu_products`
	_, err = db.Exec(q5)
	return err
}

// RowId MAKE SURE ROWS HAS NOT COPIES
func RowId(db *sql.DB, selectQuery, insertQuery string, args ...interface{}) int64 {
	row := db.QueryRow(selectQuery, args...)
	var id int64
	err := row.Scan(&id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Println(err)
		}
		result, err := db.Exec(insertQuery, args...)
		if err != nil {
			if strings.HasPrefix(err.Error(), "Error 1062") {
				return RowId(db, selectQuery, insertQuery, args...)
			}
		}
		id, err = result.LastInsertId()
		if err != nil {
			log.Fatalln(err)
		}
	}
	return id
}
