package keystone

type Token struct {
	ID             string
	ProjectID      string
	UserID         string
	ExpiresAt      string
	Roles          []Role
	ServiceCatalog []Service
}

type Role struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type Service struct {
	Name      string     `json:"name"`
	Type      string     `json:"type"`
	Endpoints []Endpoint `json:"endpoints"`
}

type Endpoint struct {
	Region    string `json:"region"`
	Interface string `json:"interface"`
	URL       string `json:"url"`
}

type Credentials struct {
	AuthURL   string
	Username  string
	Password  string
	ProjectID string
	DomainID  string
}

type Project struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}
