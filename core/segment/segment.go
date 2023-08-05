package segment

import (
	"errors"

	"github.com/draganm/follow-me/core/types"
)

var ErrNoMoreMessagesInSegment = errors.New("no more messages in segment")

type Segment interface {
	ContainsID(id types.UUID) bool
	GetMessageAfterID(prevID types.UUID) (id string, typ string, data []byte)
}
