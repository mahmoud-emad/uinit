package client

func (c *UinitClient) List() ([]ProcessInfo, error) {
	defer func() { _ = c.conn.Close() }()

	rsp, err := c.sendRequest("LIST", "")
	if err != nil {
		return nil, err
	}

	var processes []ProcessInfo

	if err := c.decodeResponse(rsp.Data, &processes); err != nil {
		return nil, err
	}
	return processes, nil
}
