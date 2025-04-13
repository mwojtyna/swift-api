package main

import (
	"encoding/csv"
	"errors"
	"io"
	"log"
	"os"
)

var logger = log.New(os.Stderr, "[CSV PARSER] ", log.LstdFlags|log.Lshortfile)

func main() {
	if len(os.Args) != 2 {
		logger.Fatal("Error: CSV filename not specified")
	}

	csv_name := os.Args[1]
	err := parseAndInsertCSV(csv_name)
	if err != nil {
		logger.Fatalf(`Error parsing file '%s': "%s"`, csv_name, err.Error())
	}
}

func parseAndInsertCSV(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	reader := csv.NewReader(file)

	for {
		record, err := reader.Read()
		if err != nil {
			if !errors.Is(err, io.EOF) {
				return err
			} else {
				break
			}
		}

		logger.Println(record)
	}

	return nil
}
