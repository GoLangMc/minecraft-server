package util

import (
	"time"

	"github.com/hako/durafmt"
)

func FormatTime(durationInSeconds int64) string {
	return durafmt.Parse(time.Second * time.Duration(durationInSeconds)).String()
}
