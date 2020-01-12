package db

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"go.uber.org/zap"

	"github.com/smeshkov/cab-data-researcher/cfg"
	"github.com/smeshkov/cab-data-researcher/ctx"
)

// NewCabDB creates new instance of cabDB.
func NewCabDB(driverName, dataSourceName string) (CabDBCloser, error) {
	var err error

	db, err := sql.Open(driverName, dataSourceName)
	if err != nil {
		return nil, fmt.Errorf("failed to open DB connection: %w", err)
	}

	// wait for when DB is ready
	ok := false
	cnt := 0
	for !ok || cnt < 3 {
		cnt++
		err = db.Ping()
		if err != nil {
			time.Sleep(time.Second)
			continue
		}
		ok = true
	}

	if err != nil {
		return nil, fmt.Errorf("failed to ping DB: %w", err)
	}

	return &cabDB{db}, nil
}

type cabDB struct {
	db *sql.DB
}

func (db *cabDB) GetTripCount(c context.Context, medallions []string, pickupDatetime time.Time) ([]*CabTripCount, error) {
	res := make([]*CabTripCount, 0)

	fromDate := cfg.FormatDate(cfg.StartOfDay(pickupDatetime))
	toDate := cfg.FormatDate(cfg.StartOfDay(pickupDatetime.Add(cfg.Day)))

	l := ctx.GetLog(c).With(
		zap.String("from_date", fromDate),
		zap.String("to_date", toDate),
	)

	l.Debug("querying DB")

	selectTripCount, err := db.db.Prepare("SELECT medallion, COUNT(*) FROM cab_trip_data WHERE medallion IN (?) AND pickup_datetime >= ? AND pickup_datetime < ? GROUP BY medallion")
	if err != nil {
		return res, fmt.Errorf("failed to create DB statement: %w", err)
	}

	rows, err := selectTripCount.Query(strings.Join(medallions, ","), fromDate, toDate)
	if err != nil {
		return res, fmt.Errorf("failed to query DB: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		tripCount := CabTripCount{}
		err = rows.Scan(&tripCount.Medallion, &tripCount.Count)
		if err != nil {
			return res, fmt.Errorf("failed to scan DB rows: %w", err)
		}
		res = append(res, &tripCount)
	}

	l.Debug(fmt.Sprintf("found %d rows in DB", len(res)))

	return res, nil
}

// Close closes the connection.
func (db *cabDB) Close() error {
	return db.db.Close()
}
