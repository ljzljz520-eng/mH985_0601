package workflow

import (
	"errors"
	"trainingdesk/internal/model"
)

func NextStep(w model.Workflow) (string, error) {
	for _, step := range w.Steps {
		completed := false
		for _, done := range w.Completed {
			if step == done {
				completed = true
				break
			}
		}
		if !completed {
			return step, nil
		}
	}
	return "", errors.New("workflow is complete")
}

func IsComplete(w model.Workflow) bool {
	_, err := NextStep(w)
	return err != nil
}
