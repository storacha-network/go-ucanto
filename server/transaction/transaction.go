package transaction

import (
	"github.com/storacha/go-ucanto/core/ipld"
	"github.com/storacha/go-ucanto/core/receipt"
	"github.com/storacha/go-ucanto/core/result"
)

// Transaction defines a result & effect pair, used by provider that wishes to
// return results that have effects.
type Transaction[O any, X any] interface {
	Out() result.Result[O, X]
	Fx() receipt.Effects
}

type transaction[O, X any] struct {
	out result.Result[O, X]
	fx  receipt.Effects
}

func (t transaction[O, X]) Out() result.Result[O, X] {
	return t.out
}

func (t transaction[O, X]) Fx() receipt.Effects {
	return t.fx
}

type effects struct {
	fork []ipld.Link
	join ipld.Link
}

func (fx effects) Fork() []ipld.Link {
	return fx.fork
}

func (fx effects) Join() ipld.Link {
	return fx.join
}

func NewEffects(fork []ipld.Link, join ipld.Link) receipt.Effects {
	return effects{fork, join}
}

// Option is an option configuring a transaction.
type Option func(cfg *txConfig)

type txConfig struct {
	fx receipt.Effects
}

// WithEffects configures the effects for the receipt.
func WithEffects(fx receipt.Effects) Option {
	return func(cfg *txConfig) {
		cfg.fx = fx
	}
}

func NewTransaction[O, X any](result result.Result[O, X], options ...Option) Transaction[O, X] {
	cfg := txConfig{}
	for _, opt := range options {
		opt(&cfg)
	}
	return transaction[O, X]{out: result, fx: cfg.fx}
}
