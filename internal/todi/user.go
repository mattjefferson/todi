package todi

import "context"

// GetUserInfo fetches the authenticated user.
func (c *Client) GetUserInfo(ctx context.Context) (User, error) {
	var user User
	if err := c.get(ctx, "/api/v1/user", nil, &user); err != nil {
		return User{}, err
	}
	return user, nil
}
