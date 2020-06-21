package dao

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"text/template"
	"time"
)

type genesis struct {
	Balances map[Account]uint `json:"balances"`
}

func loadGenesis(path string) (genesis, error) {
	content, err := ioutil.ReadFile(path)
	if err != nil {
		return genesis{}, err
	}

	var loadedGenesis genesis
	err = json.Unmarshal(content, &loadedGenesis)
	if err != nil {
		return genesis{}, err
	}

	return loadedGenesis, nil
}

func writeGenesisToDisk(genesisFile string) error {
	// get current working directory or error
	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("CWD error: %w", err)
	}

	tmplFiles := path.Join(cwd, "static", "tmpl", "*.*")

	// Parse the template directory
	// If we create more templates then this parse needs to be in an init
	// somewhere else
	tmpls, err := template.New("templates").ParseGlob(tmplFiles)
	if err != nil {
		return fmt.Errorf("Template parse error: %w", err)
	}
	f, err := os.Create(genesisFile)
	if err != nil {
		return fmt.Errorf("Cannot create gensis file '%s': %w", genesisFile, err)
	}
	err = tmpls.ExecuteTemplate(f, "genesis", map[string]interface{}{
		"genesisTime": time.Now().Format(time.RFC3339Nano),
		"chainId":     "the-refactored-blockchain-bar-ledger",
		"balances":    Balances{"andrej": 10000, "tim": 20000}, // Tim getting in early before the encryption stops him!!
	})
	f.Close()
	if err != nil {
		os.Remove(genesisFile)
		return fmt.Errorf("Template execution error: %w", err)
	}
	return nil
}
