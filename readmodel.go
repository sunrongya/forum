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
	_repository es.ReadRepository
}

func NewPostProjector(repository es.ReadRepository) *PostProjector {
	return &PostProjector{_repository: repository}
}

func (this *PostProjector) HandlePostCreatedEvent(event *PostCreatedEvent) {
	post := &RPost{
		Id:       event.GetGuid(),
		Subject:  event.Subject,
		Body:     event.Body,
		AuthorId: event.AuthorId,
	}
	this._repository.Save(post.Id, post)
}

func (this *PostProjector) HandlePostReplyStatisticInfoChangedEvent(event *PostReplyStatisticInfoChangedEvent) {
	this.do(event.GetGuid(), func(post *RPost) {
		post.PostReplyStatisticInfo = event.PostReplyStatisticInfo
	})
}

func (this *PostProjector) do(id es.Guid, assignRPostFn func(*RPost)) {
	i, err := this._repository.Find(id)
	if err != nil {
		return
	}
	post := i.(*RPost)
	assignRPostFn(post)
	this._repository.Save(id, post)
}
