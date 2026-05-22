package main

import (
	"time"
)

func get_time_from_str(s string) (retv time.Time, err error) {
	retv, err = time.Parse(time.DateTime, s)
	if err == nil {
		return
	}
	retv, err = time.Parse(time.DateOnly, s)
	return
}
