package server

import (
	"bytes"
	"encoding/json"
	"net/http"
	"reflect"
	"testing"
	"time"
)

// marshals a struct and submit it to sever with given method
// unmarshal it back to ret (pointer to target struct)
func marshalAndSubmit(data interface{}, ret interface{}, method string, t *testing.T) {
	payloadBytes, err := json.Marshal(data)
	if err != nil {
		t.Fatal(err)
	}
	body := bytes.NewReader(payloadBytes)

	req, err := http.NewRequest(method, "http://127.0.0.1:8080", body)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(ret)
	if err != nil {
		t.Fatal(err)
	}
}

// only tests normal inputs, abnormal cases are omitted
func TestIntegration(t *testing.T) {
	srv := NewHTTPServer(":8080")

	echan := make(chan error, 1)
	go func() {
		err := srv.ListenAndServe()
		if err != nil {
			echan <- err
		}
	}()

	select {
	case err := <-echan:
		t.Fatal(err)
	case <-time.After(500 * time.Millisecond):
	}

	produceData := []ProduceRequest{
		{
			Record: Record{
				Value: []byte("114"),
			},
		},

		{
			Record: Record{
				Value: []byte("514"),
			},
		},
	}

	expectedOffsets := []uint64{0, 1}

	for idx, d := range produceData {
		var respStruct ProduceResponse
		marshalAndSubmit(d, &respStruct, "POST", t)
		if respStruct.Offset != expectedOffsets[idx] {
			t.Fatal("wrong return val")
		}
	}

	consumeData := []ConsumeRequest{
		{
			Offset: &expectedOffsets[0],
		},
		{
			Offset: &expectedOffsets[1],
		},
	}

	for idx, d := range consumeData {
		var respStruct ConsumeResponse
		marshalAndSubmit(d, &respStruct, "GET", t)

		if respStruct.Record.Offset != expectedOffsets[idx] {
			t.Fatal("wrong return idx")
		}

		if !reflect.DeepEqual(respStruct.Record.Value, produceData[idx].Record.Value) {
			t.Fatalf("wrong return val, expect %q, got %q", produceData[idx].Record.Value, respStruct.Record.Value)
		}
	}
}
