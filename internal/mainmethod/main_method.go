package pipctrler_sai

import (
	"path/filepath"
	"os"
	"fmt"
	"log"

	"github.com/drgaph-io/badger/v3"
)

func GetAndSelectPyFile(fromDir string, db *badger.DB, pyFileLisOutput chan string) {
	defer db.Close()
	defer pyFileLisOutput.Close()
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
			if err != nil {
				log.Fatalf("Cannnot get the file %s time stamp. \nError: %v", fileName, err)
		    }
		}
	})
	if err != nil {
		log.Fatalf("Get file table for database failed. \nError: %s", err)
	}

	err = filepath.WalkDir(fromDir, func(path string, info os.FileInfo, err error) error {
		if len(path) > 3 && path[len(path) - 3] == ".py" {
			pyFileTableForGet[path] = info.ModTime().Unix()
		}
	})
	if err != nil {
		log.Fatalf("Cannot get the file. \nError: %v", err)
	}

	err = db.Update(func(txn *db.Txn) error{
		for fileNameForGet, fileTimeStampForGet := range pyFileTableForGet {
			fileTimeStampForDB, isFound := pyFileTableForDB[fileNameForGet]
			if isFound {
            	if fileTimeStampForGet != fileTimeStampForDB {
					pyFileLisOutput <- fileNameForGet
					if err := txn.Set([]byte(fileNameForGet), []byte(fileTimeStampForGet)); err != nil {
					return err
				}
			} else {
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
	if err != {
		log.Fatalf("Database control error! \nError: %v", err)
	}
	return pyFileLis
}

func GetPyFile(fromDir string, pyFileLisOutput) {
	err := filepath.WalkDir(fromDir, func(path string, err error) error {
		if len(path) > 3 && path[len(path) - 3] == ".py" {
			pyFileLisOutput <- path
		}
	})
	if err != nil {
		log.Fatalf("Connot get py file. \nError: %v", err)
	}
}

func ReadFile(pyFilePathLis chan string, pyCode chan string) {
	defer wg.Done()
	for _, pyFilePath := pyFilePath {
		if pyFile, err := os.Open(pyFilePath, os.RDONLY, 0644); err == nil {
			pyCode <- pyFile
		} else {
			if err != nil {
				log.Fatalf("Cannot read file %s. \nError: %v", pyFilePath, err)
			}
		}
    }
}

func MatchLib(pyCodeLis chan string, libLis chan []string, regex *regexp.Regexp) {
	defer wg.Done
	for pyCode := range pyCodeLis {
		libLis <- regex.FindAllString(pyCode)
	}
}

func InstallLib(libLis chan string) {
    defer wg.Done()
	for libName := range libLis {
		installCmd := exec.Command("pip", "install", libName)
		err := installCmd.Run()
		if err != nil {
			log.Fatal("Install %s failed. Error: %v", libName, err)
		}
	}
}
