package main

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

func formatUTCTime(ds string) string {
	var s string = ""
	var t time.Time
	var err error
	var dstr string = ds
	var year, mon, day int
	var min, hour int
	var month time.Month
	year, _ = strconv.Atoi(ds[:2])
	year += 2000
	mon, _ = strconv.Atoi(ds[2:4])
	switch mon {
	case 1:
		month = time.January
	case 2:
		month = time.February
	case 3:
		month = time.March
	case 4:
		month = time.April
	case 5:
		month = time.May
	case 6:
		month = time.June
	case 7:
		month = time.July
	case 8:
		month = time.August
	case 9:
		month = time.September
	case 10:
		month = time.October
	case 11:
		month = time.November
	case 12:
		month = time.December
	default:
		month = time.January
	}
	day, _ = strconv.Atoi(ds[4:6])
	hour, _ = strconv.Atoi(ds[6:8])
	min, _ = strconv.Atoi(ds[8:10])
	t = time.Date(year, month, day, hour, min, 0, 0, time.UTC)
	if err != nil {
		s += dstr
		fmt.Fprintf(os.Stderr, "can not parse [%s] err[%s]\n", dstr, err.Error())
	} else {
		s += t.Format(time.RFC1123Z)
	}
	return s
}

func main() {
	for _, c := range os.Args[1:] {
		fmt.Fprintf(os.Stdout, "%s => %s\n", c, formatUTCTime(c))
	}
}
