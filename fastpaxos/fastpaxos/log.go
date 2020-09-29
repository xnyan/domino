package fastpaxos

import (
	//"errors"
	"fmt"
	"strconv"
	"strings"
)

// Log index
type LogIdx struct {
	Seg    int64 // segment index
	Offset int64 // offset in a segment
}

func NewLogIdx(seg, offset int64) *LogIdx {
	return &LogIdx{seg, offset}
}

func (idx *LogIdx) IncOffset() {
	idx.Offset++
}

func (idx *LogIdx) Equals(i *LogIdx) bool {
	if idx.Offset != i.Offset || idx.Seg != i.Seg {
		return false
	}
	return true
}

func (idx *LogIdx) String() string {
	return strconv.FormatInt(idx.Seg, 10) + "-" + strconv.FormatInt(idx.Offset, 10)
}

func ParseLogIdx(s string) *LogIdx {
	ele := strings.Split(s, "-")
	if len(ele) != 2 {
		logger.Fatalf("Invalid log index format: (%s)", s)
	}

	seg, err1 := strconv.ParseInt(ele[0], 10, 64)
	if err1 != nil {
		logger.Fatalf("Invalid log index format, segment: (%s), error: %v", ele[0], err1)
	}
	offset, err2 := strconv.ParseInt(ele[1], 10, 64)
	if err2 != nil {
		logger.Fatalf("Invalid log index offset, segment: (%s), error: %v", ele[1], err2)
	}

	return &LogIdx{seg, offset}
}

type Log interface {
	Get(idx *LogIdx) (*Entry, error)

	Put(idx *LogIdx, entry *Entry) error

	// Increases the given index by one
	IncIdx(idx *LogIdx) error

	// Returns the next index
	NextIdx(idx *LogIdx) (*LogIdx, error)

	Size() string
}

type FixedLog struct {
	log []*Entry // Operation log
	cap int64
}

func NewFixedLog(size int64) *FixedLog {
	return &FixedLog{
		log: make([]*Entry, size, size),
		cap: size,
	}
}

func (l *FixedLog) Get(idx *LogIdx) (*Entry, error) {
	if idx.Offset >= l.cap {
		return nil, fmt.Errorf("FixedLog Get() idx = (%s) is out of range, capacity = (%d)",
			idx, l.cap)
	}

	return l.log[idx.Offset], nil
}

func (l *FixedLog) Put(idx *LogIdx, entry *Entry) error {
	if idx.Offset >= l.cap {
		return fmt.Errorf("FixedLog Put() idx = (%s) is out of range, capacity = (%d)", idx, l.cap)
	}

	l.log[idx.Offset] = entry
	return nil
}

func (l *FixedLog) IncIdx(idx *LogIdx) error {
	if idx.Offset+1 >= l.cap {
		return fmt.Errorf("FixedLog IncIdx() idx = (%s) is the max idx.", idx)
	}

	idx.Offset++
	return nil
}

func (l *FixedLog) NextIdx(idx *LogIdx) (*LogIdx, error) {
	if idx.Offset+1 >= l.cap {
		return nil, fmt.Errorf("FixedLog NextIdx() idx = (%s) is the max idx.", idx)

	}

	return &LogIdx{idx.Seg, idx.Offset + 1}, nil
}

func (l *FixedLog) Size() string {
	return strconv.FormatInt(l.cap, 10)
}

// A log that consists of multiple segments. Segments are in the same size. The
// log will automatically extend as needed.
type SegLog struct {
	SegNum  int64 // inital number of segments
	SegSize int64

	log [][]*Entry
}

func NewSegLog(segNum, segSize int64) *SegLog {
	l := &SegLog{
		SegNum:  segNum,
		SegSize: segSize,
	}

	l.log = make([][]*Entry, l.SegNum)
	for i, _ := range l.log {
		l.log[i] = make([]*Entry, segSize, segSize)
	}

	return l
}

func (l *SegLog) Get(idx *LogIdx) (*Entry, error) {
	if idx.Seg < 0 || idx.Offset >= l.SegSize {
		return nil, fmt.Errorf("SegLog Get() idx = (%s) is out of range, segment capacity = (%d)",
			idx, l.SegSize)
	}

	if idx.Seg >= l.SegNum {
		// The segment is not allocated yet
		return nil, nil
	}

	return l.log[idx.Seg][idx.Offset], nil
}

func (l *SegLog) Put(idx *LogIdx, entry *Entry) error {
	if idx.Offset >= l.SegSize {
		return fmt.Errorf("SegLog Put() idx = (%s) is out of range, segment capacity = (%d)",
			idx, l.SegSize)
	}

	if idx.Seg >= l.SegNum {
		num := idx.Seg - l.SegNum + 1
		l.Extend(num)
	}

	l.log[idx.Seg][idx.Offset] = entry

	return nil
}

func (l *SegLog) IncIdx(idx *LogIdx) error {
	if idx.Offset+1 >= l.SegSize {
		idx.Offset = 0
		idx.Seg++ // the segment may not be allocated yet
		// overflow
		if idx.Seg < 0 {
			return fmt.Errorf("SegLog IncIdx() overflow. No more segments can be allocated.")
		}
		return nil
	}
	idx.Offset++
	return nil
}

func (l *SegLog) NextIdx(idx *LogIdx) (*LogIdx, error) {
	if idx.Offset+1 >= l.SegSize {
		ret := &LogIdx{idx.Seg + 1, 0} // the segment may not be allocated yet
		// overflow
		if ret.Seg < 0 {
			return nil, fmt.Errorf("SegLog NextIdx() overflow. No more segments can be allocated.")
		}
		return ret, nil
	}
	return &LogIdx{idx.Seg, idx.Offset + 1}, nil

}

func (l *SegLog) Extend(num int64) {
	l.log = append(l.log, make([][]*Entry, num)...)
	var i int64
	for i = 0; i < num; i++ {
		l.log[l.SegNum] = make([]*Entry, l.SegSize)
		l.SegNum++
	}
	//for ; l.SegNum < len(l.log); l.SegNum++ {
	//	l.log[l.SegNum] = make([]*Entry, l.SegSize)
	//}
}

func (l *SegLog) Size() string {
	return strconv.FormatInt(l.SegNum, 10) + "-" + strconv.FormatInt(l.SegSize, 10)
}
