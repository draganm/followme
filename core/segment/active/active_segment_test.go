package active_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/draganm/follow-me/core/segment/active"
	"github.com/draganm/follow-me/core/types"
	"github.com/stretchr/testify/require"
)

func TestActiveSegment(t *testing.T) {
	require := require.New(t)
	td, err := os.MkdirTemp("", "")
	require.NoError(err)

	t.Cleanup(func() {
		err := os.RemoveAll(td)
		require.NoError(err)
	})

	as, err := active.Open(filepath.Join(td, "active"), 128)
	require.NoError(err)

	firstUUID, err := types.ParseUUID("00010203-0405-0607-0809-0a0b0c0d0e0f")
	require.NoError(err)

	err = as.StoreMessage(firstUUID, "test", []byte{1, 2, 3})
	require.NoError(err)

	secondUUID, err := types.ParseUUID("00010203-0405-0607-0809-0a0b0c0d0f0f")
	require.NoError(err)

	err = as.StoreMessage(secondUUID, "test2", []byte{4, 5, 6})
	require.NoError(err)

	prevUUID, err := types.ParseUUID("00010203-0405-0607-0809-0a0b0c0d0e0e")
	require.NoError(err)

	id, typ, data, err := as.GetMessageAfterID(prevUUID)
	require.NoError(err)

	require.Equal(firstUUID, id)
	require.Equal("test", typ)
	require.Equal([]byte{1, 2, 3}, data)

	id, typ, data, err = as.GetMessageAfterID(firstUUID)
	require.NoError(err)

	require.Equal(secondUUID, id)
	require.Equal("test2", typ)
	require.Equal([]byte{4, 5, 6}, data)

	err = as.Close()
	require.NoError(err)

	as, err = active.Open(filepath.Join(td, "active"), 128)
	require.NoError(err)

	id, typ, data, err = as.GetMessageAfterID(prevUUID)
	require.NoError(err)

	require.Equal(firstUUID, id)
	require.Equal("test", typ)
	require.Equal([]byte{1, 2, 3}, data)

	id, typ, data, err = as.GetMessageAfterID(firstUUID)
	require.NoError(err)

	require.Equal(secondUUID, id)
	require.Equal("test2", typ)
	require.Equal([]byte{4, 5, 6}, data)

}
