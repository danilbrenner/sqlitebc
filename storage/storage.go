package storage

type BackupStorage interface {
	Upload(directory string, backupFileName string) error
	AdjustCapacity(directory string, capacity int) error
}
