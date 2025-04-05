package nur

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
)

const Prefix = "nur/option/"

type Nur struct {
	Version  int                `json:"version"`
	Packages map[string]Package `json:"packages"`
}

type Package struct {
	Meta    Meta   `json:"meta"`
	Version string `json:"version"`
}

func (p Package) GetName() string {
	return "nixpkgs"
}

type Meta struct {
	Description     string               `json:"description"`
	LongDescription string               `json:"longDescription"`
	MainProgram     string               `json:"mainProgram"`
	Homepages       ElemOrSlice[string]  `json:"homepage"`
	Licenses        ElemOrSlice[License] `json:"license"`
	Maintainers     Maintainers          `json:"maintainers"`
	Broken          bool                 `json:"broken"`
	Unfree          bool                 `json:"unfree"`
	Name            string               `json:"name"`
	Position        string               `json:"position"`
	Platforms       FlexibleStringSlice  `json:"platforms"`
}

type (
	Maintainers         []Maintainer
	FlexibleStringSlice []string
)

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

func (m *Maintainer) UnmarshalJSON(data []byte) error {
	type Alias Maintainer
	aux := &struct {
		GithubId any `json:"githubId"`
		*Alias
	}{
		Alias: (*Alias)(m),
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	switch v := aux.GithubId.(type) {
	case float64:
		m.GithubId = int(v)
	case string:
		id, err := strconv.Atoi(v)
		if err != nil {
			return err
		}
		m.GithubId = id
	default:
		return fmt.Errorf("unexpected type for githubId: %T", v)
	}
	return nil
}

func (m *Maintainers) UnmarshalJSON(data []byte) error {
	var rawMessages []json.RawMessage
	if err := json.Unmarshal(data, &rawMessages); err != nil {
		return err
	}
	var maintainers []Maintainer
	for _, msg := range rawMessages {
		var maint Maintainer
		if err := json.Unmarshal(msg, &maint); err == nil {
			maintainers = append(maintainers, maint)
			continue
		}
		var tmp struct {
			Name string `json:"name"`
		}
		if err := json.Unmarshal(msg, &tmp); err == nil && tmp.Name != "" {
			maintainers = append(maintainers, Maintainer{Name: tmp.Name})
			continue
		}
		var s string
		if err := json.Unmarshal(msg, &s); err == nil {
			maintainers = append(maintainers, Maintainer{Name: s})
			continue
		}
		return fmt.Errorf("unable to unmarshal maintainer: %s", string(msg))
	}
	*m = maintainers
	return nil
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
