package handler

import (
	"net/http"

	"github.com/dotuan0601/go-backend-template/internal/response"
	"github.com/dotuan0601/go-backend-template/internal/service"
)

type StatisticHandler struct {
	service service.StatisticService
}

func NewStatisticHandler(
	service service.StatisticService,
) *StatisticHandler {

	return &StatisticHandler{
		service: service,
	}
}

func (h *StatisticHandler) Statistic(
	w http.ResponseWriter,
	r *http.Request,
) {

	data := h.service.GetStatistic()

	response.Success(
		w,
		data,
	)
}
