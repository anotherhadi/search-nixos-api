package darwin

const Prefix = "darwin/option/"

type Darwin struct {
	Version  int                `json:"version"`
	Packages map[string]Package `json:"packages"`
}

type Package struct {
	Type        string   `json:"type"`
	Default     string   `json:"default"`
	Example     string   `json:"example"`
	DeclaredBy  []string `json:"declarations"`
	Description string   `json:"description"`
}
