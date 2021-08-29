package define

type FileMeta struct {
	Name      string   `json:"name"`
	LocalPath string   `json:"local_path"`
	Size      int64    `json:"size"`
	Md5       string   `json:"md5"`
	Folders   []string `json:"folders"`
}
