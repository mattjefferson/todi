package todoist

type Task struct {
	ID          string   `json:"id"`
	Content     string   `json:"content"`
	Description string   `json:"description,omitempty"`
	ProjectID   string   `json:"project_id"`
	Labels      []string `json:"labels"`
	Priority    int      `json:"priority"`
	Due         *Due     `json:"due"`
}

type Due struct {
	Date     string `json:"date"`
	Datetime string `json:"datetime"`
	String   string `json:"string"`
	Timezone string `json:"timezone"`
}

type Project struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}
