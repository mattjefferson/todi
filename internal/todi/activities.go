package todi

import (
	"context"
	"strconv"
)

// ListActivities fetches a page of activity log entries.
func (c *Client) ListActivities(ctx context.Context, params map[string]string) ([]Activity, string, error) {
	var resp listResponse[Activity]
	if err := c.get(ctx, "/api/v1/activities", params, &resp); err != nil {
		return nil, "", err
	}
	return resp.Results, resp.NextCursor, nil
}

// ListActivitiesAll fetches all activity log entries across pages.
func (c *Client) ListActivitiesAll(ctx context.Context, params map[string]string) ([]Activity, error) {
	if params == nil {
		params = map[string]string{}
	}
	params["limit"] = strconv.Itoa(100)
	var all []Activity
	cursor := ""
	for {
		if cursor != "" {
			params["cursor"] = cursor
		}
		page, next, err := c.ListActivities(ctx, params)
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
