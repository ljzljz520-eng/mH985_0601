package workflow

import (
	"errors"
	"fmt"

	"trainingdesk/internal/model"
	"trainingdesk/internal/store"
)

type Engine struct {
	store *store.Store
}

func New(s *store.Store) *Engine {
	return &Engine{store: s}
}

func (e *Engine) Start(recordID, name string) (model.Workflow, error) {
	if recordID == "" || name == "" {
		return model.Workflow{}, errors.New("workflow record and name are required")
	}
	w := model.Workflow{ID: fmt.Sprintf("wf-%s", recordID), RecordID: recordID, Name: name, Stage: "registration", Steps: []string{"register", "review", "confirm", "archive"}, Completed: []string{}}
	if err := e.store.PutWorkflow(w); err != nil {
		return model.Workflow{}, err
	}
	return w, nil
}

func (e *Engine) Current(id string) (model.Workflow, error) {
	if id == "" {
		return model.Workflow{}, errors.New("workflow id is required")
	}
	return e.store.GetWorkflow(id)
}

func (e *Engine) CompleteStep(id, step string) (model.Workflow, error) {
	w, err := e.Current(id)
	if err != nil {
		return model.Workflow{}, err
	}
	if step == "" {
		return model.Workflow{}, errors.New("step is required")
	}
	for _, done := range w.Completed {
		if done == step {
			return w, nil
		}
	}
	valid := false
	for _, candidate := range w.Steps {
		if candidate == step {
			valid = true
			break
		}
	}
	if !valid {
		return model.Workflow{}, errors.New("unknown workflow step")
	}
	w.Completed = append(w.Completed, step)
	w.Stage = step
	return w, e.store.PutWorkflow(w)
}
