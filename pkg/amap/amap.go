package amap

import (
	"net/http"
)

type Amap struct {
	Key   string
	httpc *http.Client
}

type Option func(*Amap)

func WithHttpClient(c *http.Client) Option {
	return func(a *Amap) {
		a.httpc = c
	}
}

func New(key string, opts ...Option) *Amap {
	a := &Amap{
		key,
		http.DefaultClient,
	}
	for _, opt := range opts {
		opt(a)
	}
	return a
}
