package models

import "encoding/json"

// Metadata is used to manage the key rotations
type Metadata struct {
	Owner     string
	Resource  string
	Env       string
	ActiveKey string
	PrevKey   string
	DeleteKey string
}

type Credential struct {
	BootstrapURI string
	APIKey       string
	Metadata     Metadata
}

func (c *Credential) ToJSON() ([]byte, error) {
	return json.Marshal(c)
}

func CredentialsFromJSON(bytes []byte) (*Credential, error) {
	cred := Credential{}
	err := json.Unmarshal(bytes, &cred)
	if err != nil {
		return &cred, err
	}

	return &cred, nil
}
