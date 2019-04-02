package gcloud

type Project struct {
	Name      string `json:"name"`
	ProjectID string `json:"projectId"`
	State     string `json:"lifecycleState"`
}

type Projects []Project

func (ps Projects) Names() []string {
	var names []string
	for _, p := range ps {
		names = append(names, p.Name)
	}
	return names
}
