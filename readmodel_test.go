package main

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	es "github.com/sunrongya/eventsourcing"
	"testing"
	"time"
)

func TestPostReadModel(t *testing.T) {
	readRepository := es.NewMemoryReadRepository()
	postProjector := NewPostProjector(readRepository)

	// 帖子创建
	createdEvents := []*PostCreatedEvent{
		&PostCreatedEvent{WithGuid: es.WithGuid{es.NewGuid()}, Subject: "subject1", Body: "body1", AuthorId: "author1"},
		&PostCreatedEvent{WithGuid: es.WithGuid{es.NewGuid()}, Subject: "subject2", Body: "body2", AuthorId: "author2"},
	}

	for _, event := range createdEvents {
		postProjector.HandlePostCreatedEvent(event)
	}

	// 帖子创建验证
	for _, event := range createdEvents {
		i, err := readRepository.Find(event.GetGuid())
		assert.NoError(t, err, fmt.Sprintf("读取帖子创建[%s]信息错误", event.Subject))
		post := i.(*RPost)

		assert.Equal(t, event.GetGuid(), post.Id, "ID 不相等")
		assert.Equal(t, event.Subject, post.Subject, "Subject 不相等")
		assert.Equal(t, event.Body, post.Body, "Price 不相等")
		assert.Equal(t, event.AuthorId, post.AuthorId, "Price 不相等")
	}

	// 回复
	replyEvents := []*PostReplyStatisticInfoChangedEvent{
		&PostReplyStatisticInfoChangedEvent{
			WithGuid: es.WithGuid{createdEvents[0].GetGuid()},
			PostReplyStatisticInfo: PostReplyStatisticInfo{
				LastReplyId:       es.NewGuid(),
				LastReplyAuthorId: "author1",
				LastReplyTime:     time.Now(),
				ReplyCount:        2,
			},
		},
		&PostReplyStatisticInfoChangedEvent{
			WithGuid: es.WithGuid{createdEvents[1].GetGuid()},
			PostReplyStatisticInfo: PostReplyStatisticInfo{
				LastReplyId:       es.NewGuid(),
				LastReplyAuthorId: "author2",
				LastReplyTime:     time.Now(),
				ReplyCount:        3,
			},
		},
		&PostReplyStatisticInfoChangedEvent{
			WithGuid: es.WithGuid{createdEvents[1].GetGuid()},
			PostReplyStatisticInfo: PostReplyStatisticInfo{
				LastReplyId:       es.NewGuid(),
				LastReplyAuthorId: "author4",
				LastReplyTime:     time.Now(),
				ReplyCount:        5,
			},
		},
	}

	for _, event := range replyEvents {
		postProjector.HandlePostReplyStatisticInfoChangedEvent(event)
	}

	// 验证统计信息
	postReplyStatisticInfos := []PostReplyStatisticInfo{
		replyEvents[0].PostReplyStatisticInfo,
		replyEvents[2].PostReplyStatisticInfo,
	}
	for i, event := range createdEvents {
		iPost, err := readRepository.Find(event.GetGuid())
		assert.NoError(t, err, fmt.Sprintf("读取帖子创建[%s]信息错误", event.Subject))
		post := iPost.(*RPost)

		assert.Equal(t, event.GetGuid(), post.Id, "ID 不相等")
		assert.Equal(t, event.Subject, post.Subject, "Subject 不相等")
		assert.Equal(t, event.Body, post.Body, "Price 不相等")
		assert.Equal(t, event.AuthorId, post.AuthorId, "Price 不相等")
		assert.Equal(t, postReplyStatisticInfos[i], post.PostReplyStatisticInfo, "PostReplyStatisticInfo 不相等")
	}
}
