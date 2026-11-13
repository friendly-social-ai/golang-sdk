package sdk

import (
	"context"
	"fmt"
	"time"
)

type addFriendRequest struct {
	Token *FriendToken `json:"token"`
}

type friendRequestRequest struct {
	UserId     UserId          `json:"userId"`
	AccessHash *UserAccessHash `json:"userAccessHash"`
}

// GeneratedFriendToken is a FriendToken with the time the backend stops accepting it.
type GeneratedFriendToken struct {
	Token      FriendToken `json:"token"`
	Expiration time.Time   `json:"expiration"`
}

type addFriendResponse struct {
	Type string `json:"type"`
}

// FriendTokenOwner is the public profile of the user who owns a FriendToken: Nickname, Description, list of Interests
// and Avatar, with the Expiration of the FriendToken.
type FriendTokenOwner struct {
	Nickname    Nickname        `json:"nickname"`
	Description UserDescription `json:"description"`
	Interests   Interests       `json:"interests"`
	Avatar      *FileDescriptor `json:"avatar"`
	Expiration  time.Time       `json:"expiration"`
}

type friendTokenOwnerResponse struct {
	Type string `json:"type"`
	FriendTokenOwner
}

// GenerateFriendToken creates token for Authorization's user by which another users can add them.
func (c *Client) GenerateFriendToken(ctx context.Context, auth *Authorization) (GeneratedFriendToken, error) {
	var resp GeneratedFriendToken
	err := c.do(ctx, auth, "POST", "/friends/generate", nil, &resp)
	if err != nil {
		return GeneratedFriendToken{}, fmt.Errorf("failed to generate friend token: %w", err)
	}

	return resp, nil
}

// ForceGenerateFriendToken creates token for Authorization's user by which another users can add them, through the
// forced variant of token generation.
func (c *Client) ForceGenerateFriendToken(ctx context.Context, auth *Authorization) (GeneratedFriendToken, error) {
	var resp GeneratedFriendToken
	err := c.do(ctx, auth, "POST", "/friends/generate/force", nil, &resp)
	if err != nil {
		return GeneratedFriendToken{}, fmt.Errorf("failed to force generate friend token: %w", err)
	}

	return resp, nil
}

// AddFriend makes request to add the owner of provided FriendToken to Authorization's friends list.
func (c *Client) AddFriend(ctx context.Context, auth *Authorization, token FriendToken) error {
	req := addFriendRequest{
		Token: &token,
	}

	var resp addFriendResponse
	err := c.do(ctx, auth, "POST", "/friends/add", req, &resp)
	if err != nil {
		return fmt.Errorf("failed to add friend: %w", err)
	}

	if resp.Type == "FriendTokenExpired" {
		return ErrFriendTokenExpired
	}

	return nil
}

// GetFriendTokenOwner returns FriendTokenOwner of provided FriendToken without adding them, or ErrFriendTokenExpired
// when the token expired or never existed.
func (c *Client) GetFriendTokenOwner(ctx context.Context, token FriendToken) (*FriendTokenOwner, error) {
	var resp friendTokenOwnerResponse
	err := c.do(ctx, nil, "GET", "/friends/add/details/"+token.value, nil, &resp)
	if err != nil {
		return nil, fmt.Errorf("failed to get friend token owner: %w", err)
	}

	if resp.Type == "expired" {
		return nil, ErrFriendTokenExpired
	}

	return &resp.FriendTokenOwner, nil
}

// SendFriendRequest sends friend request from Authorization to user with provided ID and AccessHash.
func (c *Client) SendFriendRequest(ctx context.Context, auth *Authorization, userId UserId, accessHash UserAccessHash) error {
	req := friendRequestRequest{
		UserId:     userId,
		AccessHash: &accessHash,
	}

	err := c.do(ctx, auth, "POST", "/friends/request", req, nil)
	if err != nil {
		return fmt.Errorf("failed to send friend request: %w", err)
	}

	return nil
}

// DeclineFriendRequest declines Authorization's request to user with provided ID and AccessHash.
func (c *Client) DeclineFriendRequest(ctx context.Context, auth *Authorization, userId UserId, accessHash UserAccessHash) error {
	req := friendRequestRequest{
		UserId:     userId,
		AccessHash: &accessHash,
	}

	err := c.do(ctx, auth, "POST", "/friends/decline", req, nil)
	if err != nil {
		return fmt.Errorf("failed to decline friend request: %w", err)
	}

	return nil
}
