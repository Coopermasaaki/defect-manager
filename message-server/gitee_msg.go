package messageserver

import (
	"errors"
	"fmt"

	sdk "github.com/opensourceways/go-gitee/gitee"

	"github.com/opensourceways/defect-manager/issue"
)

const (
	msgHeaderUUID      = "X-Gitee-Timestamp"
	msgHeaderUserAgent = "User-Agent"
	msgHeaderEventType = "X-Gitee-Event"
)

type giteeEventHandler struct {
	userAgent string
	handler   issue.EventHandler
}

func (msg *giteeEventHandler) handle(payload []byte, header map[string]string) error {
	eventType, err := msg.parseRequest(header)
	if err != nil {
		return fmt.Errorf("invalid msg, err:%s", err.Error())
	}

	switch eventType {
	case sdk.EventTypeIssue:
		e, err := sdk.ConvertToIssueEvent(payload)
		if err != nil {
			return err
		}

		return msg.handler.HandleIssueEvent(&e)

	case sdk.EventTypeNote:
		e, err := sdk.ConvertToNoteEvent(payload)
		if err != nil {
			return err
		}

		return msg.handler.HandleNoteEvent(&e)

	default:
		return nil
	}
}

func (msg *giteeEventHandler) parseRequest(header map[string]string) (
	eventType string, err error,
) {
	if header == nil {
		err = errors.New("no header")

		return
	}

	if header[msgHeaderUserAgent] != msg.userAgent {
		err = errors.New("unknown " + msgHeaderUserAgent)

		return
	}

	if eventType = header[msgHeaderEventType]; eventType == "" {
		err = errors.New("missing " + msgHeaderEventType)

		return
	}

	if header[msgHeaderUUID] == "" {
		err = errors.New("missing " + msgHeaderUUID)
	}

	return
}
