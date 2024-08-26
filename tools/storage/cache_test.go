package storage

import (
	"testing"
)

func TestInsertСache(t *testing.T) {
	cache := CreateCacheObject()

	row1 := Telnum{
		Msisdn:     "79201112233",
		Region:     "msk",
		Abc:        "74951002030",
		Enabled:    true,
		ServiceKey: 3507,
	}
	key := cache.Create(row1)

	if key != MakeHash(row1) {
		t.Error("invalid create key")
	}
}

func TestCreateSelect(t *testing.T) {
	cache := CreateCacheObject()

	row1 := Telnum{
		Msisdn:     "79201112233",
		Region:     "msk",
		Abc:        "74951002030",
		Enabled:    true,
		ServiceKey: 3507,
	}
	key := cache.Create(row1)

	if key != MakeHash(row1) {
		t.Error("invalid create key")
	}

	row, ok := cache.Get(key)

	if !ok {
		t.Error("value not found")
	}

	if MakeHash(row1) != MakeHash(row) {
		t.Error("values do not match")
	}
}

func TestCreateDeleteGПри(t *testing.T) {
	cache := CreateCacheObject()

	row1 := Telnum{
		Msisdn:     "79201112233",
		Region:     "msk",
		Abc:        "74951002030",
		Enabled:    true,
		ServiceKey: 3507,
	}
	key := cache.Create(row1)

	if key != MakeHash(row1) {
		t.Error("invalid create key")
	}

	row, ok := cache.Get(key)

	if !ok {
		t.Error("value not found")
	}

	if MakeHash(row1) != MakeHash(row) {
		t.Error("values do not match")
	}
}
