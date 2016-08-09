package main

import (
	es "github.com/sunrongya/eventsourcing"
)

type PostService struct {
	es.Service
}

func NewPostService(store es.EventStore) *PostService {
	service := &PostService{
		Service: es.NewService(store, NewPost),
	}
	return service
}

func (this *PostService) CreatePost(subject, body, authorId string) es.Guid {
	guid := es.NewGuid()
	c := &CreatePostCommand{
		WithGuid: es.WithGuid{guid},
		Subject:  subject,
		Body:     body,
		AuthorId: authorId,
	}
	this.PublishCommand(c)
	return guid
}

func (this *PostService) UpdatePost(guid es.Guid, subject, body string) {
	c := &UpdatePostCommand{
		WithGuid: es.WithGuid{guid},
		Subject:  subject,
		Body:     body,
	}
	this.PublishCommand(c)
}
