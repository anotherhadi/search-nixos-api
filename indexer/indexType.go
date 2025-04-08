package indexer

type Index struct {
	Info map[string]string `json:"info"`

	// Options
	Nixos       Options `json:"nixos"`
	Homemanager Options `json:"home-manager"`
	Darwin      Options `json:"darwin"`

	// Packages
	Nixpkgs Packages `json:"nixpkgs"`
	Nur     Packages `json:"nur"`
}

type Packages map[string]Package

type Package struct {
	Source            string       `json:"source"` // nixpkgs, nur
	Name              string       `json:"name"`
	Version           string       `json:"version"`
	Description       string       `json:"description"`
	LongDescription   string       `json:"longDescription"`
	MainProgram       string       `json:"mainProgram"`
	Homepages         []string     `json:"homepages"`
	Maintainers       []Maintainer `json:"maintainers"`
	Platforms         []string     `json:"platforms"`
	PlatformsSimplify []string     `json:"platformsSimplify"`
	Position          string       `json:"position"`
	PositionUrl       string       `json:"positionUrl"`

	Broken     bool `json:"broken"`
	Vulnerable bool `json:"vulnerable"`

	Licenses []License `json:"licenses"`
	Unfree   bool      `json:"unfree"`

	Insecure             bool     `json:"insecure"`
	KnownVulnerabilities []string `json:"knownVulnerabilities"`
}

type Maintainer struct {
	Name     string `json:"name"`
	Email    string `json:"email,omitempty"`
	GitHub   string `json:"github"`
	GithubId int    `json:"githubId"`
}

type License struct {
	Free     bool   `json:"free"`
	FullName string `json:"fullName"`
	SpdxID   string `json:"spdxId"`
}

type Options map[string]Option

type Option struct {
	Source       string   `json:"source"` // nixpkgs, homemanager, darwin
	Description  string   `json:"description"`
	Example      string   `json:"example"`
	Type         string   `json:"type"`
	Declarations []string `json:"declarations"`
	Default      string   `json:"default"`
}
