package main

import (
	es "github.com/sunrongya/eventsourcing"
	"time"
)

// --------------
// Post Commands
// --------------
type CreatePostCommand struct {
	es.WithGuid
	Subject  string
	Body     string
	AuthorId string
}

type UpdatePostCommand struct {
	es.WithGuid
	Subject string
	Body    string
}

type AcceptNewReplyCommand struct {
	es.WithGuid
	ReplyId   es.Guid
	AuthorId  string
	CreatedOn time.Time
}

// --------------
// Reply Commands
// --------------
type CreateReplyCommand struct {
	es.WithGuid
	PostId    es.Guid
	ParentId  es.Guid
	AuthorId  string
	Body      string
	CreatedOn time.Time
}

type ChangeReplyBodyCommand struct {
	es.WithGuid
	Body string
}
