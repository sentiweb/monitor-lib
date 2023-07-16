package email

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/sentiweb/monitor-lib/utils"
	"gopkg.in/mail.v2"
)

// FileSender stores message in a directory in the filesystem. Only for testing purpose
type FileSender struct {
	path string
}

func NewFileSender(path string) *FileSender {
	return &FileSender{path: path}
}

func (s *FileSender) Start() error {
	_, err := os.Stat(s.path)
	if err != nil {
		log.Printf("Unable to access to '%s'", s.path)
		return err
	}
	f, err := os.Create(s.path + "/.touch")
	if err != nil {
		log.Printf("Unable to write in'%s'", s.path)
		return err
	}
	defer f.Close()
	return nil
}

func (s *FileSender) Send(ctx context.Context, msg *mail.Message) error {
	fn := fmt.Sprintf("%s/%s.eml", s.path, utils.RandomName(8))
	f, err := os.Create(fn)
	if err != nil {
		log.Printf("Unable to write file %s", fn)
		return err
	}
	defer f.Close()
	log.Printf("Sending fake to %s", fn)
	_, err = msg.WriteTo(f)
	return err
}
