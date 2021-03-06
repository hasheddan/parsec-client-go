// Copyright 2021 Contributors to the Parsec project.
// SPDX-License-Identifier: Apache-2.0

package auth

import (
	"bytes"
	"fmt"
)

// AuthenticationType identifies type of authentication to be used in a request message.
type AuthenticationType uint8

// Authentication Types
const (
	AuthNoAuth AuthenticationType = 0
	// Direct authentication
	AuthDirect AuthenticationType = 1
	// JSON Web Tokens (JWT) authentication (not currently supported)
	AuthJwt AuthenticationType = 2
	// Unix peer credentials authentication
	AuthUnixPeerCredentials AuthenticationType = 3
	// Authentication verifying a JWT SPIFFE Verifiable Identity Document
	AuthJwtSvid AuthenticationType = 4
)

func NewAuthenticationTypeFromU32(t uint32) (AuthenticationType, error) {
	if t > uint32(AuthJwtSvid) {
		return AuthNoAuth, fmt.Errorf("cannot convert value %v to AuthenticationType", t)
	}
	return AuthenticationType(t), nil
}

// AuthenticatorInfo contains information about an authenticator.
// Id is the id used to select the authenticator
// Name name of the authenticator
type AuthenticatorInfo struct {
	ID          AuthenticationType
	Description string
	VersionMaj  uint32
	VersionMin  uint32
	VersionRev  uint32
}

// Authenticator interface for an authenticator
// Info returns information about the authenticator
// NewRequestAuth creates a RequestAuthToken ready to populate a request
type Authenticator interface {
	Info() AuthenticatorInfo
	NewRequestAuth() (RequestAuthToken, error)
}

// AuthenticatorFactory Create an authenticator for the given authentication type
func AuthenticatorFactory(authType AuthenticationType) (Authenticator, error) {
	var auth Authenticator
	var err error
	switch authType {
	case AuthNoAuth:
		auth, err = newNoAuthAuthenticator()
	case AuthUnixPeerCredentials:
		auth, err = newUnixPeerAuthenticator()
	case AuthJwt, AuthDirect, AuthJwtSvid:
		err = fmt.Errorf("unsupported authenticator type %v", authType)
	default:
		err = fmt.Errorf("invalid authenticator type %v", authType)
	}
	if err != nil {
		return nil, err
	}
	return auth, nil
}

// RequestAuthToken describes interface for token to contain an authentication field in a request
type RequestAuthToken interface {
	Buffer() *bytes.Buffer
	AuthType() AuthenticationType
}

// DefaultRequestAuthToken represents a request authentication payload
type DefaultRequestAuthToken struct {
	buf      *bytes.Buffer
	authType AuthenticationType
}

// Buffer returns byte buffer with the token to be sent in a request to the server
func (a DefaultRequestAuthToken) Buffer() *bytes.Buffer {
	return a.buf
}

// AuthType returns the auth type value to put in a request header
func (a DefaultRequestAuthToken) AuthType() AuthenticationType {
	return a.authType
}
