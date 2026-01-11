package todi

import (
	"context"
	"fmt"
	"net/url"
)

// ListProjects fetches a page of projects.
func (c *Client) ListProjects(ctx context.Context, params map[string]string) ([]Project, string, error) {
	var resp listResponse[Project]
	if err := c.get(ctx, "/api/v1/projects", params, &resp); err != nil {
		return nil, "", err
	}
	return resp.Results, resp.NextCursor, nil
}

// ListProjectsAll fetches all projects across pages.
func (c *Client) ListProjectsAll(ctx context.Context) ([]Project, error) {
	params := map[string]string{"limit": "200"}
	var all []Project
	cursor := ""
	for {
		if cursor != "" {
			params["cursor"] = cursor
		}
		page, next, err := c.ListProjects(ctx, params)
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

// CreateProject creates a new project.
func (c *Client) CreateProject(ctx context.Context, body map[string]any) (Project, []byte, error) {
	var project Project
	raw, err := c.post(ctx, "/api/v1/projects", body, &project)
	return project, raw, err
}

// UpdateProject updates an existing project.
func (c *Client) UpdateProject(ctx context.Context, id string, body map[string]any) (Project, []byte, error) {
	var project Project
	raw, err := c.post(ctx, "/api/v1/projects/"+url.PathEscape(id), body, &project)
	return project, raw, err
}

// GetProject fetches a project by ID.
func (c *Client) GetProject(ctx context.Context, id string) (Project, error) {
	var project Project
	if err := c.get(ctx, "/api/v1/projects/"+url.PathEscape(id), nil, &project); err != nil {
		return Project{}, err
	}
	return project, nil
}

// DeleteProject deletes a project by ID.
func (c *Client) DeleteProject(ctx context.Context, id string) ([]byte, error) {
	return c.delete(ctx, "/api/v1/projects/"+url.PathEscape(id))
}

// ArchiveProject archives a project by ID.
func (c *Client) ArchiveProject(ctx context.Context, id string) ([]byte, error) {
	return c.post(ctx, "/api/v1/projects/"+url.PathEscape(id)+"/archive", nil, nil)
}

// UnarchiveProject unarchives a project by ID.
func (c *Client) UnarchiveProject(ctx context.Context, id string) ([]byte, error) {
	return c.post(ctx, "/api/v1/projects/"+url.PathEscape(id)+"/unarchive", nil, nil)
}

// FindProjectByName returns the project for a unique name match.
func (c *Client) FindProjectByName(ctx context.Context, name string) (Project, error) {
	projects, err := c.ListProjectsAll(ctx)
	if err != nil {
		return Project{}, err
	}
	matches := make([]Project, 0, 2)
	for _, project := range projects {
		if project.Name == name {
			matches = append(matches, project)
		}
	}
	if len(matches) == 0 {
		return Project{}, fmt.Errorf("project not found: %s", name)
	}
	if len(matches) > 1 {
		return Project{}, fmt.Errorf("project name not unique: %s", name)
	}
	return matches[0], nil
}

// FindProjectIDByName returns the project ID for a unique project name.
func (c *Client) FindProjectIDByName(ctx context.Context, name string) (string, error) {
	project, err := c.FindProjectByName(ctx, name)
	if err != nil {
		return "", err
	}
	return project.ID, nil
}
