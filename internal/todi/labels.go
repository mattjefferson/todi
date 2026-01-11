package todi

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
)

// ListLabels fetches a page of labels.
func (c *Client) ListLabels(ctx context.Context, params map[string]string) ([]Label, string, error) {
	var resp listResponse[Label]
	if err := c.get(ctx, "/api/v1/labels", params, &resp); err != nil {
		return nil, "", err
	}
	return resp.Results, resp.NextCursor, nil
}

// ListLabelsAll fetches all labels across pages.
func (c *Client) ListLabelsAll(ctx context.Context, params map[string]string) ([]Label, error) {
	if params == nil {
		params = map[string]string{}
	}
	params["limit"] = strconv.Itoa(200)
	var all []Label
	cursor := ""
	for {
		if cursor != "" {
			params["cursor"] = cursor
		}
		page, next, err := c.ListLabels(ctx, params)
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

// CreateLabel creates a new label.
func (c *Client) CreateLabel(ctx context.Context, body map[string]any) (Label, []byte, error) {
	var label Label
	raw, err := c.post(ctx, "/api/v1/labels", body, &label)
	return label, raw, err
}

// UpdateLabel updates an existing label.
func (c *Client) UpdateLabel(ctx context.Context, id string, body map[string]any) (Label, []byte, error) {
	var label Label
	raw, err := c.post(ctx, "/api/v1/labels/"+url.PathEscape(id), body, &label)
	return label, raw, err
}

// GetLabel fetches a label by ID.
func (c *Client) GetLabel(ctx context.Context, id string) (Label, error) {
	var label Label
	if err := c.get(ctx, "/api/v1/labels/"+url.PathEscape(id), nil, &label); err != nil {
		return Label{}, err
	}
	return label, nil
}

// DeleteLabel deletes a label by ID.
func (c *Client) DeleteLabel(ctx context.Context, id string) ([]byte, error) {
	return c.delete(ctx, "/api/v1/labels/"+url.PathEscape(id))
}

// FindLabelByName returns the label for a unique name match.
func (c *Client) FindLabelByName(ctx context.Context, name string) (Label, error) {
	labels, err := c.ListLabelsAll(ctx, nil)
	if err != nil {
		return Label{}, err
	}
	matches := make([]Label, 0, 2)
	for _, label := range labels {
		if label.Name == name {
			matches = append(matches, label)
		}
	}
	if len(matches) == 0 {
		return Label{}, fmt.Errorf("label not found: %s", name)
	}
	if len(matches) > 1 {
		return Label{}, fmt.Errorf("label name not unique: %s", name)
	}
	return matches[0], nil
}

// FindLabelIDByName returns the label ID for a unique name match.
func (c *Client) FindLabelIDByName(ctx context.Context, name string) (string, error) {
	label, err := c.FindLabelByName(ctx, name)
	if err != nil {
		return "", err
	}
	return label.ID, nil
}
