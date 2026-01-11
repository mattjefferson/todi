package todi

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
)

// ListTasks fetches a page of tasks.
func (c *Client) ListTasks(ctx context.Context, params map[string]string) ([]Task, string, error) {
	var resp listResponse[Task]
	if err := c.get(ctx, "/api/v1/tasks", params, &resp); err != nil {
		return nil, "", err
	}
	return resp.Results, resp.NextCursor, nil
}

// ListTasksAll fetches all tasks across pages.
func (c *Client) ListTasksAll(ctx context.Context, params map[string]string) ([]Task, error) {
	if params == nil {
		params = map[string]string{}
	}
	params["limit"] = strconv.Itoa(200)
	var all []Task
	cursor := ""
	for {
		if cursor != "" {
			params["cursor"] = cursor
		}
		page, next, err := c.ListTasks(ctx, params)
		if err != nil {
			return nil, err
		}
		all = append(all, page...)
		if next == "" {
			break
		}
		cursor = next
	}
	return all, nil
}

// CreateTask creates a new task.
func (c *Client) CreateTask(ctx context.Context, body map[string]any) (Task, []byte, error) {
	var task Task
	raw, err := c.post(ctx, "/api/v1/tasks", body, &task)
	return task, raw, err
}

// UpdateTask updates an existing task.
func (c *Client) UpdateTask(ctx context.Context, id string, body map[string]any) (Task, []byte, error) {
	var task Task
	raw, err := c.post(ctx, "/api/v1/tasks/"+url.PathEscape(id), body, &task)
	return task, raw, err
}

// GetTask fetches a task by ID.
func (c *Client) GetTask(ctx context.Context, id string) (Task, error) {
	var task Task
	if err := c.get(ctx, "/api/v1/tasks/"+url.PathEscape(id), nil, &task); err != nil {
		return Task{}, err
	}
	return task, nil
}

// DeleteTask deletes a task by ID.
func (c *Client) DeleteTask(ctx context.Context, id string) ([]byte, error) {
	return c.delete(ctx, "/api/v1/tasks/"+url.PathEscape(id))
}

// CloseTask completes a task by ID.
func (c *Client) CloseTask(ctx context.Context, id string) ([]byte, error) {
	return c.post(ctx, "/api/v1/tasks/"+url.PathEscape(id)+"/close", nil, nil)
}

// ReopenTask reopens a completed task by ID.
func (c *Client) ReopenTask(ctx context.Context, id string) ([]byte, error) {
	return c.post(ctx, "/api/v1/tasks/"+url.PathEscape(id)+"/reopen", nil, nil)
}

// QuickAdd creates a task using Todoist quick-add syntax.
func (c *Client) QuickAdd(ctx context.Context, body map[string]any) (Task, []byte, error) {
	var resp struct {
		Task Task `json:"task"`
	}
	raw, err := c.post(ctx, "/api/v1/tasks/quick", body, &resp)
	return resp.Task, raw, err
}

// FindTaskByContent returns the task for a unique content match.
func (c *Client) FindTaskByContent(ctx context.Context, title string) (Task, error) {
	tasks, err := c.ListTasksAll(ctx, nil)
	if err != nil {
		return Task{}, err
	}
	matches := make([]Task, 0, 2)
	for _, task := range tasks {
		if task.Content == title {
			matches = append(matches, task)
		}
	}
	if len(matches) == 0 {
		return Task{}, fmt.Errorf("task not found: %s", title)
	}
	if len(matches) > 1 {
		return Task{}, fmt.Errorf("task title not unique: %s", title)
	}
	return matches[0], nil
}
