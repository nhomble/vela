package server

import (
	"fmt"
	"time"
)

const COMMON_LOG_FORMAT = "02/Jan/2006:15:04:05 -0700"

type CommonLog struct {
	clientHost   string
	time         time.Time
	method       string
	url          string
	status       string
	meta         string
	responseSize int
}

func (log CommonLog) format() string {
	return fmt.Sprintf("%s - - [%s] \"%s\" %s %s %d",
		log.clientHost,
		log.time.Format(COMMON_LOG_FORMAT),
		log.url,
		log.status,
		log.meta,
		log.responseSize,
	)
}
