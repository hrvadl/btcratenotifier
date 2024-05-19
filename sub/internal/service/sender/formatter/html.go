package formatter

import (
	"fmt"
	"time"
)

// NewWithDate constructs new HTML formatter for mails.
func NewWithDate() *WithDateFormatter {
	return &WithDateFormatter{}
}

// WithDateFormatter is a HTML formatter for mails,
// which will include date in the message and format float to
// 2 point precision.
type WithDateFormatter struct{}

// Format method taking exchange rate as a argument, then
// includes date in the message and formats float to
// 2 point precision.
func (hf *WithDateFormatter) Format(r float32) string {
	return fmt.Sprintf(
		"Latest exchange rate as for %v: 1 USD worth %.2f UAH",
		time.Now().Format(time.DateTime),
		r,
	)
}
