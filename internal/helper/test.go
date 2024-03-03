package helper

import (
	"time"

	"github.com/aimerzarashi/timeslice"
)

func NewItem[T any](value T, startAt, endAt time.Time) *timeslice.Item[T] {
	item, err := timeslice.NewItem(value, startAt, endAt)
	if err != nil {
		panic(err)
	}
	return item
}