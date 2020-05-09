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

import (
	"context"
	"fmt"
	"testing"

	corev1 "go.zenithar.org/solid/api/gen/go/oidc/core/v1"
	registrationv1 "go.zenithar.org/solid/api/gen/go/oidc/registration/v1"
	"go.zenithar.org/solid/api/oidc"
	"go.zenithar.org/solid/pkg/rfcerrors"
	"go.zenithar.org/solid/pkg/storage"
	storagemock "go.zenithar.org/solid/pkg/storage/mock"
	tokenmock "go.zenithar.org/solid/pkg/token/mock"

	"github.com/golang/mock/gomock"
	"github.com/golang/protobuf/ptypes/wrappers"
	"github.com/google/go-cmp/cmp"
)

func Test_service_authorizationCode(t *testing.T) {
	type args struct {
		ctx    context.Context
		client *registrationv1.Client
		req    *corev1.TokenRequest
	}
	tests := []struct {
		name    string
		args    args
		prepare func(*storagemock.MockAuthorizationRequestReader, *tokenmock.MockAccessTokenGenerator)
		want    *corev1.TokenResponse
		wantErr bool
	}{
		{
			name: "nil client",
			args: args{
				ctx: context.Background(),
				req: &corev1.TokenRequest{
					GrantType: oidc.GrantTypeAuthorizationCode,
					Grant: &corev1.TokenRequest_AuthorizationCode{
						AuthorizationCode: &corev1.GrantAuthorizationCode{
							CodeVerifier: "dBjftJeZ4CVP-mB92K27uhbUJU1p1r_wW1gFWFOEjXk",
							RedirectUri:  "https://client.example.org/cb",
						},
					},
				},
			},
			wantErr: true,
			want: &corev1.TokenResponse{
				Error: rfcerrors.ServerError(""),
			},
		},
		{
			name: "nil request",
			args: args{
				ctx:    context.Background(),
				client: &registrationv1.Client{},
			},
			wantErr: true,
			want: &corev1.TokenResponse{
				Error: rfcerrors.ServerError(""),
			},
		},
		{
			name: "nil grant",
			args: args{
				ctx:    context.Background(),
				client: &registrationv1.Client{},
				req: &corev1.TokenRequest{
					GrantType: oidc.GrantTypeAuthorizationCode,
				},
			},
			wantErr: true,
			want: &corev1.TokenResponse{
				Error: rfcerrors.ServerError(""),
			},
		},
		{
			name: "client not support code response_type",
			args: args{
				ctx: context.Background(),
				client: &registrationv1.Client{
					GrantTypes:    []string{oidc.GrantTypeClientCredentials},
					ResponseTypes: []string{"code"},
					RedirectUris:  []string{"https://client.example.org/cb"},
				},
				req: &corev1.TokenRequest{
					Client: &corev1.ClientAuthentication{
						ClientId: "s6BhdRkqt3",
					},
					GrantType: oidc.GrantTypeAuthorizationCode,
					Grant: &corev1.TokenRequest_AuthorizationCode{
						AuthorizationCode: &corev1.GrantAuthorizationCode{
							CodeVerifier: "dBjftJeZ4CVP-mB92K27uhbUJU1p1r_wW1gFWFOEjXk",
							RedirectUri:  "https://client.example.org/cb",
						},
					},
				},
			},
			wantErr: true,
			want: &corev1.TokenResponse{
				Error: rfcerrors.UnsupportedGrantType(""),
			},
		},
		{
			name: "missing code",
			args: args{
				ctx: context.Background(),
				client: &registrationv1.Client{
					GrantTypes:    []string{oidc.GrantTypeAuthorizationCode},
					ResponseTypes: []string{"code"},
					RedirectUris:  []string{"https://client.example.org/cb"},
				},
				req: &corev1.TokenRequest{
					Client: &corev1.ClientAuthentication{
						ClientId: "s6BhdRkqt3",
					},
					GrantType: oidc.GrantTypeAuthorizationCode,
					Grant: &corev1.TokenRequest_AuthorizationCode{
						AuthorizationCode: &corev1.GrantAuthorizationCode{
							CodeVerifier: "dBjftJeZ4CVP-mB92K27uhbUJU1p1r_wW1gFWFOEjXk",
							RedirectUri:  "https://client.example.org/cb",
						},
					},
				},
			},
			wantErr: true,
			want: &corev1.TokenResponse{
				Error: rfcerrors.InvalidGrant(""),
			},
		},
		{
			name: "missing code_verifier",
			args: args{
				ctx: context.Background(),
				client: &registrationv1.Client{
					GrantTypes:    []string{oidc.GrantTypeAuthorizationCode},
					ResponseTypes: []string{"code"},
					RedirectUris:  []string{"https://client.example.org/cb"},
				},
				req: &corev1.TokenRequest{
					Client: &corev1.ClientAuthentication{
						ClientId: "s6BhdRkqt3",
					},
					GrantType: oidc.GrantTypeAuthorizationCode,
					Grant: &corev1.TokenRequest_AuthorizationCode{
						AuthorizationCode: &corev1.GrantAuthorizationCode{
							Code:        "1234567891234567890",
							RedirectUri: "https://client.example.org/cb",
						},
					},
				},
			},
			wantErr: true,
			want: &corev1.TokenResponse{
				Error: rfcerrors.InvalidGrant(""),
			},
		},
		{
			name: "missing redirect_uri",
			args: args{
				ctx: context.Background(),
				client: &registrationv1.Client{
					GrantTypes:    []string{oidc.GrantTypeAuthorizationCode},
					ResponseTypes: []string{"code"},
					RedirectUris:  []string{"https://client.example.org/cb"},
				},
				req: &corev1.TokenRequest{
					Client: &corev1.ClientAuthentication{
						ClientId: "s6BhdRkqt3",
					},
					GrantType: oidc.GrantTypeAuthorizationCode,
					Grant: &corev1.TokenRequest_AuthorizationCode{
						AuthorizationCode: &corev1.GrantAuthorizationCode{
							Code:         "1234567891234567890",
							CodeVerifier: "dBjftJeZ4CVP-mB92K27uhbUJU1p1r_wW1gFWFOEjXk",
						},
					},
				},
			},
			wantErr: true,
			want: &corev1.TokenResponse{
				Error: rfcerrors.InvalidGrant(""),
			},
		},
		{
			name: "authorization request not found",
			args: args{
				ctx: context.Background(),
				client: &registrationv1.Client{
					GrantTypes:    []string{oidc.GrantTypeAuthorizationCode},
					ResponseTypes: []string{"code"},
					RedirectUris:  []string{"https://client.example.org/cb"},
				},
				req: &corev1.TokenRequest{
					Client: &corev1.ClientAuthentication{
						ClientId: "s6BhdRkqt3",
					},
					GrantType: oidc.GrantTypeAuthorizationCode,
					Grant: &corev1.TokenRequest_AuthorizationCode{
						AuthorizationCode: &corev1.GrantAuthorizationCode{
							Code:         "1234567891234567890",
							CodeVerifier: "dBjftJeZ4CVP-mB92K27uhbUJU1p1r_wW1gFWFOEjXk",
							RedirectUri:  "https://client.example.org/cb",
						},
					},
				},
			},
			prepare: func(ar *storagemock.MockAuthorizationRequestReader, _ *tokenmock.MockAccessTokenGenerator) {
				ar.EXPECT().GetByCode(gomock.Any(), "1234567891234567890").Return(nil, storage.ErrNotFound)
			},
			wantErr: true,
			want: &corev1.TokenResponse{
				Error: rfcerrors.InvalidGrant(""),
			},
		},
		{
			name: "authorization request storage error",
			args: args{
				ctx: context.Background(),
				client: &registrationv1.Client{
					GrantTypes:    []string{oidc.GrantTypeAuthorizationCode},
					ResponseTypes: []string{"code"},
					RedirectUris:  []string{"https://client.example.org/cb"},
				},
				req: &corev1.TokenRequest{
					Client: &corev1.ClientAuthentication{
						ClientId: "s6BhdRkqt3",
					},
					GrantType: oidc.GrantTypeAuthorizationCode,
					Grant: &corev1.TokenRequest_AuthorizationCode{
						AuthorizationCode: &corev1.GrantAuthorizationCode{
							Code:         "1234567891234567890",
							CodeVerifier: "dBjftJeZ4CVP-mB92K27uhbUJU1p1r_wW1gFWFOEjXk",
							RedirectUri:  "https://client.example.org/cb",
						},
					},
				},
			},
			prepare: func(ar *storagemock.MockAuthorizationRequestReader, _ *tokenmock.MockAccessTokenGenerator) {
				ar.EXPECT().GetByCode(gomock.Any(), "1234567891234567890").Return(nil, fmt.Errorf("foo"))
			},
			wantErr: true,
			want: &corev1.TokenResponse{
				Error: rfcerrors.ServerError(""),
			},
		},
		{
			name: "redirect_uri mismatch",
			args: args{
				ctx: context.Background(),
				client: &registrationv1.Client{
					GrantTypes:    []string{oidc.GrantTypeAuthorizationCode},
					ResponseTypes: []string{"code"},
					RedirectUris:  []string{"https://client.example.org/cb"},
				},
				req: &corev1.TokenRequest{
					Client: &corev1.ClientAuthentication{
						ClientId: "s6BhdRkqt3",
					},
					GrantType: oidc.GrantTypeAuthorizationCode,
					Grant: &corev1.TokenRequest_AuthorizationCode{
						AuthorizationCode: &corev1.GrantAuthorizationCode{
							Code:         "1234567891234567890",
							CodeVerifier: "dBjftJeZ4CVP-mB92K27uhbUJU1p1r_wW1gFWFOEjXk",
							RedirectUri:  "https://client.example.org/cb12346",
						},
					},
				},
			},
			prepare: func(ar *storagemock.MockAuthorizationRequestReader, _ *tokenmock.MockAccessTokenGenerator) {
				ar.EXPECT().GetByCode(gomock.Any(), "1234567891234567890").Return(&corev1.AuthorizationRequest{
					ResponseType:        "code",
					Scope:               "openid profile email",
					ClientId:            "s6BhdRkqt3",
					State:               "af0ifjsldkj",
					RedirectUri:         "https://client.example.org/cb",
					CodeChallenge:       "E9Melhoa2OwvFrEMTJguCHaoeK1t8URWbuGJSstw-cM",
					CodeChallengeMethod: "S256",
				}, nil)
			},
			wantErr: true,
			want: &corev1.TokenResponse{
				Error: rfcerrors.InvalidGrant("af0ifjsldkj"),
			},
		},
		{
			name: "redirect_uri mismatch: client changes between request",
			args: args{
				ctx: context.Background(),
				client: &registrationv1.Client{
					GrantTypes:    []string{oidc.GrantTypeAuthorizationCode},
					ResponseTypes: []string{"code"},
					RedirectUris:  []string{"https://client.example.org/cb"},
				},
				req: &corev1.TokenRequest{
					Client: &corev1.ClientAuthentication{
						ClientId: "s6BhdRkqt3",
					},
					GrantType: oidc.GrantTypeAuthorizationCode,
					Grant: &corev1.TokenRequest_AuthorizationCode{
						AuthorizationCode: &corev1.GrantAuthorizationCode{
							Code:         "1234567891234567890",
							CodeVerifier: "dBjftJeZ4CVP-mB92K27uhbUJU1p1r_wW1gFWFOEjXk",
							RedirectUri:  "https://client.example.org/cb1",
						},
					},
				},
			},
			prepare: func(ar *storagemock.MockAuthorizationRequestReader, _ *tokenmock.MockAccessTokenGenerator) {
				ar.EXPECT().GetByCode(gomock.Any(), "1234567891234567890").Return(&corev1.AuthorizationRequest{
					ResponseType:        "code",
					Scope:               "openid profile email",
					ClientId:            "s6BhdRkqt3",
					State:               "af0ifjsldkj",
					RedirectUri:         "https://client.example.org/cb1",
					CodeChallenge:       "E9Melhoa2OwvFrEMTJguCHaoeK1t8URWbuGJSstw-cM",
					CodeChallengeMethod: "S256",
				}, nil)
			},
			wantErr: true,
			want: &corev1.TokenResponse{
				Error: rfcerrors.InvalidGrant("af0ifjsldkj"),
			},
		},
		{
			name: "invalid code_verifier",
			args: args{
				ctx: context.Background(),
				client: &registrationv1.Client{
					GrantTypes:    []string{oidc.GrantTypeAuthorizationCode},
					ResponseTypes: []string{"code"},
					RedirectUris:  []string{"https://client.example.org/cb"},
				},
				req: &corev1.TokenRequest{
					Client: &corev1.ClientAuthentication{
						ClientId: "s6BhdRkqt3",
					},
					GrantType: oidc.GrantTypeAuthorizationCode,
					Grant: &corev1.TokenRequest_AuthorizationCode{
						AuthorizationCode: &corev1.GrantAuthorizationCode{
							Code:         "1234567891234567890",
							CodeVerifier: "foo",
							RedirectUri:  "https://client.example.org/cb",
						},
					},
				},
			},
			prepare: func(ar *storagemock.MockAuthorizationRequestReader, _ *tokenmock.MockAccessTokenGenerator) {
				ar.EXPECT().GetByCode(gomock.Any(), "1234567891234567890").Return(&corev1.AuthorizationRequest{
					ResponseType:        "code",
					Scope:               "openid profile email",
					ClientId:            "s6BhdRkqt3",
					State:               "af0ifjsldkj",
					RedirectUri:         "https://client.example.org/cb",
					CodeChallenge:       "E9Melhoa2OwvFrEMTJguCHaoeK1t8URWbuGJSstw-cM",
					CodeChallengeMethod: "S256",
				}, nil)
			},
			wantErr: true,
			want: &corev1.TokenResponse{
				Error: rfcerrors.InvalidGrant("af0ifjsldkj"),
			},
		},
		{
			name: "invalid code_challenge_method",
			args: args{
				ctx: context.Background(),
				client: &registrationv1.Client{
					GrantTypes:    []string{oidc.GrantTypeAuthorizationCode},
					ResponseTypes: []string{"code"},
					RedirectUris:  []string{"https://client.example.org/cb"},
				},
				req: &corev1.TokenRequest{
					Client: &corev1.ClientAuthentication{
						ClientId: "s6BhdRkqt3",
					},
					GrantType: oidc.GrantTypeAuthorizationCode,
					Grant: &corev1.TokenRequest_AuthorizationCode{
						AuthorizationCode: &corev1.GrantAuthorizationCode{
							Code:         "1234567891234567890",
							CodeVerifier: "dBjftJeZ4CVP-mB92K27uhbUJU1p1r_wW1gFWFOEjXk",
							RedirectUri:  "https://client.example.org/cb",
						},
					},
				},
			},
			prepare: func(ar *storagemock.MockAuthorizationRequestReader, _ *tokenmock.MockAccessTokenGenerator) {
				ar.EXPECT().GetByCode(gomock.Any(), "1234567891234567890").Return(&corev1.AuthorizationRequest{
					ResponseType:        "code",
					Scope:               "openid profile email",
					ClientId:            "s6BhdRkqt3",
					State:               "af0ifjsldkj",
					RedirectUri:         "https://client.example.org/cb",
					CodeChallenge:       "E9Melhoa2OwvFrEMTJguCHaoeK1t8URWbuGJSstw-cM",
					CodeChallengeMethod: "xxx",
				}, nil)
			},
			wantErr: true,
			want: &corev1.TokenResponse{
				Error: rfcerrors.InvalidGrant("af0ifjsldkj"),
			},
		},
		// ---------------------------------------------------------------------
		{
			name: "openid: generate access token error",
			args: args{
				ctx: context.Background(),
				client: &registrationv1.Client{
					GrantTypes:    []string{oidc.GrantTypeAuthorizationCode},
					ResponseTypes: []string{"code"},
					RedirectUris:  []string{"https://client.example.org/cb"},
				},
				req: &corev1.TokenRequest{
					Client: &corev1.ClientAuthentication{
						ClientId: "s6BhdRkqt3",
					},
					GrantType: oidc.GrantTypeAuthorizationCode,
					Grant: &corev1.TokenRequest_AuthorizationCode{
						AuthorizationCode: &corev1.GrantAuthorizationCode{
							Code:         "1234567891234567890",
							CodeVerifier: "dBjftJeZ4CVP-mB92K27uhbUJU1p1r_wW1gFWFOEjXk",
							RedirectUri:  "https://client.example.org/cb",
						},
					},
				},
			},
			prepare: func(ar *storagemock.MockAuthorizationRequestReader, at *tokenmock.MockAccessTokenGenerator) {
				ar.EXPECT().GetByCode(gomock.Any(), "1234567891234567890").Return(&corev1.AuthorizationRequest{
					ResponseType:        "code",
					Scope:               "openid profile email",
					ClientId:            "s6BhdRkqt3",
					State:               "af0ifjsldkj",
					RedirectUri:         "https://client.example.org/cb",
					CodeChallenge:       "E9Melhoa2OwvFrEMTJguCHaoeK1t8URWbuGJSstw-cM",
					CodeChallengeMethod: "S256",
				}, nil)
				at.EXPECT().Generate(gomock.Any()).Return("", fmt.Errorf("foo"))
			},
			wantErr: true,
			want: &corev1.TokenResponse{
				Error: rfcerrors.ServerError("af0ifjsldkj"),
			},
		},
		{
			name: "openid: generate refresh token error",
			args: args{
				ctx: context.Background(),
				client: &registrationv1.Client{
					GrantTypes:    []string{oidc.GrantTypeAuthorizationCode},
					ResponseTypes: []string{"code"},
					RedirectUris:  []string{"https://client.example.org/cb"},
				},
				req: &corev1.TokenRequest{
					Client: &corev1.ClientAuthentication{
						ClientId: "s6BhdRkqt3",
					},
					GrantType: oidc.GrantTypeAuthorizationCode,
					Grant: &corev1.TokenRequest_AuthorizationCode{
						AuthorizationCode: &corev1.GrantAuthorizationCode{
							Code:         "1234567891234567890",
							CodeVerifier: "dBjftJeZ4CVP-mB92K27uhbUJU1p1r_wW1gFWFOEjXk",
							RedirectUri:  "https://client.example.org/cb",
						},
					},
				},
			},
			prepare: func(ar *storagemock.MockAuthorizationRequestReader, at *tokenmock.MockAccessTokenGenerator) {
				ar.EXPECT().GetByCode(gomock.Any(), "1234567891234567890").Return(&corev1.AuthorizationRequest{
					ResponseType:        "code",
					Scope:               "openid profile email offline_access",
					ClientId:            "s6BhdRkqt3",
					State:               "af0ifjsldkj",
					RedirectUri:         "https://client.example.org/cb",
					CodeChallenge:       "E9Melhoa2OwvFrEMTJguCHaoeK1t8URWbuGJSstw-cM",
					CodeChallengeMethod: "S256",
				}, nil)
				accessTokenSuccess := at.EXPECT().Generate(gomock.Any()).Return("1/fFAGRNJru1FTz70BzhT3Zg", nil)
				at.EXPECT().Generate(gomock.Any()).Return("", fmt.Errorf("foo")).After(accessTokenSuccess)
			},
			wantErr: true,
			want: &corev1.TokenResponse{
				Error: rfcerrors.ServerError("af0ifjsldkj"),
			},
		},
		// ---------------------------------------------------------------------
		{
			name: "openid: valid",
			args: args{
				ctx: context.Background(),
				client: &registrationv1.Client{
					GrantTypes:    []string{oidc.GrantTypeAuthorizationCode},
					ResponseTypes: []string{"code"},
					RedirectUris:  []string{"https://client.example.org/cb"},
				},
				req: &corev1.TokenRequest{
					Client: &corev1.ClientAuthentication{
						ClientId: "s6BhdRkqt3",
					},
					GrantType: oidc.GrantTypeAuthorizationCode,
					Grant: &corev1.TokenRequest_AuthorizationCode{
						AuthorizationCode: &corev1.GrantAuthorizationCode{
							Code:         "1234567891234567890",
							CodeVerifier: "dBjftJeZ4CVP-mB92K27uhbUJU1p1r_wW1gFWFOEjXk",
							RedirectUri:  "https://client.example.org/cb",
						},
					},
				},
			},
			prepare: func(ar *storagemock.MockAuthorizationRequestReader, at *tokenmock.MockAccessTokenGenerator) {
				ar.EXPECT().GetByCode(gomock.Any(), "1234567891234567890").Return(&corev1.AuthorizationRequest{
					ResponseType:        "code",
					Scope:               "openid profile email offline_access",
					ClientId:            "s6BhdRkqt3",
					State:               "af0ifjsldkj",
					RedirectUri:         "https://client.example.org/cb",
					CodeChallenge:       "E9Melhoa2OwvFrEMTJguCHaoeK1t8URWbuGJSstw-cM",
					CodeChallengeMethod: "S256",
				}, nil)
				accessTokenSuccess := at.EXPECT().Generate(gomock.Any()).Return("1/fFAGRNJru1FTz70BzhT3Zg", nil)
				at.EXPECT().Generate(gomock.Any()).Return("5ZsdF6h/sQAghJFRD", nil).After(accessTokenSuccess)
			},
			wantErr: false,
			want: &corev1.TokenResponse{
				Error: nil,
				Openid: &corev1.OpenIDToken{
					AccessToken:  "1/fFAGRNJru1FTz70BzhT3Zg",
					ExpiresIn:    3600,
					RefreshToken: &wrappers.StringValue{Value: "5ZsdF6h/sQAghJFRD"},
					TokenType:    "Bearer",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			// Arm mocks
			authorizationRequests := storagemock.NewMockAuthorizationRequestReader(ctrl)
			accessTokens := tokenmock.NewMockAccessTokenGenerator(ctrl)
			idTokens := tokenmock.NewMockIDTokenGenerator(ctrl)

			// Prepare them
			if tt.prepare != nil {
				tt.prepare(authorizationRequests, accessTokens)
			}

			s := &service{
				accessTokenGenerator:  accessTokens,
				idTokenGenerator:      idTokens,
				authorizationRequests: authorizationRequests,
			}
			got, err := s.authorizationCode(tt.args.ctx, tt.args.client, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("service.authorizationCode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if diff := cmp.Diff(got, tt.want, cmpOpts...); diff != "" {
				t.Errorf("service.authorizationCode() res = %s", diff)
			}
		})
	}
}
