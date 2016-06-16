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

func (p *Reply) ApplyEvents(events []es.Event) {
	for _, event := range events {
		switch e := event.(type) {
		case *ReplyCreatedEvent:
			p._postId, p._parentId, p._authorId, p._body, p._createdOn = e.PostId, e.ParentId, e.AuthorId, e.Body, e.CreatedOn
		case *ReplyBodyChangedEvent:
			p._body = e.Body
		default:
			panic(fmt.Errorf("Unknown event %#v", e))
		}
	}
	p.SetVersion(len(events))
}

func (p *Reply) ProcessCommand(command es.Command) []es.Event {
	var event es.Event
	switch c := command.(type) {
	case *CreateReplyCommand:
		event = p.processCreateReplyCommand(c)
	case *ChangeReplyBodyCommand:
		event = p.processChangeReplyBodyCommand(c)
	default:
		panic(fmt.Errorf("Unknown command %#v", c))
	}
	event.SetGuid(command.GetGuid())
	return []es.Event{event}
}

func (p *Reply) processCreateReplyCommand(command *CreateReplyCommand) es.Event {
	if len(command.Body) > 4000 {
		panic(fmt.Errorf("回复内容长度不能超过4000"))
	}
	fmt.Println("processCreateReplyCommand")
	return &ReplyCreatedEvent{
		PostId:    command.PostId,
		ParentId:  command.ParentId,
		AuthorId:  command.AuthorId,
		Body:      command.Body,
		CreatedOn: command.CreatedOn,
	}
}

func (p *Reply) processChangeReplyBodyCommand(command *ChangeReplyBodyCommand) es.Event {
	if len(command.Body) > 4000 {
		panic(fmt.Errorf("回复内容长度不能超过4000"))
	}
	return &ReplyBodyChangedEvent{Body: command.Body}
}

func NewReply() es.Aggregate {
	return &Reply{}
}
