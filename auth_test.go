package sdk

import (
	"context"
	"testing"

	"github.com/h2non/gock"
	"github.com/stretchr/testify/require"
)

func TestRegister_Success(t *testing.T) {
	defer gock.Off()

	gock.New("https://api.getfriend.ly").
		Post("/auth/generate").
		JSON(`{"nickname":"atennop", "description":"bio","interests":["programming"],"avatar":{"id":10,"accessHash":"hash"},"socialLink":"https://github.com/Atennop1"}`).
		Reply(200).
		JSON(`{"id":1,"token":"token","accessHash":"hash"}`)

	client := NewClient()
	auth, err := client.Register(context.Background(),
		MockNickname("atennop"),
		MockUserDescription("bio"),
		MockInterests([]Interest{MockInterest("programming")}),
		&FileDescriptor{Id: MockFileId(10), AccessHash: MockFileAccessHash("hash")},
		MockSocialLink("https://github.com/Atennop1"))

	require.NoError(t, err)
	require.Equal(t, &Authorization{
		Id:         MockUserId(1),
		Token:      MockToken("token"),
		AccessHash: MockUserAccessHash("hash"),
	}, auth)
}

func TestRegister_Failed(t *testing.T) {
	defer gock.Off()

	gock.New("https://api.getfriend.ly").
		Post("/auth/generate").
		Reply(400)

	client := NewClient()
	_, err := client.Register(context.Background(),
		MockNickname("atennop"),
		MockUserDescription("bio"),
		MockInterests([]Interest{MockInterest("programming")}),
		&FileDescriptor{Id: MockFileId(10), AccessHash: MockFileAccessHash("hash")},
		MockSocialLink("https://github.com/Atennop1"))

	require.Error(t, err)
}

func TestSendLoginRequest_Success(t *testing.T) {
	defer gock.Off()

	gock.New("https://api.getfriend.ly").
		Post("/auth/email").
		MatchHeader("X-Locale", "ru").
		JSON(`{"email": "example@example.com"}`).
		Reply(200)

	client := NewClient()
	err := client.SendLoginRequest(context.Background(), MockEmail("example@example.com"), MockEmailLocale("ru"))
	require.NoError(t, err)
}

func TestSendLoginRequest_Failed(t *testing.T) {
	defer gock.Off()

	gock.New("https://api.getfriend.ly").
		Post("/auth/email").
		Reply(400)

	client := NewClient()
	err := client.SendLoginRequest(context.Background(), MockEmail("example@example.com"), MockEmailLocale("ru"))
	require.Error(t, err)
}

func TestSendLoginRequest_NewRequestFailed(t *testing.T) {
	client := NewClient()
	err := client.SendLoginRequest(nil, MockEmail("example@example.com"), MockEmailLocale("ru")) //nolint:staticcheck
	require.Error(t, err)
}

func TestConfirmLogin_Success(t *testing.T) {
	defer gock.Off()

	gock.New("https://api.getfriend.ly").
		Post("/auth/login").
		JSON(`{"email":"example@example.com", "code":11111111}`).
		Reply(200).
		JSON(`{"id":1,"token":"token","accessHash":"hash"}`)

	client := NewClient()
	auth, err := client.ConfirmLogin(context.Background(), MockEmail("example@example.com"), MockEmailCode(11111111))

	require.NoError(t, err)
	require.Equal(t, &Authorization{
		Id:         MockUserId(1),
		Token:      MockToken("token"),
		AccessHash: MockUserAccessHash("hash"),
	}, auth)
}

func TestConfirmLogin_Failed(t *testing.T) {
	defer gock.Off()

	gock.New("https://api.getfriend.ly").
		Post("/auth/login").
		Reply(400)

	client := NewClient()
	_, err := client.ConfirmLogin(context.Background(), MockEmail("example@example.com"), MockEmailCode(11111111))
	require.Error(t, err)
}

func TestRegister_NoSocialLink(t *testing.T) {
	defer gock.Off()

	gock.New("https://api.getfriend.ly").
		Post("/auth/generate").
		JSON(`{"nickname":"atennop","description":"bio","interests":["programming"],"avatar":null,"socialLink":null}`).
		Reply(200).
		JSON(`{"id":1,"accessHash":"hash","token":"token"}`)

	client := NewClient()
	_, err := client.Register(context.Background(), MockNickname("atennop"), MockUserDescription("bio"),
		MockInterests([]Interest{MockInterest("programming")}), nil, SocialLink{})
	require.NoError(t, err)
}

func TestRegisterFirebaseToken_Success(t *testing.T) {
	defer gock.Off()

	gock.New("https://api.getfriend.ly").
		Post("/auth/firebase").
		MatchHeader("X-User-Id", "1").
		MatchHeader("X-Token", "token").
		JSON(`{"firebaseToken":"fcm"}`).
		Reply(200)

	client := NewClient()
	auth := &Authorization{Id: MockUserId(1), Token: MockToken("token")}
	err := client.RegisterFirebaseToken(context.Background(), auth, MockFirebaseToken("fcm"))
	require.NoError(t, err)
}

func TestRegisterFirebaseToken_Failed(t *testing.T) {
	defer gock.Off()

	gock.New("https://api.getfriend.ly").
		Post("/auth/firebase").
		Reply(401)

	client := NewClient()
	err := client.RegisterFirebaseToken(context.Background(), nil, MockFirebaseToken("fcm"))
	require.Error(t, err)
}
