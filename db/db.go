package db

import (
	"database/sql"
	"time"
)

// NewCabDB creates new instance of cabDB.
func NewCabDB(driverName, dataSourceName string) (CabDBCloser, error) {
	db, err := sql.Open(driverName, dataSourceName)
	if err != nil {
		return nil, err
	}
	return &cabDB{db}, nil
}

type cabDB struct {
	db *sql.DB
}

func (cdb *cabDB) GetTripCount(medallions []string, pickupDatetime time.Time) (int, error) {
	return 0, nil
}

// Close closes the connection.
func (cdb *cabDB) Close() error {
	return cdb.db.Close()
}
