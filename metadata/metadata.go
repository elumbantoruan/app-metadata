package metadata

import (
	"fmt"
	"regexp"
)

// ApplicationMetadata represents a metadata for an application
type ApplicationMetadata struct {
	ApplicationID string       `yaml:"applicationID"`
	Title         string       `yaml:"title"`
	Version       string       `yaml:"version"`
	Maintainers   []Maintainer `yaml:"maintainers"`
	Company       string       `yaml:"company"`
	Website       string       `yaml:"website"`
	Source        string       `yaml:"source"`
	License       string       `yaml:"license"`
	Description   string       `yaml:"description"`
}

// Maintainer contains the information of application maintainer.
type Maintainer struct {
	Name  string `yaml:"name"`
	Email string `yaml:"email"`
}

// ValidationMessage returns a validation message such as error description
type ValidationMessage struct {
	Description string `yaml:""`
}

// IsValid validates the ApplicationMetadata
func (am ApplicationMetadata) IsValid() (valid bool, desc *ValidationMessage) {
	valid, desc = am.isValidVersion()
	if !valid {
		return valid, desc
	}
	return am.isValidEmail()
}

func (am ApplicationMetadata) isValidVersion() (valid bool, desc *ValidationMessage) {
	if len(am.Version) == 0 {
		vm := ValidationMessage{
			Description: "version is empty",
		}
		return false, &vm
	}
	return true, nil
}

func (am ApplicationMetadata) isValidEmail() (valid bool, desc *ValidationMessage) {
	re := regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	var vm ValidationMessage
	for _, m := range am.Maintainers {
		if len(m.Email) == 0 {
			vm.Description = "maintainer's email is empty"
			return false, &vm
		}
		if !re.MatchString(m.Email) {
			vm.Description = fmt.Sprintf("%s is not a valid email address", m.Email)
			return false, &vm
		}
	}
	return true, nil
}
