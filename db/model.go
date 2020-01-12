package db

import (
	"context"
	"io"
	"time"
)

// CabDB ...
type CabDB interface {
	GetTripCount(c context.Context, medallions []string, pickupDatetime time.Time) ([]*CabTripCount, error)
}

// CabDBCloser ...
type CabDBCloser interface {
	CabDB
	io.Closer
}

// CabTripCount ...
type CabTripCount struct {
	Medallion string
	Count     int
}
