package aws

import (
	"bytes"
	"encoding/base64"

	"gopkg.in/yaml.v2"
)

type CloudConfig struct {
	Users []User `yaml:"users"`
}

type User struct {
	Name              string   `yaml:"name"`
	Sudo              string   `yaml:"sudo"`
	Groups            []string `yaml:"groups"`
	SSHAuthorizedKeys []string `yaml:"ssh_authorized_keys"`
}

// ConvertToUserData converts the CloudConfig to a base64-encoded string suitable for use as user data in an AWS EC2 instance.
func (c *CloudConfig) ConvertToUserData() (string, error) {
	buf := bytes.NewBufferString("#cloud-config\n")
	err := yaml.NewEncoder(buf).Encode(c)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(buf.Bytes()), nil
}
