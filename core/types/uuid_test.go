package types_test

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/draganm/follow-me/core/types"
	"github.com/stretchr/testify/require"
)

func TestUUIDFormatting(t *testing.T) {
	require := require.New(t)

	uuid := types.UUID{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15}

	require.Equal("00010203-0405-0607-0809-0a0b0c0d0e0f", uuid.String())

	require.Equal("|00010203-0405-0607-0809-0a0b0c0d0e0f|", fmt.Sprintf("|%s|", uuid))

	d, err := json.Marshal(uuid)
	require.NoError(err)

	require.Equal(`"00010203-0405-0607-0809-0a0b0c0d0e0f"`, string(d))
}

func TestParsing(t *testing.T) {
	require := require.New(t)

	parsed, err := types.ParseUUID("00010203-0405-0607-0809-0a0b0c0d0e0f")
	require.NoError(err)

	uuid := types.UUID{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15}

	require.Equal(uuid, parsed)

}
