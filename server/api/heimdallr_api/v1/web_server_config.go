package v1

type GroupOpts struct {
	Label    string     `json:"label"`
	Value    uint       `json:"value"`
	Children []HostOpts `json:"children"`
}

type HostOpts struct {
	Label    string    `json:"label"`
	Value    uint      `json:"value"`
	Children []SrvOpts `json:"children"`
}

type SrvOpts struct {
	Label string `json:"label"`
	Value string `json:"value"`
}
