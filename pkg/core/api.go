package core

import (
	"database/sql"
	"errors"
	"fmt"
)

// ошибки - это тоже часть API
var ErrInvalidPass = errors.New("invalid password")

type QueryError struct { // alt + enter
	Query string
	Err error
}

type DbError struct {
	Err error
}

type Products struct {
	ID int64
	Name string
	Price string
	Qty int64
}

func (receiver *QueryError) Unwrap() error {
	return receiver.Err
}

func (receiver *QueryError) Error() string {
	return fmt.Sprintf("can't execute query %s: %s", loginSQL, receiver.Err.Error())
}

func queryError(query string, err error) *QueryError {
	return &QueryError{Query: query, Err: err}
}


func (reciver DbError) Error() string {
	return fmt.Sprintf("can't handle db operation: %v", reciver.Err.Error())
}

func (reciever DbError) Unwrap() error {
	return reciever.Err
}

func dbError(err error) *DbError {
	return &DbError{Err: err}
}

const loginSQL = `SELECT login, password FROM managers WHERE login = ?`

func Login(login, password string, db *sql.DB) (bool, error) {
	var dbLogin, dbPassword string

	// QueryRow если такого логина нет?
	err := db.QueryRow(
		loginSQL,
		login).Scan(&dbLogin, &dbPassword)

	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}

		return false, queryError(loginSQL, err)
	}

	if dbPassword != password {
		return false, ErrInvalidPass
	}

	return true, nil
}

func LoadManagerByLogin(login string, db *sql.DB) {

}

func QueryData(db *sql.DB, query string, mapRow func(rows *sql.Rows)) {
	// mapping -> отображение одних данных в другие
	// map
}

func ShowAllProducts(db *sql.DB) ([]Products, error) {
	const getProductsSQL = `select name, price, qty from products`
	products := make([]Products, 0)
	rows, err := db.Query(getProductsSQL)
	if err != nil {
		return nil, queryError(getProductsSQL, err)
	}

	defer func() {
		if innerErr := rows.Close(); innerErr != nil {
			products, err = nil, dbError(innerErr)
		}
	}()

	for rows.Next() {
		p := new(Products)
		err := rows.Scan(&p.Name, &p.Price, &p.Qty)
		if err != nil {
			return  nil, dbError(err)
		}
		products = append(products, *p)
	}

	if rows.Err() != nil {
		return nil, dbError(rows.Err())
	}

	return products, nil
}

// TODO: add manager_id
func Sale(product_id, product_qty int64, db *sql.DB) (err error) {
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	var (
		currentPrice int64
		currentQty int64
	)
	err = tx.QueryRow(`SELECT price, qty from products WHERE id = ?`, product_id).Scan(
		currentPrice, currentQty)
	if err != nil {
		return err
	}

	_, err  = tx.Exec(`INSERT into sasles(manager_id, product_id, price, qty) VALUES (:manager_id, :product_id, :price, :qty)`,
		sql.Named("manager_id", 1),
		sql.Named("product_id", product_id),
		sql.Named("price", currentPrice),
		sql.Named("qty", product_qty))

	if err != nil {
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}



	return nil
}