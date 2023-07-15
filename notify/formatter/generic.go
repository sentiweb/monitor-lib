package formatter

import(
	"fmt"
	"time"
	"github.com/sentiweb/monitor-lib/notify/types"
)

// GenericFormatter provides very simple formatting for notification.
type GenericFormatter struct {

}

func (g *GenericFormatter) Title(n types.Notification) string {
	return  n.Label()
}

func (g *GenericFormatter) Text(n types.Notification) string {
	var what string
	if n.Status() == types.NotificationStatusUp {
		what = "is Online"
	} else {
		what = "is Offline"
	}
	return fmt.Sprintf("%s %s from %s", n.Label(), what, n.FromTime().Format(time.RFC822))
}