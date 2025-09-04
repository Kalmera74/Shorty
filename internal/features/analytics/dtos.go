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
