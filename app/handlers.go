package app

import (
	"encoding/json"
	"fmt"
	"net/http"

	"go.uber.org/zap"

	"github.com/smeshkov/cab-data-researcher/app/handlers"
	"github.com/smeshkov/cab-data-researcher/cache"
	"github.com/smeshkov/cab-data-researcher/cfg"
	"github.com/smeshkov/cab-data-researcher/ctx"
	"github.com/smeshkov/cab-data-researcher/db"
)

type tripCountReq struct {
	Medallions []string `json:"medallions"`
	PickupDate string   `json:"pickupDate"`
	NoCache    bool     `json:"noCache,omitempty"`
}

type tripCountsResp struct {
	Counts []*medallionTripCount `json:"counts"`
}

type medallionTripCount struct {
	Medallion string `json:"medallion"`
	Value     int    `json:"value"`
}

func getTripCount(cdb cache.CabDBCache) func(rw http.ResponseWriter, req *http.Request) *handlers.AppError {
	return func(w http.ResponseWriter, r *http.Request) *handlers.AppError {
		defer cfg.LogSync()

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

		pickupDate, err := cfg.ParseDate(tripCount.PickupDate)
		if err != nil {
			return &handlers.AppError{
				Err:     err,
				Message: fmt.Sprintf("wrong pickupDate time format, e.g. \"%s\"", cfg.DateFormat),
				Code:    http.StatusBadRequest,
			}
		}

		l := zap.L().With(
			zap.String("medallions", fmt.Sprintf("%v", tripCount.Medallions)),
			zap.Int("medallion_count", len(tripCount.Medallions)),
			zap.Time("pickup_date", pickupDate),
		)

		var store db.CabDB = cdb

		if tripCount.NoCache {
			store = cdb.Delegate()
			l.Debug("get trip count not cached")
		}

		res, err := store.GetTripCount(ctx.WithLogReq(r, l), tripCount.Medallions, pickupDate)
		if err != nil {
			return &handlers.AppError{
				Err:     err,
				Message: "failed to get trip count",
				Code:    http.StatusInternalServerError,
				Context: ctx.WithLog(r.Context(), l),
			}
		}

		counts := make([]*medallionTripCount, 0)
		for _, v := range res {
			counts = append(counts, &medallionTripCount{Medallion: v.Medallion, Value: v.Count})
		}

		return handlers.Write(ctx.WithLogReq(r, l), w, &tripCountsResp{Counts: counts})
	}
}

func clearCache(cdb cache.CabCache) func(rw http.ResponseWriter, req *http.Request) *handlers.AppError {
	return func(w http.ResponseWriter, r *http.Request) *handlers.AppError {
		defer cfg.LogSync()
		l := zap.L()
		cdb.Clear(ctx.WithLogReq(r, l))
		l.Debug("cache cleared")
		w.WriteHeader(http.StatusNoContent)
		return nil
	}
}
