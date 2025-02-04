package backup

import (
	"database/sql"
	"fmt"
	"os"
	"time"
)

func ToBackupFileName(date time.Time) string {
	datePart := date.Format("05041502012006")
	backupFileName := fmt.Sprintf("%s.bcp", datePart)
	return backupFileName
}

func PerformBackupDb(fileToBackup string, backupFileName string) error {
	_, err := os.Stat(backupFileName)
	if !os.IsNotExist(err) {
		os.Remove(backupFileName)
	}

	db, err := sql.Open("sqlite3", fileToBackup)
	if err != nil {
		return err
	}
	defer db.Close()

	sql := fmt.Sprintf("VACUUM INTO '%s';", backupFileName)
	_, err = db.Exec(sql)
	if err != nil {
		return err
	}

	return nil
}
