package sdk

import (
	"context"
	"fmt"
	"time"
)

// Activity represents an event for Authorization's user. Type "reply" means Post replied to one of user's posts.
type Activity struct {
	Type    string         `json:"type"`
	Id      ActivityId     `json:"id"`
	Instant time.Time      `json:"instant"`
	IsRead  bool           `json:"isRead"`
	Post    *CommunityPost `json:"post"`
}

// ListActivity returns page of Authorization's activity starting from cursor, or the first page when cursor is nil.
func (c *Client) ListActivity(ctx context.Context, auth *Authorization, cursor *CursorId) (*Cursor[Activity], error) {
	var resp Cursor[Activity]
	err := c.do(ctx, auth, "GET", cursorPath("/activity/list", cursor), nil, &resp)
	if err != nil {
		return nil, fmt.Errorf("failed to list activity: %w", err)
	}

	return &resp, nil
}

// ReadActivity marks Authorization's activity as read.
func (c *Client) ReadActivity(ctx context.Context, auth *Authorization, id ActivityId) error {
	err := c.do(ctx, auth, "POST", fmt.Sprintf("/activity/read/%d", id.value), nil, nil)
	if err != nil {
		return fmt.Errorf("failed to read activity: %w", err)
	}

	return nil
}
