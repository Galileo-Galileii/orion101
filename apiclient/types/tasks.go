package types

type Task struct {
	Metadata
	TaskManifest
}

type TaskList List[Task]

type TaskManifest struct {
	Name        string     `json:"name"`
	Description string     `json:"description"`
	Steps       []TaskStep `json:"steps"`
}

type TaskStep struct {
	ID   string  `json:"id,omitempty"`
	If   *TaskIf `json:"if,omitempty"`
	Step string  `json:"step,omitempty"`
}

type TaskIf struct {
	Condition string     `json:"condition,omitempty"`
	Steps     []TaskStep `json:"steps,omitempty"`
	Else      []TaskStep `json:"else,omitempty"`
}
