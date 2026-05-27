package main

import (
	"database/sql"
)

type ParcelStore struct {
	db *sql.DB
}

func NewParcelStore(db *sql.DB) ParcelStore {
	return ParcelStore{db: db}
}

func (s ParcelStore) Add(p Parcel) (int, error) {
	// adding new row into DB using parcel data
	res, err := s.db.Exec("INSERT INTO parcel (client, status, address, created_at) VALUES (:client, :status, :address, :created_at)",
		sql.Named("client", p.Client),
		sql.Named("status", p.Status),
		sql.Named("address", p.Address),
		sql.Named("created_at", p.CreatedAt))
	if err != nil {
		return 0, err
	}
	// returning id of last added row
	addedNumber, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}
	return int(addedNumber), nil
}

func (s ParcelStore) Get(number int) (Parcel, error) {
	// variable type Parcel for returning
	p := Parcel{}
	row := s.db.QueryRow("SELECT * FROM parcel WHERE number = :number", sql.Named("number", number))
	// add values from got row to empty variable
	err := row.Scan(&p.Number, &p.Client, &p.Status, &p.Address, &p.CreatedAt)
	if err != nil {
		// returning Parcel{} variable in case of error
		return Parcel{}, err
	}
	return p, nil
}

func (s ParcelStore) GetByClient(client int) ([]Parcel, error) {
	// slice of Parcels for returning
	var res []Parcel
	rows, err := s.db.Query("SELECT number, client, status, address, created_at FROM parcel WHERE client = :client", sql.Named("client", client))
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	// adding parcel into slice for each got row from DB
	for rows.Next() {
		p := Parcel{}
		// adding valies into new variable for the slice
		err := rows.Scan(&p.Number, &p.Client, &p.Status, &p.Address, &p.CreatedAt)
		if err != nil {
			return nil, err
		}
		res = append(res, p)
	}
	// checking if rows.Next() didn't return error
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (s ParcelStore) SetStatus(number int, status string) error {
	// changing status of parcel in DB
	_, err := s.db.Exec("UPDATE parcel SET status = :status WHERE number = :number",
		sql.Named("status", status),
		sql.Named("number", number))
	if err != nil {
		return err
	}
	return nil
}

func (s ParcelStore) SetAddress(number int, address string) error {
	// setting new address for the parcel in DB
	// setting is possible only if parcel's status is "registered"
	_, err := s.db.Exec("UPDATE parcel SET address = :address WHERE number = :number AND status = :status",
		sql.Named("address", address),
		sql.Named("number", number),
		sql.Named("status", ParcelStatusRegistered))
	if err != nil {
		return err
	}
	return nil
}

func (s ParcelStore) Delete(number int) error {
	// deleting row from the DB
	// deleting is possible only if parcel's status is "registred"
	_, err := s.db.Exec("DELETE FROM parcel WHERE number = :number AND status = :status",
		sql.Named("number", number),
		sql.Named("status", ParcelStatusRegistered))
	if err != nil {
		return err
	}

	return nil
}
