package main

import (
	"database/sql"
	"fmt"
)

type ParcelStore struct {
	db *sql.DB
}

func NewParcelStore(db *sql.DB) ParcelStore {
	return ParcelStore{db: db}
}

func (s ParcelStore) Add(p Parcel) (int, error) {
	res, err := s.db.Exec("INSERT INTO parcel (client, status, address, created_at) VALUES (:client, :status, :address, :created_at)",
		sql.Named("client", p.Client),
		sql.Named("status", p.Status),
		sql.Named("address", p.Address),
		sql.Named("created_at", p.CreatedAt))
	if err != nil {
		return 0, fmt.Errorf("exec failed: %w", err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("last insert id failed: %w", err)
	}

	return int(id), nil
}

func (s ParcelStore) Get(number int) (Parcel, error) {
	row := s.db.QueryRow("SELECT number, client, status, address, created_at FROM parcel WHERE number = :n",
		sql.Named("n", number))

	p := Parcel{}
	err := row.Scan(&p.Number, &p.Client, &p.Status, &p.Address, &p.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return Parcel{}, fmt.Errorf("no rows found: %w", err)
		}
		return Parcel{}, fmt.Errorf("scan failed: %w", err)
	}

	return p, nil
}

func (s ParcelStore) GetByClient(client int) ([]Parcel, error) {
	var res []Parcel

	rows, err := s.db.Query("SELECT number, client, status, address, created_at FROM parcel WHERE client = :n",
		sql.Named("n", client))
	if err != nil {
		return nil, fmt.Errorf("query failed: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		p := Parcel{}

		err := rows.Scan(&p.Number, &p.Client, &p.Status, &p.Address, &p.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("scan failed: %w", err)
		}

		res = append(res, p)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}

	return res, nil
}

func (s ParcelStore) SetStatus(number int, status string) error {
	_, err := s.db.Exec("UPDATE parcel SET status = :status WHERE number = :n",
		sql.Named("status", status),
		sql.Named("n", number))
	if err != nil {
		return fmt.Errorf("exec failed: %w", err)
	}

	return nil
}

func (s ParcelStore) SetAddress(number int, address string) error {
	_, err := s.db.Exec("UPDATE parcel SET address = :address WHERE number = :n AND status = :s",
		sql.Named("address", address),
		sql.Named("n", number),
		sql.Named("s", "registered"))
	if err != nil {
		return fmt.Errorf("exec failed: %w", err)
	}

	return nil
}

func (s ParcelStore) Delete(number int) error {
	_, err := s.db.Exec("DELETE FROM parcel WHERE number = :n AND status = :s",
		sql.Named("n", number),
		sql.Named("s", "registered"))
	if err != nil {
		return fmt.Errorf("exec failed: %w", err)
	}
	return nil
}
