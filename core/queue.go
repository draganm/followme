package core

import (
	"fmt"
	"sync"

	"github.com/draganm/follow-me/core/segment/active"
	"github.com/draganm/follow-me/core/segment/archived"
)

type Queue struct {
	mu       *sync.Mutex
	archived []*archived.Segment
	active   *active.Segment
}

func Open(path string, segmentSize uint64, maxSegments uint) (*Queue, error) {
	return nil, fmt.Errorf("not yet implemented")
}
