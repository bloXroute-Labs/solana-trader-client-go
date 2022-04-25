package connections

import (
	"fmt"
	"github.com/gorilla/websocket"
)

type ConnectionManager struct {
	connectionMap map[int]*websocket.Conn
	address       string
	id            int
}

func NewConnectionManager(address string) ConnectionManager {
	return ConnectionManager{
		connectionMap: make(map[int]*websocket.Conn),
		address:       address,
		id:            0,
	}
}

func (c *ConnectionManager) Next() (*websocket.Conn, int, error) {
	conn, _, err := websocket.DefaultDialer.Dial(c.address, nil)
	if err != nil {
		return nil, 0, err
	}
	if conn == nil {
		return nil, 0, fmt.Errorf("connection to %s was nil", c.address)
	}

	c.id++
	c.connectionMap[c.id] = conn

	return conn, c.id, nil
}

func (c *ConnectionManager) RemoveConnection(id int) error {
	conn, ok := c.connectionMap[id]
	if !ok {
		return fmt.Errorf("conn with id %v not found", id)
	}

	err := conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
	if err != nil {
		return fmt.Errorf("error writing close msg -  %v", err)
	}

	delete(c.connectionMap, id)

	return nil
}
