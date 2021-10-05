package eventsub_framework

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	esb "github.com/dnsge/twitch-eventsub-bindings"
	"net/http"
	"time"
)

const EventSubSubscriptionsEndpoint = "https://api.twitch.tv/helix/eventsub/subscriptions"

type Credentials interface {
	ClientID() (string, error)
	AppToken() (string, error)
}

type SubRequest struct {
	Type      string
	Condition interface{}
	Callback  string
	Secret    string
}

type SubClient struct {
	httpClient  http.Client
	credentials Credentials
}

func NewSubClient(credentials Credentials) *SubClient {
	return &SubClient{
		httpClient: http.Client{
			Timeout: time.Second * 3,
		},
		credentials: credentials,
	}
}

func (s *SubClient) Do(req *http.Request) (*http.Response, error) {
	clientID, err := s.credentials.ClientID()
	if err != nil {
		return nil, fmt.Errorf("get client id: %w", err)
	}

	appToken, err := s.credentials.AppToken()
	if err != nil {
		return nil, fmt.Errorf("get app token: %w", err)
	}

	req.Header.Set("Client-ID", clientID)
	req.Header.Set("Authorization", "Bearer "+appToken)
	if req.Header.Get("Content-Type") == "" {
		req.Header.Set("Content-Type", "application/json")
	}
	return s.httpClient.Do(req)
}

func (s *SubClient) Subscribe(ctx context.Context, srq *SubRequest) (*esb.RequestStatus, error) {
	reqJSON := esb.Request{
		Type:      srq.Type,
		Version:   "1",
		Condition: srq.Condition,
		Transport: esb.Transport{
			Method:   "webhook",
			Callback: srq.Callback,
			Secret:   srq.Secret,
		},
	}

	buf := new(bytes.Buffer)
	err := json.NewEncoder(buf).Encode(reqJSON)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", EventSubSubscriptionsEndpoint, buf)
	if err != nil {
		return nil, err
	}
	res, err := s.Do(req)

	var statusResponse esb.RequestStatus
	if err := json.NewDecoder(res.Body).Decode(&statusResponse); err != nil {
		return nil, err
	}
	_ = res.Body.Close()

	return &statusResponse, nil
}
