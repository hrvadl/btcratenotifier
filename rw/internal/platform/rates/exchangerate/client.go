package exchangerate

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

const (
	operation = "usd to uah rate"
)

const (
	usd = "USD"
	uah = "UAH"
)

type usdUahResponse struct {
	ConversionRate float32 `json:"conversion_rate"`
	BaseCode       string  `json:"base_code"`
	TargetCode     string  `json:"target_code"`
	UpdatedAt      int     `json:"time_last_update_unix"`
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
	res := new(usdUahResponse)
	if err := c.getRate(ctx, res, usd, uah); err != nil {
		return 0, fmt.Errorf("%s: %w", operation, err)
	}

	return res.ConversionRate, nil
}

func (c Client) getRate(
	ctx context.Context,
	response any,
	from string,
	to string,
) error {
	url, err := url.Parse(fmt.Sprintf("%s/%s/pair/%s/%s", c.url, c.token, from, to))
	if err != nil {
		return fmt.Errorf("failed to parse url: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url.String(), nil)
	if err != nil {
		return fmt.Errorf("failed to construct request: %w", err)
	}

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
