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

package jarm

import (
	"context"
	"fmt"
	"time"

	"zntr.io/solid/pkg/jwk"

	corev1 "zntr.io/solid/api/gen/go/oidc/core/v1"

	"github.com/square/go-jose/v3"
	"github.com/square/go-jose/v3/jwt"
)

// -----------------------------------------------------------------------------

// JWTEncoder builds a JWT Response encoder instance.
func JWTEncoder(alg jose.SignatureAlgorithm, keyProvider jwk.KeyProviderFunc) ResponseEncoder {
	return &jwtEncoder{
		alg:         alg,
		keyProvider: keyProvider,
	}
}

type jwtEncoder struct {
	alg         jose.SignatureAlgorithm
	keyProvider jwk.KeyProviderFunc
}

func (d *jwtEncoder) Encode(ctx context.Context, issuer string, resp *corev1.AuthorizationCodeResponse) (string, error) {
	// Check arguments
	if resp == nil {
		return "", fmt.Errorf("unable to process nil response")
	}
	if issuer == "" {
		return "", fmt.Errorf("unable to process empty issuer")
	}

	// Retrieve key from provider
	k, err := d.keyProvider(ctx)
	if err != nil {
		return "", fmt.Errorf("unable to retrieve signing key from provider: %w", err)
	}
	if k == nil {
		return "", fmt.Errorf("key provider returned nil key")
	}
	if k.IsPublic() {
		return "", fmt.Errorf("key provider returned a public key")
	}

	// Prepare response claims
	var claims *jwtResponseClaims
	if resp.Error != nil {
		claims = &jwtResponseClaims{
			State:            resp.State,
			Error:            resp.Error.Err,
			ErrorDescription: resp.Error.ErrorDescription,
		}
	} else {
		claims = &jwtResponseClaims{
			State:     resp.State,
			Issuer:    issuer,
			Audience:  resp.ClientId,
			Code:      resp.Code,
			ExpiresAt: uint64(time.Now().Add(time.Duration(resp.ExpiresIn) * time.Second).Unix()),
		}
	}

	// Prepare a signer
	sig, err := jose.NewSigner(jose.SigningKey{
		Algorithm: d.alg,
		Key:       k,
	}, nil)
	if err != nil {
		return "", fmt.Errorf("unable to prepare signer: %w", err)
	}

	// Generate the final proof
	raw, err := jwt.Signed(sig).Claims(claims).CompactSerialize()
	if err != nil {
		return "", fmt.Errorf("unable to generate JARM: %w", err)
	}

	// No error
	return raw, nil
}