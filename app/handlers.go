package app

import (
	"encoding/json"
	"net/http"
	"time"

	"go.uber.org/zap"

	"github.com/smeshkov/cab-data-researcher/app/handlers"
	"github.com/smeshkov/cab-data-researcher/cache"
	"github.com/smeshkov/cab-data-researcher/db"
)

type tripCountReq struct {
	Medallions []string `json:"medallions"`
	PickupDate string   `json:"pickupDate"`
	NoCache    bool     `json:"noCache,omitempty"`
}

type tripCountResp struct {
	Count int `json:"count"`
}

func getTripCount(cdb cache.CabDBCache) func(rw http.ResponseWriter, req *http.Request) *handlers.AppError {
	return func(w http.ResponseWriter, r *http.Request) *handlers.AppError {

		var tripCount tripCountReq
		err := json.NewDecoder(r.Body).Decode(&tripCount)
		if err != nil {
			return &handlers.AppError{
				Err:     err,
				Message: "wrong request",
				Code:    http.StatusBadRequest,
			}
		}
		defer r.Body.Close()

		l := zap.L().With(zap.Int("medallion_count", len(tripCount.Medallions)))

		t, err := time.Parse(time.RFC3339, tripCount.PickupDate)
		if err != nil {
			return &handlers.AppError{
				Err:     err,
				Message: "wrong time format",
				Code:    http.StatusBadRequest,
			}
		}

		l = l.With(zap.Time("pickup_date", t))

		var store db.CabDB = cdb

		if tripCount.NoCache {
			store = cdb.Delegate()
			l.Debug("get trip count not cached")
		}

		count, err := store.GetTripCount(tripCount.Medallions, t)
		if err != nil {
			return &handlers.AppError{
				Err:     err,
				Message: "failed to get trip count",
				Code:    http.StatusInternalServerError,
				Log:     l,
			}
		}

		return handlers.WriteWithLog(w, &tripCountResp{Count: count}, l)
	}
}

func clearCache(cdb cache.CabCache) func(rw http.ResponseWriter, req *http.Request) *handlers.AppError {
	return func(w http.ResponseWriter, r *http.Request) *handlers.AppError {
		cdb.Clear()
		zap.L().Debug("cache cleared")
		return nil
	}
}
