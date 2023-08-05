package active

import (
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"math"
	"os"

	"github.com/draganm/follow-me/core/segment"
	"github.com/draganm/follow-me/core/types"
	"github.com/edsrzf/mmap-go"
	"github.com/minio/highwayhash"
	"golang.org/x/exp/slices"
)

type indexEntry struct {
	id     types.UUID
	offset int
}

type Segment struct {
	f                  *os.File
	mm                 mmap.MMap
	index              []indexEntry
	nextFreeByteOffset int
}

// ErrNotEnoughSpace is a predefined error returned when attempting to add a new
// message to an active segment that lacks sufficient storage space.
// This error typically signifies the need to begin using a new segment for message storage.
var ErrNotEnoughSpace = errors.New("not enough space in active segment")

// Open opens a file at the given path and returns a pointer to a Segment.
// If the file does not exist, Open will create it.
// The function also ensures that the size of the file matches the specified size.
// If the existing file is smaller, it will be resized.
// If the existing file is larger, an error will be returned.
// The file is memory-mapped for read and write operations.
//
// Example usage:
//
//	seg, err := Open("/path/to/file", 1024)
//	if err != nil {
//	    log.Fatalf("failed to open file: %v", err)
//	}
//
// Parameters:
//
//	path (string): The path to the file.
//	size (uint64): The minimum size of the file.
//
// Returns:
//
//	s (*Segment): A pointer to the Segment representing the opened file.
//	err (error): An error object that describes the error, if any occurred. nil otherwise.
//
// Errors:
//
// Returns an error if the function fails to open the file, get file stats,
// resize the file, or if the existing file is larger than the specified size,
// or if the file fails to be memory-mapped.
func Open(path string, size uint64) (s *Segment, err error) {
	f, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		return nil, fmt.Errorf("could not open %s: %w", path, err)
	}

	defer func() {
		if err != nil {
			err = errors.Join(
				err,
				f.Close(),
			)
		}
	}()

	st, err := f.Stat()
	if err != nil {
		return nil, fmt.Errorf("could not stat %s: %w", path, err)
	}

	if st.Size() < int64(size) {
		err = f.Truncate(int64(size))
		if err != nil {
			return nil, fmt.Errorf("could not resize %s to %d: %w", path, size, err)
		}
	}

	if st.Size() > int64(size) {
		return nil, fmt.Errorf("active segment %s has length %d, which is larger than requested size %d", path, st.Size(), size)
	}

	mm, err := mmap.Map(f, mmap.RDWR, 0)
	if err != nil {
		return nil, fmt.Errorf("could not mmap %s: %w", path, err)
	}
	pos := 0
	index := []indexEntry{}
	{
		for pos < len(mm) {
			if len(mm)-pos < 16 {
				return nil, fmt.Errorf("internal error: active segment corrupted - not enough bytes to get uuid")
			}
			id, _, _, ln, err := getMessage(mm, pos)
			if err != nil {
				return nil, fmt.Errorf("could not regenerate index: %w", err)
			}

			if id == types.ZeroUUID {
				break
			}

			index = append(index, indexEntry{
				id:     id,
				offset: pos,
			})
			pos += ln

		}

	}

	return &Segment{
		f:                  f,
		mm:                 mm,
		index:              index,
		nextFreeByteOffset: pos,
	}, nil

}

var randomKey []byte

func init() {
	randomKey, _ = hex.DecodeString(`70d7994cf728e46e23b7df19e0980a34c22ed448393da3581dcec378c67db2a3`)
}

// StoreMessage attempts to store a new message with the provided ID, type, and body in the current active segment.
// It adds the message to the segment's existing messages, preserving the order in which messages are stored.
// If the active segment does not have enough space to accommodate the new message,
// StoreMessage returns an ErrNotEnoughSpace error.
//
// Parameters:
//
//	id: A unique identifier for the message (uuidv6).
//
//	typ: The type of the message being stored.
//
//	body: The actual content of the message, provided as a byte slice.
//
// Returns:
//
//	nil if the message is successfully stored, ErrNotEnoughSpace if the active segment cannot accommodate the message.
func (s *Segment) StoreMessage(id types.UUID, typ string, body []byte) error {

	// TODO: check if the ID is larger than previous stored

	if len(typ) > 255 {
		return errors.New("type of the message is too long")
	}

	if len(body) > math.MaxInt32 {
		return errors.New("body of the message is too large")
	}

	spaceNeeded := len(id) + len(typ) + 1 + len(body) + 4 + 8
	spaceAvailable := len(s.mm) - s.nextFreeByteOffset

	if spaceNeeded > spaceAvailable {
		return ErrNotEnoughSpace
	}

	pos := s.nextFreeByteOffset
	copy(s.mm[pos:], id[:])
	pos += len(id)

	s.mm[pos] = uint8(len(typ))
	pos++
	copy(s.mm[pos:], []byte(typ))
	pos += len(typ)

	binary.BigEndian.PutUint32(s.mm[pos:], uint32(len(body)))
	pos += 4

	copy(s.mm[pos:], body)
	pos += len(body)

	hash := highwayhash.Sum64(s.mm[s.nextFreeByteOffset:pos], randomKey)

	binary.BigEndian.PutUint64(s.mm[pos:], hash)
	pos += 8

	s.index = append(s.index, indexEntry{
		id:     id,
		offset: s.nextFreeByteOffset,
	})
	s.nextFreeByteOffset = pos

	err := s.mm.Flush()
	if err != nil {
		return fmt.Errorf("could not flush memory map: %w", err)
	}

	return nil
}

func getMessage(d []byte, offset int) (msgID types.UUID, typ string, data []byte, ln int, err error) {

	dat := d[offset:]

	msgID, err = types.UUIDFromBytes(dat)
	if err != nil {
		return types.UUID{}, "", nil, 0, err
	}

	if msgID == types.ZeroUUID {
		return types.ZeroUUID, "", nil, 0, nil
	}

	defer func() {
		if err != nil {
			err = fmt.Errorf("message id %s: %w", msgID, err)
		}
	}()

	dat = dat[16:]
	ln += 16

	if len(dat) < 1 {
		return msgID, "", nil, 0, errors.New("internal error: message malformed: can't get type length")
	}

	typLen := int(dat[0])

	if len(dat) < typLen {
		return msgID, "", nil, 0, errors.New("internal error: message malformed: can't get types")
	}
	dat = dat[1:]

	ln++

	typ = string(dat[:typLen])

	dat = dat[typLen:]
	ln += typLen

	if len(dat) < 4 {
		return msgID, "", nil, 0, errors.New("internal error: message malformed: can't get message length")
	}

	msgLen := int(binary.BigEndian.Uint32(dat[:4]))

	dat = dat[4:]
	ln += 4

	if len(dat) < msgLen {
		return msgID, "", nil, 0, errors.New("internal error: message malformed: can't get message length")
	}

	data = dat[:msgLen]
	dat = dat[msgLen:]
	ln += msgLen

	calculatedHash := highwayhash.Sum64(d[offset:offset+ln], randomKey)

	if len(dat) < 8 {
		return msgID, "", nil, 0, errors.New("internal error: message malformed: not enough bytes for hash")
	}

	storedHash := binary.BigEndian.Uint64(dat)
	ln += 8

	if calculatedHash != storedHash {
		return msgID, "", nil, 0, errors.New("internal error: message malformed: hash mismatch")
	}

	return msgID, typ, data, ln, nil

}

func (s *Segment) GetMessageAfterID(id types.UUID) (msgID types.UUID, typ string, data []byte, err error) {
	idx, found := slices.BinarySearchFunc(s.index, id, func(e indexEntry, id types.UUID) int {
		return e.id.Compare(id)
	})

	if found {
		idx++
	}

	if idx >= len(s.index) {
		return types.UUID{}, "", nil, segment.ErrNoMoreMessagesInSegment
	}

	pos := s.index[idx].offset

	msgID, typ, data, _, err = getMessage(s.mm, pos)

	return msgID, typ, data, err

}

func (s *Segment) Close() error {
	return errors.Join(
		s.mm.Unmap(),
		s.f.Close(),
	)
}
