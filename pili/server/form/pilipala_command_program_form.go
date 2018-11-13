package form

type CreateCommandForm struct {
	Title         string  `json:"title"         binding:"required"`
	FileName      string  `json:"file_name"     binding:"required"`
	TmpFileName   string  `json:"tmp_file_name" binding:"required"`
	HaveDedicate  int64   `json:"have_dedicate"`
	Params        string  `json:"params"`
	DedicateHosts []int64 `json:"dedicate_hosts"`
}

type UpdateCommandForm struct {
	Id            int64   `json:"id"            binding:"required"`
	Title         string  `json:"title"         binding:"required"`
	FileName      string  `json:"file_name"     binding:"required"`
	TmpFileName   string  `json:"tmp_file_name"`
	HaveDedicate  int64   `json:"have_dedicate"`
	Params        string  `json:"params"`
	DedicateHosts []int64 `json:"dedicate_hosts"`
}
