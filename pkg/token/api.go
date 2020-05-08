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

package token

import "context"

//go:generate mockgen -destination mock/access_token_generator.gen.go -package mock go.zenithar.org/solid/pkg/token AccessTokenGenerator

// AccessTokenGenerator describes accessToken generator contract.
type AccessTokenGenerator interface {
	Generate(ctx context.Context) (string, error)
}

//go:generate mockgen -destination mock/phantom_token_generator.gen.go -package mock go.zenithar.org/solid/pkg/token PhantomTokenGenerator

// PhantomTokenGenerator describes phantom token generator contract.
type PhantomTokenGenerator interface {
	Generate(ctx context.Context) (string, error)
}

//go:generate mockgen -destination mock/id_token_generator.gen.go -package mock go.zenithar.org/solid/pkg/token IDTokenGenerator

// IDTokenGenerator describes idToken generator contract.
type IDTokenGenerator interface {
	Generate(ctx context.Context) (string, error)
}
