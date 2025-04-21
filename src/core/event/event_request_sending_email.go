package event

import (
	"github.com/KhaiHust/email-notification-service/core/event/message"
	"github.com/golibs-starter/golib/web/event"
)

type EventRequestSendingEmail struct {
	*event.AbstractEvent
	PayloadData *message.EmailRequestSendingMessage `json:"payload"`
}

func (a EventRequestSendingEmail) Payload() interface{} {
	return a.PayloadData
}
func (a EventRequestSendingEmail) String() string {
	return a.ToString(a)
}
