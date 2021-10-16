package eventsub_framework

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	esb "github.com/dnsge/twitch-eventsub-bindings"
	"net/http"
	"net/url"
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

type Status string

const (
	StatusAny                  Status = ""
	StatusEnabled              Status = "enabled"
	StatusVerificationPending  Status = "webhook_callback_verification_pending"
	StatusVerificationFailed   Status = "webhook_callback_verification_failed"
	StatusFailuresExceeded     Status = "notification_failures_exceeded"
	StatusAuthorizationRevoked Status = "authorization_revoked"
	StatusUserRemoved          Status = "user_removed"
)

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

func (s *SubClient) do(req *http.Request) (*http.Response, error) {
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
	res, err := s.do(req)
	if err != nil {
		return nil, err
	}

	var statusResponse esb.RequestStatus
	if err := json.NewDecoder(res.Body).Decode(&statusResponse); err != nil {
		return nil, err
	}
	_ = res.Body.Close()

	return &statusResponse, nil
}

func (s *SubClient) Unsubscribe(ctx context.Context, subscriptionID string) error {
	u, err := url.Parse(EventSubSubscriptionsEndpoint)
	if err != nil {
		return fmt.Errorf("unsubscribe: parse EventSubSubscriptionsEndpoint url: %w", err)
	}

	q := u.Query()
	q.Set("id", subscriptionID)
	u.RawQuery = q.Encode()

	req, err := http.NewRequestWithContext(ctx, "DELETE", u.String(), nil)
	if err != nil {
		return err
	}
	res, err := s.do(req)
	if err != nil {
		return err
	}

	if res.StatusCode < 200 || res.StatusCode >= 300 {
		return fmt.Errorf("unsubscribe: bad status code %d (%s)", res.StatusCode, http.StatusText(res.StatusCode))
	}

	return nil
}

// GetSubscriptions returns all EventSub subscriptions.
// If statusFilter != StatusAny, it will apply the filter to the query.
func (s *SubClient) GetSubscriptions(ctx context.Context, statusFilter Status) (*esb.RequestStatus, error) {
	var reqUrl string
	if statusFilter == StatusAny {
		reqUrl = EventSubSubscriptionsEndpoint
	} else {
		u, err := url.Parse(EventSubSubscriptionsEndpoint)
		if err != nil {
			return nil, fmt.Errorf("get subscriptions: parse EventSubSubscriptionsEndpoint url: %w", err)
		}
		q := u.Query()
		q.Set("status", string(statusFilter))
		u.RawQuery = q.Encode()
		reqUrl = u.String()
	}

	req, err := http.NewRequestWithContext(ctx, "GET", reqUrl, nil)
	if err != nil {
		return nil, err
	}
	res, err := s.do(req)
	if err != nil {
		return nil, err
	}

	var subscriptionsResponse esb.RequestStatus
	if err := json.NewDecoder(res.Body).Decode(&subscriptionsResponse); err != nil {
		return nil, err
	}
	_ = res.Body.Close()

	return &subscriptionsResponse, nil
}
