package main

import (
	es "github.com/sunrongya/eventsourcing"
)

type EventHandler struct {
	_postChan chan<- es.Command
}

func (this *EventHandler) HandleReplyCreatedEvent(event *ReplyCreatedEvent) {
	c := &AcceptNewReplyCommand{
		WithGuid:  es.WithGuid{Guid: event.PostId},
		ReplyId:   event.GetGuid(),
		AuthorId:  event.AuthorId,
		CreatedOn: event.CreatedOn,
	}
	this._postChan <- c
}

func NewEventHandler(postChan chan<- es.Command) *EventHandler {
	return &EventHandler{
		_postChan: postChan,
	}
}
