package helpfunc

import (
	"fmt"
	"github.com/rs/xid"
	"time"
)

func CreateRequestId() string {
	return fmt.Sprintf("%s:%s", xid.New().String(), time.Now().Format("20060102150405"))
}
