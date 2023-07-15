notify

import (
	"context"
	"fmt"
	"log"
	"os"

	"gopkg.in/mail.v2"
	"github.com/sentiweb/monitor-lib/utils"
)

type FileSender struct {
	path string
}

func NewFileSender(path string) (*FileSender) {
	return &FileSender{path: path}, nil
}

func (s* FileSender) Start() error {
	_, err := os.Stat(s.path)
	if(err != nil) {
		log.Println(fmt.Sprintf("Unable to access to '%s'", s.path))
		return err
	}
	f, err := os.Create(path + "/.touch")
	if(err != nil) {
		log.Println(fmt.Sprintf("Unable to write in'%s'", s.path))
		return err
	}
	defer f.Close()
	return nil
}


func (s* FileSender) Send(ctx context.Context, msg *mail.Message) error {
	fn := fmt.Sprintf("%s/%s.eml", s.path, utils.RandomName(8) )
	f, err := os.Create(fn)
	defer f.Close()
	if err != nil {
		log.Println(fmt.Sprintf("Unable to write file %s", fn))
		return err
	}
	log.Println(fmt.Sprintf("Sending fake to %s", fn))
	_, err = msg.WriteTo(f)
	return err
}