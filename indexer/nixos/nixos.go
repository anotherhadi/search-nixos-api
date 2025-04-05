package nixos

const (
	Prefix = "nixpkgs/option/"
)

type Package struct {
	Example struct {
		Text string `json:"text"`
	} `json:"example"`
	Type         string   `json:"type"`
	Description  string   `json:"description"`
	Declarations []string `json:"declarations"`
	Default      struct {
		Text string `json:"text"`
	} `json:"default"`
}
