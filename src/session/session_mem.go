package session

import (
	"container/list"
	"sync"
	"time"
)

type MemProvider struct {
	lock        sync.RWMutex
	maxlifetime int64
	list        *list.List
	sessions    map[string]*list.Element
}

func (provider *MemProvider) update(sid string) {
	provider.lock.RLock()
	defer provider.lock.RUnlock()
	if ele, ok := provider.sessions[sid]; ok {
		ele.Value.(*MemSession).timeAccessed = time.Now()
		provider.list.MoveToFront(ele)
	}
}

func (provider *MemProvider) ProviderInit(maxlifetime int64) error {
	provider.maxlifetime = maxlifetime
	return nil
}

func (provider *MemProvider) SessionRead(sid string) (Session, error) {
	provider.lock.RLock()
	if ele, ok := provider.sessions[sid]; ok {
		go provider.update(sid)
		provider.lock.RUnlock()
		return ele.Value.(*MemSession), nil
	} else {
		provider.lock.RUnlock()
		provider.lock.Lock()
		session := &MemSession{sid: sid, timeAccessed: time.Now(), value: make(map[interface{}]interface{})}
		ele := provider.list.PushBack(session)
		provider.sessions[sid] = ele
		provider.lock.Unlock()
		return session, nil
	}
	return nil, nil
}

func (provider *MemProvider) SessionDestroy(sid string) error {
	provider.lock.Lock()
	defer provider.lock.Unlock()
	if ele, ok := provider.sessions[sid]; ok {
		delete(provider.sessions, sid)
		provider.list.Remove(ele)
	}
	return nil
}

func (provider *MemProvider) SessionGC() {
	provider.lock.RLock()
	defer provider.lock.RUnlock()
	for {
		ele := provider.list.Back()
		if ele == nil {
			break
		}
		if ele.Value.(*MemSession).timeAccessed.Unix()+provider.maxlifetime < time.Now().Unix() {
			provider.lock.RUnlock()
			provider.lock.Lock()
			provider.list.Remove(ele)
			delete(provider.sessions, ele.Value.(*MemSession).sid)
			provider.lock.Unlock()
			provider.lock.RLock()
		} else {
			break
		}
	}
}

type MemSession struct {
	sid          string
	timeAccessed time.Time
	value        map[interface{}]interface{}
	lock         sync.RWMutex
}

func (session *MemSession) Set(key, value interface{}) {
	session.lock.Lock()
	defer session.lock.Unlock()
	session.value[key] = value
}

func (session *MemSession) Get(key interface{}) interface{} {
	session.lock.Lock()
	defer session.lock.Unlock()
	if v, ok := session.value[key]; ok {
		return v
	}
	return nil
}

func (session *MemSession) Delete(key interface{}) {
	session.lock.Lock()
	defer session.lock.Unlock()
	delete(session.value, key)
}

func (session *MemSession) SessionId() string {
	return session.sid
}

//var provider = &MemProvider{list: list.New(), sessions: make(map[string]*list.Element)}

func init() {
	Register("memory", &MemProvider{list: list.New(), sessions: make(map[string]*list.Element)})
}
