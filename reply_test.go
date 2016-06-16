package main

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	es "github.com/sunrongya/eventsourcing"
	"strings"
	"testing"
	"time"
)

func TestReplyRestore(t *testing.T) {
	postId, parentId := es.NewGuid(), es.NewGuid()
	now := time.Now()
	reply := &Reply{}
	reply.ApplyEvents([]es.Event{
		&ReplyCreatedEvent{
			PostId:    postId,
			ParentId:  parentId,
			AuthorId:  "sry",
			Body:      "body1",
			CreatedOn: now,
		},
		&ReplyBodyChangedEvent{
			Body: "bodyttttt",
		},
	})
	assert.Equal(t, 2, reply.Version(), "version error")
	assert.Equal(t, postId, reply._postId, "postId 错误")
	assert.Equal(t, parentId, reply._parentId, "parentId 错误")
	assert.Equal(t, "sry", reply._authorId, "authorId 错误")
	assert.Equal(t, "bodyttttt", reply._body, "body 错误")
	assert.Equal(t, now, reply._createdOn, "createdOn 错误")
}

func TestReplyRestoreForErrorEvent(t *testing.T) {
	assert.Panics(t, func() {
		NewReply().ApplyEvents([]es.Event{&struct{ es.WithGuid }{}})
	}, "restore error event must panic error")
}

func TestCheckReplyApplyEvents(t *testing.T) {
	events := []es.Event{
		&ReplyCreatedEvent{},
		&ReplyBodyChangedEvent{},
	}
	assert.NotPanics(t, func() { NewReply().ApplyEvents(events) }, "Check Process All Event")
}

func TestReplyCommand(t *testing.T) {
	guid, postId, parentId := es.NewGuid(), es.NewGuid(), es.NewGuid()
	now := time.Now()
	tests := []struct {
		reply   *Reply
		command es.Command
		event   es.Event
	}{
		{
			&Reply{},
			&CreateReplyCommand{WithGuid: es.WithGuid{Guid: guid}, ParentId: parentId, PostId: postId, Body: "body1", AuthorId: "author1", CreatedOn: now},
			&ReplyCreatedEvent{WithGuid: es.WithGuid{Guid: guid}, ParentId: parentId, PostId: postId, Body: "body1", AuthorId: "author1", CreatedOn: now},
		},
		{
			&Reply{},
			&ChangeReplyBodyCommand{WithGuid: es.WithGuid{Guid: guid}, Body: "body2"},
			&ReplyBodyChangedEvent{WithGuid: es.WithGuid{Guid: guid}, Body: "body2"},
		},
	}

	for _, v := range tests {
		assert.Equal(t, []es.Event{v.event}, v.reply.ProcessCommand(v.command))
	}
}

func TestReplyCommand_Panic(t *testing.T) {
	tests := []struct {
		reply   *Reply
		command es.Command
	}{
		{
			&Reply{},
			&struct{ es.WithGuid }{},
		},
		{
			&Reply{},
			&CreateReplyCommand{ParentId: es.NewGuid(), PostId: es.NewGuid(), Body: strings.Repeat("s", 4001), AuthorId: "author1"},
		},
		{
			&Reply{},
			&ChangeReplyBodyCommand{Body: strings.Repeat("s", 4001)},
		},
	}

	for _, v := range tests {
		assert.Panics(t, func() { v.reply.ProcessCommand(v.command) }, fmt.Sprintf("test panics error: command:%v", v.command))
	}
}
