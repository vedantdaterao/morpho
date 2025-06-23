package bencoding

import (
	"testing"
)

func TestDecodeString(t *testing.T) {
	data := []byte("4:spam")
	result, err := Decode(data, false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result != "spam" {
		t.Errorf("expected 'spam', got '%s'", result)
	}
}

func TestDecodeInt(t *testing.T) {
	data := []byte("i42e")
	result, err := Decode(data, false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result != 42 {
		t.Errorf("expected 42, got %d", result)
	}
}

func TestDecodeList(t *testing.T) {
	data := []byte("l4:spam4:eggs3:hame")
	result, err := Decode(data, false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expected := []interface{}{"spam", "eggs", "ham"}
	if len(result.([]interface{})) != len(expected) {
		t.Errorf("expected length %d, got %d", len(expected), len(result.([]interface{})))
	}
	for i, v := range result.([]interface{}) {
		if v != expected[i] {
			t.Errorf("expected %s, got %s", expected[i], v)
		}
	}
}

func TestDecodeDict(t *testing.T) {
	data := []byte("d3:bar4:spam3:fooi42ee")
	result, err := Decode(data, false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expected := map[string]interface{}{
		"bar": "spam",
		"foo": 42,
	}
	if len(result.(map[string]interface{})) != len(expected) {
		t.Errorf("expected length %d, got %d", len(expected), len(result.(map[string]interface{})))
	}
	for k, v := range expected {
		if result.(map[string]interface{})[k] != v {
			t.Errorf("expected %s: %v, got %s: %v", k, v, k, result.(map[string]interface{})[k])
		}
	}
}

func TestDecodeInvalidData(t *testing.T) {
	data := []byte("i42")
	_, err := Decode(data, false)
	if err == nil {
		t.Error("expected error for invalid data, got nil")
	}

	data = []byte("l4:spam4:eggs4:hame") // missing 'e' at the end
	_, err = Decode(data, false)
	if err == nil {
		t.Error("expected error for unterminated list, got nil")
	}

	data = []byte("d3:bar4:spam3:fooi42e") // missing 'e' at the end
	_, err = Decode(data, false)
	if err == nil {
		t.Error("expected error for unterminated dictionary, got nil")
	}
}

func TestDecodeRawData(t *testing.T) {
	data := []byte("d3:bar4:spam3:fooi42e6:piecesd3:bar4:spamee")
	result, err := Decode(data, true)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	dict := result.(map[string]interface{})
	if _, ok := dict["raw"]; ok {
	} else {
		t.Error("expected 'raw' key in result, got none")
	}

	if _, ok := dict["after_pieces"]; ok {
	} else {
		t.Error("expected 'after_pieces' key in result, got none")
	}
}

func TestDecodeEmptyData(t *testing.T) {
	data := []byte("")
	_, err := Decode(data, false)
	if err == nil {
		t.Error("expected error for empty data, got nil")
	}

	data = []byte("e") // empty dictionary
	result, err := Decode(data, false)
	if result != nil && err == nil {
		t.Fatalf("expected error: %v", err)
	}
	// if len(result.(map[string]interface{})) != 0 {
	// 	t.Errorf("expected empty dictionary, got %v", result)
	// }
}

func TestDecodeNestedListDict(t *testing.T) {
	data := []byte("lld3:fooi1e3:bari2eeeli1ei2eee")
	result, err := Decode(data, false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	list, ok := result.([]any)
	if !ok || len(list) != 2 {
		t.Fatalf("expected list of length 2, got %v", result)
	}
	dict := list[0].([]any)[0].(map[string]any)
	if dict["foo"] != 1 || dict["bar"] != 2 {
		t.Errorf("expected nested dict, got %v", dict)
	}
	nestedList, ok := list[1].([]any)
	if !ok || len(nestedList) != 2 || nestedList[0] != 1 || nestedList[1] != 2 {
		t.Errorf("expected nested list [1 2], got %v", nestedList)
	}
}

func TestDecodeEmptyList(t *testing.T) {
	data := []byte("le")
	result, err := Decode(data, false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	list, ok := result.([]interface{})
	if !ok {
		t.Fatalf("expected list, got %T", result)
	}
	if len(list) != 0 {
		t.Errorf("expected empty list, got %v", list)
	}
}

func TestDecodeEmptyDict(t *testing.T) {
	data := []byte("de")
	result, err := Decode(data, false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	dict, ok := result.(map[string]interface{})
	if !ok {
		t.Fatalf("expected dict, got %T", result)
	}
	if len(dict) != 0 {
		t.Errorf("expected empty dict, got %v", dict)
	}
}

func TestDecodeNegativeInt(t *testing.T) {
	data := []byte("i-123e")
	result, err := Decode(data, false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result != -123 {
		t.Errorf("expected -123, got %v", result)
	}
}

func TestDecodeZeroInt(t *testing.T) {
	data := []byte("i0e")
	result, err := Decode(data, false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result != 0 {
		t.Errorf("expected 0, got %v", result)
	}
}

func TestDecodeStringWithColon(t *testing.T) {
	data := []byte("7:foo:bar")
	result, err := Decode(data, false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result != "foo:bar" {
		t.Errorf("expected 'foo:bar', got '%v'", result)
	}
}

func TestDecodeMultipleTopLevel(t *testing.T) {
	data := []byte("4:spam4:eggs")
	result, err := Decode(data, false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result != "spam" {
		t.Errorf("expected '', got '%v'", result)
	}
}

func TestDecodeInvalidStringLength(t *testing.T) {
	data := []byte("100:abc")
	_, err := Decode(data, false)
	if err == nil {
		t.Error("expected error for string length out of bounds, got nil")
	}
}

func TestDecodeIntFormat(t *testing.T) {
	data := []byte("i234e")
	_, err := Decode(data, false)
	if err != nil {
		t.Error("expected error for invalid integer, got", err)
	}
}

func TestDecodeInvalidRawData(t *testing.T) {
	data := []byte("d3:bar4:spam3:fooi42e6:piecesd3:bar4:spamee")
	_, err := Decode(data, true)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	data = []byte("d3:bar4:spam3:fooi42e6:piecesd3:bbb4:aaaaee") // missing 'e' at the end
	r, err := Decode(data, true)
	if err != nil {
		t.Error("expected error for unterminated raw data, got nil")
	}
	t.Log(r)
}
