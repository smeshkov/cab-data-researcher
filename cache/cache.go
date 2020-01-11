package cache

import (
	"fmt"
	"sync"
	"time"

	"github.com/smeshkov/cab-data-researcher/db"
)

// NewCabCache creates new instance of CabDB cache.
func NewCabCache(delegate db.CabDB) CabDBCache {
	return &cabCache{delegate: delegate}
}

func (cc *cabCache) GetTripCount(medallions []string, pickupDatetime time.Time) (int, error) {
	var res int

	df := formatPickupDate(pickupDatetime)

	misses := make([]string, 0)

	for _, m := range medallions {
		key := toCacheKey(m, df)
		if v, ok := cc.sm.Load(key); ok {
			res += v.(int)
		} else {
			misses = append(misses, m)
		}
	}

	if len(misses) > 0 {
		cnt, err := cc.Delegate().GetTripCount(misses, pickupDatetime)
		if err != nil {
			return res, err
		}
		res += cnt
	}

	return res, nil
}

func (cc *cabCache) Delegate() db.CabDB {
	return cc.delegate
}

func (cc *cabCache) Clear() {
	cc.sm.Range(func (key, value interface{}) bool {
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

func formatPickupDate(pickupDatetime time.Time) string {
	return pickupDatetime.Format("02/01/2006")
}
