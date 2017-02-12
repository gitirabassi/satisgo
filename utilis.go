package satisgo

import (
	"strconv"
	"time"

	uuid "github.com/satori/go.uuid"
)

func putTime(t time.Time) string {
	return t.Format("2006-01-02T15:04:05.0000Z")
}

func getTime(s string) *time.Time {
	t, err := time.Parse("2006-01-02T15:04:05.0000Z", s)
	if err != nil {
		return nil
	}
	return &t
}

func putUnix(t time.Time) string {
	n := t.UnixNano()
	n = n / 1000000
	s := strconv.FormatInt(n, 10)
	return s
}

func putMoney(m float64) uint64 {
	return uint64(m * 100)
}

func getMoney(n uint64) float64 {
	var val float64
	val = float64(n) / 100
	return val
}

func generateUUID() string {
	u := uuid.NewV4()
	return u.String()
}
