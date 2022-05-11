package env

import (
	"fmt"
	"os"
	"strconv"
)

//Config struct holds data parsed from environment
type Config struct {
	//Add some field if i need to
	PostgresURI string
	Threads     int
}

//NewConfig parses env to struct
func NewConfig() (*Config, error) {
	// postgresql://[userspec@][hostspec][/dbname][?paramspec]
	postgresURI, ok := os.LookupEnv("POSTGRESURI")
	if !ok {
		return nil, fmt.Errorf("no POSTGRESURI env variable")
	}

	threads, ok := os.LookupEnv("THREADS")
	if !ok {
		return nil, fmt.Errorf("no THREADS env variable")
	}
	threadsInt, err := strconv.ParseInt(threads, 10, 32)
	if err != nil {
		return nil, fmt.Errorf("wrong THREADS variable format")
	}

	return &Config{
		PostgresURI: postgresURI,
		Threads:     int(threadsInt),
	}, nil
}
