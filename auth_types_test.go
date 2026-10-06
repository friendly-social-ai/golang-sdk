package sdk

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func MockUserId(i int64) UserId {
	return UserId{value: i}
}

func MockUserAccessHash(s string) UserAccessHash {
	return UserAccessHash{value: s}
}

func MockToken(s string) Token {
	return Token{value: s}
}

func MockFirebaseToken(s string) FirebaseToken {
	return FirebaseToken{value: s}
}

func TestAuthTypes(t *testing.T) {
	t.Run("FirebaseToken", func(t *testing.T) {
		token, err := NewFirebaseToken("fcm")
		require.NoError(t, err)
		require.Equal(t, MockFirebaseToken("fcm"), token)
		require.Equal(t, "fcm", token.Value())

		_, err = NewFirebaseToken(strings.Repeat("1", 4097))
		require.Error(t, err)
		require.ErrorIs(t, err, ErrTooLongFirebaseToken)

		data, err := json.Marshal(token)
		require.NoError(t, err)
		require.Equal(t, `"fcm"`, string(data))

		var loaded FirebaseToken
		require.NoError(t, json.Unmarshal(data, &loaded))
		require.Equal(t, token, loaded)
	})

	t.Run("UserId", func(t *testing.T) {
		id := NewUserId(123)
		require.Equal(t, UserId{value: 123}, id)
		require.Equal(t, int64(123), id.Value())

		data, err := json.Marshal(id)
		require.NoError(t, err)
		require.Equal(t, `123`, string(data))

		var loadedId UserId
		err = json.Unmarshal(data, &loadedId)
		require.NoError(t, err)
		require.Equal(t, id, loadedId)
	})

	t.Run("Token", func(t *testing.T) {
		token, err := NewToken(strings.Repeat("1", 256))
		require.NoError(t, err)
		require.Equal(t, MockToken(strings.Repeat("1", 256)), token)
		require.Equal(t, strings.Repeat("1", 256), token.Value())

		_, err = NewToken("1")
		require.Error(t, err)
		require.ErrorIs(t, err, ErrTokenLengthMustBe256)

		data, err := json.Marshal(token)
		require.NoError(t, err)
		require.Equal(t, `"`+strings.Repeat("1", 256)+`"`, string(data))

		var loadedToken Token
		err = json.Unmarshal(data, &loadedToken)
		require.NoError(t, err)
		require.Equal(t, token, loadedToken)
	})

	t.Run("UserAccesssHash", func(t *testing.T) {
		hash, err := NewUserAccessHash(strings.Repeat("1", 256))
		require.NoError(t, err)
		require.Equal(t, MockUserAccessHash(strings.Repeat("1", 256)), hash)
		require.Equal(t, strings.Repeat("1", 256), hash.Value())

		_, err = NewUserAccessHash("1")
		require.Error(t, err)
		require.ErrorIs(t, err, ErrUserAccessHashLengthMustBe256)

		data, err := json.Marshal(hash)
		require.NoError(t, err)
		require.Equal(t, `"`+strings.Repeat("1", 256)+`"`, string(data))

		var loadedHash UserAccessHash
		err = json.Unmarshal(data, &loadedHash)
		require.NoError(t, err)
		require.Equal(t, hash, loadedHash)
	})
}
