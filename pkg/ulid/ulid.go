// Copyright (c) Mainflux
// SPDX-License-Identifier: Apache-2.0

// Package ulid provides a ULID identity provider.
package ulid

import (
	"time"

	"git.willowglen.ca/sq/third-party/mainflux"
	"git.willowglen.ca/sq/third-party/mainflux/pkg/errors"
	"github.com/oklog/ulid/v2"

	mathrand "math/rand"
)

// ErrGeneratingID indicates error in generating ULID
var ErrGeneratingID = errors.New("generating id failed")

var _ mainflux.IDProvider = (*ulidProvider)(nil)

type ulidProvider struct {
	entropy *mathrand.Rand
}

// New instantiates a ULID provider.
func New() mainflux.IDProvider {
	seed := time.Now().UnixNano()
	source := mathrand.NewSource(seed)
	return &ulidProvider{
		entropy: mathrand.New(source),
	}
}

func (up *ulidProvider) ID() (string, error) {
	id, err := ulid.New(ulid.Timestamp(time.Now()), up.entropy)
	if err != nil {
		return "", err
	}

	return id.String(), nil
}
