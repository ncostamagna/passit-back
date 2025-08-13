package types

import "encoding/json"

type Secret struct {
	Expiration int32  `json:"expiration,omitempty"`
	Message    string `json:"message"`
	OneTime    bool   `json:"one_time,omitempty"`
}

func (s *Secret) ToJSON() ([]byte, error) {
	return json.Marshal(&s)
}

func (s *Secret) FromJSON(data []byte) error {
	return json.Unmarshal(data, s)
}