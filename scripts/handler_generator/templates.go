package main

import (
	"bytes"
	"fmt"
	"go/format"
	"io"
	"strings"
	"text/template"
)

var (
	outputTemplate = template.Must(template.New("output").Parse(`
// Code generated by handler_generator. DO NOT EDIT.

package eventsub

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/dnsge/twitch-eventsub-framework/v2/bindings"
)

func deserializeAndCallHandler[EventType any](
	h *bindings.NotificationHeaders, 
	notification bindings.EventNotification,
	handler EventHandler[EventType],
) error {
	if handler == nil {
		return nil
	}

	var data EventType
	if err := json.Unmarshal(notification.Event, &data); err != nil {
		return err
	}

	go handler(*h, notification.Subscription, data)
	return nil
}

func (h *Handler) handleNotification(
	ctx context.Context,
	w http.ResponseWriter,
	bodyBytes []byte,
	headers *bindings.NotificationHeaders,
) {
	var notification bindings.EventNotification
	if err := json.Unmarshal(bodyBytes, &notification); err != nil {
		http.Error(w, "Invalid JSON body", http.StatusBadRequest)
		return
	}

	if h.BeforeHandleEvent != nil {
		h.BeforeHandleEvent(ctx, headers, &notification)
	}

	var err error
	selector := headers.SubscriptionType + "_" + headers.SubscriptionVersion
	switch selector { {{ .SwitchCases }}
	default:
		http.Error(w, "Unsupported notification type and version", http.StatusBadRequest)
		return
	}

	if err != nil {
		http.Error(w, "Invalid notification", http.StatusBadRequest)
		return
	}

	writeEmptyOK(w)
}
`))

	switchCaseTemplate = template.Must(template.New("switch").Parse(`
case "{{ .EventsubMessageType }}_{{ .EventsubMessageVersion }}":
	err = deserializeAndCallHandler(headers, notification, h.{{ .HandlerFieldName }});
`))
)

func renderOutputFile(cases []switchCase, writer io.Writer) error {
	var buf bytes.Buffer

	var caseStrings []string
	for _, handlerCase := range cases {
		var caseBuf bytes.Buffer
		err := switchCaseTemplate.Execute(&caseBuf, handlerCase)
		if err != nil {
			return fmt.Errorf("execute switch: %w", err)
		}
		caseStrings = append(caseStrings, caseBuf.String())
	}

	err := outputTemplate.Execute(&buf, struct {
		SwitchCases string
	}{
		SwitchCases: strings.Join(caseStrings, ""),
	})
	if err != nil {
		return fmt.Errorf("execute output: %w", err)
	}

	formatted, err := format.Source(buf.Bytes())
	if err != nil {
		return fmt.Errorf("format: %w", err)
	}

	_, err = writer.Write(formatted)
	return err
}
