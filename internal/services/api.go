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

package services

import (
	"context"

	corev1 "go.zenithar.org/solid/api/gen/go/oidc/core/v1"
)

// CodeGenerator is the function contract used by authorization_code generator.
type CodeGenerator func(context.Context) (string, error)

// Authorization describes authorization request processor.
type Authorization interface {
	// Authorize a request.
	Authorize(ctx context.Context, req *corev1.AuthorizationRequest) (*corev1.AuthorizationResponse, error)
	// Register a request.
	Register(ctx context.Context, req *corev1.RegistrationRequest) (*corev1.RegistrationResponse, error)
}

// Token describes token requet processor.
type Token interface {
	// Token handles token retrieval.
	Token(ctx context.Context, req *corev1.TokenRequest) (*corev1.TokenResponse, error)
}
