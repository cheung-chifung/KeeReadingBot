package main

import (
	"net/url"
	"sync"
)

type CommentState int

const (
	CommentStateWaiting CommentState = iota + 1
	CommentStateURLEntered
	CommentStateCommentEntered
	CommentStateCommentConfirm
)

type CommentManager struct {
	url     *url.URL
	comment string
	state   CommentState
	*sync.Mutex
}

func NewCommentManager() CommentManager {
	return CommentManager{
		state: CommentStateWaiting,
		Mutex: new(sync.Mutex),
	}
}

func (cm *CommentManager) State() CommentState {
	cm.Lock()
	defer cm.Unlock()

	return cm.state
}

func (cm *CommentManager) Reset() {
	cm.Lock()
	cm.state = CommentStateWaiting
	cm.url = nil
	cm.comment = ""
	cm.Unlock()
}

func (cm *CommentManager) URL() string {
	cm.Lock()
	defer cm.Unlock()

	return cm.url.String()
}

func (cm *CommentManager) Comment() string {
	cm.Lock()
	defer cm.Unlock()

	return cm.comment
}

func (cm *CommentManager) InputURL(url *url.URL) {
	cm.Lock()
	if cm.state == CommentStateWaiting {
		cm.url = url
		cm.state = CommentStateURLEntered
	}
	cm.Unlock()
}

func (cm *CommentManager) InputComment(input string) {
	cm.Lock()
	if cm.state == CommentStateURLEntered {
		cm.comment = input
		cm.state = CommentStateCommentEntered
	}
	cm.Unlock()
}
