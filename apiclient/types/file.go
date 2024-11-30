package types

type File struct {
	Name string `json:"name,omitempty"`
}

type FileList List[File]

type FolderSet map[string]Item

type Item struct{}
