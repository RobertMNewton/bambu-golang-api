package report

import "encoding/json"

type ReportPayload struct {
	SequenceID string
	Command    string
	Result     string
	Reason     string
	Params     map[string]interface{}
}

func (r *ReportPayload) UnmarshalJSON(data []byte) error {
	var raw map[string]interface{}
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	// Extract known fields
	if seqID, ok := raw["sequence_id"].(string); ok {
		r.SequenceID = seqID
	}
	if cmd, ok := raw["command"].(string); ok {
		r.Command = cmd
	}
	if res, ok := raw["result"].(string); ok {
		r.Result = res
	}
	if reason, ok := raw["reason"].(string); ok {
		r.Reason = reason
	}

	delete(raw, "sequence_id")
	delete(raw, "command")
	delete(raw, "result")
	delete(raw, "reason")
	r.Params = raw

	return nil
}
