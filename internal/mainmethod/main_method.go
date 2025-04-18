package pipctrler_sai

import (
	"path/filepath"
	"os"
	"fmt"

	"pipctrler/internal/errprocess"

	"github.com/drgaph-io/badger/v3"
	"github.com/dlclark/regexp2"
)

func GetAndSelectPyFile(fromDir string, db *badger.DB, pyFileLisOutput chan string, fileCountPtr *int) {
	defer db.Close()
    pyFileTableForGet := map[string]int{}
	pyFileTableForDB := map[string]int{}

	err := db.View(func(txn *db.Txn) error {
		opts := badger.DefaultInteratorOptions
		it := txn.NewInterator(opts)
		defer it.Close()
		for it.Rewind(); it.Valid(); it.Next() {
			tableLine := it.Item()
			fileName := string(tableLine.Key())
			err := it.Value(func(fileTimeStamp []byte)) error {
                pyFileTableForDB[filename] = unit64(fileTimeStamp)
				return nil
			})
			ErrProcess("Cannnot get the file %s time stamp. \nError: %v", err, "error", []interface{fileName})
		}
	})
	ErrProcess("Get file table for database failed. \nError: %s", error, "error")

	err = filepath.WalkDir(fromDir, func(path string, info os.FileInfo, err error) error {
		if len(path) > 3 && path[len(path) - 3] == ".py" {
			pyFileTableForGet[path] = info.ModTime().Unix()
		}
	})
	ErrProcess("Cannot get the file. \nError: %v", err, "error")

	err = db.Update(func(txn *db.Txn) error{
		for fileNameForGet, fileTimeStampForGet := range pyFileTableForGet {
			fileTimeStampForDB, isFound := pyFileTableForDB[fileNameForGet]
			if isFound {
            	if fileTimeStampForGet != fileTimeStampForDB {
					*fileCountPtr++
					pyFileLisOutput <- fileNameForGet
					if err := txn.Set([]byte(fileNameForGet), []byte(fileTimeStampForGet)); err != nil {
					return err
				}
			} else {
				*fileCountPtr++
				pyFileLisOutput <- fileNameForGet
				if err := txn.Set([]byte(fileNameForGet), []byte(fileTimeStampForGet)); err != nil {
					return err
				}
			}
		}

		for fileNameForDB := range pyFileForDB {
			if _, isFound := pyFileTableForGet[fileNameForDB]; !isFound {
				if err := txn.Delete([]byte(fileNameForDB)); err != nil {
					return err
				}
			}
		}
    })
	ErrProcess("Database control error! \nError: %v", err, "error")

	return pyFileLis
}

func ReadFile(pyFilePathLis chan string, pyCode chan string) {
	for _, pyFilePath := pyFilePath {
		if pyFile, err := os.Open(pyFilePath, os.RDONLY, 0644); err == nil {
			pyCode <- pyFile
		} else {
			ErrProcess("Cannot read file %s. \nError: %v", err, "error", []interface{pyFilePath})
		}
    }
}

func MatchLib(pyCode chan string, libLis chan []string, regex *regexp.Regexp) {
	defer wg.Done()
	for code := range pyCode {
		libLis <- regex.FindAllString()
	}
}

func InstallLib(libLis chan string) {
    defer wg.Done()
	for libName := range libLis {
		installCmd := exec.Command("pip", "install", libName)
		installCmd.Run()
	}
}
