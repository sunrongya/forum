package main

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	es "github.com/sunrongya/eventsourcing"
	"strings"
	"testing"
	"time"
)

func TestPostRestore(t *testing.T) {
	post := &Post{}
	post.ApplyEvents([]es.Event{
		&PostCreatedEvent{
			Subject:  "subject1",
			Body:     "bodysssss",
			AuthorId: "sry",
		},
		&PostUpdatedEvent{
			Subject: "subject2",
			Body:    "bodyttttt",
		},
	})
	assert.Equal(t, 2, post.Version(), "version error")
	assert.Equal(t, "subject2", post._subject, "subject 错误")
	assert.Equal(t, "bodyttttt", post._body, "body 错误")
	assert.Equal(t, "sry", post._authorId, "authorId 错误")
}

func TestPostRestoreForErrorEvent(t *testing.T) {
	assert.Panics(t, func() {
		NewPost().ApplyEvents([]es.Event{&struct{ es.WithGuid }{}})
	}, "restore error event must panic error")
}

func TestCheckPostApplyEvents(t *testing.T) {
	events := []es.Event{
		&PostCreatedEvent{},
		&PostUpdatedEvent{},
		&PostReplyStatisticInfoChangedEvent{},
		&RepeatPostReplyChangedEvent{},
	}
	assert.NotPanics(t, func() { NewPost().ApplyEvents(events) }, "Check Process All Event")
}

func TestPostCommand(t *testing.T) {
	guid, replyId, newReplyId := es.NewGuid(), es.NewGuid(), es.NewGuid()
	now := time.Now()
	tests := []struct {
		post    *Post
		command es.Command
		event   es.Event
	}{
		{
			&Post{},
			&CreatePostCommand{WithGuid: es.WithGuid{Guid: guid}, Subject: "subject1", Body: "body1", AuthorId: "author1"},
			&PostCreatedEvent{WithGuid: es.WithGuid{Guid: guid}, Subject: "subject1", Body: "body1", AuthorId: "author1"},
		},
		{
			&Post{},
			&UpdatePostCommand{WithGuid: es.WithGuid{Guid: guid}, Subject: "subject2", Body: "body2"},
			&PostUpdatedEvent{WithGuid: es.WithGuid{Guid: guid}, Subject: "subject2", Body: "body2"},
		},
		{
			&Post{_replyIds: map[es.Guid]bool{}},
			&AcceptNewReplyCommand{
				WithGuid:  es.WithGuid{Guid: guid},
				ReplyId:   replyId,
				AuthorId:  "author1",
				CreatedOn: now,
			},
			&PostReplyStatisticInfoChangedEvent{
				WithGuid: es.WithGuid{Guid: guid},
				ReplyId:  replyId,
				PostReplyStatisticInfo: PostReplyStatisticInfo{
					LastReplyId:       replyId,
					LastReplyAuthorId: "author1",
					LastReplyTime:     now,
					ReplyCount:        1,
				},
			},
		},
		{
			&Post{
				_replyIds: map[es.Guid]bool{},
				_replyStatisticInfo: PostReplyStatisticInfo{
					LastReplyId:       replyId,
					LastReplyAuthorId: "author2",
					LastReplyTime:     now,
					ReplyCount:        3,
				},
			},
			&AcceptNewReplyCommand{
				WithGuid:  es.WithGuid{Guid: guid},
				ReplyId:   newReplyId,
				AuthorId:  "author1",
				CreatedOn: now.Add(1 * time.Second),
			},
			&PostReplyStatisticInfoChangedEvent{
				WithGuid: es.WithGuid{Guid: guid},
				ReplyId:  newReplyId,
				PostReplyStatisticInfo: PostReplyStatisticInfo{
					LastReplyId:       newReplyId,
					LastReplyAuthorId: "author1",
					LastReplyTime:     now.Add(1 * time.Second),
					ReplyCount:        4,
				},
			},
		},
		{
			&Post{
				_replyIds: map[es.Guid]bool{},
				_replyStatisticInfo: PostReplyStatisticInfo{
					LastReplyId:       replyId,
					LastReplyAuthorId: "author2",
					LastReplyTime:     now,
					ReplyCount:        3,
				},
			},
			&AcceptNewReplyCommand{
				WithGuid:  es.WithGuid{Guid: guid},
				ReplyId:   newReplyId,
				AuthorId:  "author1",
				CreatedOn: now.Add(-1 * time.Second),
			},
			&PostReplyStatisticInfoChangedEvent{
				WithGuid: es.WithGuid{Guid: guid},
				ReplyId:  newReplyId,
				PostReplyStatisticInfo: PostReplyStatisticInfo{
					LastReplyId:       replyId,
					LastReplyAuthorId: "author2",
					LastReplyTime:     now,
					ReplyCount:        4,
				},
			},
		},
		{
			&Post{_replyIds: map[es.Guid]bool{replyId: true}},
			&AcceptNewReplyCommand{WithGuid: es.WithGuid{Guid: guid}, ReplyId: replyId, AuthorId: "author1", CreatedOn: now},
			&RepeatPostReplyChangedEvent{WithGuid: es.WithGuid{Guid: guid}, ReplyId: replyId},
		},
	}

	for _, v := range tests {
		assert.Equal(t, []es.Event{v.event}, v.post.ProcessCommand(v.command))
	}
}

func TestPostCommand_Panic(t *testing.T) {
	tests := []struct {
		post    *Post
		command es.Command
	}{
		{
			&Post{},
			&struct{ es.WithGuid }{},
		},
		{
			&Post{},
			&CreatePostCommand{Subject: strings.Repeat("s", 257), Body: "body1", AuthorId: "author1"},
		},
		{
			&Post{},
			&CreatePostCommand{Subject: "subject1", Body: strings.Repeat("s", 4001), AuthorId: "author1"},
		},
		{
			&Post{},
			&UpdatePostCommand{Subject: strings.Repeat("s", 257), Body: "body1"},
		},
		{
			&Post{},
			&UpdatePostCommand{Subject: "subject1", Body: strings.Repeat("s", 4001)},
		},
	}

	for _, v := range tests {
		assert.Panics(t, func() { v.post.ProcessCommand(v.command) }, fmt.Sprintf("test panics error: command:%v", v.command))
	}
}
