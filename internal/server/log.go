// Q1: why offset is also kept inside of Record (btw it's set when appended)

package server

import (
	"errors"
	"sync"
)

type Record struct {
	Value  []byte `json:"value"`
	Offset uint64 `json:"offset"`
}

// Log is thread-safe
type Log struct {
	mu      sync.Mutex
	records []Record
}

// Offset requested exceeds max offset number.
var ErrOffsetNotFound = errors.New("offset not found")

func NewLog() *Log {
	return &Log{records: make([]Record, 0)}
}

// Append a record to the Log. Record's offset will be set.
func (l *Log) Append(record Record) (offset uint64, err error) {
	l.mu.Lock()
	defer l.mu.Unlock()

	record.Offset = uint64(len(l.records))
	l.records = append(l.records, record)

	return record.Offset, nil
}

// Read the record of a offset, raise error if record not exist.
func (l *Log) Read(offset uint64) (Record, error) {
	l.mu.Lock()
	defer l.mu.Unlock()

	if offset >= uint64(len(l.records)) {
		return Record{}, ErrOffsetNotFound
	}

	return l.records[offset], nil
}
