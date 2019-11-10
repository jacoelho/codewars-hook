package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"

	"github.com/jacoelho/codewars/internal/user"
)

const apiEndpoint = "https://www.codewars.com/api/v1/users/"

type Repo struct {
	Client   *http.Client
	Endpoint string

	values map[string]user.User
	mu     sync.RWMutex
}

func New(client *http.Client) *Repo {
	return &Repo{
		Client:   client,
		Endpoint: apiEndpoint,
		values:   make(map[string]user.User),
	}
}

func (r *Repo) fetchFromAPI(ctx context.Context, id string) (user.User, error) {
	var fetchedUser user.User

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, r.Endpoint+id, nil)
	if err != nil {
		return fetchedUser, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "Content-Type: application/json")

	resp, err := r.Client.Do(req)
	if err != nil {
		return fetchedUser, fmt.Errorf("failed to send request: %w", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fetchedUser, fmt.Errorf("request failed: %d", resp.StatusCode)
	}

	err = json.NewDecoder(resp.Body).Decode(&fetchedUser)

	return fetchedUser, err
}

func (r *Repo) GetUserByID(ctx context.Context, id string) (user.User, error) {
	r.mu.RLock()

	u, ok := r.values[id]
	if ok {
		r.mu.RUnlock()
		return u, nil
	}

	r.mu.RUnlock()

	u, err := r.fetchFromAPI(ctx, id)
	if err != nil {
		return u, err
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	r.values[id] = u

	return u, nil
}
