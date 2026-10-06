package sdk

import (
	"context"
	"fmt"
)

// UserDetails represents complete information about some user: ID, AccessHash, Nickname, Email, Description, list of
// Interests, Avatar, SocialLink and Friendship with Authorization's user. Email is nil when none is bound or hidden.
type UserDetails struct {
	Id          UserId          `json:"id"`
	AccessHash  UserAccessHash  `json:"accessHash"`
	Nickname    Nickname        `json:"nickname"`
	Email       *Email          `json:"email"`
	Description UserDescription `json:"description"`
	Interests   Interests       `json:"interests"`
	Avatar      *FileDescriptor `json:"avatar"`
	SocialLink  SocialLink      `json:"socialLink"`
	Friendship  Friendship      `json:"friendship"`
}

// UserProfile represents UserDetails of some user together with friends in common with Authorization's user.
type UserProfile struct {
	User          UserDetails   `json:"user"`
	CommonFriends []UserDetails `json:"commonFriends"`
}

type editAccountOption func(*editAccountRequest)

type editAccountValue[T any] struct {
	Value T `json:"value"`
}

type editAccountRequest struct {
	Nickname    *editAccountValue[Nickname]        `json:"nickname,omitempty"`
	Description *editAccountValue[UserDescription] `json:"description,omitempty"`
	Interests   *editAccountValue[Interests]       `json:"interests,omitempty"`
	Avatar      *editAccountValue[*FileDescriptor] `json:"avatar,omitempty"`
	SocialLink  *editAccountValue[SocialLink]      `json:"socialLink,omitempty"`
}

// GetSelfDetails returns UserDetails structure for provided Authorization data.
func (c *Client) GetSelfDetails(ctx context.Context, auth *Authorization) (*UserDetails, error) {
	var profile UserProfile
	err := c.do(ctx, auth, "GET", "/users/details2", nil, &profile)
	if err != nil {
		return nil, fmt.Errorf("failed to get self details: %w", err)
	}

	return &profile.User, nil
}

// GetUserDetails returns UserProfile for provided user's ID and AccessHash from provided Authorization's perspective.
func (c *Client) GetUserDetails(ctx context.Context, auth *Authorization, userId UserId, accessHash UserAccessHash) (*UserProfile, error) {
	var profile UserProfile
	err := c.do(ctx, auth, "GET", fmt.Sprintf("/users/details2/%d/%s", userId.value, accessHash.value), nil, &profile)
	if err != nil {
		return nil, fmt.Errorf("failed to get user details: %w", err)
	}

	return &profile, nil
}

// EditNicknameOption applies new Nickname for editing account request.
func EditNicknameOption(nickname Nickname) editAccountOption {
	return func(r *editAccountRequest) {
		r.Nickname = &editAccountValue[Nickname]{nickname}
	}
}

// EditDescriptionOption applies new UserDescription for editing account request.
func EditDescriptionOption(description UserDescription) editAccountOption {
	return func(r *editAccountRequest) {
		r.Description = &editAccountValue[UserDescription]{description}
	}
}

// EditInterestsOption applies new Interests for editing account request.
func EditInterestsOption(interests Interests) editAccountOption {
	return func(r *editAccountRequest) {
		r.Interests = &editAccountValue[Interests]{interests}
	}
}

// EditAvatarOption applies new Avatar for editing account request.
func EditAvatarOption(avatar *FileDescriptor) editAccountOption {
	return func(r *editAccountRequest) {
		r.Avatar = &editAccountValue[*FileDescriptor]{avatar}
	}
}

// EditSocialLinkOption applies new SocialLink for editing account request.
func EditSocialLinkOption(link SocialLink) editAccountOption {
	return func(r *editAccountRequest) {
		r.SocialLink = &editAccountValue[SocialLink]{link}
	}
}

// EditAccount makes request for editing Authorization's account using provided options.
func (c *Client) EditAccount(ctx context.Context, auth *Authorization, opts ...editAccountOption) error {
	if len(opts) == 0 {
		return nil
	}

	req := editAccountRequest{}
	for _, opt := range opts {
		opt(&req)
	}

	err := c.do(ctx, auth, "PATCH", "/users/edit", req, nil)
	if err != nil {
		return fmt.Errorf("failed to edit account: %w", err)
	}

	return nil
}
