package todi

import (
	"context"
	"net/url"
	"strconv"
)

// ListComments fetches a page of comments.
func (c *Client) ListComments(ctx context.Context, params map[string]string) ([]Comment, string, error) {
	var resp listResponse[Comment]
	if err := c.get(ctx, "/api/v1/comments", params, &resp); err != nil {
		return nil, "", err
	}
	return resp.Results, resp.NextCursor, nil
}

// ListCommentsAll fetches all comments across pages.
func (c *Client) ListCommentsAll(ctx context.Context, params map[string]string) ([]Comment, error) {
	if params == nil {
		params = map[string]string{}
	}
	params["limit"] = strconv.Itoa(200)
	var all []Comment
	cursor := ""
	for {
		if cursor != "" {
			params["cursor"] = cursor
		}
		page, next, err := c.ListComments(ctx, params)
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

// CreateComment creates a new comment.
func (c *Client) CreateComment(ctx context.Context, body map[string]any) (Comment, []byte, error) {
	var comment Comment
	raw, err := c.post(ctx, "/api/v1/comments", body, &comment)
	return comment, raw, err
}

// UpdateComment updates an existing comment.
func (c *Client) UpdateComment(ctx context.Context, id string, body map[string]any) (Comment, []byte, error) {
	var comment Comment
	raw, err := c.post(ctx, "/api/v1/comments/"+url.PathEscape(id), body, &comment)
	return comment, raw, err
}

// GetComment fetches a comment by ID.
func (c *Client) GetComment(ctx context.Context, id string) (Comment, error) {
	var comment Comment
	if err := c.get(ctx, "/api/v1/comments/"+url.PathEscape(id), nil, &comment); err != nil {
		return Comment{}, err
	}
	return comment, nil
}

// DeleteComment deletes a comment by ID.
func (c *Client) DeleteComment(ctx context.Context, id string) ([]byte, error) {
	return c.delete(ctx, "/api/v1/comments/"+url.PathEscape(id))
}
