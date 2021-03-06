// Licensed to SolID under one or more contributor
// license agreements. See the NOTICE file distributed with
// this work for additional information regarding copyright
// ownership. SolID licenses this file to you under
// the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.

package main

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/dchest/uniuri"

	"github.com/square/go-jose/v3"
	jwt "github.com/square/go-jose/v3/jwt"
)

var clientPrivateKey = []byte(`{
    "kty": "EC",
    "d": "Uwq56PhVB6STB8MvLQWcOsKQlZbBvWFQba8D6Uhb2qDunpzqvoNyFsnAHKS_AkQB",
    "use": "sig",
    "crv": "P-384",
    "x": "m2NDaWfRRGlCkUa4FK949uLtMqitX1lYgi8UCIMtsuR60ux3d00XBlsC6j_YDOTe",
    "y": "6vxuUq3V1aoWi4FQ_h9ZNwUsmcGP8Uuqq_YN5dhP0U8lchdmZJbLF9mPiimo_6p4",
    "alg": "ES384"
}`)

type privateJWTClaims struct {
	JTI      string `json:"jti"`
	Subject  string `json:"sub"`
	Issuer   string `json:"iss"`
	Audience string `json:"aud"`
	Expires  uint64 `json:"exp"`
	IssuedAt uint64 `json:"iat"`
}

func generateAssertion(claims *privateJWTClaims) string {
	var privateKey jose.JSONWebKey
	// Decode JWK
	err := json.Unmarshal(clientPrivateKey, &privateKey)
	if err != nil {
		panic(err)
	}

	// Prepare a signer
	sig, err := jose.NewSigner(jose.SigningKey{Algorithm: jose.ES384, Key: privateKey}, (&jose.SignerOptions{}).WithType("JWT"))
	if err != nil {
		panic(err)
	}

	raw, err := jwt.Signed(sig).Claims(claims).CompactSerialize()
	if err != nil {
		panic(err)
	}

	// Assertion
	return raw
}

func main() {
	fmt.Printf("%s\n", generateAssertion(&privateJWTClaims{
		JTI:      uniuri.NewLen(8),
		Subject:  "6779ef20e75817b79602",
		Issuer:   "6779ef20e75817b79602",
		Audience: "http://127.0.0.1:8080",
		Expires:  uint64(time.Now().Add(2 * time.Hour).Unix()),
		IssuedAt: uint64(time.Now().Unix()),
	}))
}
