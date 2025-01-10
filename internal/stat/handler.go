package stat

import (
	"GoAdvanced/configs"
	"GoAdvanced/pkg/middleware"
	"GoAdvanced/pkg/res"
	"net/http"
	"time"
)

const (
	GroupByDay   = "day"
	GroupByMonth = "month"
)

type StatHandlerDeps struct {
	StatRepository *StatRepository
	Config         *configs.Config
}

type StatHandler struct {
	StatRepository *StatRepository
}

func NewStatHandler(router *http.ServeMux, Deps StatHandlerDeps) {
	handler := &StatHandler{
		StatRepository: Deps.StatRepository,
	}
	router.Handle("GET /stat", middleware.IsAuthed(handler.GetStat(), Deps.Config))
}

func (handler *StatHandler) GetStat() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		from, err := time.Parse(time.DateOnly, r.URL.Query().Get("from"))
		if err != nil {
			http.Error(w, "invalid date", http.StatusBadRequest)
			return
		}
		to, err := time.Parse(time.DateOnly, r.URL.Query().Get("to"))
		if err != nil {
			http.Error(w, "invalid date", http.StatusBadRequest)
			return
		}
		by := r.URL.Query().Get("by")
		if by != GroupByDay && by != GroupByMonth {
			http.Error(w, "invalid date", http.StatusBadRequest)
			return
		}
		stats := handler.StatRepository.GetStat(by, from, to)
		res.Json(w, http.StatusOK, stats)
	}
}
