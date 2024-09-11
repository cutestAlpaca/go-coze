package coze

import (
	"context"
	"testing"
)

func TestWorkFlowRun(t *testing.T) {
	newWorkflow := &Workflow{
		PersonalAccessToken: "",
		WorkflowId:          "",
		Parameters: map[string]any{
			"main":     "读书",
			"subtitle": "读了活着第一章",
			"question": "这本书写的什么意思",
		},
	}

	workflow, err := NewWorkflow(newWorkflow).Request(context.Background())
	if err != nil {
		panic(err)
	}

	t.Log(workflow)

}