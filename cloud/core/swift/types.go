package swift

type Container struct {
	Name  string `json:"name"`
	Count int    `json:"count"`
	Bytes int64  `json:"bytes"`
}

type Object struct {
	Name        string `json:"name"`
	Bytes       int64  `json:"bytes"`
	ContentType string `json:"content_type"`
}
