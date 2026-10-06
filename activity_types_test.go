package sdk

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

func MockActivityId(i int64) ActivityId {
	return ActivityId{value: i}
}

func TestActivityTypes(t *testing.T) {
	t.Run("ActivityId", func(t *testing.T) {
		id := NewActivityId(123)
		require.Equal(t, MockActivityId(123), id)
		require.Equal(t, int64(123), id.Value())

		data, err := json.Marshal(id)
		require.NoError(t, err)
		require.Equal(t, `123`, string(data))

		var loadedId ActivityId
		err = json.Unmarshal(data, &loadedId)
		require.NoError(t, err)
		require.Equal(t, id, loadedId)
	})
}
