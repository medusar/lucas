package store

import "time"

var (
	id       = int64(1)
	slowReqs = make([]*slowReq, 1024)
)

//https://redis.io/commands/slowlog
type slowReq struct {
	//A unique progressive identifier for every slow log entry.
	//The ID is never reset in the course of the Redis server execution, only a server restart will reset it.
	id int64
	//The unix timestamp at which the logged command was processed
	timestamp time.Time
	//The amount of time needed for its execution, in microseconds.
	timeTake int64
	//The array composing the arguments of the command.
	args []string
}

//AddSlowLog is used to save a new slow log record
func AddSlowLog(start time.Time, args []string, timeTake int64) {
	id++
	slowReqs[id%1024] = &slowReq{id: id, timestamp: start, timeTake: timeTake, args: args}
}
