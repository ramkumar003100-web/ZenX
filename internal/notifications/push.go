package notifications

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"time"
)

type PushClient struct {
	HTTP *http.Client
}

func NewPushClient() *PushClient {
	return &PushClient{HTTP: &http.Client{Timeout: 10 * time.Second}}
}

func (p *PushClient) SendWebhook(ctx context.Context, endpoint string, payload any, retries int) error {
	body, _ := json.Marshal(payload)
	var err error
	for i := 0; i <= retries; i++ {
		req, _ := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		resp, doErr := p.HTTP.Do(req)
		if doErr == nil && resp.StatusCode < 300 {
			_ = resp.Body.Close()
			return nil
		}
		if resp != nil {
			_ = resp.Body.Close()
		}
		err = doErr
		time.Sleep(time.Duration(1<<i) * time.Second)
	}
	return err
}
