package report

import (
	"encoding/json"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type Report struct {
	Type    string
	Payload ReportPayload
}

func (r *Report) UnmarshalJSON(data []byte) error {
	var raw map[string]ReportPayload
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	for k, v := range raw {
		r.Type = k
		r.Payload = v
	}

	return nil
}

func FromMessage(msg mqtt.Message) (Report, error) {
	var raw map[string]ReportPayload
	if err := json.Unmarshal(msg.Payload(), &raw); err != nil {
		return Report{}, nil
	}

	var r Report
	for k, v := range raw {
		r.Type = k
		r.Payload = v
	}

	return r, nil
}
