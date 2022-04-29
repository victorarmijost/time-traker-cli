package bairestt

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
	"varmijo/time-tracker/utils"
)

type Bairestt struct {
	client *client
	email  string
	cache  *cache
}

type cache struct {
	Token string
	Tasks []TaskInfo
}

const tmpPath = ".tmp"
const employessPortalUrl = "https://employees.bairesdev.com"
const cacheFileName = "tt_cahe.json"

func NewService(email string) *Bairestt {
	return &Bairestt{
		client: newClient(employessPortalUrl),
		email:  email,
		cache: &cache{
			Token: "",
			Tasks: []TaskInfo{},
		},
	}
}

func init() {
	err := os.MkdirAll(utils.GeAppPath(tmpPath), os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}
}

func (t *Bairestt) saveCache() error {
	cacheData, err := json.Marshal(t.cache)
	if err != nil {
		return err
	}

	return os.WriteFile(utils.GeAppPath(tmpPath)+"/"+cacheFileName, cacheData, 0644)

}

func (t *Bairestt) loadCache() error {
	cacheData, err := os.ReadFile(utils.GeAppPath(tmpPath) + "/" + cacheFileName)

	if err != nil {
		return err
	}

	err = json.Unmarshal(cacheData, t.cache)

	if err != nil {
		return err
	}

	return nil
}

func (t *Bairestt) StartWithPass(ctx context.Context, pass string) error {
	token, err := t.emulate_login(ctx, pass)

	if err != nil {
		return fmt.Errorf("login fail, %w", err)
	}

	t.client.SetToken(token)

	err = t.refresh(ctx)

	if err != nil {
		return fmt.Errorf("error validating token, %w", err)
	}

	go t.autoRefresh()

	return nil
}

//Login and validates login token, start auto refresh process
func (t *Bairestt) Start(ctx context.Context) error {
	err := t.loadCache()

	if err != nil {
		return err
	}

	t.client.SetToken(t.cache.Token)
	_ = t.refresh(ctx)

	if !t.IsActive() {
		return fmt.Errorf("token expired")

	}

	go t.autoRefresh()

	return nil

}

const refreshPath = "/api/apsi/auth/employees/refresh"

//Calls to the refresh service, with generates a new token
func (t *Bairestt) refresh(ctx context.Context) error {
	res, err := t.client.Post(refreshPath, struct{}{})

	if err != nil {
		return fmt.Errorf("error sending refresh, %w", err)
	}

	if res.StatusCode != http.StatusOK {
		t.client.SetToken("")
		return fmt.Errorf("refresh fail with code: %d", res.StatusCode)
	}

	token := res.Header.Get("x-bairesdev-token")

	t.client.SetToken(token)

	if token == "" {
		return fmt.Errorf("missing refresh token")
	}

	t.cache.Token = token
	err = t.saveCache()

	if err != nil {
		return err
	}

	return nil
}

func (t *Bairestt) autoRefresh() {
	for {
		time.Sleep(30 * time.Minute)

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		err := t.refresh(ctx)

		if err != nil {
			log.Println("Token expired, login again")
			break
		}
	}

	os.Exit(0)
}

func (t *Bairestt) IsActive() bool {
	return t.client.IsToken()
}

const getTasksUrl = "/api/v1/employees/taskdescriptions"

func (t *Bairestt) GetTasks(ctx context.Context) ([]TaskInfo, error) {
	if len(t.cache.Tasks) > 0 {
		return t.cache.Tasks, nil
	}

	fm, lm := getFirstAndLastOfMonth()

	req := map[string]time.Time{
		"fromDate": fm,
		"toDate":   lm,
	}

	res, err := t.client.Put(getTasksUrl, req)

	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("get tasks fail with code: %d", res.StatusCode)
	}

	body, err := ioutil.ReadAll(res.Body)

	if err != nil {
		return nil, err
	}

	tasks := &TasksResponse{}
	err = json.Unmarshal(body, tasks)

	if err != nil {
		return nil, err
	}

	t.cache.Tasks = tasks.Data

	return t.cache.Tasks, nil
}

func getFirstAndLastOfMonth() (time.Time, time.Time) {
	now := time.Now()
	currentYear, currentMonth, _ := now.Date()
	currentLocation := now.Location()

	firstOfMonth := time.Date(currentYear, currentMonth, 1, 0, 0, 0, 0, currentLocation)
	lastOfMonth := firstOfMonth.AddDate(0, 1, -1)

	return firstOfMonth, lastOfMonth
}

const addRecordPath = "/api/v1/employees/timetracker-record-upsert"

func (t *Bairestt) AddRecord(ctx context.Context, record *TimeRecord) (*RecordResponse, error) {
	//Just take the date in the time record
	record.Date = adjustTime(record.Date)

	res, err := t.client.Put(addRecordPath, record)

	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("add record fail with code: %d", res.StatusCode)
	}

	body, err := ioutil.ReadAll(res.Body)

	if err != nil {
		return nil, fmt.Errorf("error reading response")
	}

	respRecord := &RecordResponse{}
	err = json.Unmarshal(body, respRecord)

	if err != nil {
		return nil, fmt.Errorf("error unmarshalling response,%w", err)
	}

	return respRecord, nil
}

func adjustTime(it time.Time) time.Time {
	ot, _ := time.Parse("2006-01-02", it.Format("2006-01-02"))

	return ot
}

const getFocalUrl = "/api/v1/employees/focalpoints"

func (t *Bairestt) GetFocalPoints(ctx context.Context) ([]FocalPointInfo, error) {
	fm, lm := getFirstAndLastOfMonth()

	req := map[string]time.Time{
		"fromDate": fm,
		"toDate":   lm,
	}

	res, err := t.client.Put(getFocalUrl, req)

	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("get focal point fail with code: %d", res.StatusCode)
	}

	body, err := ioutil.ReadAll(res.Body)

	if err != nil {
		return nil, err
	}

	fp := &FocalPointsResponse{}
	err = json.Unmarshal(body, fp)

	if err != nil {
		return nil, err
	}

	return fp.Data, nil
}

const getProjectUrl = "/api/v1/employees/projects"

func (t *Bairestt) GetProjects(ctx context.Context) ([]ProjectInfo, error) {
	fm, lm := getFirstAndLastOfMonth()

	req := map[string]time.Time{
		"fromDate": fm,
		"toDate":   lm,
	}

	res, err := t.client.Put(getProjectUrl, req)

	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("get projects fail with code: %d", res.StatusCode)
	}

	body, err := ioutil.ReadAll(res.Body)

	if err != nil {
		return nil, err
	}

	p := &ProjectsResponse{}
	err = json.Unmarshal(body, p)

	if err != nil {
		return nil, err
	}

	return p.Data, nil
}
