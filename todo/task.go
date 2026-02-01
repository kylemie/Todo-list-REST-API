package todo

import "time"

type Task struct {
	Title       string
	Description string
	Completed   bool
	CreatedAt   time.Time
	CompletedAd *time.Time
}

func NewTask(title string, description string) Task {
	return Task{
		Title:       title,
		Description: description,
		Completed:   false,
		CreatedAt:   time.Now(),
		CompletedAd: nil,
	}
}

func (t *Task) CompleteTask() {
	completeTime := time.Now()
	t.Completed = true
	t.CompletedAd = &completeTime
}

func (t *Task) UnCompleteTask() {
	t.Completed = false
	t.CompletedAd = nil
}
