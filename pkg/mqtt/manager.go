package mqtt

import (
	"net"
	"strings"
	"sync"
	"time"
)

type Manager struct {
	Provider Provider
	CxtMutex sync.RWMutex
	IdToConn map[uint64]*Connection
}

func NewManager(p Provider) *Manager {
	m := &Manager{
		Provider: p,
		IdToConn: make(map[uint64]*Connection),
	}

	go m.CleanWorker()

	return m
}

func (m *Manager) NewConn(conn net.Conn) {
	NewConnection(conn, m)
}

func (m *Manager) AddConn(id uint64, c *Connection) {
	m.CxtMutex.Lock()
	oldSub, exist := m.IdToConn[id]
	if exist {
		oldSub.Close()
	}

	m.IdToConn[id] = c
	m.CxtMutex.Unlock()
}

func (m *Manager) DelConn(id uint64) {
	m.CxtMutex.Lock()
	_, exist := m.IdToConn[id]

	if exist {
		delete(m.IdToConn, id)
	}
	m.CxtMutex.Unlock()
}

func (m *Manager) GetToken(deviceid uint64) ([]byte, error) {
	m.CxtMutex.RLock()
	con, exist := m.IdToConn[deviceid]
	m.CxtMutex.RUnlock()
	if !exist {
		return nil, errorf("device not exist: %v[%v]", deviceid, deviceid)
	}

	return con.Token, nil
}

func (m *Manager) PublishMessage2Device(deviceid uint64, msg *Publish, timeout time.Duration) error {
	m.CxtMutex.RLock()
	con, exist := m.IdToConn[deviceid]
	m.CxtMutex.RUnlock()
	if !exist {
		return errorf("device not exist: %v", deviceid)
	}

	err := <-con.Publish(msg, timeout)
	return err
}

func (m *Manager) PublishMessage2Server(deviceid uint64, msg *Publish) error {
	topic := msg.TopicName
	msgtype := strings.Join(strings.Split(topic, "/")[1:], "/")

	payload := msg.Payload.(BytesPayload)

	m.Provider.OnDeviceMessage(deviceid, msgtype, payload)
	return nil
}

func (m *Manager) CleanWorker() {
	for {
		curTime := time.Now().Unix()

		for _, con := range m.IdToConn {
			if con.KeepAlive == 0 {
				continue
			}

			if uint16(curTime-con.LastHbTime) > uint16(2*con.KeepAlive/2) {
				con.Close()
			}
		}

		time.Sleep(60 * time.Second)
	}
}
