package cache

import (
	"context"
	
	"github.com/smeshkov/cab-data-researcher/db"
)

// CabCache ...
type CabCache interface {
	// Delegate allows to bypass cache.
	Delegate() db.CabDB

	// Clear clears cache.
	Clear(c context.Context)
}

// CabDBCache ...
type CabDBCache interface {
	db.CabDB
	CabCache
}
