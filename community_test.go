package sdk

import (
	"context"
	"testing"
	"time"

	"github.com/h2non/gock"
	"github.com/stretchr/testify/require"
)

func mockCommunityPost(id int64, owner int64, text string) CommunityPost {
	postText := MockCommunityPostText(text)
	return CommunityPost{
		Type:          "plain",
		Id:            MockCommunityPostId(id),
		AccessHash:    MockCommunityPostAccessHash("hash"),
		Instant:       time.Date(2026, 10, 6, 8, 40, 32, 308976000, time.UTC),
		ReplyPreviews: []UserDetails{},
		Text:          &postText,
		Owner:         &UserDetails{Id: MockUserId(owner)},
	}
}

func TestPostCommunity_Success(t *testing.T) {
	defer gock.Off()

	gock.New("https://api.getfriend.ly").
		Post("/community").
		MatchHeader("Content-Type", "application/json").
		MatchHeader("X-User-Id", "1").
		MatchHeader("X-Token", "token").
		JSON(`{"text":"hello"}`).
		Reply(200).
		JSON(`{"id":5,"accessHash":"hash"}`)

	client := NewClient()
	auth := &Authorization{Id: MockUserId(1), Token: MockToken("token")}
	post, err := client.PostCommunity(context.Background(), auth, MockCommunityPostText("hello"), nil)

	require.NoError(t, err)
	require.Equal(t, &CommunityPostDescriptor{Id: MockCommunityPostId(5), AccessHash: MockCommunityPostAccessHash("hash")}, post)
}

func TestPostCommunity_Reply(t *testing.T) {
	defer gock.Off()

	gock.New("https://api.getfriend.ly").
		Post("/community").
		JSON(`{"text":"hello","replyTo":{"id":4,"accessHash":"parent"}}`).
		Reply(200).
		JSON(`{"id":5,"accessHash":"hash"}`)

	client := NewClient()
	auth := &Authorization{Id: MockUserId(1), Token: MockToken("token")}
	replyTo := &CommunityPostDescriptor{Id: MockCommunityPostId(4), AccessHash: MockCommunityPostAccessHash("parent")}
	_, err := client.PostCommunity(context.Background(), auth, MockCommunityPostText("hello"), replyTo)

	require.NoError(t, err)
}

func TestPostCommunity_Failed(t *testing.T) {
	defer gock.Off()

	gock.New("https://api.getfriend.ly").
		Post("/community").
		Reply(400)

	client := NewClient()
	_, err := client.PostCommunity(context.Background(), nil, MockCommunityPostText("hello"), nil)
	require.Error(t, err)
}

func TestListCommunity_Success(t *testing.T) {
	defer gock.Off()

	gock.New("https://api.getfriend.ly").
		Get("/community/list").
		MatchHeader("Content-Type", "application/json").
		MatchHeader("X-User-Id", "1").
		MatchHeader("X-Token", "token").
		Reply(200).
		JSON(`{"data":[{"id":7,"accessHash":"hash","instant":"2026-10-06T08:40:32.308976Z","replyPreviews":[],"text":"hello","owner":{"id":2},"edited":false}],"nextId":"next"}`)

	client := NewClient()
	auth := &Authorization{Id: MockUserId(1), Token: MockToken("token")}
	page, err := client.ListCommunity(context.Background(), auth, nil)

	require.NoError(t, err)

	post := mockCommunityPost(7, 2, "hello")
	post.Type = ""
	next := MockCursorId("next")
	require.Equal(t, &Cursor[CommunityPost]{Data: []CommunityPost{post}, NextId: &next}, page)
	require.False(t, page.Data[0].Deleted())
}

func TestListCommunity_Cursor(t *testing.T) {
	defer gock.Off()

	gock.New("https://api.getfriend.ly").
		Get("/community/list/next").
		Reply(200).
		JSON(`{"data":[],"nextId":null}`)

	client := NewClient()
	cursor := MockCursorId("next")
	page, err := client.ListCommunity(context.Background(), nil, &cursor)

	require.NoError(t, err)
	require.Equal(t, &Cursor[CommunityPost]{Data: []CommunityPost{}}, page)
}

func TestListCommunity_Failed(t *testing.T) {
	defer gock.Off()

	gock.New("https://api.getfriend.ly").
		Get("/community/list").
		Reply(400)

	client := NewClient()
	_, err := client.ListCommunity(context.Background(), nil, nil)
	require.Error(t, err)
}

func TestGetCommunityPost_Success(t *testing.T) {
	defer gock.Off()

	gock.New("https://api.getfriend.ly").
		Get("/community/2/7/hash").
		MatchHeader("Content-Type", "application/json").
		MatchHeader("X-User-Id", "1").
		MatchHeader("X-Token", "token").
		Reply(200).
		JSON(`{
			"post":{"type":"plain","id":7,"accessHash":"hash","instant":"2026-10-06T08:40:32.308976Z","replyPreviews":[],"text":"hello","owner":{"id":2},"edited":false},
			"replies":{"data":[
				{"type":"single","post":{"type":"plain","id":8,"accessHash":"hash","instant":"2026-10-06T08:40:32.308976Z","replyPreviews":[],"text":"single","owner":{"id":3},"edited":false}},
				{"type":"thread","thread":[
					{"type":"deleted","id":9,"accessHash":"hash","instant":"2026-10-06T08:40:32.308976Z","replyPreviews":[]},
					{"type":"plain","id":10,"accessHash":"hash","instant":"2026-10-06T08:40:32.308976Z","replyPreviews":[],"text":"thread","owner":{"id":2},"edited":true}
				]}
			],"nextId":null},
			"upstream":[]
		}`)

	client := NewClient()
	auth := &Authorization{Id: MockUserId(1), Token: MockToken("token")}
	desc := CommunityPostDescriptor{Id: MockCommunityPostId(7), AccessHash: MockCommunityPostAccessHash("hash")}
	details, err := client.GetCommunityPost(context.Background(), auth, desc)

	require.NoError(t, err)

	single := mockCommunityPost(8, 3, "single")
	deleted := CommunityPost{
		Type:          "deleted",
		Id:            MockCommunityPostId(9),
		AccessHash:    MockCommunityPostAccessHash("hash"),
		Instant:       time.Date(2026, 10, 6, 8, 40, 32, 308976000, time.UTC),
		ReplyPreviews: []UserDetails{},
	}
	edited := mockCommunityPost(10, 2, "thread")
	edited.Edited = true

	require.Equal(t, &CommunityPostDetails{
		Post: mockCommunityPost(7, 2, "hello"),
		Replies: Cursor[CommunityPostReply]{Data: []CommunityPostReply{
			{Type: "single", Post: &single},
			{Type: "thread", Thread: []CommunityPost{deleted, edited}},
		}},
		Upstream: []CommunityPost{},
	}, details)
	require.Equal(t, desc, details.Post.Descriptor())
	require.Equal(t, []CommunityPost{single}, details.Replies.Data[0].Posts())
	require.Equal(t, []CommunityPost{deleted, edited}, details.Replies.Data[1].Posts())
	require.True(t, deleted.Deleted())
}

func TestGetCommunityPost_Failed(t *testing.T) {
	defer gock.Off()

	gock.New("https://api.getfriend.ly").
		Get("/community/2/7/hash").
		Reply(404)

	client := NewClient()
	desc := CommunityPostDescriptor{Id: MockCommunityPostId(7), AccessHash: MockCommunityPostAccessHash("hash")}
	_, err := client.GetCommunityPost(context.Background(), nil, desc)
	require.Error(t, err)
}

func TestGetCommunityReplies_Success(t *testing.T) {
	defer gock.Off()

	gock.New("https://api.getfriend.ly").
		Get("/community/7/hash/replies2/next").
		MatchHeader("X-User-Id", "1").
		MatchHeader("X-Token", "token").
		Reply(200).
		JSON(`{"data":[{"type":"single","post":{"type":"plain","id":8,"accessHash":"hash","instant":"2026-10-06T08:40:32.308976Z","replyPreviews":[],"text":"single","owner":{"id":3},"edited":false}}],"nextId":null}`)

	client := NewClient()
	auth := &Authorization{Id: MockUserId(1), Token: MockToken("token")}
	desc := CommunityPostDescriptor{Id: MockCommunityPostId(7), AccessHash: MockCommunityPostAccessHash("hash")}
	cursor := MockCursorId("next")
	page, err := client.GetCommunityReplies(context.Background(), auth, desc, &cursor)

	require.NoError(t, err)

	single := mockCommunityPost(8, 3, "single")
	require.Equal(t, &Cursor[CommunityPostReply]{Data: []CommunityPostReply{{Type: "single", Post: &single}}}, page)
}

func TestGetCommunityReplies_Failed(t *testing.T) {
	defer gock.Off()

	gock.New("https://api.getfriend.ly").
		Get("/community/7/hash/replies2").
		Reply(400)

	client := NewClient()
	desc := CommunityPostDescriptor{Id: MockCommunityPostId(7), AccessHash: MockCommunityPostAccessHash("hash")}
	_, err := client.GetCommunityReplies(context.Background(), nil, desc, nil)
	require.Error(t, err)
}

func TestEditCommunityPost_Success(t *testing.T) {
	defer gock.Off()

	gock.New("https://api.getfriend.ly").
		Post("/community/7/edit").
		MatchHeader("Content-Type", "application/json").
		MatchHeader("X-User-Id", "1").
		MatchHeader("X-Token", "token").
		JSON(`{"text":{"value":"updated"}}`).
		Reply(200)

	client := NewClient()
	auth := &Authorization{Id: MockUserId(1), Token: MockToken("token")}
	err := client.EditCommunityPost(context.Background(), auth, MockCommunityPostId(7), MockCommunityPostText("updated"))
	require.NoError(t, err)
}

func TestEditCommunityPost_Failed(t *testing.T) {
	defer gock.Off()

	gock.New("https://api.getfriend.ly").
		Post("/community/7/edit").
		Reply(404)

	client := NewClient()
	err := client.EditCommunityPost(context.Background(), nil, MockCommunityPostId(7), MockCommunityPostText("updated"))
	require.Error(t, err)
}

func TestDeleteCommunityPost_Success(t *testing.T) {
	defer gock.Off()

	gock.New("https://api.getfriend.ly").
		Post("/community/7/delete").
		MatchHeader("X-User-Id", "1").
		MatchHeader("X-Token", "token").
		Reply(200)

	client := NewClient()
	auth := &Authorization{Id: MockUserId(1), Token: MockToken("token")}
	err := client.DeleteCommunityPost(context.Background(), auth, MockCommunityPostId(7))
	require.NoError(t, err)
}

func TestDeleteCommunityPost_Failed(t *testing.T) {
	defer gock.Off()

	gock.New("https://api.getfriend.ly").
		Post("/community/7/delete").
		Reply(404)

	client := NewClient()
	err := client.DeleteCommunityPost(context.Background(), nil, MockCommunityPostId(7))
	require.Error(t, err)
}
