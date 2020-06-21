package dao

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
)

func initDataDirIfNotExists(dataDir string) error {
	// Ensure the database directory exist
	if err := os.MkdirAll(getDatabaseDirPath(dataDir), os.ModePerm); err != nil {
		return fmt.Errorf("Error creating data directory: %w", err)
	}

	// Ensure the genesis file exists
	if !fileExist(getGenesisJsonFilePath(dataDir)) {
		if err := writeGenesisToDisk(getGenesisJsonFilePath(dataDir)); err != nil {
			return fmt.Errorf("Could not create genesis file: %w", err)
		}
	}

	// Ensure the Block file exists
	if !fileExist(getBlocksDbFilePath(dataDir)) {
		if err := writeEmptyBlocksDbToDisk(getBlocksDbFilePath(dataDir)); err != nil {
			return fmt.Errorf("Could not create empty block file: %w", err)
		}
	}

	return nil
}

func getDatabaseDirPath(dataDir string) string {
	return filepath.Join(dataDir, "db")
}

func getGenesisJsonFilePath(dataDir string) string {
	return filepath.Join(getDatabaseDirPath(dataDir), "genesis.json")
}

func getBlocksDbFilePath(dataDir string) string {
	return filepath.Join(getDatabaseDirPath(dataDir), "block.db")
}

func fileExist(filePath string) bool {
	//fmt.Print(filePath, " :")
	_, err := os.Stat(filePath)
	if err != nil && os.IsNotExist(err) {
		//fmt.Println("Missing")
		return false
	}
	//fmt.Println("Exists")
	return true
}

func dirExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}

	return true, err
}

func writeEmptyBlocksDbToDisk(path string) error {
	return ioutil.WriteFile(path, []byte(""), os.ModePerm)
}

// Read text file
func readTextFile(dir string, filename string) ([]byte, error) {
	// get current working directory or error
	cwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	//Determine the filepath
	textFile := path.Join(cwd, dir, filename)
	res, err := ioutil.ReadFile(textFile)
	if err != nil {
		return nil, fmt.Errorf("Error loading text file %s:%w", textFile, err)
	}
	return res, nil
}
