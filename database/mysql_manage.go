package database

import (
	"sync"
)

var MysqlManage *MysqlClientManage

func init() {
	MysqlManage = &MysqlClientManage{
		m: make(map[string]*MysqlClient),
	}
}

type MysqlClientManage struct {
	mu sync.RWMutex
	m  map[string]*MysqlClient

	// TODO: 链表存个list
}

func (m *MysqlClientManage) IsExist(uniqueDBName string) (*MysqlClient, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	client, ok := m.m[uniqueDBName]
	return client, ok
}

func (m *MysqlClientManage) AddClient(uniqueDBName string, db *MysqlClient) error {
	if err := db.TestConnection(); err != nil {
		return err
	}

	m.mu.Lock()
	defer m.mu.Unlock()
	//if _, ok := m.IsExist(uniqueDBName); ok {
	//	return fmt.Errorf("%s name is exist", uniqueDBName)
	//}
	m.m[uniqueDBName] = db
	return nil
}

func (m *MysqlClientManage) GetClientList() []string {
	m.mu.RLock()
	defer m.mu.RUnlock()
	res := make([]string, 0, len(m.m))
	for name := range m.m {
		res = append(res, name)
	}
	return res
}
