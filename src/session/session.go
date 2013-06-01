package session

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

type Session interface {
	Set(key, value interface{})
	Get(key interface{}) interface{}
	Delete(key interface{})
	SessionId() string
}

type Provider interface {
	ProviderInit(maxlifetime int64) error
	SessionRead(sid string) (Session, error)
	SessionDestroy(sid string) error
	SessionGC()
}

type SessionManager struct {
	cookieName  string
	provider    Provider
	maxlifetime int64
}

func (manager *SessionManager) sessionId() string {
	b := make([]byte, 24)
	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		panic("session:create sessionId error")
	}
	return base64.URLEncoding.EncodeToString(b)
}

func (manager *SessionManager) SessionStart(w http.ResponseWriter, r *http.Request) (session Session) {
	cookie, err := r.Cookie(manager.cookieName)
	if err != nil || cookie.Value == "" {
		sessionId := manager.sessionId()
		session, _ = manager.provider.SessionRead(sessionId)
		cookie := http.Cookie{Name: manager.cookieName, Value: url.QueryEscape(sessionId), Path: "/"}
		cookie.Expires = time.Now().Add(time.Duration(manager.maxlifetime) * time.Second)
		http.SetCookie(w, &cookie)
		r.AddCookie(&cookie)
	} else {
		cookie.Expires = time.Now().Add(time.Duration(manager.maxlifetime) * time.Second)
		cookie.Path = "/"
		http.SetCookie(w, cookie)
		sid, _ := url.QueryUnescape(cookie.Value)
		session, _ = manager.provider.SessionRead(sid)
	}
	return
}

func (manager *SessionManager) SessionDestroy(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie(manager.cookieName)
	if err != nil || cookie.Value == "" {
		return
	} else {
		manager.provider.SessionDestroy(cookie.Value)
		expires := time.Now()
		cookie := http.Cookie{Name: manager.cookieName, Path: "/", Expires: expires, MaxAge: -1}
		http.SetCookie(w, &cookie)
		return
	}
}

func (manager *SessionManager) GC() {
	manager.provider.SessionGC()
	time.AfterFunc(time.Duration(manager.maxlifetime)*time.Second, func() { manager.GC() })
}

func NewManager(providerName string, cookieName string, maxlifetime int64) (*SessionManager, error) {
	provider, ok := providerMap[providerName]
	if !ok {
		return nil, fmt.Errorf("session:unknown provider %v ", providerName)
	}
	err := provider.ProviderInit(maxlifetime)
	if err != nil {
		return nil, fmt.Errorf("session:session init error in provider %v", providerName)
	}
	return &SessionManager{cookieName: cookieName, provider: provider, maxlifetime: maxlifetime}, nil
}

var providerMap = make(map[string]Provider)

func Register(name string, provider Provider) {
	if provider == nil {
		panic("session:Register provider is nil")
	}
	if _, dum := providerMap[name]; dum {
		panic("session:Register called twice for provider " + name)
	}
	providerMap[name] = provider
}
