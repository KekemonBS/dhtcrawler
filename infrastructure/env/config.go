package env

import (
	"fmt"
	"os"
)

//Config struct holds data parsed from environment
type Config struct {
	//Add some field if i need to
	PostgresURI string
	ImportCSV   bool
}

//NewConfig parses env to struct
func NewConfig() (*Config, error) {
	// postgresql://[userspec@][hostspec][/dbname][?paramspec]
	postgresURI, ok := os.LookupEnv("POSTGRESURI")
	if !ok {
		return nil, fmt.Errorf("no POSTGRESURI env variable")
	}

	//[YES/NO]
	importCSVAns, ok := os.LookupEnv("IMPORTCSV")
	if !ok {
		return nil, fmt.Errorf("no IMPORTCSV env variable")
	}
	var importCSV bool
	if importCSVAns == "YES" {
		importCSV = true
	} else {
		importCSV = false
	}

	return &Config{
		PostgresURI: postgresURI,
		ImportCSV:   importCSV,
	}, nil
}
