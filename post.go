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

func (p *Post) ApplyEvents(events []es.Event) {
	for _, event := range events {
		switch e := event.(type) {
		case *PostCreatedEvent:
			p._subject, p._body, p._authorId = e.Subject, e.Body, e.AuthorId
		case *PostUpdatedEvent:
			p._subject, p._body = e.Subject, e.Body
		case *PostReplyStatisticInfoChangedEvent:
			p._replyIds[e.ReplyId] = true
			p._replyStatisticInfo = e.PostReplyStatisticInfo
		case *RepeatPostReplyChangedEvent:
		default:
			panic(fmt.Errorf("Unknown event %#v", e))
		}
	}
	p.SetVersion(len(events))
}

func (p *Post) ProcessCommand(command es.Command) []es.Event {
	var event es.Event
	switch c := command.(type) {
	case *CreatePostCommand:
		event = p.processCreatePostCommand(c)
	case *UpdatePostCommand:
		event = p.processUpdatePostCommand(c)
	case *AcceptNewReplyCommand:
		event = p.processAcceptNewReplyCommand(c)
	default:
		panic(fmt.Errorf("Unknown command %#v", c))
	}
	event.SetGuid(command.GetGuid())
	return []es.Event{event}
}

func (p *Post) processCreatePostCommand(command *CreatePostCommand) es.Event {
	if len(command.Subject) > 256 {
		panic(fmt.Errorf("帖子标题长度不能超过256"))
	}
	if len(command.Body) > 4000 {
		panic(fmt.Errorf("帖子内容长度不能超过4000"))
	}

	return &PostCreatedEvent{
		Subject:  command.Subject,
		Body:     command.Body,
		AuthorId: command.AuthorId,
	}
}

func (p *Post) processUpdatePostCommand(command *UpdatePostCommand) es.Event {
	if len(command.Subject) > 256 {
		panic(fmt.Errorf("帖子标题长度不能超过256"))
	}
	if len(command.Body) > 4000 {
		panic(fmt.Errorf("帖子内容长度不能超过4000"))
	}

	return &PostUpdatedEvent{
		Subject: command.Subject,
		Body:    command.Body,
	}
}

func (p *Post) processAcceptNewReplyCommand(command *AcceptNewReplyCommand) es.Event {
	if _, ok := p._replyIds[command.ReplyId]; ok {
		return &RepeatPostReplyChangedEvent{ReplyId: command.ReplyId}
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
	return &PostReplyStatisticInfoChangedEvent{
		ReplyId:                command.ReplyId,
		PostReplyStatisticInfo: replyStatisticInfo,
	}
}

func NewPost() es.Aggregate {
	return &Post{
		_replyIds: make(map[es.Guid]bool),
	}
}
