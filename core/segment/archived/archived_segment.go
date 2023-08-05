package archived

import (
	"fmt"
	"os"

	"github.com/edsrzf/mmap-go"
)

type Segment struct {
	mm mmap.MMap
}

func Open(path string) (*Segment, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("could not open %s: %w", path, err)
	}

	mm, err := mmap.Map(f, mmap.RDONLY, 0)
	if err != nil {
		return nil, fmt.Errorf("could not mmap %s: %w", path, err)
	}

	return &Segment{mm: mm}, nil
}
