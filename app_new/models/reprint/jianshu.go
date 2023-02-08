package reprint

type ParamJianShu struct {
	Id              string `json:"id"`
	AutoSaveControl int    `json:"autosave_control"`
	Title           string `json:"title"`
	Content         string `json:"content"`
}

type RetrieveJianShu struct {
	Id                int    `json:"id"`
	ContentUpdatedAt  int    `json:"content_updated_at"`
	ContentSizeStatus string `json:"content_size_status"`
	LastCompiledAt    int    `json:"last_compiled_at"`
}
