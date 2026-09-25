package nova

type Server struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Status string `json:"status"`
}

type FlavorRef string
type ImageRef string

type CreateServerRequest struct {
	Name      string
	FlavorRef string
	ImageRef  string
	NetworkID string
}
