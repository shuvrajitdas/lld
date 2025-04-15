package stackoverflow

import (
	"fmt"
	"sync"
	"time"
)

type Answer struct {
	ID           string
	Content      string
	Author       *User
	Question     *Question
	isAccepted   bool
	CreationDate time.Time
	comments     []*Comment
	votes        []*Vote
	mu           sync.RWMutex
}

func NewAnswer(author *User, question *Question, content string) *Answer {
	return &Answer{
		ID:           generateID(),
		Author:       author,
		Question:     question,
		Content:      content,
		CreationDate: time.Now(),
		comments:     make([]*Comment, 0),
		votes:        make([]*Vote, 0),
	}
}

func (a *Answer) Vote(user *User, value int) error {
	a.mu.Lock()
	defer a.mu.Unlock()

	if value != 1 && value != -1 {
		return fmt.Errorf("vote value must be either 1 or -1")
	}

	for i, v := range a.votes {
		if v.User.ID == user.ID {
			a.votes = append(a.votes[:i], a.votes[i+1:]...)
			break
		}
	}

	a.votes = append(a.votes, &Vote{User: user, Value: value})
	a.Author.UpdateReputation(value * 10)
	return nil
}

func (a *Answer) GetVoteCount() int {
	a.mu.Lock()
	defer a.mu.Unlock()

	count := 0
	for _, vote := range a.votes {
		count += vote.Value
	}
	return count
}

func (a *Answer) AddComment(comment *Comment) error {
	a.mu.Lock()
	defer a.mu.Unlock()

	a.comments = append(a.comments, comment)
	return nil
}

func (a *Answer) GetComments() []*Comment {
	a.mu.Lock()
	defer a.mu.Unlock()

	return a.comments
}

func (a *Answer) MarkAsAccepted() error {
	a.mu.Lock()
	defer a.mu.Unlock()

	if a.isAccepted {
		return fmt.Errorf("answer is already accepted")
	}

	a.isAccepted = true
	a.Author.UpdateReputation(15)
	return nil
}

func (a *Answer) IsAccepted() bool {
	a.mu.Lock()
	defer a.mu.Unlock()

	return a.isAccepted
}
