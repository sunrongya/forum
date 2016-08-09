package main

import (
	"fmt"
	es "github.com/sunrongya/eventsourcing"
	"time"
)

type Reply struct {
	es.BaseAggregate
	_postId    es.Guid
	_parentId  es.Guid
	_authorId  string
	_body      string
	_createdOn time.Time
}

var _ es.Aggregate = (*Reply)(nil)

func NewReply() es.Aggregate {
	return &Reply{}
}

func (this *Reply) ProcessCreateReplyCommand(command *CreateReplyCommand) []es.Event {
	if len(command.Body) > 4000 {
		panic(fmt.Errorf("回复内容长度不能超过4000"))
	}
	return []es.Event{
		&ReplyCreatedEvent{
			PostId:    command.PostId,
			ParentId:  command.ParentId,
			AuthorId:  command.AuthorId,
			Body:      command.Body,
			CreatedOn: command.CreatedOn,
		},
	}
}

func (this *Reply) ProcessChangeReplyBodyCommand(command *ChangeReplyBodyCommand) []es.Event {
	if len(command.Body) > 4000 {
		panic(fmt.Errorf("回复内容长度不能超过4000"))
	}
	return []es.Event{
		&ReplyBodyChangedEvent{Body: command.Body},
	}
}

func (this *Reply) HandleReplyCreatedEvent(event *ReplyCreatedEvent) {
	this._postId = event.PostId
	this._parentId = event.ParentId
	this._authorId = event.AuthorId
	this._body = event.Body
	this._createdOn = event.CreatedOn
}

func (this *Reply) HandleReplyBodyChangedEvent(event *ReplyBodyChangedEvent) {
	this._body = event.Body
}
