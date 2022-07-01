package csvimport

import (
	"encoding/base64"
	"encoding/csv"
	"encoding/hex"
	"io"
	"os"
	"path/filepath"
)

//DHTCrawler describes needed struct that can process read infohash
type DHTCrawler interface {
	AddIH(ih string)
}

//ImportCSV writes to dhtc storage parsed infohashes
//that will be handled
func ImportCSV(dirName string, dhtc DHTCrawler) error {
	dirFile, err := os.Open(dirName)
	if err != nil {
		return err
	}
	files, err := dirFile.ReadDir(0)
	if err != nil {
		return err
	}
	fp, err := filepath.Abs(dirName)
	for _, v := range files {
		f, err := os.Open(filepath.Join(fp, v.Name()))
		if err != nil {
			return err
		}

		reader := csv.NewReader(f)
		reader.Comma = ';'
		reader.Comment = '#'
		for {
			record, err := reader.Read()
			if err == io.EOF {
				break
			} else if err != nil {
				return err
			}
			decoded, err := base64.StdEncoding.DecodeString(record[1]) //second element is info_hash
			if err != nil {
				return err
			}
			encodedString := hex.EncodeToString(decoded)
			//fmt.Println(encodedString)
			dhtc.AddIH(encodedString)
		}
	}
	return nil

}
