package command

import (
	"github.com/medusar/lucas/protocol"
	"github.com/medusar/lucas/store"
	"time"
)

//WithTime will monitor the time a request costs,
// if it costs more than 10 microseconds it will be added to the slow log log
func WithTime(realFunc cmdFunc) cmdFunc {
	return func(args []string, r *protocol.RedisConn) error {
		defer timeTrack(time.Now(), args)
		return realFunc(args, r)
	}
}

func timeTrack(start time.Time, args []string) {
	takes := time.Since(start)
	if takes > 10*time.Microsecond {
		store.AddSlowLog(start, args, int64(takes/time.Microsecond))
	}
}
