package email

import (
	"context"

	"gopkg.in/mail.v2"
)

// MemorySender stores message in memory. Only for testing purpose
type MemorySender struct {
	messages []*mail.Message
}

func NewMemorySender() *MemorySender {
	messages := make([]*mail.Message, 0)
	return &MemorySender{messages: messages}
}

func (s *MemorySender) Start() error {
	return nil
}

func (s *MemorySender) Send(ctx context.Context, msg *mail.Message) error {
	s.messages = append(s.messages, msg)
	return nil
}
