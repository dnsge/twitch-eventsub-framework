package eventsub_framework

import (
	"bytes"
	"context"
	esb "github.com/dnsge/twitch-eventsub-bindings"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

const secret = `hey this is really secret`

func TestSubHandler_ServeHTTP_VerificationBasic(t *testing.T) {
	handler := NewSubHandler(true, []byte(secret))
	res := handleVerificationRequest(handler, newVerificationRequest)
	body, _ := io.ReadAll(res.Body)
	_ = res.Body.Close()

	assert.True(t, isOK(res.StatusCode))
	assert.Equal(t, body, []byte("olYc8-klwIH9BthhWWhTU-AhJQ0eatixVF2y6x3G5kk"))
}

func TestSubHandler_ServeHTTP_VerificationInvalidSignature(t *testing.T) {
	handler := NewSubHandler(true, []byte(secret))
	res := handleVerificationRequest(handler, newBadVerificationRequest)
	_ = res.Body.Close()

	assert.Equal(t, res.StatusCode, http.StatusForbidden)
}

func TestSubHandler_ServeHTTP_VerificationDynamic(t *testing.T) {
	handler := NewSubHandler(false, nil)
	handler.VerifyChallenge = func(h *esb.ResponseHeaders, chal *esb.SubscriptionChallenge) bool {
		return h.SubscriptionType == "channel.update"
	}

	res := handleVerificationRequest(handler, newVerificationRequestWithType("channel.update"))
	body, _ := io.ReadAll(res.Body)
	_ = res.Body.Close()

	assert.True(t, isOK(res.StatusCode))
	assert.Equal(t, body, []byte("olYc8-klwIH9BthhWWhTU-AhJQ0eatixVF2y6x3G5kk"))

	res = handleVerificationRequest(handler, newVerificationRequestWithType("channel.follow"))
	body, _ = io.ReadAll(res.Body)
	_ = res.Body.Close()

	assert.Equal(t, res.StatusCode, http.StatusBadRequest)
}

func TestSubHandler_ServeHTTP_IDTracker(t *testing.T) {
	handler := NewSubHandler(false, nil)
	tracker := &wrapper{m: NewMapTracker(), DuplicateSeen: false}
	handler.IDTracker = tracker

	// First request should succeed
	res := handleVerificationRequest(handler, newVerificationRequest)
	_ = res.Body.Close()

	assert.True(t, isOK(res.StatusCode))
	assert.False(t, tracker.DuplicateSeen)

	// Second request should fail as duplicate
	res = handleVerificationRequest(handler, newVerificationRequest)
	_ = res.Body.Close()

	assert.True(t, isOK(res.StatusCode)) // should still give 2xx code
	assert.True(t, tracker.DuplicateSeen)
}

func handleVerificationRequest(handler *SubHandler, reqFactory func() *http.Request) *http.Response {
	verificationReq := reqFactory()
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, verificationReq)
	return w.Result()
}

func newVerificationRequest() *http.Request {
	bodyData := []byte(`{"subscription":{"id":"ef7e8fba-6c32-4ead-965d-61f21660d095","status":"webhook_callback_verification_pending","type":"channel.update","version":"1","condition":{"broadcaster_user_id":"132532813"},"transport":{"method":"webhook","callback":"https://testing.proxy.b.dnsge.org/webhooks"},"created_at":"2023-03-09T04:44:48.057734342Z","cost":0},"challenge":"olYc8-klwIH9BthhWWhTU-AhJQ0eatixVF2y6x3G5kk"}`)

	req := httptest.NewRequest("POST", "/", bytes.NewReader(bodyData))
	req.Header = http.Header{
		"Content-Type":                                     {"application/json"},
		"Twitch-Eventsub-Message-Id":                       {"e7f8151c-5849-48d5-8db9-234618442877"},
		"Twitch-Eventsub-Message-Retry":                    {"1"},
		"Twitch-Eventsub-Message-Signature":                {"sha256=9afe337a0526eda98c12cd5c6892f5eec5c86a2f0fde7e0655764382d464bce8"},
		"Twitch-Eventsub-Message-Timestamp":                {"2023-03-09T04:44:48.062323705Z"},
		"Twitch-Eventsub-Message-Type":                     {"webhook_callback_verification"},
		"Twitch-Eventsub-Subscription-Is-Batching-Enabled": {"false"},
		"Twitch-Eventsub-Subscription-Type":                {"channel.update"},
		"Twitch-Eventsub-Subscription-Version":             {"1"},
	}

	return req
}

func newBadVerificationRequest() *http.Request {
	req := newVerificationRequest()
	// overwrite header with invalid signature
	req.Header.Set("Twitch-Eventsub-Message-Signature", "sha256=9afe337a0526eda98c12cd5c6892f5eec5c86a2f0fde7e0655764382d464bce9")
	return req
}

func newVerificationRequestWithType(typ string) func() *http.Request {
	return func() *http.Request {
		req := newVerificationRequest()
		req.Header.Set("Twitch-Eventsub-Subscription-Type", typ)
		return req
	}
}

func isOK(statusCode int) bool {
	return 200 <= statusCode && statusCode < 300
}

type wrapper struct {
	m             *MapTracker
	DuplicateSeen bool
}

func (w *wrapper) AddAndCheckIfDuplicate(ctx context.Context, id string) (bool, error) {
	a, e := w.m.AddAndCheckIfDuplicate(ctx, id)
	if a {
		w.DuplicateSeen = true
	}
	return a, e
}
