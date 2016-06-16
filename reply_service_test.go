package main

import (
	es "github.com/sunrongya/eventsourcing"
	"github.com/sunrongya/eventsourcing/utiltest"
	"testing"
	"time"
)

func TestReplyServiceDoCreateReply(t *testing.T) {
	utiltest.TestServicePublishCommand(t, func(service es.Service) es.Command {
		gs := ReplyService{Service: service}
		postId, parentId := es.NewGuid(), es.NewGuid()
		body, authorId := "body111", "author1"
		now := time.Now()
		guid := gs.CreateReply(postId, parentId, body, authorId, now)
		return &CreateReplyCommand{WithGuid: es.WithGuid{guid}, PostId: postId, ParentId: parentId, AuthorId: authorId, Body: body, CreatedOn: now}
	})
}

func TestReplyServiceDoUpdateReply(t *testing.T) {
	utiltest.TestServicePublishCommand(t, func(service es.Service) es.Command {
		gs := ReplyService{Service: service}
		guid := es.NewGuid()
		body := "body111"
		gs.UpdateReply(guid, body)
		return &ChangeReplyBodyCommand{WithGuid: es.WithGuid{guid}, Body: body}
	})
}
