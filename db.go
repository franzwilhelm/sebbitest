package main

import (
	"database/sql"
	"fmt"
	"os"
	"strconv"

	_ "github.com/mattn/go-sqlite3"
)

// DBHandler ...
type DBHandler struct {
	DB *sql.DB
}

// NewDBHandler ...
func NewDBHandler() (DBHandler, error) {
	db, err := sql.Open("sqlite3", "./vouchers.db")
	if err != nil {
		return DBHandler{}, err
	}
	return DBHandler{DB: db}, nil
}

// CreateDatabaseFile ...
func CreateDatabaseFile() error {
	file, err := os.Create("vouchers.db") // Create SQLite file
	if err != nil {
		return err
	}
	file.Close()
	return nil
}

// Init ...
func (handler *DBHandler) Init() error {
	// SQL Statement for Create Table
	createStudentTableSQL :=
		`CREATE TABLE vouchers (
  		"id" integer NOT NULL PRIMARY KEY AUTOINCREMENT,
  		"code" TEXT,
  		"sent" integer,
			"email" TEXT
	  );`

	statement, err := handler.DB.Prepare(createStudentTableSQL) // Prepare SQL Statement
	if err != nil {
		return err
	}
	_, err = statement.Exec() // Execute SQL Statements
	if err != nil {
		return err
	}
	return nil
}

// CreateVoucher ...
func (handler *DBHandler) CreateVoucher(v Voucher) error {
	insertStudentSQL := "INSERT INTO vouchers(code, sent) VALUES (?, ?)"

	statement, err := handler.DB.Prepare(insertStudentSQL)
	if err != nil {
		return err
	}

	_, err = statement.Exec(v.Code, v.Sent)
	if err != nil {
		return err
	}

	return nil
}

// GetVouchers ...
func (handler *DBHandler) GetAvailableVouchers(quantityVouchers int) ([]Voucher, error) {
	query := "SELECT * FROM vouchers where sent=0 LIMIT " + strconv.Itoa(quantityVouchers)
	row, err := handler.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer row.Close()

	var vouchers []Voucher
	for row.Next() { // Iterate and fetch the records from result cursor
		voucher := Voucher{}
		row.Scan(&voucher.ID, &voucher.Code, &voucher.Sent, &voucher.Email)
		vouchers = append(vouchers, voucher)
	}
	if len(vouchers) != quantityVouchers {
		return nil, fmt.Errorf("Not enough vouchers left. Got %v, needed %v", len(vouchers), quantityVouchers)
	}

	return vouchers, nil
}

func (handler *DBHandler) MarkVouchersSent(vouchers []Voucher, sentToEmail string) error {
	codeArray := "("
	for i, voucher := range vouchers {
		codeArray += "'" + voucher.Code + "'"
		if i != len(vouchers)-1 {
			codeArray += ", "
		}
	}
	codeArray += ")"

	insertStudentSQL := "UPDATE vouchers SET email = ?,sent=1 WHERE code IN " + codeArray
	statement, err := handler.DB.Prepare(insertStudentSQL)
	if err != nil {
		return err
	}

	_, err = statement.Exec(sentToEmail)
	return err
}
