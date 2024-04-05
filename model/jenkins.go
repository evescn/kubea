package model

type Jenkins struct {
	GroupName   string `json:"group_name"`
	Name        string `json:"project_name"`
	CopyJobName string `json:"copy_job_name"`
}
