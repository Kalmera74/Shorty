package analytics

import (
	"time"

	"github.com/Kalmera74/Shorty/internal/types"
)

type ClickEvent struct {
	ShortID   types.ShortId `json:"short_id"`
	Ip        string        `json:"ip"`
	UserAgent string        `json:"user_agent"`
	TimeStamp time.Time     `json:"time_stamp"`
}

type Analysis struct {
	ShortUrl     string
	OriginalUrl  string
	UsageDetails []Usage
}

type Usage struct {
	ClickTimes time.Time
	IpAddress  string
	UserAgents string
}
type PaginatedAnalytics struct {
	Total     int        `json:"total"`
	Page      int        `json:"page"`
	PageSize  int        `json:"page_size"`
	Analytics []Analysis `json:"analytics"`
}
type PaginatedAnalysis struct {
	Total   int      `json:"total"`
	Page    int      `json:"page"`
	Limit   int      `json:"limit"`
	Results Analysis `json:"results"`
}

type PaginatedClicks struct {
	Total  int          `json:"total"`
	Page   int          `json:"page"`
	Limit  int          `json:"limit"`
	Clicks []ClickEvent `json:"clicks"`
}
