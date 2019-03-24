package gcloud

type Project struct {
	Name      string `json:"name"`
	ProjectID string `json:"projectId"`
	State     string `json:"lifecycleState"`
}
