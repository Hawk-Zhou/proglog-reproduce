package server

import (
	"reflect"
	"strconv"
	"testing"
)

func TestLog_Append_Read(t *testing.T) {

	log := NewLog()

	if _, err := log.Read(0); err != ErrOffsetNotFound {
		t.Error("should get err")
	}

	records := []Record{
		{
			Value: []byte("1"),
		},

		{
			Value: []byte("2"),
		},
	}

	log.Append(records[0])

	log.Append(records[1])

	records[0].Offset = 0
	records[1].Offset = 1

	for idx, dataRec := range records {
		t.Run(strconv.Itoa(idx), func(t *testing.T) {
			rec, err := log.Read(uint64(idx))

			if err != nil {
				t.Error("unexpected err")
			}

			if rec.Offset != dataRec.Offset {
				t.Errorf("wrong Offset expect: %d, got: %d", dataRec.Offset, rec.Offset)
			}

			if !reflect.DeepEqual(rec.Value, dataRec.Value) {
				t.Error("wrong Value")
			}
		})
	}

	if _, err := log.Read(2); err != ErrOffsetNotFound {
		t.Error("should get err")
	}
}
