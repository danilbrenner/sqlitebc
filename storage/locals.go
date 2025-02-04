package storage

import "time"

type backupItem struct {
	name     string
	modified time.Time
}

func partition(arr []backupItem) ([]backupItem, []backupItem) {
	if len(arr) == 0 {
		return []backupItem{}, []backupItem{}
	}

	pivot := arr[0]
	var left, right []backupItem
	for _, val := range arr[1:] {
		if pivot.modified.Before(val.modified) {
			left = append(left, val)
		} else {
			right = append(right, val)
		}
	}
	return append(left, pivot), right
}

func sortOutItemsToRemove(capacity int, blobs *[]backupItem, blobsToRemove []backupItem) []backupItem {
	if capacity > len(*blobs) {
		return blobsToRemove
	}
	if capacity < 1 {
		return *blobs
	}

	left, right := partition(*blobs)
	if len(left) == capacity {
		return append(blobsToRemove, right...)
	}
	if len(left) < capacity {
		return sortOutItemsToRemove(capacity-len(left), &right, blobsToRemove)
	}
	if len(left) > capacity {
		return sortOutItemsToRemove(capacity, &left, append(blobsToRemove, right...))
	}
	return blobsToRemove
}
