package entity

type FsFileInfo struct {
	Name      string       `json:"name"`      // file name
	Path      string       `json:"path"`      // file path
	FullPath  string       `json:"full_path"` // file full path
	Extension string       `json:"extension"` // file extension
	Md5       string       `json:"md5"`       // MD5 hash
	IsDir     bool         `json:"is_dir"`    // whether it is directory
	FileSize  int64        `json:"file_size"` // file size (bytes)
	Children  []FsFileInfo `json:"children"`  // children for sub-directory
}
