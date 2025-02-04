package storage

import (
	"sort"
	"testing"
	"time"
)

func TestPartition(t *testing.T) {
	now := time.Now()

	items := []backupItem{
		{name: "A", modified: now.Add(-3 * time.Hour)},
		{name: "B", modified: now.Add(-1 * time.Hour)},
		{name: "C", modified: now.Add(-2 * time.Hour)},
		{name: "D", modified: now.Add(-4 * time.Hour)},
	}

	left, right := partition(items)

	for _, item := range left {
		if item != items[0] && !item.modified.After(items[0].modified) {
			t.Errorf("Expected left partition items to be after pivot, but got %v", item.modified)
		}
	}

	for _, item := range right {
		if !items[0].modified.After(item.modified) {
			t.Errorf("Expected right partition items to be before or equal to pivot, but got %v", item.modified)
		}
	}
}

func sameElements(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	sort.Strings(a)
	sort.Strings(b)
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func TestSortOutItemsToRemove(t *testing.T) {
	now := time.Now()

	blobs := []backupItem{
		{name: "A", modified: now.Add(-3 * time.Hour)},
		{name: "B", modified: now.Add(-1 * time.Hour)},
		{name: "C", modified: now.Add(-2 * time.Hour)},
		{name: "D", modified: now.Add(-4 * time.Hour)},
		{name: "E", modified: now.Add(-5 * time.Hour)},
	}

	tests := []struct {
		name           string
		capacity       int
		expectedRemove []string
	}{
		{"Keep all blobs", 5, []string{}},
		{"Remove one blob", 4, []string{"E"}},
		{"Remove two blobs", 3, []string{"E", "D"}},
		{"Remove three blobs", 2, []string{"E", "D", "A"}},
		{"Remove four blobs", 1, []string{"E", "D", "A", "C"}},
		{"Remove all blobs", 0, []string{"E", "D", "A", "C", "B"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			blobsCopy := append([]backupItem{}, blobs...)
			removed := sortOutItemsToRemove(tt.capacity, &blobsCopy, []backupItem{})

			var removedNames []string
			for _, item := range removed {
				removedNames = append(removedNames, item.name)
			}

			if !sameElements(removedNames, tt.expectedRemove) {
				t.Errorf("Expected removed items %v, but got %v", tt.expectedRemove, removedNames)
			}
		})
	}
}
