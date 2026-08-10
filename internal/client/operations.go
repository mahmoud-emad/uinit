package client

import (
	"bufio"
	"encoding/json"

	"github.com/uinit/internal/manager"
)

func (c *UinitClient) sendRequest(
	action string,
	processName string,
) (manager.Response, error) {
	req := manager.Request{
		Action:  action,
		Process: processName,
	}

	rsp := manager.Response{}

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
