package coze

import (
	"bytes"
	"context"
	"fmt"
	jsoniter "github.com/json-iterator/go"
	"io"
	"net/http"
	"time"
)

/*
*
coze官方文档地址
https://www.coze.cn/docs/developer_guides/workflow_run
*/
const (
	urlWorkflowRun = "https://api.coze.cn/v1/workflow/run"
)

type Workflow struct {
	PersonalAccessToken string        `json:"-"`
	Timeout             time.Duration `json:"-"`

	// 必选 执行的 Workflow ID，此工作流应已发布。
	WorkflowId string `json:"workflow_id"`

	// 工作流开始节点的输入参数及取值，你可以在指定工作流的编排页面查看参数列表。
	Parameters map[string]any `json:"parameters"`

	// 需要关联的 Bot ID。 部分工作流执行时需要指定关联的 Bot，例如存在数据库节点、变量节点等节点的工作流。
	BotID string `json:"bot_id"`

	// 用于指定一些额外的字段，以 Map[String][String] 格式传入。例如某些插件 会隐式用到的经纬度等字段。
	ext ext `json:"ext"`

	// 使用启用流式返回。
	Stream bool `json:"stream"`
}

type ext struct {
	Latitude  string `json:"latitude"`
	Longitude string `json:"longitude"`
	UserId    int    `json:"user_id"`
}

func NewWorkflow(workflow *Workflow) *Workflow {
	return &Workflow{
		PersonalAccessToken: workflow.PersonalAccessToken,
		WorkflowId:          workflow.WorkflowId,
		Parameters:          workflow.Parameters,
		BotID:               workflow.BotID,
		ext:                 workflow.ext,
	}
}

func (c *Workflow) Request(ctx context.Context) (*WorkFlowRunResponse, error) {
	if c.Stream {
		return nil, fmt.Errorf("stream request not supported")
	}
	body, err := jsoniter.Marshal(c)
	if err != nil {
		return nil, err
	}
	resp := new(WorkFlowRunResponse)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, urlWorkflowRun, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.Header.Add(HeaderContentType, HeaderApplicationJson)
	req.Header.Add(HeaderAuthorization, fmt.Sprintf("Bearer %s", c.PersonalAccessToken))
	req.Header.Add(HeaderConnection, HeaderKeepAlive)
	req.Header.Add(HeaderAccept, HeaderAcceptAll)
	client := http.DefaultClient
	if c.Timeout != 0 {
		client.Timeout = c.Timeout
	}

	httpResp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer httpResp.Body.Close()
	data, err := io.ReadAll(httpResp.Body)
	if err != nil {
		return resp, err
	}

	if httpResp.StatusCode != http.StatusOK {
		return nil, &HttpErrorResponse{
			Status:     httpResp.Status,
			StatusCode: httpResp.StatusCode,
			Body:       data,
		}
	}

	if err = jsoniter.Unmarshal(data, &resp); err != nil {
		return resp, err
	}

	return resp, nil
}