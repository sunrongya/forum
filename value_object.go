package main

import (
	es "github.com/sunrongya/eventsourcing"
	"time"
)

type PostReplyStatisticInfo struct {
	LastReplyId       es.Guid
	LastReplyAuthorId string
	LastReplyTime     time.Time
	ReplyCount        int
}
