package main

import (
	es "github.com/sunrongya/eventsourcing"
)

type RPost struct {
	Id       es.Guid
	Subject  string
	Body     string
	AuthorId string
	PostReplyStatisticInfo
}

type PostProjector struct {
	repository es.ReadRepository
}

func NewPostProjector(repository es.ReadRepository) *PostProjector {
	return &PostProjector{repository: repository}
}

func (g *PostProjector) HandlePostCreatedEvent(event *PostCreatedEvent) {
	post := &RPost{
		Id:       event.GetGuid(),
		Subject:  event.Subject,
		Body:     event.Body,
		AuthorId: event.AuthorId,
	}
	g.repository.Save(post.Id, post)
}

func (g *PostProjector) HandlePostReplyStatisticInfoChangedEvent(event *PostReplyStatisticInfoChangedEvent) {
	g.do(event.GetGuid(), func(post *RPost) {
		post.PostReplyStatisticInfo = event.PostReplyStatisticInfo
	})
}

func (g *PostProjector) do(id es.Guid, assignRPostFn func(*RPost)) {
	i, err := g.repository.Find(id)
	if err != nil {
		return
	}
	post := i.(*RPost)
	assignRPostFn(post)
	g.repository.Save(id, post)
}
