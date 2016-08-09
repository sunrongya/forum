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

func (this *ReplyService) CreateReply(postId, parentId es.Guid, body, authorId string, createdOn time.Time) es.Guid {
	guid := es.NewGuid()
	c := &CreateReplyCommand{
		WithGuid:  es.WithGuid{guid},
		PostId:    postId,
		ParentId:  parentId,
		Body:      body,
		AuthorId:  authorId,
		CreatedOn: createdOn,
	}
	this.PublishCommand(c)
	return guid
}

func (this *ReplyService) UpdateReply(guid es.Guid, body string) {
	c := &ChangeReplyBodyCommand{
		WithGuid: es.WithGuid{guid},
		Body:     body,
	}
	this.PublishCommand(c)
}
