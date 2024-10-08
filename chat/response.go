package chat

import (
	"fmt"

	"github.com/cutestAlpaca/go-coze/common/response"
)

type BaseResponse struct {
	// The ID of the code.
	// 0 represents a successful call.
	// 状态码。
	// 0 代表调用成功。
	Code int `json:"code"`
	// 状态信息。API 调用失败时可通过此字段查看详细错误信息。
	Msg string `json:"msg"`
}

type DataResponse struct {
	BaseResponse
	Data *response.Chat `json:"data"`
}

type StreamingResponse struct {
	BaseResponse
	Event   string
	Chat    *response.Chat
	Message *response.Message
}

type MessageResponse struct {
	BaseResponse
	Data []response.Message `json:"data"`
}

type HttpErrorResponse struct {
	Status     string // e.g. "200 OK"
	StatusCode int    `json:"status_code"` // http 状态码
	Body       []byte `json:"body"`        // http 响应体
}

func (h *HttpErrorResponse) Error() string {
	return fmt.Sprintf("response error: statusCode: %d, status: %s", h.StatusCode, h.Status)
}