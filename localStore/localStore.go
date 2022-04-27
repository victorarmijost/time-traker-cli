package localStore

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/google/uuid"
)

const storePath = "local"

func init() {
	err := os.MkdirAll(storePath, os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}
}

func Save(record *Record) error {
	folder := getDateString(&record.Date)

	return saveToFile(folder, record)
}

func SaveToPool(record *Record) error {
	return saveToFile("pool", record)
}

func DeleteRecord(date *time.Time, name string) error {
	folder := getDateString(date)
	return os.Remove(getFilePath(folder, name))
}

func PourePool(date *time.Time) error {
	files, err := ListByStatus(nil, StatusPool)
	if err != nil {
		return err
	}

	if len(files) == 0 {
		return fmt.Errorf("nothing to poure")
	}

	newFolder := getDateString(date)

	err = os.MkdirAll(getBasePath(newFolder), os.ModePerm)
	if err != nil {
		return err
	}

	for _, name := range files {

		err = os.Rename(getFilePath("pool", name), getFilePath(newFolder, name))
		if err != nil {
			return err
		}

		file, err := Get(date, name)

		if err != nil {
			return err
		}

		file.Date = getDate(date)

		err = Save(file)

		if err != nil {
			return err
		}
	}

	return nil

}

func saveToFile(folder string, record *Record) error {

	err := os.MkdirAll(getBasePath(folder), os.ModePerm)
	if err != nil {
		return err
	}

	if record.Id == "" {
		id := uuid.New()
		record.Id = id.String()
	}

	fileName := getFilePath(folder, record.Id)

	data, err := json.MarshalIndent(record, "", "\t")
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(fileName, data, 0644)
	if err != nil {
		return err
	}

	return nil
}

func ListByStatus(date *time.Time, status string) ([]string, error) {
	var folder string

	if status == StatusPool {
		folder = "pool"
		status = StatusStored
	} else {
		folder = getDateString(date)
	}

	files, err := ioutil.ReadDir(getBasePath(folder))

	if errors.Is(err, os.ErrNotExist) {
		return []string{}, nil
	}

	if err != nil {
		return nil, err
	}

	list := []string{}
	for _, f := range files {
		fileName := f.Name()
		ext := filepath.Ext(fileName)
		name := strings.TrimSuffix(fileName, ext)
		commited := IsCommitted(name)
		if !f.IsDir() && ext == ".json" && (status == StatusStored || (status == StatusPending && !commited || status == StatusCommited && commited)) {
			list = append(list, name)
		}
	}

	return list, nil
}

func getDate(date *time.Time) time.Time {
	if date == nil {
		return time.Now()
	}
	return *date
}

func getDateString(date *time.Time) string {
	return getDate(date).Format("2006-01-02")
}

func Get(date *time.Time, name string) (*Record, error) {
	sdate := getDateString(date)
	return getRecord(sdate, name)
}

func GetFromPool(name string) (*Record, error) {
	return getRecord("pool", name)
}

func getRecord(folder, name string) (*Record, error) {
	fileName := getFilePath(folder, name)

	data, err := ioutil.ReadFile(fileName)
	if err != nil {
		return nil, err
	}

	record := &Record{}
	err = json.Unmarshal(data, record)
	if err != nil {
		return nil, err
	}

	if record.Id == "" {
		record.Id = name
	}

	return record, nil
}

func SetCommit(date *time.Time, name string) error {
	if IsCommitted(name) {
		return fmt.Errorf("record already commited")
	}

	sdate := getDateString(date)

	return os.Rename(getFilePath(sdate, name), getCommitedFilePath(sdate, name))
}

func IsCommitted(name string) bool {
	parts := strings.Split(name, ".")

	return len(parts) == 2
}

func GetAllByStatus(date *time.Time, status string) ([]*Record, error) {
	list := []*Record{}

	files, err := ListByStatus(date, status)

	if err != nil {
		return nil, err
	}

	for _, f := range files {
		record, err := func(f string) (*Record, error) {
			if status == StatusPool {
				return GetFromPool(f)
			} else {
				return Get(date, f)
			}
		}(f)

		if err != nil {
			return nil, err
		}

		list = append(list, record)
	}

	sort.Slice(list, func(i, j int) bool {
		return list[i].Date.After(list[j].Date)
	})

	return list, nil
}

func GetTimeByStatus(date *time.Time, status string) float32 {
	records, err := GetAllByStatus(date, status)

	if err != nil {
		return 0
	}

	Worked := float32(0.0)
	for _, r := range records {
		Worked += r.Hours
	}

	return Worked
}

func getFilePath(folder, name string) string {
	return fmt.Sprintf("%s/%s.json", getBasePath(folder), name)
}

func getCommitedFilePath(folder, name string) string {
	return fmt.Sprintf("%s/%s.com.json", getBasePath(folder), name)
}

func getBasePath(folder string) string {
	return fmt.Sprintf("%s/%s", storePath, folder)
}
