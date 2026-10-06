package sdk

import (
	"context"
	"testing"
	"time"

	"github.com/h2non/gock"
	"github.com/stretchr/testify/require"
)

func TestListActivity_Success(t *testing.T) {
	defer gock.Off()

	gock.New("https://api.getfriend.ly").
		Get("/activity/list").
		MatchHeader("X-User-Id", "1").
		MatchHeader("X-Token", "token").
		Reply(200).
		JSON(`{"data":[{"type":"reply","id":5,"instant":"2026-10-06T09:30:47.408013Z","isRead":false,"post":{"id":7,"accessHash":"hash","instant":"2026-10-06T09:30:47.408013Z","replyPreviews":[],"text":"hello","owner":{"id":2},"edited":false}}],"nextId":"next"}`)

	client := NewClient()
	auth := &Authorization{Id: MockUserId(1), Token: MockToken("token")}
	page, err := client.ListActivity(context.Background(), auth, nil)

	require.NoError(t, err)

	instant := time.Date(2026, 10, 6, 9, 30, 47, 408013000, time.UTC)
	text := MockCommunityPostText("hello")
	next := MockCursorId("next")
	require.Equal(t, &Cursor[Activity]{
		Data: []Activity{{
			Type:    "reply",
			Id:      MockActivityId(5),
			Instant: instant,
			Post: &CommunityPost{
				Id:            MockCommunityPostId(7),
				AccessHash:    MockCommunityPostAccessHash("hash"),
				Instant:       instant,
				ReplyPreviews: []UserDetails{},
				Text:          &text,
				Owner:         &UserDetails{Id: MockUserId(2)},
			},
		}},
		NextId: &next,
	}, page)
}

func TestListActivity_Cursor(t *testing.T) {
	defer gock.Off()

	gock.New("https://api.getfriend.ly").
		Get("/activity/list/next").
		Reply(200).
		JSON(`{"data":[],"nextId":null}`)

	client := NewClient()
	cursor := MockCursorId("next")
	page, err := client.ListActivity(context.Background(), nil, &cursor)

	require.NoError(t, err)
	require.Equal(t, &Cursor[Activity]{Data: []Activity{}}, page)
}

func TestListActivity_Failed(t *testing.T) {
	defer gock.Off()

	gock.New("https://api.getfriend.ly").
		Get("/activity/list").
		Reply(400)

	client := NewClient()
	_, err := client.ListActivity(context.Background(), nil, nil)
	require.Error(t, err)
}

func TestReadActivity_Success(t *testing.T) {
	defer gock.Off()

	gock.New("https://api.getfriend.ly").
		Post("/activity/read/5").
		MatchHeader("X-User-Id", "1").
		MatchHeader("X-Token", "token").
		Reply(200)

	client := NewClient()
	auth := &Authorization{Id: MockUserId(1), Token: MockToken("token")}
	err := client.ReadActivity(context.Background(), auth, MockActivityId(5))
	require.NoError(t, err)
}

func TestReadActivity_Failed(t *testing.T) {
	defer gock.Off()

	gock.New("https://api.getfriend.ly").
		Post("/activity/read/5").
		Reply(404)

	client := NewClient()
	err := client.ReadActivity(context.Background(), nil, MockActivityId(5))
	require.Error(t, err)
}
