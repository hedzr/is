/*
 * Copyright © 2021 Hedzr Yeh.
 */

package timing_test

import (
	"github.com/hedzr/log/timing"
	"math/big"
	"testing"
)

func vFactorial(t *testing.T, x *big.Int) *big.Int {
	defer timing.New(timing.WithWriter(t.Logf)).Duration()

	v := big.NewInt(1)
	for one := big.NewInt(1); x.Sign() > 0; x.Sub(x, one) {
		v.Mul(v, x)
	}

	return x.Set(v)
}

func TestNew(t *testing.T) {
	t.Logf("v = %v", vFactorial(t, big.NewInt(19)))
}

func vFactorial2(t *testing.T, x *big.Int) *big.Int {
	defer timing.New(
		timing.WithMsgFormat("vFactorial2 takes %v"),
		timing.WithWriter(t.Logf),
	).Duration()

	v := big.NewInt(1)
	for one := big.NewInt(1); x.Sign() > 0; x.Sub(x, one) {
		v.Mul(v, x)
	}

	return x.Set(v)
}

func TestNewT2(t *testing.T) {
	t.Logf("v = %v", vFactorial2(t, big.NewInt(19)))
}
