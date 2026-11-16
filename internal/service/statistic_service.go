package service

import (
	"runtime"
	"time"
)

type StatisticService interface {
	GetStatistic() map[string]interface{}
}

type statisticService struct {
	startTime time.Time
}

func NewStatisticService() StatisticService {

	return &statisticService{
		startTime: time.Now(),
	}
}

func (s *statisticService) GetStatistic() map[string]interface{} {

	var mem runtime.MemStats

	runtime.ReadMemStats(
		&mem,
	)

	return map[string]interface{}{

		"uptime": time.Since(
			s.startTime,
		).String(),

		"goroutines": runtime.NumGoroutine(),

		"memory": map[string]interface{}{

			"alloc_mb": mem.Alloc / 1024 / 1024,

			"total_alloc_mb": mem.TotalAlloc / 1024 / 1024,

			"sys_mb": mem.Sys / 1024 / 1024,

			"heap_alloc_mb": mem.HeapAlloc / 1024 / 1024,

			"heap_objects": mem.HeapObjects,
		},

		"gc": map[string]interface{}{

			"num_gc": mem.NumGC,

			"last_gc": time.Unix(
				0,
				int64(mem.LastGC),
			),
		},
	}
}
