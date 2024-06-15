package gw

import (
	"fmt"
	"net/http"
	"time"

	vegeta "github.com/tsenart/vegeta/v12/lib"
	"github.com/tsenart/vegeta/v12/lib/prom"
)

const frequency = 100

func AttackGetRate(pm *prom.Metrics, url string) {
	rate := vegeta.Rate{Freq: frequency, Per: time.Second}
	duration := 1 * time.Minute
	targeter := vegeta.NewStaticTargeter(vegeta.Target{
		Method: http.MethodGet,
		URL:    fmt.Sprintf("http://%s/api/rate", url),
	})

	attacker := vegeta.NewAttacker()

	for res := range attacker.Attack(targeter, rate, duration, "Big Bang!") {
		pm.Observe(res)
	}
}
