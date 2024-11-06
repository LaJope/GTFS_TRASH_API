package main

import (
	"archive/zip"
	"bytes"
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"slices"
	"strconv"
	"strings"
)

const (
	FEED_DIR = "./gtfs_csv"
)

// func getGTFSFeed() {
func getGTFSFeed(w http.ResponseWriter, req *http.Request) {
	_, err := os.Stat(FEED_DIR)
	if errors.Is(err, os.ErrNotExist) {
		_ = os.Mkdir(FEED_DIR, os.ModePerm)
	}

	feedZip := getFeedZipArchive()

	reader := bytes.NewReader(feedZip)
	zipReader, err := zip.NewReader(reader, int64(len(feedZip)))
	if err != nil {
		log.Fatal(err)
	}

	for _, f := range zipReader.File {
		path := filepath.Join(FEED_DIR, f.Name)

		if f.FileInfo().IsDir() {
			_ = os.MkdirAll(path, os.ModePerm)
			continue
		}

		if err := os.MkdirAll(filepath.Dir(path), os.ModePerm); err != nil {
			log.Fatal(err)
		}

		dstFile, err :=
			os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			log.Fatal(err)
		}

		fileInArchive, err := f.Open()
		if err != nil {
			log.Fatal(err)
		}

		if _, err := io.Copy(dstFile, fileInArchive); err != nil {
			log.Fatal(err)
		}

		dstFile.Close()
		fileInArchive.Close()
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode("Hello")
}

func parseCSVFiles(fileName string, ignorePrefix string,
	columns ...string) []map[string]interface{} {

	file, err := os.Open(fmt.Sprintf("%s/%s", FEED_DIR, fileName))
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := file.Close(); err != nil {
			log.Fatal(err)
		}
	}()

  csvReader := csv.NewReader(file)


  colInds := make([]int, 0)
  colStrings := make([]string, 0)

  firstRow, err := csvReader.Read()
  if err != nil {
    log.Fatal(err)
  }
  for ind, elem := range firstRow {
    if slices.Index(columns, elem) != -1 {
      colInds = append(colInds, ind)
      colStrings = append(colStrings, strings.TrimPrefix(elem, ignorePrefix))
    }
  }

  data, err := csvReader.ReadAll()
  if err != nil {
    log.Fatal(err)
  }

  csvData := make([]map[string]interface{}, len(data))

  for rowInd, row := range data {
    elemInfo := make(map[string]interface{})

    for ind, elem := range row {
      colIndex := slices.Index(colInds, ind)
      if colIndex == -1 {
        continue
      }
      var elemData interface{} = elem
      if tmp, err := strconv.ParseFloat(elem, 32); err == nil {
        elemData = tmp
      }
      if tmp, err := strconv.ParseInt(elem, 10, 0); err == nil {
        elemData = tmp
      }
      elemInfo[colStrings[colIndex]] = elemData
    }
    csvData[rowInd] = elemInfo
  }

	return csvData
}
