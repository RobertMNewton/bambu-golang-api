package request

import "encoding/json"

type Request struct {
	Type    string
	Payload RequestPayload
}

func (r *Request) ToMessage() ([]byte, error) {
	return json.Marshal(map[string]RequestPayload{r.Type: r.Payload})
}

func CreateRequest(requestType, command, sequenceID string, params map[string]interface{}) Request {
	return Request{
		Type: requestType,
		Payload: RequestPayload{
			SequenceID: sequenceID,
			Command:    command,
			Params:     params,
		},
	}
}
