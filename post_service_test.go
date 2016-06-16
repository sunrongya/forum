package main

import (
	es "github.com/sunrongya/eventsourcing"
	"github.com/sunrongya/eventsourcing/utiltest"
	"testing"
)

func TestPostServiceDoCreatePost(t *testing.T) {
	utiltest.TestServicePublishCommand(t, func(service es.Service) es.Command {
		gs := PostService{Service: service}
		subject, body, authorId := "title1", "body111", "author1"
		guid := gs.CreatePost(subject, body, authorId)
		return &CreatePostCommand{WithGuid: es.WithGuid{guid}, Subject: subject, Body: body, AuthorId: authorId}
	})
}

func TestPostServiceDoUpdatePost(t *testing.T) {
	utiltest.TestServicePublishCommand(t, func(service es.Service) es.Command {
		gs := PostService{Service: service}
		guid := es.NewGuid()
		subject, body := "title1", "body111"
		gs.UpdatePost(guid, subject, body)
		return &UpdatePostCommand{WithGuid: es.WithGuid{guid}, Subject: subject, Body: body}
	})
}
