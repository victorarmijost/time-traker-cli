package localStore

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
	"varmijo/time-tracker/bairestt"

	"github.com/google/uuid"
)

const storePath = "local"

func init() {
	err := os.MkdirAll(storePath, os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}
}

func Save(record *bairestt.TimeRecord) error {
	sdate := getDate(&record.Date)
	path := storePath + "/" + sdate

	return saveToFile(path, record)
}

func SaveToPool(record *bairestt.TimeRecord) error {
	path := storePath + "/pool"

	return saveToFile(path, record)
}

func saveToFile(path string, record *bairestt.TimeRecord) error {
	err := os.MkdirAll(path, os.ModePerm)
	if err != nil {
		return err
	}

	id := uuid.New()

	fileName := fmt.Sprintf("%s/%s.json", path, id)

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
	var path string

	if status == StatusPool {
		path = storePath + "/pool"
	} else {
		sdate := getDate(date)
		path = storePath + "/" + sdate
	}

	files, err := ioutil.ReadDir(path)

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

func getDate(date *time.Time) string {
	if date == nil {
		return time.Now().Format("2006-01-02")
	}
	return date.Format("2006-01-02")
}

func Get(date *time.Time, name string) (*bairestt.TimeRecord, error) {
	sdate := getDate(date)
	path := storePath + "/" + sdate

	ext := ".json"
	fileName := path + "/" + name + ext

	data, err := ioutil.ReadFile(fileName)
	if err != nil {
		return nil, err
	}

	record := &bairestt.TimeRecord{}
	err = json.Unmarshal(data, record)
	if err != nil {
		return nil, err
	}

	return record, nil
}

func SetCommit(date *time.Time, name string) error {
	if IsCommitted(name) {
		return fmt.Errorf("record already commited")
	}

	sdate := getDate(date)
	path := storePath + "/" + sdate
	newName := name + ".com.json"

	name = name + ".json"

	return os.Rename(path+"/"+name, path+"/"+newName)
}

func IsCommitted(name string) bool {
	parts := strings.Split(name, ".")

	return len(parts) == 2
}

func GetAllByStatus(date *time.Time, status string) ([]*bairestt.TimeRecord, error) {
	list := []*bairestt.TimeRecord{}

	files, err := ListByStatus(date, status)

	if err != nil {
		return nil, err
	}

	for _, f := range files {
		record, err := Get(date, f)

		if err != nil {
			return nil, err
		}

		list = append(list, record)
	}

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
