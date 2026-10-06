package sdk

import (
	"encoding/json"
	"fmt"
	"unicode/utf8"
)

var (
	ErrCommunityPostAccessHashLengthMustBe256 = fmt.Errorf("community post access hash must be 256 characters long")
	ErrEmptyCommunityPostText                 = fmt.Errorf("community post text must not be empty")
	ErrTooLongCommunityPostText               = fmt.Errorf("community post text must be at most 4096 characters long")
)

// --- COMMUNITY POST ID ---

// CommunityPostId represents the unique identifier of community post.
type CommunityPostId struct {
	value int64
}

// NewCommunityPostId creates new CommunityPostId from int64.
func NewCommunityPostId(i int64) CommunityPostId {
	return CommunityPostId{value: i}
}

// Value returns CommunityPostId as a plain int64.
func (i CommunityPostId) Value() int64 {
	return i.value
}

func (i CommunityPostId) MarshalJSON() ([]byte, error) {
	return json.Marshal(i.value)
}

func (i *CommunityPostId) UnmarshalJSON(bytes []byte) error {
	return json.Unmarshal(bytes, &i.value)
}

// --- COMMUNITY POST ACCESS HASH ---

// CommunityPostAccessHash represents the unique hash associated with community post. Works in pair with CommunityPostId.
type CommunityPostAccessHash struct {
	value string
}

// Value returns CommunityPostAccessHash as a plain string.
func (h CommunityPostAccessHash) Value() string {
	return h.value
}

// NewCommunityPostAccessHash creates new CommunityPostAccessHash or returns an error if hash length isn't 256.
func NewCommunityPostAccessHash(s string) (CommunityPostAccessHash, error) {
	if len(s) != 256 {
		return CommunityPostAccessHash{}, fmt.Errorf("length is %d: %w", len(s), ErrCommunityPostAccessHashLengthMustBe256)
	}

	return CommunityPostAccessHash{value: s}, nil
}

func (h CommunityPostAccessHash) MarshalJSON() ([]byte, error) {
	return json.Marshal(h.value)
}

func (h *CommunityPostAccessHash) UnmarshalJSON(bytes []byte) error {
	return json.Unmarshal(bytes, &h.value)
}

// --- COMMUNITY POST TEXT ---

// CommunityPostText represents the text content of community post.
type CommunityPostText struct {
	value string
}

// Value returns CommunityPostText as a plain string.
func (t CommunityPostText) Value() string {
	return t.value
}

// NewCommunityPostText creates new CommunityPostText or returns an error if it is empty or longer than 4096 characters.
func NewCommunityPostText(s string) (CommunityPostText, error) {
	if s == "" {
		return CommunityPostText{}, ErrEmptyCommunityPostText
	}

	if n := utf8.RuneCountInString(s); n > 4096 {
		return CommunityPostText{}, fmt.Errorf("length is %d: %w", n, ErrTooLongCommunityPostText)
	}

	return CommunityPostText{value: s}, nil
}

func (t CommunityPostText) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.value)
}

func (t *CommunityPostText) UnmarshalJSON(bytes []byte) error {
	return json.Unmarshal(bytes, &t.value)
}

// --- CURSOR ID ---

// CursorId points to the next page of paginated results.
type CursorId struct {
	value string
}

// Value returns CursorId as a plain string.
func (i CursorId) Value() string {
	return i.value
}

func (i CursorId) MarshalJSON() ([]byte, error) {
	return json.Marshal(i.value)
}

func (i *CursorId) UnmarshalJSON(bytes []byte) error {
	return json.Unmarshal(bytes, &i.value)
}
