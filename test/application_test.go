package test

import (
	. "github.com/Sister20/if3230-tubes-dark-syster/lib"
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
	result := kv.Get("any")
	if result != "" {
		t.Fatalf("Err get when key doesn't exist")
	}
}

func TestKVSetGetCorrectKey(t *testing.T) {
	kv := NewKVStore()
	kv.Set("azero", "azimah")
	result := kv.Get("azero")
	if result != "azimah" {
		t.Fatalf("Fail to set then get via correct key")
	}
}

func TestKVSetGetWrongKey(t *testing.T) {
	kv := NewKVStore()
	kv.Set("abdillah", "rizky")
	result := kv.Get("abdillah")
	if result != "rizky" {
		t.Fatalf("Err set then get via wrong key")
	}
}

func TestKVChangeKey(t *testing.T) {
	kv := NewKVStore()
	kv.Set("rizky", "syaban")
	result := kv.Get("rizky")
	if result != "syaban" {
		t.Fatalf("Err set then get via wrong key")
	}
	kv.Set("rizky", "abdillah")
	result = kv.Get("rizky")
	if result != "abdillah" {
		t.Fatalf("Err set new value to a key")
	}
}

func TestKVSetThenDeleteCorrectly(t *testing.T) {
	kv := NewKVStore()
	kv.Set("ab", "cd")
	kv.Delete("ab")
	result := kv.Get("ab")
	if result != "" {
		t.Fatalf("Error deleting data")
	}
}

func TestKVStrLen(t *testing.T) {
	kv := NewKVStore()
	kv.Set("ab", "yan")
	length := kv.Strlen("ab")
	if(length != 3) {
		t.Fatalf("Err determining value length")
	}
}

func TestKVSetThenDeleteIncorrectly(t *testing.T) {
	kv := NewKVStore()
	kv.Delete("ab")
	result := kv.Get("ab")
	if result != "" {
		t.Fatalf("Error deleting data")
	}
}

func TestKVAppend(t *testing.T) {
	kv := NewKVStore()
	kv.Set("miya", " BI")
	kv.Append("miya", "segsual")
	result := kv.Get("miya")
	if( result != " BIsegsual") {
		t.Fatalf("Error appending value")
	}
}