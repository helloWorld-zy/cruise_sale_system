package service

const (
	ChannelSMS    = "sms"
	ChannelWechat = "wechat"
	ChannelInbox  = "inbox"
)

type Notifier interface {
	Send(channel, template, payload string) error
}

type NotifyService struct{ n Notifier }

func NewNotifyService(n Notifier) *NotifyService { return &NotifyService{n: n} }

func (s *NotifyService) Send(channel, template, payload string) error {
	return s.n.Send(channel, template, payload)
}
