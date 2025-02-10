package pubsub

import (
	"context"
	"ozon-tesk-task/internal/transport/graph/model"
	"sync"
)

type PubSub struct {
	commentSubscriptions map[int32][]chan *model.Comment
	lock                 sync.Mutex
}

func New() *PubSub {
	return &PubSub{
		commentSubscriptions: make(map[int32][]chan *model.Comment),
		lock:                 sync.Mutex{},
	}
}

func (p *PubSub) Subscribe(ctx context.Context, postId int32) <-chan *model.Comment {
	p.lock.Lock()
	defer p.lock.Unlock()

	ch := make(chan *model.Comment, 1)
	p.commentSubscriptions[postId] = append(p.commentSubscriptions[postId], ch)

	return ch
}

func (p *PubSub) Publish(ctx context.Context, comment *model.Comment) {
	go func() {
		p.lock.Lock()
		defer p.lock.Unlock()

		if subscribers, ok := p.commentSubscriptions[comment.PostID]; ok {
			for _, ch := range subscribers {
				ch <- comment
			}
		}
	}()
}

func (p *PubSub) Unsubscribe(ctx context.Context, postId int32, ch chan *model.Comment) {
	p.lock.Lock()
	defer p.lock.Unlock()

	var newSubscribers []chan *model.Comment
	for _, sub := range p.commentSubscriptions[postId] {
		if sub != ch {
			newSubscribers = append(newSubscribers, sub)
		}
	}

	p.commentSubscriptions[postId] = newSubscribers

	close(ch)
}
