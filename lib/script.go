package lib

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"
)

/*
 * Script is the specific up or down script.
 */
type Script struct {
	// Hash is the git commit hash for this migration
	Hash string

	Batch string

	// Path is the absolute path of the migration script
	Path string

	// MergedDate is the date the migration was committed to the repo
	MergedDate time.Time

	// CreateDate is the date the migration was created with the make command.
	// This is the date the is part of the directory where the scripts reside
	CreateDate time.Time

	Author    string
	direction int
}

/*
 * Execute runs a single migration script against the database.  If we are
 * executing an up script each script will have a row added to the goosey table
 * and if it's the last migration script the marker column will be set to true
 * to indicate the end of a batch.
 *
 * If we execute a down script, its corresponding row from the goosey table
 * will be removed.  If this is the last down script in the batch then the last
 * row in the goosey table will have its marker set to true.
 */
func (s Script) Execute(db *DB) error {
	script, err := ioutil.ReadFile(s.Path)
	if err != nil {
		return err
	}

	err = db.RunScript(string(script))
	if err != nil {
		// if there was a transaction and it failed then we need to rollback
		db.RunScript(`rollback`)
		//if !isErrorAcceptable(s.Path, err) {
		//    return fmt.Errorf("err: execute script %s %s: %s", s.Hash, s.Path, err)
		//}
	}

	if s.direction == Up {
		err = db.InsertLastMigration(s)
	} else {
		err = db.DeleteLastMigration(s.Hash)
	}
	return err
}

func isErrorAcceptable(file string, err error) bool {
	fmt.Println("")
	yellow(file)
	fmt.Println("")
	red(err.Error())
	fmt.Printf("\nwould you like to continue without applying this file? y/n ")
	reader := bufio.NewReader(os.Stdin)
	answer, err := reader.ReadString('\n')
	if err != nil {
		return false
	}
	if strings.ToLower(strings.TrimSpace(answer)) == "y" {
		return true
	}
	return false
}
