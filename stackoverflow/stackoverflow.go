package stackoverflow

import (
	"fmt"
	"strings"
	"sync"
)

type StackOverflow struct {
	users     map[string]*User
	questions map[string]*Question
	answers   map[string]*Answer
	tags      map[string]*Tag
	mu        sync.RWMutex
}

func NewStackOverflow() *StackOverflow {
	return &StackOverflow{
		users:     make(map[string]*User),
		questions: make(map[string]*Question),
		answers:   make(map[string]*Answer),
		tags:      make(map[string]*Tag),
	}
}

func (so *StackOverflow) CreateUser(username, email string) *User {
	so.mu.Lock()
	defer so.mu.Unlock()

	id := generateID()
	user := NewUser(id, username, email)

	so.users[id] = user
	return user
}

func (so *StackOverflow) AskQuestion(user *User, title, content string, tags []string) (*Question, error) {
	if user == nil {
		return nil, fmt.Errorf("user cannot be nil")
	}

	question := NewQuestion(user, title, content, tags)

	so.mu.Lock()
	defer so.mu.Unlock()

	so.questions[question.ID] = question
	user.AddQuestion(question)

	for _, tag := range question.GetTags() {
		so.tags[tag.Name] = tag
	}
	return question, nil
}

func (so *StackOverflow) AnswerQuestion(user *User, question *Question, content string) (*Answer, error) {
	if user == nil || question == nil {
		return nil, fmt.Errorf("user and question cannot be nil")
	}

	answer := NewAnswer(user, question, content)

	so.mu.Lock()
	defer so.mu.Unlock()

	so.answers[answer.ID] = answer
	question.AddAnswer(answer)
	user.AddAnswer(answer)

	return answer, nil
}

func (so *StackOverflow) AddComment(user *User, target Commentable, content string) (*Comment, error) {
	if user == nil || target == nil {
		return nil, fmt.Errorf("user and target cannot be nil")
	}

	comment := NewComment(user, content)
	if err := target.AddComment(comment); err != nil {
		return nil, err
	}

	user.AddComment(comment)
	return comment, nil
}

func (so *StackOverflow) Vote(user *User, target Votable, value int) error {
	if user == nil || target == nil {
		return fmt.Errorf("user and target cannot be nil")
	}

	return target.Vote(user, value)
}

func (so *StackOverflow) AcceptAnswer(answer *Answer) error {
	if answer == nil {
		return fmt.Errorf("answer cannot be nil")
	}

	return answer.MarkAsAccepted()
}

func (so *StackOverflow) SearchQuestions(query string) []*Question {
	so.mu.Lock()
	defer so.mu.Unlock()

	query = strings.ToLower(query)
	var results []*Question

	for _, q := range so.questions {
		if strings.Contains(strings.ToLower(q.Content), query) || strings.Contains(strings.ToLower(q.Title), query) {
			results = append(results, q)
			continue
		}

		for _, tag := range q.tags {
			if strings.EqualFold(tag.Name, query) {
				results = append(results, q)
				break
			}
		}
	}
	return results
}

func (so *StackOverflow) GetQuestionsByUser(user *User) []*Question {
	if user == nil {
		return nil
	}
	return user.GetQuestions()
}
