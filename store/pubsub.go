package store

import (
	"github.com/medusar/lucas/protocol"
	"log"
)

//https://redis.io/topics/pubsub

var (
	channelReceivers = make(map[string][]protocol.RedisRW)
	receiverChannels = make(map[protocol.RedisRW]map[string]bool)
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
		if rw.IsClosed() {
			removeSubscriber(rw)
			continue
		}
		if err := rw.WriteArray(ret); err != nil {
			failedCount++
		}
	}
	channelReceivers[channel] = receivers
	return len(receivers) - failedCount
}

func Subscribe(rw protocol.RedisRW, channels []string) {
	if !rw.IsSubscriber() {
		rw.SetSubscriber(true)
	}

	subChans, ok := receiverChannels[rw]
	if !ok {
		subChans = make(map[string]bool)
	}
	ret := make([]*protocol.Resp, 3)
	ret[0] = protocol.NewBulk("subscribe")

	for _, channel := range channels {
		if _, ok := subChans[channel]; !ok {
			subChans[channel] = true
			receivers, ok := channelReceivers[channel]
			if !ok {
				receivers = make([]protocol.RedisRW, 0)
			}
			receivers = append(receivers, rw)
			channelReceivers[channel] = receivers
		}
		ret[1] = protocol.NewBulk(channel)
		ret[2] = protocol.NewInt(len(subChans))
		rw.WriteArray(ret)
	}
	receiverChannels[rw] = subChans
}

func removeSubscriber(rw protocol.RedisRW) {
	log.Println("remove closed connection:", rw)

	channels, ok := receiverChannels[rw]
	if !ok {
		return
	}
	for channel := range channels {
		rws, ok := channelReceivers[channel]
		if !ok {
			continue
		}
		for i, crw := range rws {
			if crw == rw {
				rws = append(rws[:i], rws[i+1:]...)
			}
		}
		channelReceivers[channel] = rws
	}
	delete(receiverChannels, rw)
}
