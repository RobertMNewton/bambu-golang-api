package request

import (
	"encoding/json"
)

type RequestPayload struct {
	SequenceID string
	Command    string
	Params     map[string]interface{}
}

func (r RequestPayload) MarshalJSON() ([]byte, error) {
	data := map[string]interface{}{
		"sequence_id": r.SequenceID,
		"command":     r.Command,
	}

	for k, v := range r.Params {
		data[k] = v
	}
	return json.Marshal(data)
}
