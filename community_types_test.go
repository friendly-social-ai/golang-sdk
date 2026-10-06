package sdk

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func MockCommunityPostId(i int64) CommunityPostId {
	return CommunityPostId{value: i}
}

func MockCommunityPostAccessHash(s string) CommunityPostAccessHash {
	return CommunityPostAccessHash{value: s}
}

func MockCommunityPostText(s string) CommunityPostText {
	return CommunityPostText{value: s}
}

func MockCursorId(s string) CursorId {
	return CursorId{value: s}
}

func TestCommunityTypes(t *testing.T) {
	t.Run("CommunityPostId", func(t *testing.T) {
		id := NewCommunityPostId(123)
		require.Equal(t, MockCommunityPostId(123), id)
		require.Equal(t, int64(123), id.Value())

		data, err := json.Marshal(id)
		require.NoError(t, err)
		require.Equal(t, `123`, string(data))

		var loadedId CommunityPostId
		err = json.Unmarshal(data, &loadedId)
		require.NoError(t, err)
		require.Equal(t, id, loadedId)
	})

	t.Run("CommunityPostAccessHash", func(t *testing.T) {
		valid := strings.Repeat("a", 256)
		hash, err := NewCommunityPostAccessHash(valid)
		require.NoError(t, err)
		require.Equal(t, MockCommunityPostAccessHash(valid), hash)
		require.Equal(t, valid, hash.Value())

		_, err = NewCommunityPostAccessHash("short")
		require.Error(t, err)
		require.ErrorIs(t, err, ErrCommunityPostAccessHashLengthMustBe256)

		data, err := json.Marshal(hash)
		require.NoError(t, err)
		require.Equal(t, `"`+valid+`"`, string(data))

		var loadedHash CommunityPostAccessHash
		err = json.Unmarshal(data, &loadedHash)
		require.NoError(t, err)
		require.Equal(t, hash, loadedHash)
	})

	t.Run("CommunityPostText", func(t *testing.T) {
		text, err := NewCommunityPostText("hello")
		require.NoError(t, err)
		require.Equal(t, MockCommunityPostText("hello"), text)
		require.Equal(t, "hello", text.Value())

		_, err = NewCommunityPostText(strings.Repeat("😀", 4096))
		require.NoError(t, err)

		_, err = NewCommunityPostText(strings.Repeat("1", 4097))
		require.Error(t, err)
		require.ErrorIs(t, err, ErrTooLongCommunityPostText)

		_, err = NewCommunityPostText("")
		require.Error(t, err)
		require.ErrorIs(t, err, ErrEmptyCommunityPostText)

		data, err := json.Marshal(text)
		require.NoError(t, err)
		require.Equal(t, `"hello"`, string(data))

		var loadedText CommunityPostText
		err = json.Unmarshal(data, &loadedText)
		require.NoError(t, err)
		require.Equal(t, text, loadedText)
	})

	t.Run("CursorId", func(t *testing.T) {
		var id CursorId
		err := json.Unmarshal([]byte(`"next"`), &id)
		require.NoError(t, err)
		require.Equal(t, MockCursorId("next"), id)
		require.Equal(t, "next", id.Value())

		data, err := json.Marshal(id)
		require.NoError(t, err)
		require.Equal(t, `"next"`, string(data))
	})
}
