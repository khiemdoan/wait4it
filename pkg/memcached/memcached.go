package memcached

import (
	"context"
	"errors"

	"wait4it/pkg/model"

	"github.com/bradfitz/gomemcache/memcache"
)

func (m *MemcachedConnection) BuildContext(cx model.CheckContext) {
	m.Host = cx.Host
	m.Port = cx.Port
}

func (m *MemcachedConnection) Validate() error {
	if len(m.Host) == 0 {
		return errors.New("Host can't be empty")
	}

	if m.Port < 1 || m.Port > 65535 {
		return errors.New("Invalid port range for Memcached")
	}

	return nil
}

func (m *MemcachedConnection) Check(_ context.Context) (bool, bool, error) {
	// TODO: is it possible to handle ping using context?
	mc := memcache.New(m.BuildConnectionString())

	if err := mc.Ping(); err != nil {
		return false, false, err
	}

	return true, true, nil
}
