package cache

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/smeshkov/cab-data-researcher/cfg"
	"github.com/smeshkov/cab-data-researcher/ctx"
	"github.com/smeshkov/cab-data-researcher/db"
)

// NewCabCache creates new instance of CabDB read-through cache.
func NewCabCache(delegate db.CabDB) CabDBCache {
	return &cabCache{delegate: delegate}
}

func (cc *cabCache) GetTripCount(c context.Context, medallions []string, pickupDatetime time.Time) ([]*db.CabTripCount, error) {
	res := make([]*db.CabTripCount, 0)

	l := ctx.GetLog(c)

	df := cfg.FormatDate(pickupDatetime.UTC())

	// cache misses
	misses := make([]string, 0)
	for _, m := range medallions {
		key := toCacheKey(m, df)
		if v, ok := cc.sm.Load(key); ok {
			res = append(res, &db.CabTripCount{Medallion: m, Count: v.(int)})
		} else {
			misses = append(misses, m)
		}
	}

	l.Debug(fmt.Sprintf("number of cached entries found: %d", len(res)))

	if len(misses) > 0 {
		l.Debug(fmt.Sprintf("fetching from DB, number of cache misses: %d", len(misses)))
		cnts, err := cc.Delegate().GetTripCount(c, misses, pickupDatetime)
		if err != nil {
			return res, err
		}

		l.Debug(fmt.Sprintf("number of DB entries found: %d", len(cnts)))

		// Update cache entries
		for _, v := range cnts {
			key := toCacheKey(v.Medallion, df)
			cc.sm.Store(key, v.Count)
		}

		res = append(res, cnts...)
	}

	return res, nil
}

func (cc *cabCache) Delegate() db.CabDB {
	return cc.delegate
}

func (cc *cabCache) Clear(c context.Context) {

	cc.sm.Range(func(key, value interface{}) bool {
		cc.sm.Delete(key)
		return true
	})
}

type cabCache struct {
	sm       sync.Map
	delegate db.CabDB
}

type cabCacheKey string

func toCacheKey(medallion, pickupDatetime string) cabCacheKey {
	return cabCacheKey(fmt.Sprintf("%s|%s", medallion, pickupDatetime))
}
