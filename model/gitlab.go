package model

type GitLab struct {
	GroupName   string `json:"group_name"`
	ProjectName string `json:"project_name"`
	Visibility  string `json:"visibility"`
	Description string `json:"desc"`
	GroupID     int    `json:"group_id"`
}

type TempInfo []struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type ProjectInfo struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Path        string `json:"path"`
	NameSpaceID int    `json:"namespace_id"`
	Visibility  string `json:"visibility"`
	ImportUrl   string `json:"import_url"`
}
