package gw

import (
	"fmt"
	"net/http"
	"time"

	vegeta "github.com/tsenart/vegeta/v12/lib"
	"github.com/tsenart/vegeta/v12/lib/prom"
)

const (
	frequency = 1000
	duration  = 5 * time.Minute
)

func AttackGetRate(pm *prom.Metrics, url string) {
	rate := vegeta.Rate{Freq: frequency, Per: time.Second}
	targeter := vegeta.NewStaticTargeter(vegeta.Target{
		Method: http.MethodGet,
		URL:    fmt.Sprintf("http://%s/api/rate", url),
	})

	attacker := vegeta.NewAttacker()

	for res := range attacker.Attack(targeter, rate, duration, "Big Bang!") {
		pm.Observe(res)
	}
}
