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

type newUnreadActivityResponse struct {
	Type       string     `json:"type"`
	ActivityId ActivityId `json:"activityId"`
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

// GetNewUnreadActivity returns ID of Authorization's newest unacknowledged unread activity, or nil when there is none.
func (c *Client) GetNewUnreadActivity(ctx context.Context, auth *Authorization) (*ActivityId, error) {
	var resp newUnreadActivityResponse
	err := c.do(ctx, auth, "GET", "/activity/unreads/new", nil, &resp)
	if err != nil {
		return nil, fmt.Errorf("failed to get new unread activity: %w", err)
	}

	if resp.Type != "some" {
		return nil, nil
	}

	return &resp.ActivityId, nil
}

// AcknowledgeUnreadActivity marks Authorization's new unread activity as seen without reading it.
func (c *Client) AcknowledgeUnreadActivity(ctx context.Context, auth *Authorization) error {
	err := c.do(ctx, auth, "POST", "/activity/unreads/acknowledge", nil, nil)
	if err != nil {
		return fmt.Errorf("failed to acknowledge unread activity: %w", err)
	}

	return nil
}
