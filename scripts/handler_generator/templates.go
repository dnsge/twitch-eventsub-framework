package main

import (
	"bytes"
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
	"encoding/json"
	"net/http"

	"github.com/dnsge/twitch-eventsub-framework/v2/bindings"
)

func (s *SubHandler) handleNotification(w http.ResponseWriter, bodyBytes []byte, h *bindings.NotificationHeaders) {
	var notification bindings.EventNotification
	if err := json.Unmarshal(bodyBytes, &notification); err != nil {
		http.Error(w, "Invalid JSON body", http.StatusBadRequest)
		return
	}

	switch h.SubscriptionType { {{ .SwitchCases }}
	default:
		http.Error(w, "Unknown notification type", http.StatusBadRequest)
		return
	}

	writeEmptyOK(w)
}
`))

	switchCaseTemplate = template.Must(template.New("switch").Parse(`
case "{{ .EventsubMessageType }}":
	var data {{ .EventType }}
	if err := json.Unmarshal(notification.Event, &data); err != nil {
		http.Error(w, "Invalid JSON body", http.StatusBadRequest)
		return
	}
	if s.{{ .HandlerFieldName }} != nil {
		go s.{{ .HandlerFieldName }}(h, &data)
	}
`))
)

type switchCase struct {
	EventsubMessageType string
	EventType           string
	HandlerFieldName    string
}

func renderOutputFile(cases []switchCase, writer io.Writer) error {
	var buf bytes.Buffer

	var caseStrings []string
	for _, handlerCase := range cases {
		var caseBuf bytes.Buffer
		err := switchCaseTemplate.Execute(&caseBuf, handlerCase)
		if err != nil {
			return err
		}
		caseStrings = append(caseStrings, caseBuf.String())
	}

	err := outputTemplate.Execute(&buf, struct {
		SwitchCases string
	}{
		SwitchCases: strings.Join(caseStrings, ""),
	})
	if err != nil {
		return err
	}

	formatted, err := format.Source(buf.Bytes())
	if err != nil {
		return err
	}

	_, err = writer.Write(formatted)
	return err
}