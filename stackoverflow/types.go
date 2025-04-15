package stackoverflow

import (
	"fmt"
	"sync"
	"time"
)

const (
	QuestionReputation = 5
	AnswerReputation   = 10
	CommentReputation  = 2
)

type Commentable interface {
	AddComment(comment *Comment) error
	GetComments() []*Comment
}

type Votable interface {
	Vote(user *User, value int) error
	GetVoteCount() int
}

type Vote struct {
	User  *User
	Value int
}

type User struct {
	ID         string
	Username   string
	email      string
	reputation int
	questions  []*Question
	answers    []*Answer
	comments   []*Comment
	mu         sync.RWMutex
}

func NewUser(id, username, email string) *User {
	return &User{
		ID:        id,
		Username:  username,
		email:     email,
		questions: make([]*Question, 0),
		answers:   make([]*Answer, 0),
		comments:  make([]*Comment, 0),
	}
}

func (u *User) AddQuestion(q *Question) {
	u.mu.Lock()
	defer u.mu.Unlock()

	u.questions = append(u.questions, q)
	u.UpdateReputation(QuestionReputation)
}

func (u *User) AddAnswer(a *Answer) {
	u.mu.Lock()
	defer u.mu.Unlock()

	u.answers = append(u.answers, a)
	u.UpdateReputation(AnswerReputation)
}

func (u *User) AddComment(comment *Comment) {
	u.mu.Lock()
	defer u.mu.Unlock()

	u.comments = append(u.comments, comment)
	u.UpdateReputation(CommentReputation)
}

func (u *User) GetReputation() int {
	u.mu.RLock()
	defer u.mu.RUnlock()
	return u.reputation
}

func (u *User) UpdateReputation(reputation int) {
	u.reputation += reputation
	if u.reputation < 0 {
		u.reputation = 0
	}
}

func (u *User) GetQuestions() []*Question {
	u.mu.RLock()
	defer u.mu.RUnlock()
	questions := make([]*Question, len(u.questions))
	copy(questions, u.questions)
	return questions
}

type Question struct {
	ID           string
	Title        string
	Content      string
	Author       *User
	CreationDate time.Time
	answers      []*Answer
	comments     []*Comment
	tags         []*Tag
	votes        []*Vote
	mu           sync.RWMutex
}

func NewQuestion(author *User, title, content string, tagNames []string) *Question {
	q := &Question{
		ID:           generateID(),
		Title:        title,
		Content:      content,
		Author:       author,
		CreationDate: time.Now(),
		answers:      make([]*Answer, 0),
		comments:     make([]*Comment, 0),
		tags:         make([]*Tag, 0),
		votes:        make([]*Vote, 0),
	}
	for _, tagName := range tagNames {
		q.tags = append(q.tags, NewTag(tagName))
	}
	return q
}

func (q *Question) AddAnswer(answer *Answer) error {
	q.mu.Lock()
	defer q.mu.Unlock()

	for _, a := range q.answers {
		if a.ID == answer.ID {
			return nil
		}
	}
	q.answers = append(q.answers, answer)
	return nil
}

func (q *Question) Vote(user *User, value int) error {
	q.mu.Lock()
	defer q.mu.Unlock()

	if value != 1 && value != -1 {
		return fmt.Errorf("vote value must be either 1 or -1")
	}

	for i, v := range q.votes {
		if v.User.ID == user.ID {
			q.votes = append(q.votes[:i], q.votes[i+1:]...)
			break
		}
	}

	q.votes = append(q.votes, &Vote{User: user, Value: value})
	q.Author.UpdateReputation(value * 5)
	return nil
}

func (q *Question) GetVoteCount() int {
	q.mu.Lock()
	defer q.mu.Unlock()

	count := 0
	for _, vote := range q.votes {
		count += vote.Value
	}
	return count
}

func (q *Question) AddComment(comment *Comment) error {
	q.mu.Lock()
	defer q.mu.Unlock()
	q.comments = append(q.comments, comment)
	return nil
}

func (q *Question) GetComments() []*Comment {
	q.mu.RLock()
	defer q.mu.RUnlock()
	comments := make([]*Comment, len(q.comments))
	copy(comments, q.comments)
	return comments
}

type Tag struct {
	ID   string
	Name string
}

func NewTag(name string) *Tag {
	return &Tag{
		ID:   generateID(),
		Name: name,
	}
}

type Comment struct {
	ID           string
	Content      string
	Author       *User
	CreationDate time.Time
	mu           sync.RWMutex
}

func NewComment(author *User, content string) *Comment {
	return &Comment{
		ID:           generateID(),
		Content:      content,
		Author:       author,
		CreationDate: time.Now(),
	}
}

func (q *Question) GetTags() []*Tag {
	q.mu.Lock()
	defer q.mu.Unlock()

	tags := make([]*Tag, len(q.tags))
	copy(tags, q.tags)
	return tags
}

func generateID() string {
	return fmt.Sprintf("%d", time.Now().UnixNano())
}
