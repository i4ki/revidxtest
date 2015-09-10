package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

type FieldType int

const (
	StringType FieldType = iota
	NumberType
)

func main() {
	if len(os.Args) != 4 {
		fmt.Fprintf(os.Stderr, "Usage: csv2json headerfile txtfile jsonfile\n")
		os.Exit(1)
	}

	var (
		headerFilename string = os.Args[1]
		txtFilename    string = os.Args[2]
		jsonFilename   string = os.Args[3]
		separator      string = ""
		headerFields   int    = 0
		headerNames    []string
		headerTypes    []FieldType
		txtLine        uint64 = 0
	)

	headerData, err := ioutil.ReadFile(headerFilename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Cannor read headerfile '%s': %s\n", headerFilename, err.Error())
		os.Exit(1)
	}

	headerLines := strings.Split(string(headerData), "\n")

	if len(headerLines) < 2 || headerLines[0][0] != ':' {
		fmt.Fprintf(os.Stderr, "Invalid headerfile '%s' format\n", headerFilename)
		os.Exit(1)
	}

	for _, headerTag := range strings.Split(headerLines[0][1:], " ") {
		switch headerTag {
		case "tsv":
			separator = "\t"
		case "csv":
			separator = ","
		}
	}

	if separator == "" {
		fmt.Fprintf(os.Stderr, "Missing separator type on headerfile '%s'\n", headerFilename)
		os.Exit(1)
	}

	headerNames = make([]string, len(headerLines)-1)
	headerTypes = make([]FieldType, len(headerLines)-1)

	for _, headerLine := range headerLines[1:] {
		parts := strings.Split(headerLine, ":")
		if len(parts) != 2 {
			fmt.Fprintf(os.Stderr, "Invalid header line %d, on headerfile '%s': %s\n",
				(headerFields + 1), headerFilename, headerLine)
			os.Exit(1)
		}

		headerNames[headerFields] = strings.TrimSpace(parts[0])
		fieldType := strings.ToLower(strings.TrimSpace(parts[1]))

		switch fieldType {
		case "string":
			headerTypes[headerFields] = StringType
		case "number":
			headerTypes[headerFields] = NumberType
		default:
			fmt.Fprintf(os.Stderr, "Invalid header type '%s', line %d on headerfile '%s'\n",
				fieldType, (headerFields + 1), headerFilename)
			os.Exit(1)
		}

		headerFields++
	}

	txtFile, err := os.Open(txtFilename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Cannot read txtfile '%s': %s\n", txtFilename, err.Error())
		os.Exit(1)
	}
	defer txtFile.Close()

	jsonFile, err := os.Create(jsonFilename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Cannot create jsonfile '%s': %s\n", jsonFilename, err.Error())
		os.Exit(1)
	}
	defer jsonFile.Close()

	txtScanner := bufio.NewScanner(txtFile)

	for txtScanner.Scan() {
		parts := strings.Split(txtScanner.Text(), separator)
		if len(parts) != headerFields {
			fmt.Fprintf(os.Stderr, "Invalid column count %d on line %d of txtfile '%s'\n",
				len(parts), (txtLine + 1), txtFilename)
			os.Exit(1)
		}

		data := map[string]interface{}{}

		for i := 0; i < headerFields; i++ {
			value := strings.TrimSpace(parts[i])
			if value == "" {
				continue
			}

			switch headerTypes[i] {
			case StringType:
				data[headerNames[i]] = value
			case NumberType:
				number, err := strconv.ParseFloat(value, 64)
				if err != nil {
					fmt.Fprintf(os.Stderr, "Invalid number '%s' on line %d of txtfile '%s': %s\n",
						parts[i], (txtLine + 1), txtFilename, err.Error())
					os.Exit(1)
				}
				data[headerNames[i]] = number
			}
		}

		jsonData, err := json.Marshal(data)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Cannot encode json %+v on line %d of txtfile '%s'\n",
				data, (txtLine + 1), txtFilename)
			os.Exit(1)
		}

		_, err = jsonFile.Write(jsonData)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Cannot write jsonfile '%s'\n", jsonFilename, err.Error())
			os.Exit(1)
		}

		_, err = jsonFile.WriteString("\n")
		if err != nil {
			fmt.Fprintf(os.Stderr, "Cannot write jsonfile '%s'\n", jsonFilename, err.Error())
			os.Exit(1)
		}

		txtLine++
	}

	fmt.Printf("%d lines processed from '%s'\n", txtLine, txtFilename)

	err = txtScanner.Err()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error readind txtfile '%s': %s\n", txtFilename, err.Error())
		os.Exit(1)
	}
}
