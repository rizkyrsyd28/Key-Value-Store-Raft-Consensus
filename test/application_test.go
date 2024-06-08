package test

import (
	. "github.com/Sister20/if3230-tubes-dark-syster/lib/app"
	"testing"
)

func TestKVPing(t *testing.T) {
	kv := NewKVStore()
	result := kv.Ping()
	if result != "PONG" {
		t.Fatalf("Pong err")
	}
}

func TestKVGetEmpty(t *testing.T) {
	kv := NewKVStore()
	result, err := kv.Get("any")
	if result != "" && err != nil {
		t.Fatalf("Err get when key doesn't exist")
	}
}

func TestKVSetGetCorrectKey(t *testing.T) {
	kv := NewKVStore()
	kv.Set("azero", "azimah")
	result, err := kv.Get("azero")
	if err != nil {
		t.Fatalf("Err get when key exist")
	}
	if result != "azimah" {
		t.Fatalf("Fail to set then get via correct key")
	}
}

func TestKVSetGetWrongKey(t *testing.T) {
	kv := NewKVStore()
	kv.Set("abdillah", "rizky")
	result, err := kv.Get("abdillah")
	if err != nil {
		t.Fatalf("Err get when key exist")
	}
	if result != "rizky" {
		t.Fatalf("Err set then get via wrong key")
	}
}

func TestKVChangeKey(t *testing.T) {
	kv := NewKVStore()
	kv.Set("rizky", "syaban")
	result, err := kv.Get("rizky")
	if err != nil {
		t.Fatalf("Err get value when key exist")
	}
	if result != "syaban" {
		t.Fatalf("Err set then get via wrong key")
	}
	kv.Set("rizky", "abdillah")
	result, err = kv.Get("rizky")
	if err != nil {
		t.Fatalf("Err get value when key exist")
	}
	if result != "abdillah" {
		t.Fatalf("Err set new value to a key")
	}
}

func TestKVSetThenDeleteCorrectly(t *testing.T) {
	kv := NewKVStore()
	kv.Set("ab", "cd")
	kv.Delete("ab")
	result, err := kv.Get("ab")
	if result != "" && err != nil {
		t.Fatalf("Error deleting data")
	}
	if result != "" {
		t.Fatalf("Error deleting data")
	}
}

func TestKVStrLen(t *testing.T) {
	kv := NewKVStore()
	kv.Set("ab", "yan")
	length, err := kv.Strlen("ab")
	if err != nil {
		t.Fatalf(err.Error())
	}
	if length != 3 {
		t.Fatalf("Err determining value length")
	}
}

func TestKVSetThenDeleteIncorrectly(t *testing.T) {
	kv := NewKVStore()
	kv.Delete("ab")
	result, err := kv.Get("ab")
	if result != "" && err != nil {
		t.Fatalf("Error deleting data")
	}
	if result != "" {
		t.Fatalf("Error deleting data")
	}
}

func TestKVAppend(t *testing.T) {
	kv := NewKVStore()
	kv.Set("miya", " BI")
	kv.Append("miya", "segsual")
	result, err := kv.Get("miya")
	if err != nil {
		t.Fatalf(err.Error())
	}
	if result != " BIsegsual" {
		t.Fatalf("Error appending value")
	}
}
