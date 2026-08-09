package supervisor

import (
	"bufio"
	"encoding/json"
	"net"
)

type Request struct {
	Action  string `json:"action"`
	Service string `json:"service,omitempty"`
}

type Response struct {
	OK      bool          `json:"ok"`
	Message string        `json:"message,omitempty"`
	Data    []ServiceInfo `json:"data,omitempty"`
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

	_, err = c.conn.Write(encoded)
	if err != nil {
		return rsp, err
	}

	reader := bufio.NewReader(c.conn)

	line, err := reader.ReadBytes('\n')
	if err != nil {
		return rsp, err
	}

	if err := json.Unmarshal(line, &rsp); err != nil {
		return rsp, err
	}

	return rsp, nil
}

func (s *Supervisor) sendResponse(conn net.Conn, rsp Response) error {
	encoded, err := json.Marshal(rsp)
	if err != nil {
		return err
	}

	encoded = append(encoded, '\n')
	_, err = conn.Write(encoded)
	if err != nil {
		return err
	}
	return nil
}
