package main

import (
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
	replyCreatedEvent := &ReplyCreatedEvent{
		PostId:    postId,
		ParentId:  parentId,
		AuthorId:  "sry",
		Body:      "body1",
		CreatedOn: now,
	}
	replyBodyChangedEvent := &ReplyBodyChangedEvent{
		Body: "bodyttttt",
	}
	reply.HandleReplyCreatedEvent(replyCreatedEvent)
	reply.HandleReplyBodyChangedEvent(replyBodyChangedEvent)

	assert.Equal(t, postId, reply._postId, "postId 错误")
	assert.Equal(t, parentId, reply._parentId, "parentId 错误")
	assert.Equal(t, "sry", reply._authorId, "authorId 错误")
	assert.Equal(t, "bodyttttt", reply._body, "body 错误")
	assert.Equal(t, now, reply._createdOn, "createdOn 错误")
}

func TestCreateReplyCommand(t *testing.T) {
	postId, parentId := es.NewGuid(), es.NewGuid()
	now := time.Now()
	command := &CreateReplyCommand{ParentId: parentId, PostId: postId, Body: "body1", AuthorId: "author1", CreatedOn: now}
	events := []es.Event{&ReplyCreatedEvent{ParentId: parentId, PostId: postId, Body: "body1", AuthorId: "author1", CreatedOn: now}}

	assert.Equal(t, events, new(Reply).ProcessCreateReplyCommand(command), "处理CreateReplyCommand命令返回的事件错误")
}

func TestChangeReplyBodyCommand(t *testing.T) {
	command := &ChangeReplyBodyCommand{Body: "body2"}
	events := []es.Event{&ReplyBodyChangedEvent{Body: "body2"}}

	assert.Equal(t, events, new(Reply).ProcessChangeReplyBodyCommand(command), "处理ChangeReplyBodyCommand命令返回的事件错误")
}

func TestCreateReplyCommand_Panic(t *testing.T) {
	command := &CreateReplyCommand{ParentId: es.NewGuid(), PostId: es.NewGuid(), Body: strings.Repeat("s", 4001), AuthorId: "author1"}
	assert.Panics(t, func() { new(Reply).ProcessCreateReplyCommand(command) }, "创建回复Body长度大于4000应该报错")
}

func TestChangeReplyBodyCommand_Panic(t *testing.T) {
	command := &ChangeReplyBodyCommand{Body: strings.Repeat("s", 4001)}
	assert.Panics(t, func() { new(Reply).ProcessChangeReplyBodyCommand(command) }, "修改回复Body长度大于4000应该报错")
}
