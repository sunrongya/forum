package main

import (
	es "github.com/sunrongya/eventsourcing"
	"time"
)

type ReplyService struct {
	es.Service
}

func NewReplyService(store es.EventStore) *ReplyService {
	service := &ReplyService{
		Service: es.NewService(store, NewReply),
	}
	return service
}

func (p *ReplyService) CreateReply(postId, parentId es.Guid, body, authorId string, createdOn time.Time) es.Guid {
	guid := es.NewGuid()
	c := &CreateReplyCommand{
		WithGuid:  es.WithGuid{guid},
		PostId:    postId,
		ParentId:  parentId,
		Body:      body,
		AuthorId:  authorId,
		CreatedOn: createdOn,
	}
	p.PublishCommand(c)
	return guid
}

func (p *ReplyService) UpdateReply(guid es.Guid, body string) {
	c := &ChangeReplyBodyCommand{
		WithGuid: es.WithGuid{guid},
		Body:     body,
	}
	p.PublishCommand(c)
}
