package supervisor

import (
	"encoding/json"
	"log"
)

type Request struct {
	Action  string `json:"action"`
	Service string `json:"service,omitempty"`
}

type Response struct {
	OK      bool        `json:"ok"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

func (c *UinitClient) sendRequest(action string, service string) (Response, error) {
	req := Request{
		Action:  action,
		Service: service,
	}

	rsp := Response{}

	encoded, err := json.Marshal(req)
	if err != nil {
		return rsp, err
	}

	encoded = append(encoded, '\n')

	n, err := c.conn.Write(encoded)
	if err != nil {
		return rsp, err
	}

	log.Println("n is: ", n)
	return rsp, nil
}
