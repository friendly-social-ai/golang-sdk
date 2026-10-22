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

func TestListActivity_IndirectReplies(t *testing.T) {
	defer gock.Off()

	gock.New("https://api.getfriend.ly").
		Get("/activity/list").
		Reply(200).
		JSON(`{"data":[{"type":"indirect_replies","id":5,"instant":"2026-10-06T09:30:47.408013Z","isRead":false,"post":{"type":"plain","id":7,"accessHash":"root","instant":"2026-10-06T09:30:47.408013Z","replyPreviews":[],"newUnreadId":9,"text":"root","owner":{"id":2},"edited":false},"commonDenominator":{"type":"plain","id":8,"accessHash":"hash","instant":"2026-10-06T09:30:47.408013Z","replyPreviews":[],"newUnreadId":9,"text":"reply","owner":{"id":3},"edited":false}}],"nextId":null}`)

	client := NewClient()
	page, err := client.ListActivity(context.Background(), nil, nil)

	require.NoError(t, err)

	require.Equal(t, CommunityPostDescriptor{Id: MockCommunityPostId(8), AccessHash: MockCommunityPostAccessHash("hash")}, page.Data[0].CommonDenominator.Descriptor())
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

func TestGetNewUnreadActivity_Some(t *testing.T) {
	defer gock.Off()

	gock.New("https://api.getfriend.ly").
		Get("/activity/unreads/new").
		MatchHeader("X-User-Id", "1").
		MatchHeader("X-Token", "token").
		Reply(200).
		JSON(`{"type":"some","activityId":5}`)

	client := NewClient()
	auth := &Authorization{Id: MockUserId(1), Token: MockToken("token")}
	id, err := client.GetNewUnreadActivity(context.Background(), auth)

	require.NoError(t, err)
	expected := MockActivityId(5)
	require.Equal(t, &expected, id)
}

func TestGetNewUnreadActivity_None(t *testing.T) {
	defer gock.Off()

	gock.New("https://api.getfriend.ly").
		Get("/activity/unreads/new").
		Reply(200).
		JSON(`{"type":"none"}`)

	client := NewClient()
	id, err := client.GetNewUnreadActivity(context.Background(), nil)

	require.NoError(t, err)
	require.Nil(t, id)
}

func TestGetNewUnreadActivity_Failed(t *testing.T) {
	defer gock.Off()

	gock.New("https://api.getfriend.ly").
		Get("/activity/unreads/new").
		Reply(401)

	client := NewClient()
	_, err := client.GetNewUnreadActivity(context.Background(), nil)
	require.Error(t, err)
}

func TestAcknowledgeUnreadActivity_Success(t *testing.T) {
	defer gock.Off()

	gock.New("https://api.getfriend.ly").
		Post("/activity/unreads/acknowledge").
		MatchHeader("X-User-Id", "1").
		MatchHeader("X-Token", "token").
		Reply(200)

	client := NewClient()
	auth := &Authorization{Id: MockUserId(1), Token: MockToken("token")}
	err := client.AcknowledgeUnreadActivity(context.Background(), auth)
	require.NoError(t, err)
}

func TestAcknowledgeUnreadActivity_Failed(t *testing.T) {
	defer gock.Off()

	gock.New("https://api.getfriend.ly").
		Post("/activity/unreads/acknowledge").
		Reply(401)

	client := NewClient()
	err := client.AcknowledgeUnreadActivity(context.Background(), nil)
	require.Error(t, err)
}
