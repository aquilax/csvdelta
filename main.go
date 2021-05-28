package main

import (
	"encoding/csv"
	"flag"
	"io"
	"os"
	"strconv"
	"strings"
)

type options struct {
	columns []int
}

func main() {
	var cFlag = flag.String("c", "", "columns to calculate delta on comma (separated and 0 indexed)")

	flag.Parse()

	columns, err := getColumns(*cFlag)
	if err != nil {
		panic(err)
	}

	in := os.Stdin
	out := os.Stdout
	err = process(options{columns}, in, out)
	if err != nil {
		panic(err)
	}
}

func process(o options, r io.Reader, w io.Writer) error {
	csvReader := csv.NewReader(r)
	csvWriter := csv.NewWriter(w)

	defer csvWriter.Flush()

	var err error
	var record []string
	var newRecord []string

	buffer := make([]string, len(o.columns))
	for j := 0; j < len(buffer); j++ {
		buffer[j] = "0"
	}

	for {
		record, err = csvReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		newRecord, err = getRecord(o, buffer, record)
		if err != nil {
			return err
		}

		if err = csvWriter.Write(newRecord); err != nil {
			return err
		}
	}
	return nil
}

func getRecord(o options, buffer []string, record []string) ([]string, error) {
	var colIndex int
	var bufferInder int
	var newValue string
	var err error
	for bufferInder, colIndex = range o.columns {
		newValue, err = getDiff(record[colIndex], buffer[bufferInder])
		buffer[bufferInder] = record[colIndex]
		record[colIndex] = newValue
	}
	return record, err
}

func getDiff(new, old string) (string, error) {
	var err error
	var newI, oldI int
	if newI, err = strconv.Atoi(new); err != nil {
		return "", err
	}
	if oldI, err = strconv.Atoi(old); err != nil {
		return "", err
	}
	return strconv.Itoa(newI - oldI), nil
}

func getColumns(c string) ([]int, error) {
	var err error
	var i, value int
	var v string
	if c == "" {
		return []int{}, nil
	}
	colStr := strings.Split(c, ",")
	columns := make([]int, len(colStr))
	for i, v = range colStr {
		value, err = strconv.Atoi(v)
		if err != nil {
			break
		}
		columns[i] = value
	}
	return columns, err
}
