package bairestt

import "time"

type TasksResponse struct {
	Data []TaskInfo `json:"data"`
}

type TaskInfo struct {
	Id       int64  `json:"id"`
	Name     string `json:"name"`
	Category string `json:"categoryName"`
}

type TimeRecord struct {
	ProjectId     int64     `json:"projectId"`
	RecordTypeId  int64     `json:"recordTypeId"`
	FocalPointId  int64     `json:"focalPointId"`
	Date          time.Time `json:"date"`
	DescriptionId int64     `json:"descriptionId"`
	Hours         float32   `json:"hours"`
	Comments      string    `json:"comments"`
}

type RecordResponse struct {
	Data []ExtTimeRecord `json:"data"`
}

type ExtTimeRecord struct {
	ProjectId        int64   `json:"projectId"`
	RecordTypeId     int64   `json:"recordTypeId"`
	FocalPointId     int64   `json:"focalPointId"`
	DescriptionId    int64   `json:"descriptionId"`
	Hours            float32 `json:"hours"`
	Comments         string  `json:"comments"`
	DescriptionName  string  `json:"descriptionName"`
	TaskCategoryName string  `json:"taskCategoryName"`
	RecordTypeName   string  `json:"recordTypeName"`
	ProjectName      string  `json:"projectName"`
}

type FocalPointsResponse struct {
	Data []FocalPointInfo `json:"data"`
}

type FocalPointInfo struct {
	Id   int64  `json:"id"`
	Name string `json:"name"`
}

type ProjectsResponse struct {
	Data []ProjectInfo `json:"data"`
}

type ProjectInfo struct {
	Id   int64  `json:"id"`
	Name string `json:"name"`
}
