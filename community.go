package sdk

import (
	"context"
	"fmt"
	"time"
)

// CommunityPostDescriptor identifies community post.
type CommunityPostDescriptor struct {
	Id         CommunityPostId         `json:"id"`
	AccessHash CommunityPostAccessHash `json:"accessHash"`
}

// Cursor represents single page of paginated results. NextId is nil on the last page.
type Cursor[T any] struct {
	Data   []T       `json:"data"`
	NextId *CursorId `json:"nextId"`
}

// CommunityPost represents community post. Text and Owner are nil when post is deleted.
type CommunityPost struct {
	Type          string                  `json:"type"`
	Id            CommunityPostId         `json:"id"`
	AccessHash    CommunityPostAccessHash `json:"accessHash"`
	Instant       time.Time               `json:"instant"`
	ReplyPreviews []UserDetails           `json:"replyPreviews"`
	Text          *CommunityPostText      `json:"text"`
	Owner         *UserDetails            `json:"owner"`
	Edited        bool                    `json:"edited"`
}

// Deleted reports whether post was deleted by its owner.
func (p CommunityPost) Deleted() bool {
	return p.Type == "deleted"
}

// Descriptor returns CommunityPostDescriptor of post.
func (p CommunityPost) Descriptor() CommunityPostDescriptor {
	return CommunityPostDescriptor{Id: p.Id, AccessHash: p.AccessHash}
}

// CommunityPostReply represents reply to community post: either single post (Type "single") or thread of posts (Type "thread").
type CommunityPostReply struct {
	Type   string          `json:"type"`
	Post   *CommunityPost  `json:"post"`
	Thread []CommunityPost `json:"thread"`
}

// Posts returns all posts of reply in order regardless of its type.
func (r CommunityPostReply) Posts() []CommunityPost {
	if r.Post != nil {
		return []CommunityPost{*r.Post}
	}

	return r.Thread
}

// CommunityPostDetails represents community post with its first page of replies and chain of posts it replies to.
type CommunityPostDetails struct {
	Post     CommunityPost              `json:"post"`
	Replies  Cursor[CommunityPostReply] `json:"replies"`
	Upstream []CommunityPost            `json:"upstream"`
}

type communityPostRequest struct {
	Text    CommunityPostText        `json:"text"`
	ReplyTo *CommunityPostDescriptor `json:"replyTo,omitempty"`
}

type communityEditRequest struct {
	Text editAccountValue[CommunityPostText] `json:"text"`
}

func cursorPath(path string, cursor *CursorId) string {
	if cursor == nil {
		return path
	}

	return path + "/" + cursor.value
}

// PostCommunity creates community post with provided text, as a reply when replyTo is not nil.
func (c *Client) PostCommunity(ctx context.Context, auth *Authorization, text CommunityPostText, replyTo *CommunityPostDescriptor) (*CommunityPostDescriptor, error) {
	req := communityPostRequest{
		Text:    text,
		ReplyTo: replyTo,
	}

	var resp CommunityPostDescriptor
	err := c.do(ctx, auth, "POST", "/community", req, &resp)
	if err != nil {
		return nil, fmt.Errorf("failed to post community: %w", err)
	}

	return &resp, nil
}

// ListCommunity returns page of community posts starting from cursor, or the first page when cursor is nil.
func (c *Client) ListCommunity(ctx context.Context, auth *Authorization, cursor *CursorId) (*Cursor[CommunityPost], error) {
	var resp Cursor[CommunityPost]
	err := c.do(ctx, auth, "GET", cursorPath("/community/list", cursor), nil, &resp)
	if err != nil {
		return nil, fmt.Errorf("failed to list community: %w", err)
	}

	return &resp, nil
}

// GetCommunityPost returns CommunityPostDetails for provided descriptor.
func (c *Client) GetCommunityPost(ctx context.Context, auth *Authorization, post CommunityPostDescriptor) (*CommunityPostDetails, error) {
	var resp CommunityPostDetails
	err := c.do(ctx, auth, "GET", fmt.Sprintf("/community/2/%d/%s", post.Id.value, post.AccessHash.value), nil, &resp)
	if err != nil {
		return nil, fmt.Errorf("failed to get community post: %w", err)
	}

	return &resp, nil
}

// GetCommunityReplies returns page of replies to provided post starting from cursor, or the first page when cursor is nil.
func (c *Client) GetCommunityReplies(ctx context.Context, auth *Authorization, post CommunityPostDescriptor, cursor *CursorId) (*Cursor[CommunityPostReply], error) {
	path := fmt.Sprintf("/community/%d/%s/replies2", post.Id.value, post.AccessHash.value)

	var resp Cursor[CommunityPostReply]
	err := c.do(ctx, auth, "GET", cursorPath(path, cursor), nil, &resp)
	if err != nil {
		return nil, fmt.Errorf("failed to get community replies: %w", err)
	}

	return &resp, nil
}

// EditCommunityPost replaces text of Authorization's community post.
func (c *Client) EditCommunityPost(ctx context.Context, auth *Authorization, id CommunityPostId, text CommunityPostText) error {
	req := communityEditRequest{
		Text: editAccountValue[CommunityPostText]{Value: text},
	}

	err := c.do(ctx, auth, "POST", fmt.Sprintf("/community/%d/edit", id.value), req, nil)
	if err != nil {
		return fmt.Errorf("failed to edit community post: %w", err)
	}

	return nil
}

// DeleteCommunityPost deletes Authorization's community post.
func (c *Client) DeleteCommunityPost(ctx context.Context, auth *Authorization, id CommunityPostId) error {
	err := c.do(ctx, auth, "POST", fmt.Sprintf("/community/%d/delete", id.value), nil, nil)
	if err != nil {
		return fmt.Errorf("failed to delete community post: %w", err)
	}

	return nil
}
