package formatter

import (
	"fmt"
	"time"
)

func NewWithDate() *WithDateFormatter {
	return &WithDateFormatter{}
}

type WithDateFormatter struct{}

func (hf *WithDateFormatter) Format(r float32) string {
	return fmt.Sprintf(
		"Latest exchange rate as for %v: 1 USD worth %.2f UAH",
		time.Now().Format(time.DateTime),
		r,
	)
}
