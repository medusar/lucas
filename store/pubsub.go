package store

import (
	"github.com/medusar/lucas/protocol"
)

//https://redis.io/topics/pubsub

var (
	channelReceivers = make(map[string][]protocol.RedisRW)
	receiverChannels = make(map[string][]string)
)

func Publish(channel, message string) int {
	receivers, ok := channelReceivers[channel]
	if !ok {
		return 0
	}

	ret := make([]*protocol.Resp, 3)
	ret[0] = protocol.NewBulk("message")
	ret[1] = protocol.NewBulk(channel)
	ret[2] = protocol.NewBulk(message)

	failedCount := 0
	for _, rw := range receivers {
		if err := rw.WriteArray(ret); err == nil {
			failedCount++
		}
	}
	return len(receivers) - failedCount
}

func Subscribe(r protocol.RedisRW, channels []string) {
	//TODO:
}
