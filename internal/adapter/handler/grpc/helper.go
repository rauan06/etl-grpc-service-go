package handler

import (
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"
)

func parseTime(ts string) *timestamppb.Timestamp {
	t, err := time.Parse(time.RFC3339, ts)
	if err != nil {
		return nil // or use timestamppb.Now() or log error
	}
	return timestamppb.New(t)
}
