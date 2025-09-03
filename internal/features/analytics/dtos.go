package analytics

import (
	"time"

	"github.com/Kalmera74/Shorty/internal/types"
)

type ClickRecord struct {
	ShortID    types.ShortId
	ClickTimes time.Time
	IpAddress  string
	UserAgents string
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
