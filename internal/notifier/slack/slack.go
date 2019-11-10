package slack

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"sync"
)

// Webhook sends notifications to webhook
type Webhook struct {
	Client   *http.Client
	Endpoint string
}

const payload = `{
  "attachments":[
    {
       "title":"%s",
       "pretext":"Much programing, such honor",
       "text":"%s"
    }
  ]
}`

var bufferPool = sync.Pool{
	New: func() interface{} {
		return new(bytes.Buffer)
	},
}

// Notify sends webhook message
func (w *Webhook) Notify(ctx context.Context, user, message string) error {
	buf := bufferPool.Get().(*bytes.Buffer)

	defer func() {
		buf.Reset()
		bufferPool.Put(buf)
	}()

	fmt.Fprintf(buf, payload, user, message)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, w.Endpoint, buf)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "Content-Type: application/json")

	resp, err := w.Client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("request failed: %d", resp.StatusCode)
	}

	return nil
}
