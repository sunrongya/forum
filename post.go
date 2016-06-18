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

func (p *Post) ProcessCreatePostCommand(command *CreatePostCommand) []es.Event {
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

func (p *Post) ProcessUpdatePostCommand(command *UpdatePostCommand) []es.Event {
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

func (p *Post) ProcessAcceptNewReplyCommand(command *AcceptNewReplyCommand) []es.Event {
	if _, ok := p._replyIds[command.ReplyId]; ok {
		return []es.Event{&RepeatPostReplyChangedEvent{ReplyId: command.ReplyId}}
	}
	var replyStatisticInfo PostReplyStatisticInfo
	if p._replyStatisticInfo.ReplyCount == 0 {
		replyStatisticInfo = PostReplyStatisticInfo{
			LastReplyId:       command.ReplyId,
			LastReplyAuthorId: command.AuthorId,
			LastReplyTime:     command.CreatedOn,
			ReplyCount:        1,
		}
	} else if p._replyStatisticInfo.LastReplyTime.After(command.CreatedOn) {
		p._replyStatisticInfo.ReplyCount += 1
		replyStatisticInfo = p._replyStatisticInfo
	} else {
		replyStatisticInfo = PostReplyStatisticInfo{
			LastReplyId:       command.ReplyId,
			LastReplyAuthorId: command.AuthorId,
			LastReplyTime:     command.CreatedOn,
			ReplyCount:        p._replyStatisticInfo.ReplyCount + 1,
		}
	}
	return []es.Event{
		&PostReplyStatisticInfoChangedEvent{
			ReplyId:                command.ReplyId,
			PostReplyStatisticInfo: replyStatisticInfo,
		},
	}
}

func (p *Post) HandlePostCreatedEvent(event *PostCreatedEvent) {
	p._subject, p._body, p._authorId = event.Subject, event.Body, event.AuthorId
}

func (p *Post) HandlePostUpdatedEvent(event *PostUpdatedEvent) {
	p._subject, p._body = event.Subject, event.Body
}

func (p *Post) HandlePostReplyStatisticInfoChangedEvent(event *PostReplyStatisticInfoChangedEvent) {
	p._replyIds[event.ReplyId] = true
	p._replyStatisticInfo = event.PostReplyStatisticInfo
}

func (p *Post) HandleRepeatPostReplyChangedEvent(event *RepeatPostReplyChangedEvent) {
}
