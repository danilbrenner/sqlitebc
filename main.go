package main

import (
	_ "github.com/mattn/go-sqlite3"
	"log"
	"os"
	"sqlitebc/backup"
	storage2 "sqlitebc/storage"
	"strconv"
	"time"
)

type BackupSettings struct {
	azAccount    string
	azAccountKey string
	azContainer  string
	fileToBackup string
	directory    string
	capacity     int
}

const usageMessage = `
Configuration

Set the following environment variables to authenticate and store backups in Azure Blob Storage:

- AZ_ACCOUNT – Your Azure Storage account name.
- AZ_ACCOUNT_KEY – Your Azure Storage account key.
- AZ_CONTAINER – The Azure Blob Storage container to store backups.

Usage

Run the tool with the following arguments:

    sqlitebc <database_file_path> <azure_directory> [storage_capacity]

- <database_file_path> – Path to the SQLite database file.
- <azure_directory> – Directory in Azure Blob Storage where backups will be stored.
- [storage_capacity] - Optional. Number of files to keep in blob storage (default: 5). Must be an integer.

`

func getSettings() BackupSettings {
	settings := BackupSettings{capacity: 5}

	settings.azAccount = os.Getenv("AZ_ACCOUNT")
	if settings.azAccount == "" {
		log.Print(usageMessage)
		log.Fatal("Missing required environment variable: AZ_ACCOUNT")
	}

	settings.azAccountKey = os.Getenv("AZ_ACCOUNT_KEY")
	if settings.azAccountKey == "" {
		log.Print(usageMessage)
		log.Fatal("Missing required environment variable: AZ_ACCOUNT_KEY")
	}

	settings.azContainer = os.Getenv("AZ_CONTAINER")
	if settings.azContainer == "" {
		log.Print(usageMessage)
		log.Fatal("Missing required environment variable: AZ_CONTAINER")
	}

	if len(os.Args) < 3 || len(os.Args) > 4 {
		log.Print(usageMessage)
		log.Fatal("Missing or invalid arguments.")
	}

	settings.fileToBackup = os.Args[1]
	settings.directory = os.Args[2]
	if len(os.Args) == 4 {
		var err error
		settings.capacity, err = strconv.Atoi(os.Args[3])
		if err != nil {
			log.Print(usageMessage)
			log.Fatal(err)
		}
	}

	return settings
}

func main() {
	settings := getSettings()

	log.Println("Starting backup!")

	backupFileName := backup.ToBackupFileName(time.Now())

	log.Println("Backup filename: ", backupFileName)

	err := backup.PerformBackupDb(settings.fileToBackup, backupFileName)
	if err != nil {
		log.Fatal(err)
	}

	storage := storage2.NewAzureBackupStorage(settings.azAccount, settings.azAccountKey, settings.azContainer)
	err = storage.Upload(settings.directory, backupFileName)
	if err != nil {
		log.Fatal(err)
	}

	err = storage.AdjustCapacity(settings.directory, settings.capacity)
	if err != nil {
		log.Fatal(err)
	}

	err = os.Remove(backupFileName)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Backup complete!")
}
