package cookie

import (
	"grok-proxy/config"
	"math/rand"
	"sync"
)

// Manager manages cookies and user agents
type Manager struct {
	cookies     []string
	userAgents  []string
	cookieIndex int
	mutex       sync.Mutex
}

// NewManager creates a new cookie manager
func NewManager() (*Manager, error) {
	cfg, err := config.GetInstance()
	if err != nil {
		return nil, err
	}

	return &Manager{
		cookies:     cfg.Cookies,
		userAgents:  cfg.UserAgent,
		cookieIndex: 0,
	}, nil
}

// GetUserAgent returns a random user agent
func (m *Manager) GetUserAgent() string {
	if len(m.userAgents) == 0 {
		return "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/123.0.0.0 Safari/537.36"
	}
	return m.userAgents[rand.Intn(len(m.userAgents))]
}

// GetCookie returns the next cookie in rotation
func (m *Manager) GetCookie() string {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if len(m.cookies) == 0 {
		return ""
	}

	cookie := m.cookies[m.cookieIndex]
	m.cookieIndex = (m.cookieIndex + 1) % len(m.cookies)
	return cookie
}

// CookieCount returns the total number of cookies
func (m *Manager) CookieCount() int {
	return len(m.cookies)
}

// CurrentCookieIndex returns the current cookie index
func (m *Manager) CurrentCookieIndex() int {
	return m.cookieIndex
}
