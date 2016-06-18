package main

import (
	"github.com/stretchr/testify/assert"
	es "github.com/sunrongya/eventsourcing"
	"strings"
	"testing"
	"time"
)

func TestPostRestore(t *testing.T) {
	post := &Post{}
	createdEvent := &PostCreatedEvent{Subject: "subject1", Body: "bodysssss", AuthorId: "sry"}
	updatedEvent := &PostUpdatedEvent{Subject: "subject2", Body: "bodyttttt"}
	post.HandlePostCreatedEvent(createdEvent)
	post.HandlePostUpdatedEvent(updatedEvent)

	assert.Equal(t, "subject2", post._subject, "subject 错误")
	assert.Equal(t, "bodyttttt", post._body, "body 错误")
	assert.Equal(t, "sry", post._authorId, "authorId 错误")
}

func TestCreatePostCommand(t *testing.T) {
	command := &CreatePostCommand{Subject: "subject1", Body: "body1", AuthorId: "author1"}
	events := []es.Event{&PostCreatedEvent{Subject: "subject1", Body: "body1", AuthorId: "author1"}}

	assert.Equal(t, events, new(Post).ProcessCreatePostCommand(command), "处理CreatePostCommand命令返回的事件错误")
}

func TestUpdatePostCommand(t *testing.T) {
	command := &UpdatePostCommand{Subject: "subject2", Body: "body2"}
	events := []es.Event{&PostUpdatedEvent{Subject: "subject2", Body: "body2"}}

	assert.Equal(t, events, new(Post).ProcessUpdatePostCommand(command), "处理UpdatePostCommand命令返回的事件错误")
}

func TestAcceptNewReplyCommandOfFirst(t *testing.T) {
	replyId, now := es.NewGuid(), time.Now()
	post := &Post{_replyIds: map[es.Guid]bool{}}
	command := &AcceptNewReplyCommand{ReplyId: replyId, AuthorId: "author1", CreatedOn: now}
	events := []es.Event{
		&PostReplyStatisticInfoChangedEvent{
			ReplyId: replyId,
			PostReplyStatisticInfo: PostReplyStatisticInfo{
				LastReplyId:       replyId,
				LastReplyAuthorId: "author1",
				LastReplyTime:     now,
				ReplyCount:        1,
			},
		},
	}

	assert.Equal(t, events, post.ProcessAcceptNewReplyCommand(command), "处理AcceptNewReplyCommand命令返回的事件错误")
}

func TestAcceptNewReplyCommandOfLast(t *testing.T) {
	replyId, newReplyId, now := es.NewGuid(), es.NewGuid(), time.Now()
	post := &Post{
		_replyIds: map[es.Guid]bool{},
		_replyStatisticInfo: PostReplyStatisticInfo{
			LastReplyId:       replyId,
			LastReplyAuthorId: "author2",
			LastReplyTime:     now,
			ReplyCount:        3,
		},
	}
	command := &AcceptNewReplyCommand{
		ReplyId:   newReplyId,
		AuthorId:  "author1",
		CreatedOn: now.Add(1 * time.Second),
	}
	events := []es.Event{
		&PostReplyStatisticInfoChangedEvent{
			ReplyId: newReplyId,
			PostReplyStatisticInfo: PostReplyStatisticInfo{
				LastReplyId:       newReplyId,
				LastReplyAuthorId: "author1",
				LastReplyTime:     now.Add(1 * time.Second),
				ReplyCount:        4,
			},
		},
	}

	assert.Equal(t, events, post.ProcessAcceptNewReplyCommand(command), "处理AcceptNewReplyCommand命令返回的事件错误")
}

func TestAcceptNewReplyCommandOfBefore(t *testing.T) {
	replyId, newReplyId, now := es.NewGuid(), es.NewGuid(), time.Now()
	post := &Post{
		_replyIds: map[es.Guid]bool{},
		_replyStatisticInfo: PostReplyStatisticInfo{
			LastReplyId:       replyId,
			LastReplyAuthorId: "author2",
			LastReplyTime:     now,
			ReplyCount:        3,
		},
	}
	command := &AcceptNewReplyCommand{
		ReplyId:   newReplyId,
		AuthorId:  "author1",
		CreatedOn: now.Add(-1 * time.Second),
	}
	events := []es.Event{
		&PostReplyStatisticInfoChangedEvent{
			ReplyId: newReplyId,
			PostReplyStatisticInfo: PostReplyStatisticInfo{
				LastReplyId:       replyId,
				LastReplyAuthorId: "author2",
				LastReplyTime:     now,
				ReplyCount:        4,
			},
		},
	}

	assert.Equal(t, events, post.ProcessAcceptNewReplyCommand(command), "处理AcceptNewReplyCommand命令返回的事件错误")
}

func TestAcceptNewReplyCommandOfRepeatPostReply(t *testing.T) {
	replyId, now := es.NewGuid(), time.Now()
	post := &Post{_replyIds: map[es.Guid]bool{replyId: true}}
	command := &AcceptNewReplyCommand{ReplyId: replyId, AuthorId: "author1", CreatedOn: now}
	events := []es.Event{&RepeatPostReplyChangedEvent{ReplyId: replyId}}

	assert.Equal(t, events, post.ProcessAcceptNewReplyCommand(command), "处理AcceptNewReplyCommand命令返回的事件错误")
}

func TestCreatePostCommandOfSubject_Panic(t *testing.T) {
	command := &CreatePostCommand{Subject: strings.Repeat("s", 257), Body: "body1", AuthorId: "author1"}
	assert.Panics(t, func() { new(Post).ProcessCreatePostCommand(command) }, "帖子主题长度大于256应该报错")
}

func TestCreatePostCommandOfBody_Panic(t *testing.T) {
	command := &CreatePostCommand{Subject: "subject1", Body: strings.Repeat("s", 4001), AuthorId: "author1"}
	assert.Panics(t, func() { new(Post).ProcessCreatePostCommand(command) }, "帖子内容长度大于4000应该报错")
}

func TestUpdatePostCommandOfSubject_Panic(t *testing.T) {
	command := &UpdatePostCommand{Subject: strings.Repeat("s", 257), Body: "body1"}
	assert.Panics(t, func() { new(Post).ProcessUpdatePostCommand(command) }, "更改帖子主题长度大于256应该报错")
}

func TestUpdatePostCommandOfBody_Panic(t *testing.T) {
	command := &UpdatePostCommand{Subject: "subject1", Body: strings.Repeat("s", 4001)}
	assert.Panics(t, func() { new(Post).ProcessUpdatePostCommand(command) }, "更改帖子内容长度大于4000应该报错")
}
