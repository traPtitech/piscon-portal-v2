// Code generated by BobGen mysql v0.29.0. DO NOT EDIT.
// This file is meant to be re-generated in place and/or deleted at any time.

package factory

import (
	"context"

	models "github.com/traPtitech/piscon-portal-v2/server/models"
)

type contextKey string

var (
	sessionCtx = newContextual[*models.Session]("session")
	teamCtx    = newContextual[*models.Team]("team")
	userCtx    = newContextual[*models.User]("user")
)

// Contextual is a convienience wrapper around context.WithValue and context.Value
type contextual[V any] struct {
	key contextKey
}

func newContextual[V any](key string) contextual[V] {
	return contextual[V]{key: contextKey(key)}
}

func (k contextual[V]) WithValue(ctx context.Context, val V) context.Context {
	return context.WithValue(ctx, k.key, val)
}

func (k contextual[V]) Value(ctx context.Context) (V, bool) {
	v, ok := ctx.Value(k.key).(V)
	return v, ok
}