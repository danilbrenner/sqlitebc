package main

import (
	_ "github.com/mattn/go-sqlite3"
	"log"
	"os"
	"sqlitebc/backup"
	storage2 "sqlitebc/storage"
	"time"
)

func main() {
	log.Println("Starting backup!")

	azAccount := os.Getenv("AZ_ACCOUNT")
	if azAccount == "" {
		log.Fatal("AZ_ACCOUNT")
	}

	azAccountKey := os.Getenv("AZ_ACCOUNT_KEY")
	if azAccountKey == "" {
		log.Fatal("AZ_ACCOUNT_KEY")
	}

	azContainer := os.Getenv("AZ_CONTAINER")
	if azContainer == "" {
		log.Fatal("AZ_CONTAINER")
	}

	if len(os.Args) < 2 {
		log.Fatal("First argument:", os.Args[1])
	}

	fileToBackup := os.Args[0]
	directory := os.Args[1]

	backupFileName := backup.ToBackupFileName(time.Now())

	log.Println("Backup filename: ", backupFileName)

	err := backup.PerformBackupDb(fileToBackup, backupFileName)
	if err != nil {
		log.Fatal(err)
	}

	storage := storage2.NewAzureBackupStorage(azAccount, azAccountKey, backupFileName)
	err = storage.Upload(directory, backupFileName)
	if err != nil {
		log.Fatal(err)
	}

	err = storage.AdjustCapacity(directory, 5)
	if err != nil {
		log.Fatal(err)
	}

	err = os.Remove(backupFileName)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Backup complete!")
}
