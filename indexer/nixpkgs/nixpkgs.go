package nixpkgs

import (
	"encoding/json"
	"errors"
	"fmt"
)

const Prefix = "nixpkgs/package/"

type Nixpkgs struct {
	Version  int                `json:"version"`
	Packages map[string]Package `json:"packages"`
}

type Package struct {
	Meta    Meta   `json:"meta"`
	Version string `json:"version"`
}

type Meta struct {
	Description          string               `json:"description"`
	LongDescription      string               `json:"longDescription"`
	MainProgram          string               `json:"mainProgram"`
	Homepages            ElemOrSlice[string]  `json:"homepage"`
	Licenses             ElemOrSlice[License] `json:"license"`
	Maintainers          []Maintainer         `json:"maintainers"`
	Broken               bool                 `json:"broken"`
	Unfree               bool                 `json:"unfree"`
	Insecure             bool                 `json:"insecure"`
	Name                 string               `json:"name"`
	Position             string               `json:"position"`
	Platforms            FlexibleStringSlice  `json:"platforms"`
	KnownVulnerabilities []string             `json:"knownVulnerabilities"`
}

type FlexibleStringSlice []string

func (f *FlexibleStringSlice) UnmarshalJSON(data []byte) error {
	var arr []string
	if err := json.Unmarshal(data, &arr); err == nil {
		*f = arr
		return nil
	}
	var s string
	if err := json.Unmarshal(data, &s); err == nil {
		*f = []string{s}
		return nil
	}

	// Sometimes, the package has a weird platforms field:
	// {"_type":"kernel","execFormat":{"_type":"exec-format","name":"elf"},"families":{},"name":"linux"}},{"cpu":{"bits":64,"family":"x86"},"kernel":{"_type":"kernel","execFormat":{"_type":"exec-format","name":"elf"},"families":{},"name":"linux"}},{"cpu":{"family":"power"},"kernel":{"_type":"kernel","execFormat":{"_type":"exec-format","name":"elf"},"families":{},"name":"linux"}},{"cpu":{"bits":64,"family":"arm"},"kernel":{"_type":"kernel","execFormat":{"_type":"exec-format","name":"elf"},"families":{},"name":"linux"}},{"cpu":{"family":"sparc"},"kernel":{"_type":"kernel","execFormat":{"_type":"exec-format","name":"elf"},"families":{},"name":"linux"}},{"cpu":{"arch":"armv7-a"},"kernel":{"_type":"kernel","execFormat":{"_type":"exec-format","name":"elf"},"families":{},"name":"linux"}},{"cpu":{"arch":"armv7"},"kernel":{"_type":"kernel","execFormat":{"_type":"exec-format","name":"elf"},"families":{},"name":"linux"}},{"cpu":{"arch":"armv7-m"},"kernel":{"_type":"kernel","execFormat":{"_type":"exec-format","name":"elf"},"families":{},"name":"linux"}},{"cpu":{"arch":"armv7-r"},"kernel":{"_type":"kernel","execFormat":{"_type":"exec-format","name":"elf"},"families":{},"name":"linux"}}]
	*f = []string{}

	return nil
}

type Maintainer struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	GitHub   string `json:"github"`
	GithubId int    `json:"githubId"`
}

type License struct {
	Free     bool   `json:"free"`
	FullName string `json:"fullName"`
	SpdxID   string `json:"spdxId"`
}

type LicenseNoUnmarshal struct {
	Free     bool   `json:"free"`
	FullName string `json:"fullName"`
	SpdxID   string `json:"spdxId"`
}

func (l *License) UnmarshalJSON(data []byte) error {
	if len(data) == 0 {
		return errors.New("input data is empty")
	}

	switch {
	case data[0] == '"':
		s := ""
		err := json.Unmarshal(data, &s)
		if err != nil {
			return fmt.Errorf("unmarshal string: %w", err)
		}
		(*l).FullName = s

	default:
		lu := LicenseNoUnmarshal{}
		err := json.Unmarshal(data, &lu)
		if err != nil {
			return fmt.Errorf("unmarshal struct: %w", err)
		}
		*l = License(lu)
	}

	return nil
}

type ElemOrSlice[T any] []T

func (eos *ElemOrSlice[T]) UnmarshalJSON(data []byte) error {
	if len(data) == 0 {
		return errors.New("input data is empty")
	}

	switch {
	case data[0] == '[':
		s := []T{}
		err := json.Unmarshal(data, &s)
		if err != nil {
			return fmt.Errorf("unmarshal slice: %w", err)
		}
		*eos = s

	default:
		var e T
		err := json.Unmarshal(data, &e)
		if err != nil {
			return fmt.Errorf("unmarshal element: %w", err)
		}
		*eos = []T{e}
	}

	return nil
}
