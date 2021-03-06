// Tideland GoREST - JSON Web Token - Unit Tests
//
// Copyright (C) 2016 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

package jwt_test

//--------------------
// IMPORTS
//--------------------

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"strings"
	"testing"
	"time"

	"github.com/tideland/golib/audit"

	"github.com/tideland/gorest/jwt"
)

//--------------------
// TESTS
//--------------------

var (
	esTests = []jwt.Algorithm{jwt.ES256, jwt.ES384, jwt.ES512}
	hsTests = []jwt.Algorithm{jwt.HS256, jwt.HS384, jwt.HS512}
	psTests = []jwt.Algorithm{jwt.PS256, jwt.PS384, jwt.PS512}
	rsTests = []jwt.Algorithm{jwt.RS256, jwt.RS384, jwt.RS512}
)

// TestESAlgorithms tests the ECDSA algorithms for the
// JWT signature.
func TestESAlgorithms(t *testing.T) {
	assert := audit.NewTestingAssertion(t, true)
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	assert.Nil(err)
	claims := initClaims()
	for _, test := range esTests {
		assert.Logf("testing algorithm %q", test)
		// Encode.
		jwtEncode, err := jwt.Encode(claims, privateKey, test)
		assert.Nil(err)
		parts := strings.Split(jwtEncode.String(), ".")
		assert.Length(parts, 3)
		// Verify.
		jwtVerify, err := jwt.Verify(jwtEncode.String(), privateKey.Public())
		assert.Nil(err)
		assert.Equal(jwtEncode.Algorithm(), jwtVerify.Algorithm())
		assert.Equal(jwtEncode.String(), jwtVerify.String())
		testClaims(assert, jwtVerify.Claims())
	}
}

// TestHSAlgorithms tests the HMAC algorithms for the
// JWT signature.
func TestHSAlgorithms(t *testing.T) {
	assert := audit.NewTestingAssertion(t, true)
	key := []byte("secret")
	claims := initClaims()
	for _, test := range hsTests {
		assert.Logf("testing algorithm %q", test)
		// Encode.
		jwtEncode, err := jwt.Encode(claims, key, test)
		assert.Nil(err)
		parts := strings.Split(jwtEncode.String(), ".")
		assert.Length(parts, 3)
		// Verify.
		jwtVerify, err := jwt.Verify(jwtEncode.String(), key)
		assert.Nil(err)
		assert.Equal(jwtEncode.Algorithm(), jwtVerify.Algorithm())
		assert.Equal(jwtEncode.String(), jwtVerify.String())
		testClaims(assert, jwtVerify.Claims())
	}
}

// TestPSAlgorithms tests the RSAPSS algorithms for the
// JWT signature.
func TestPSAlgorithms(t *testing.T) {
	assert := audit.NewTestingAssertion(t, true)
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	assert.Nil(err)
	claims := initClaims()
	for _, test := range psTests {
		assert.Logf("testing algorithm %q", test)
		// Encode.
		jwtEncode, err := jwt.Encode(claims, privateKey, test)
		assert.Nil(err)
		parts := strings.Split(jwtEncode.String(), ".")
		assert.Length(parts, 3)
		// Verify.
		jwtVerify, err := jwt.Verify(jwtEncode.String(), privateKey.Public())
		assert.Nil(err)
		assert.Equal(jwtEncode.Algorithm(), jwtVerify.Algorithm())
		assert.Equal(jwtEncode.String(), jwtVerify.String())
		testClaims(assert, jwtVerify.Claims())
	}
}

// TestRSAlgorithms tests the RSA algorithms for the
// JWT signature.
func TestRSAlgorithms(t *testing.T) {
	assert := audit.NewTestingAssertion(t, true)
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	assert.Nil(err)
	claims := initClaims()
	for _, test := range rsTests {
		assert.Logf("testing algorithm %q", test)
		// Encode.
		jwtEncode, err := jwt.Encode(claims, privateKey, test)
		assert.Nil(err)
		parts := strings.Split(jwtEncode.String(), ".")
		assert.Length(parts, 3)
		// Verify.
		jwtVerify, err := jwt.Verify(jwtEncode.String(), privateKey.Public())
		assert.Nil(err)
		assert.Equal(jwtEncode.Algorithm(), jwtVerify.Algorithm())
		assert.Equal(jwtEncode.String(), jwtVerify.String())
		testClaims(assert, jwtVerify.Claims())
	}
}

// TestNoneAlgorithm tests the none algorithm for the
// JWT signature.
func TestNoneAlgorithm(t *testing.T) {
	assert := audit.NewTestingAssertion(t, true)
	assert.Logf("testing algorithm \"none\"")
	// Encode.
	claims := initClaims()
	jwtEncode, err := jwt.Encode(claims, "", jwt.NONE)
	assert.Nil(err)
	parts := strings.Split(jwtEncode.String(), ".")
	assert.Length(parts, 3)
	assert.Equal(parts[2], "")
	// Verify.
	jwtVerify, err := jwt.Verify(jwtEncode.String(), "")
	assert.Nil(err)
	assert.Equal(jwtEncode.Algorithm(), jwtVerify.Algorithm())
	assert.Equal(jwtEncode.String(), jwtVerify.String())
	testClaims(assert, jwtVerify.Claims())
}

// TestNotMatchingAlgorithm
func TestNotMatchingAlgorithm(t *testing.T) {
	assert := audit.NewTestingAssertion(t, true)
	esPrivateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	esPublicKey := esPrivateKey.Public()
	assert.Nil(err)
	hsKey := []byte("secret")
	rsPrivateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	rsPublicKey := rsPrivateKey.Public()
	assert.Nil(err)
	noneKey := ""
	claims := initClaims()
	errorMatch := ".* combination of algorithm .* and key type .*"
	tests := []struct {
		description string
		algorithm   jwt.Algorithm
		key         jwt.Key
		encodeKeys  []jwt.Key
		verifyKeys  []jwt.Key
	}{
		{"ECDSA", jwt.ES512, esPrivateKey,
			[]jwt.Key{hsKey, rsPrivateKey, noneKey}, []jwt.Key{hsKey, rsPublicKey, noneKey}},
		{"HMAC", jwt.HS512, hsKey,
			[]jwt.Key{esPrivateKey, rsPrivateKey, noneKey}, []jwt.Key{esPublicKey, rsPublicKey, noneKey}},
		{"RSA", jwt.RS512, rsPrivateKey,
			[]jwt.Key{esPrivateKey, hsKey, noneKey}, []jwt.Key{esPublicKey, hsKey, noneKey}},
		{"RSAPSS", jwt.PS512, rsPrivateKey,
			[]jwt.Key{esPrivateKey, hsKey, noneKey}, []jwt.Key{esPublicKey, hsKey, noneKey}},
		{"none", jwt.NONE, noneKey,
			[]jwt.Key{esPrivateKey, hsKey, rsPrivateKey}, []jwt.Key{esPublicKey, hsKey, rsPublicKey}},
	}
	// Run the tests.
	for _, test := range tests {
		assert.Logf("testing %q algorithm key type mismatch", test.description)
		for _, key := range test.encodeKeys {
			_, err = jwt.Encode(claims, key, test.algorithm)
			assert.ErrorMatch(err, errorMatch)
		}
		jwtEncode, err := jwt.Encode(claims, test.key, test.algorithm)
		assert.Nil(err)
		for _, key := range test.verifyKeys {
			_, err = jwt.Verify(jwtEncode.String(), key)
			assert.ErrorMatch(err, errorMatch)
		}
	}
}

// TestESTools tests the tools for the reading of PEM encoded
// ECDSA keys.
func TestESTools(t *testing.T) {
	assert := audit.NewTestingAssertion(t, true)
	assert.Logf("testing \"ECDSA\" tools")
	// Generate keys and PEMs.
	privateKeyIn, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	assert.Nil(err)
	privateBytes, err := x509.MarshalECPrivateKey(privateKeyIn)
	assert.Nil(err)
	privateBlock := pem.Block{
		Type:  "EC PRIVATE KEY",
		Bytes: privateBytes,
	}
	privatePEM := pem.EncodeToMemory(&privateBlock)
	publicBytes, err := x509.MarshalPKIXPublicKey(privateKeyIn.Public())
	assert.Nil(err)
	publicBlock := pem.Block{
		Type:  "EC PUBLIC KEY",
		Bytes: publicBytes,
	}
	publicPEM := pem.EncodeToMemory(&publicBlock)
	assert.NotNil(publicPEM)
	// Now read them.
	buf := bytes.NewBuffer(privatePEM)
	privateKeyOut, err := jwt.ReadECPrivateKey(buf)
	assert.Nil(err)
	buf = bytes.NewBuffer(publicPEM)
	publicKeyOut, err := jwt.ReadECPublicKey(buf)
	assert.Nil(err)
	// And as a last step check if they are correctly usable.
	claims := initClaims()
	jwtEncode, err := jwt.Encode(claims, privateKeyOut, jwt.ES512)
	assert.Nil(err)
	parts := strings.Split(jwtEncode.String(), ".")
	assert.Length(parts, 3)
	jwtVerify, err := jwt.Verify(jwtEncode.String(), publicKeyOut)
	assert.Nil(err)
	assert.Equal(jwtEncode.Algorithm(), jwtVerify.Algorithm())
	assert.Equal(jwtEncode.String(), jwtVerify.String())
	testClaims(assert, jwtVerify.Claims())
}

// TestRSTools tests the tools for the reading of PEM encoded
// RSA keys.
func TestRSTools(t *testing.T) {
	assert := audit.NewTestingAssertion(t, true)
	assert.Logf("testing \"RSA\" tools")
	// Generate keys and PEMs.
	privateKeyIn, err := rsa.GenerateKey(rand.Reader, 2048)
	assert.Nil(err)
	privateBytes := x509.MarshalPKCS1PrivateKey(privateKeyIn)
	privateBlock := pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: privateBytes,
	}
	privatePEM := pem.EncodeToMemory(&privateBlock)
	publicBytes, err := x509.MarshalPKIXPublicKey(privateKeyIn.Public())
	assert.Nil(err)
	publicBlock := pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: publicBytes,
	}
	publicPEM := pem.EncodeToMemory(&publicBlock)
	assert.NotNil(publicPEM)
	// Now read them.
	buf := bytes.NewBuffer(privatePEM)
	privateKeyOut, err := jwt.ReadRSAPrivateKey(buf)
	assert.Nil(err)
	buf = bytes.NewBuffer(publicPEM)
	publicKeyOut, err := jwt.ReadRSAPublicKey(buf)
	assert.Nil(err)
	// And as a last step check if they are correctly usable.
	claims := initClaims()
	jwtEncode, err := jwt.Encode(claims, privateKeyOut, jwt.RS512)
	assert.Nil(err)
	parts := strings.Split(jwtEncode.String(), ".")
	assert.Length(parts, 3)
	jwtVerify, err := jwt.Verify(jwtEncode.String(), publicKeyOut)
	assert.Nil(err)
	assert.Equal(jwtEncode.Algorithm(), jwtVerify.Algorithm())
	assert.Equal(jwtEncode.String(), jwtVerify.String())
	testClaims(assert, jwtVerify.Claims())
}

// TestDecode tests the decoding without verifying the signature.
func TestDecode(t *testing.T) {
	assert := audit.NewTestingAssertion(t, true)
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	assert.Nil(err)
	claims := initClaims()
	assert.Logf("testing decoding without verifying")
	// Encode.
	jwtEncode, err := jwt.Encode(claims, privateKey, jwt.RS512)
	assert.Nil(err)
	parts := strings.Split(jwtEncode.String(), ".")
	assert.Length(parts, 3)
	// Decode.
	jwtDecode, err := jwt.Decode(jwtEncode.String())
	assert.Nil(err)
	assert.Equal(jwtEncode.Algorithm(), jwtDecode.Algorithm())
	key, err := jwtDecode.Key()
	assert.Nil(key)
	assert.ErrorMatch(err, ".*no key available, only after encoding or verifying.*")
	assert.Equal(jwtEncode.String(), jwtDecode.String())
	testClaims(assert, jwtDecode.Claims())
}

// TestIsValid checks the time validation of a token.
func TestIsValid(t *testing.T) {
	assert := audit.NewTestingAssertion(t, true)
	assert.Logf("testing time validation")
	now := time.Now()
	leeway := time.Minute
	key := []byte("secret")
	// Create token with no times set, encode, decode, validate ok.
	claims := jwt.NewClaims()
	jwtEncode, err := jwt.Encode(claims, key, jwt.HS512)
	assert.Nil(err)
	jwtDecode, err := jwt.Decode(jwtEncode.String())
	assert.Nil(err)
	ok := jwtDecode.IsValid(leeway)
	assert.True(ok)
	// Now a token with a long timespan, still valid.
	claims = jwt.NewClaims()
	claims.SetNotBefore(now.Add(-time.Hour))
	claims.SetExpiration(now.Add(time.Hour))
	jwtEncode, err = jwt.Encode(claims, key, jwt.HS512)
	assert.Nil(err)
	jwtDecode, err = jwt.Decode(jwtEncode.String())
	assert.Nil(err)
	ok = jwtDecode.IsValid(leeway)
	assert.True(ok)
	// Now a token with a long timespan in the past, not valid.
	claims = jwt.NewClaims()
	claims.SetNotBefore(now.Add(-2 * time.Hour))
	claims.SetExpiration(now.Add(-time.Hour))
	jwtEncode, err = jwt.Encode(claims, key, jwt.HS512)
	assert.Nil(err)
	jwtDecode, err = jwt.Decode(jwtEncode.String())
	assert.Nil(err)
	ok = jwtDecode.IsValid(leeway)
	assert.False(ok)
	// And at last a token with a long timespan in the future, not valid.
	claims = jwt.NewClaims()
	claims.SetNotBefore(now.Add(time.Hour))
	claims.SetExpiration(now.Add(2 * time.Hour))
	jwtEncode, err = jwt.Encode(claims, key, jwt.HS512)
	assert.Nil(err)
	jwtDecode, err = jwt.Decode(jwtEncode.String())
	assert.Nil(err)
	ok = jwtDecode.IsValid(leeway)
	assert.False(ok)
}

// EOF
