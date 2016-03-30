package maint

import (
	"io"
	"os"
	"strings"
	"time"

	"github.com/boltdb/bolt"
	"github.com/yosssi/boltstore/reaper"
)

var db *bolt.DB

// Run Maint Schedule
func RunSchedule() (bool, error) {
	ok, err := BackupDB()
	defer reaper.Quit(reaper.Run(db, reaper.Options{}))
	return ok, err
}

// Close down and backup our DB
func BackupDB() (bool, error) {
	now := time.Now()
	now.Format(time.UnixDate)
	nowtime := strings.Replace(now.String(), " ", "", -1)
	err := copyFileContents("database.db", "backups/latest.db")
	if err != nil {
		return false, err
	}
	err = copyFileContents("database.db", "backups/database.db"+string(nowtime))
	if err != nil {
		return false, err
	}
	return true, nil
}

// Transfer from CC to Bank
func BankTransfer() (bool, error) {
	return true, nil
}

// Sync Bitcoin account with Bank and sync Blockchain
func BitcoinSync() (bool, error) {
	return true, nil
}

// Backup function
// TODO: make better.
func copyFileContents(src, dst string) (err error) {
	in, err := os.Open(src)

	if err != nil {
		return
	}
	defer in.Close()
	out, err := os.Create(dst)
	if err != nil {
		return
	}
	defer func() {
		cerr := out.Close()
		if err == nil {
			err = cerr
		}
	}()
	if _, err = io.Copy(out, in); err != nil {
		return
	}
	err = out.Sync()
	return
}
