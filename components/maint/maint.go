package maint

import (
	"io"
	"os"
	"strings"
	"time"

	//"github.com/aerth/ndjinn/components/database"
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
	_ = os.Mkdir("backups", 0700)
	err := os.Chmod("backups/latest.db", 0700)
	if err != nil {
		return false, err
	}
	err = copyFileContents("database/database.db", "backups/latest.db")
	if err != nil {
		return false, err
	}
	err = os.Chmod("backups/latest.db", 0600)
	if err != nil {
		return false, err
	}
	err = copyFileContents("database/database.db", "backups/database.db"+string(nowtime))
	if err != nil {
		return false, err
	}
	err = os.Chmod("backups/database.db"+string(nowtime), 0600)
	if err != nil {
		return false, err
	}
	return true, nil
}

// Bank check in
func BankTransfer() (bool, error) {
	return true, nil
}

// Bitcoin check in
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
	out, err := CreateFile(dst)
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
	_ = out.Chmod(0600)

	err = out.Sync()


	return
}

func CreateFile(name string) (*os.File, error) {
   		return os.OpenFile(name, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
}
