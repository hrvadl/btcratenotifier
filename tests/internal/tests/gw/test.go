package gw

import (
	"fmt"
	"net/http"
	"net/url"
	"time"

	vegeta "github.com/tsenart/vegeta/v12/lib"
	"github.com/tsenart/vegeta/v12/lib/prom"
)

const (
	frequency = 1000
	duration  = 5 * time.Minute
)

func NewLoadTest(pm *prom.Metrics, url string) *LoadTest {
	return &LoadTest{
		prometheus: pm,
		url:        url,
	}
}

type LoadTest struct {
	prometheus *prom.Metrics
	url        string
}

func (lt *LoadTest) GetRate() {
	rate := vegeta.Rate{Freq: frequency, Per: time.Second}
	targeter := vegeta.NewStaticTargeter(vegeta.Target{
		Method: http.MethodGet,
		URL:    fmt.Sprintf("http://%s/api/rate", lt.url),
	})

	attacker := vegeta.NewAttacker()

	for res := range attacker.Attack(targeter, rate, duration, "Big Bang!") {
		lt.prometheus.Observe(res)
	}
}

func (lt *LoadTest) Subscribe() {
	rate := vegeta.Rate{Freq: frequency, Per: time.Second}
	targeter := newSenderTargeter(lt.url)

	attacker := vegeta.NewAttacker()

	for res := range attacker.Attack(targeter, rate, duration, "Big Bang!") {
		lt.prometheus.Observe(res)
	}
}

func newSenderTargeter(addr string) vegeta.Targeter {
	return func(tgt *vegeta.Target) error {
		if tgt == nil {
			return vegeta.ErrNilTarget
		}

		mail := fmt.Sprintf("mail%v@mail.com", time.Now().Nanosecond())
		tgt.Body = []byte(url.Values{"email": {mail}}.Encode())
		header := http.Header{}
		header.Set("Content-Type", "application/x-www-form-urlencoded")
		tgt.Header = header
		tgt.Method = http.MethodPost
		tgt.URL = fmt.Sprintf("http://%s/api/subscribe", addr)

		return nil
	}
}
