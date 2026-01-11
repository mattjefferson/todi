package todi

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
)

// ListSections fetches a page of sections.
func (c *Client) ListSections(ctx context.Context, params map[string]string) ([]Section, string, error) {
	var resp listResponse[Section]
	if err := c.get(ctx, "/api/v1/sections", params, &resp); err != nil {
		return nil, "", err
	}
	return resp.Results, resp.NextCursor, nil
}

// ListSectionsAll fetches all sections across pages.
func (c *Client) ListSectionsAll(ctx context.Context, params map[string]string) ([]Section, error) {
	if params == nil {
		params = map[string]string{}
	}
	params["limit"] = strconv.Itoa(200)
	var all []Section
	cursor := ""
	for {
		if cursor != "" {
			params["cursor"] = cursor
		}
		page, next, err := c.ListSections(ctx, params)
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

// CreateSection creates a new section.
func (c *Client) CreateSection(ctx context.Context, body map[string]any) (Section, []byte, error) {
	var section Section
	raw, err := c.post(ctx, "/api/v1/sections", body, &section)
	return section, raw, err
}

// UpdateSection updates an existing section.
func (c *Client) UpdateSection(ctx context.Context, id string, body map[string]any) (Section, []byte, error) {
	var section Section
	raw, err := c.post(ctx, "/api/v1/sections/"+url.PathEscape(id), body, &section)
	return section, raw, err
}

// GetSection fetches a section by ID.
func (c *Client) GetSection(ctx context.Context, id string) (Section, error) {
	var section Section
	if err := c.get(ctx, "/api/v1/sections/"+url.PathEscape(id), nil, &section); err != nil {
		return Section{}, err
	}
	return section, nil
}

// DeleteSection deletes a section by ID.
func (c *Client) DeleteSection(ctx context.Context, id string) ([]byte, error) {
	return c.delete(ctx, "/api/v1/sections/"+url.PathEscape(id))
}

// FindSectionByName returns the section for a unique name match.
func (c *Client) FindSectionByName(ctx context.Context, name, projectID string) (Section, error) {
	params := map[string]string{}
	if projectID != "" {
		params["project_id"] = projectID
	}
	sections, err := c.ListSectionsAll(ctx, params)
	if err != nil {
		return Section{}, err
	}
	matches := make([]Section, 0, 2)
	for _, section := range sections {
		if section.Name == name {
			matches = append(matches, section)
		}
	}
	if len(matches) == 0 {
		return Section{}, fmt.Errorf("section not found: %s", name)
	}
	if len(matches) > 1 {
		return Section{}, fmt.Errorf("section name not unique: %s", name)
	}
	return matches[0], nil
}

// FindSectionIDByName returns the section ID for a unique name match.
func (c *Client) FindSectionIDByName(ctx context.Context, name, projectID string) (string, error) {
	section, err := c.FindSectionByName(ctx, name, projectID)
	if err != nil {
		return "", err
	}
	return section.ID, nil
}
