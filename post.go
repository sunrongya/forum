package main

import (
	"fmt"
	es "github.com/sunrongya/eventsourcing"
)

type Post struct {
	es.BaseAggregate
	_subject            string
	_body               string
	_authorId           string
	_replyIds           map[es.Guid]bool
	_replyStatisticInfo PostReplyStatisticInfo
}

var _ es.Aggregate = (*Post)(nil)

func NewPost() es.Aggregate {
	return &Post{
		_replyIds: make(map[es.Guid]bool),
	}
}

func (this *Post) ProcessCreatePostCommand(command *CreatePostCommand) []es.Event {
	if len(command.Subject) > 256 {
		panic(fmt.Errorf("帖子标题长度不能超过256"))
	}
	if len(command.Body) > 4000 {
		panic(fmt.Errorf("帖子内容长度不能超过4000"))
	}

	return []es.Event{
		&PostCreatedEvent{
			Subject:  command.Subject,
			Body:     command.Body,
			AuthorId: command.AuthorId,
		},
	}
}

func (this *Post) ProcessUpdatePostCommand(command *UpdatePostCommand) []es.Event {
	if len(command.Subject) > 256 {
		panic(fmt.Errorf("帖子标题长度不能超过256"))
	}
	if len(command.Body) > 4000 {
		panic(fmt.Errorf("帖子内容长度不能超过4000"))
	}

	return []es.Event{
		&PostUpdatedEvent{
			Subject: command.Subject,
			Body:    command.Body,
		},
	}
}

func (this *Post) ProcessAcceptNewReplyCommand(command *AcceptNewReplyCommand) []es.Event {
	if _, ok := this._replyIds[command.ReplyId]; ok {
		return []es.Event{&RepeatPostReplyChangedEvent{ReplyId: command.ReplyId}}
	}
	var replyStatisticInfo PostReplyStatisticInfo
	if this._replyStatisticInfo.ReplyCount == 0 {
		replyStatisticInfo = PostReplyStatisticInfo{
			LastReplyId:       command.ReplyId,
			LastReplyAuthorId: command.AuthorId,
			LastReplyTime:     command.CreatedOn,
			ReplyCount:        1,
		}
	} else if this._replyStatisticInfo.LastReplyTime.After(command.CreatedOn) {
		this._replyStatisticInfo.ReplyCount += 1
		replyStatisticInfo = this._replyStatisticInfo
	} else {
		replyStatisticInfo = PostReplyStatisticInfo{
			LastReplyId:       command.ReplyId,
			LastReplyAuthorId: command.AuthorId,
			LastReplyTime:     command.CreatedOn,
			ReplyCount:        this._replyStatisticInfo.ReplyCount + 1,
		}
	}
	return []es.Event{
		&PostReplyStatisticInfoChangedEvent{
			ReplyId:                command.ReplyId,
			PostReplyStatisticInfo: replyStatisticInfo,
		},
	}
}

func (this *Post) HandlePostCreatedEvent(event *PostCreatedEvent) {
	this._subject, this._body, this._authorId = event.Subject, event.Body, event.AuthorId
}

func (this *Post) HandlePostUpdatedEvent(event *PostUpdatedEvent) {
	this._subject, this._body = event.Subject, event.Body
}

func (this *Post) HandlePostReplyStatisticInfoChangedEvent(event *PostReplyStatisticInfoChangedEvent) {
	this._replyIds[event.ReplyId] = true
	this._replyStatisticInfo = event.PostReplyStatisticInfo
}

func (this *Post) HandleRepeatPostReplyChangedEvent(event *RepeatPostReplyChangedEvent) {
}
