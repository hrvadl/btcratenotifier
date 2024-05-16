package cryptocompare

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"strings"
)

const (
	operation = "btc to uah rate"
)

const (
	apiKeyQueryParam = "api_key"
	fromQueryParam   = "fsym"
	toQueryParam     = "tsyms"
)

const (
	btc = "BTC"
	uah = "UAH"
)

type btcUahResponse struct {
	UAH float32 `json:"UAH"`
}

type Client struct {
	token string
	url   string
}

func NewClient(token, url string) Client {
	return Client{
		token: token,
		url:   url,
	}
}

func (c Client) Convert(ctx context.Context) (float32, error) {
	res := new(btcUahResponse)
	if err := c.getRate(ctx, res, btc, uah); err != nil {
		return 0, fmt.Errorf("%s: %w", operation, err)
	}

	return res.UAH, nil
}

func (c Client) getRate(
	ctx context.Context,
	response any,
	from string,
	to ...string,
) error {
	url, err := url.Parse(fmt.Sprintf("%s/data/price", c.url))
	if err != nil {
		return fmt.Errorf("failed to parse url: %w", err)
	}

	q := url.Query()
	q.Set(fromQueryParam, from)
	q.Set(toQueryParam, strings.Join(to, ","))
	q.Set(apiKeyQueryParam, c.token)
	url.RawQuery = q.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url.String(), nil)
	if err != nil {
		return fmt.Errorf("failed to construct request: %w", err)
	}

	slog.Info("Constructed the URL", "url", url.String())
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}

	defer res.Body.Close()
	bytes, err := io.ReadAll(res.Body)
	if err != nil {
		return fmt.Errorf("failed to read body bytes: %w", err)
	}

	if err := json.Unmarshal(bytes, &response); err != nil {
		return fmt.Errorf("failed to parse response body: %w", err)
	}

	return nil
}
