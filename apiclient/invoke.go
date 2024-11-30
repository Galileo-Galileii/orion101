package apiclient

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/orion101-ai/orion101/apiclient/types"
)

type InvokeOptions struct {
	ThreadID       string
	WorkflowStepID string
	Async          bool
}

func (c *Client) PromptResponse(ctx context.Context, resp types.PromptResponse) error {
	_, _, err := c.postJSON(ctx, "/prompt", resp)
	return err
}

func (c *Client) Invoke(ctx context.Context, agentID string, input string, opts InvokeOptions) (*types.InvokeResponse, error) {
	url := fmt.Sprintf("/invoke/%s?async=%v&step=%s", agentID, opts.Async, opts.WorkflowStepID)
	if opts.ThreadID != "" {
		url = fmt.Sprintf("/invoke/%s/threads/%s?async=%v&step=%s", agentID, opts.ThreadID, opts.Async, opts.WorkflowStepID)
	}

	_, resp, err := c.doRequest(ctx, http.MethodPost, url, bytes.NewBuffer([]byte(input)), "Accept", "text/event-stream")
	if err != nil {
		return nil, err
	}

	if opts.Async {
		_, _ = io.Copy(io.Discard, resp.Body)
		resp.Body.Close()
		return &types.InvokeResponse{
			ThreadID: resp.Header.Get("X-Orion101-Thread-Id"),
		}, nil
	}

	return &types.InvokeResponse{
		Events:   toStream[types.Progress](resp),
		ThreadID: resp.Header.Get("X-Orion101-Thread-Id"),
	}, nil
}
