package main

import (
	"github.com/stretchr/testify/assert"
	es "github.com/sunrongya/eventsourcing"
	"testing"
	"time"
)

func TestHandleVoteRecordCreatedEvent(t *testing.T) {
	ch := make(chan es.Command)
	event := &ReplyCreatedEvent{
		WithGuid:  es.WithGuid{Guid: es.NewGuid()},
		PostId:    es.NewGuid(),
		ParentId:  es.NewGuid(),
		AuthorId:  "author1",
		Body:      "body111",
		CreatedOn: time.Now(),
	}
	command := &AcceptNewReplyCommand{
		WithGuid:  es.WithGuid{Guid: event.PostId},
		ReplyId:   event.GetGuid(),
		AuthorId:  event.AuthorId,
		CreatedOn: event.CreatedOn,
	}

	handler := NewEventHandler(ch)
	go handler.HandleReplyCreatedEvent(event)
	select {
	case c := <-ch:
		assert.Equal(t, c, command, "HandleReplyCreatedEvent command error")
	case <-time.After(1 * time.Second):
		t.Error("HandleReplyCreatedEvent error")
	}
}
