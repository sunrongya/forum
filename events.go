package main

import (
	es "github.com/sunrongya/eventsourcing"
	"time"
)

// --------------
// Post Events
// --------------

type PostCreatedEvent struct {
	es.WithGuid
	Subject  string
	Body     string
	AuthorId string
}

type PostUpdatedEvent struct {
	es.WithGuid
	Subject string
	Body    string
}

type PostReplyStatisticInfoChangedEvent struct {
	es.WithGuid
	ReplyId es.Guid
	PostReplyStatisticInfo
}

type RepeatPostReplyChangedEvent struct {
	es.WithGuid
	ReplyId es.Guid
}

// ------------------
// Reply Events
// ------------------
type ReplyCreatedEvent struct {
	es.WithGuid
	PostId    es.Guid
	ParentId  es.Guid
	AuthorId  string
	Body      string
	CreatedOn time.Time
}

type ReplyBodyChangedEvent struct {
	es.WithGuid
	Body string
}
