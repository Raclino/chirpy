package config

import "sync/atomic"

type ApiConfig struct {
	fileserverHits atomic.Int32
}
