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
	"go.zenithar.org/solid/api/oidc"
	"go.zenithar.org/solid/pkg/rfcerrors"
	"go.zenithar.org/solid/pkg/storage"
	storagemock "go.zenithar.org/solid/pkg/storage/mock"
	tokenmock "go.zenithar.org/solid/pkg/token/mock"

	"github.com/golang/mock/gomock"
	"github.com/google/go-cmp/cmp"
)

func Test_service_Introspect(t *testing.T) {
	type args struct {
		ctx context.Context
		req *corev1.TokenIntrospectionRequest
	}
	tests := []struct {
		name    string
		args    args
		prepare func(*storagemock.MockClientReader, *storagemock.MockAuthorizationRequestReader, *tokenmock.MockAccessTokenGenerator, *storagemock.MockSession, *storagemock.MockToken)
		want    *corev1.TokenIntrospectionResponse
		wantErr bool
	}{
		{
			name: "nil request",
			args: args{
				ctx: context.Background(),
			},
			wantErr: true,
			want: &corev1.TokenIntrospectionResponse{
				Error: rfcerrors.InvalidRequest(""),
			},
		},
		{
			name: "nil client authentication",
			args: args{
				ctx: context.Background(),
				req: &corev1.TokenIntrospectionRequest{},
			},
			wantErr: true,
			want: &corev1.TokenIntrospectionResponse{
				Error: rfcerrors.InvalidClient(""),
			},
		},
		{
			name: "nil token",
			args: args{
				ctx: context.Background(),
				req: &corev1.TokenIntrospectionRequest{
					Client: &corev1.Client{},
				},
			},
			wantErr: true,
			want: &corev1.TokenIntrospectionResponse{
				Error: rfcerrors.InvalidRequest(""),
			},
		},
		{
			name: "empty token",
			args: args{
				ctx: context.Background(),
				req: &corev1.TokenIntrospectionRequest{
					Client: &corev1.Client{},
					Token:  "",
				},
			},
			wantErr: true,
			want: &corev1.TokenIntrospectionResponse{
				Error: rfcerrors.InvalidRequest(""),
			},
		},
		// ---------------------------------------------------------------------
		{
			name: "client not found",
			args: args{
				ctx: context.Background(),
				req: &corev1.TokenIntrospectionRequest{
					Client: &corev1.Client{
						ClientId: "s6BhdRkqt3",
					},
					Token: "cwE.HcbVtkyQCyCUfjxYvjHNODfTbVpSlmyo",
				},
			},
			prepare: func(clients *storagemock.MockClientReader, _ *storagemock.MockAuthorizationRequestReader, _ *tokenmock.MockAccessTokenGenerator, _ *storagemock.MockSession, tokens *storagemock.MockToken) {
				clients.EXPECT().Get(gomock.Any(), "s6BhdRkqt3").Return(nil, storage.ErrNotFound)
			},
			wantErr: true,
			want: &corev1.TokenIntrospectionResponse{
				Error: rfcerrors.InvalidClient(""),
			},
		},
		{
			name: "client storage error",
			args: args{
				ctx: context.Background(),
				req: &corev1.TokenIntrospectionRequest{
					Client: &corev1.Client{
						ClientId: "s6BhdRkqt3",
					},
					Token: "cwE.HcbVtkyQCyCUfjxYvjHNODfTbVpSlmyo",
				},
			},
			prepare: func(clients *storagemock.MockClientReader, _ *storagemock.MockAuthorizationRequestReader, _ *tokenmock.MockAccessTokenGenerator, _ *storagemock.MockSession, tokens *storagemock.MockToken) {
				clients.EXPECT().Get(gomock.Any(), "s6BhdRkqt3").Return(nil, fmt.Errorf("foo"))
			},
			wantErr: true,
			want: &corev1.TokenIntrospectionResponse{
				Error: rfcerrors.ServerError(""),
			},
		},
		// ---------------------------------------------------------------------
		{
			name: "token not found",
			args: args{
				ctx: context.Background(),
				req: &corev1.TokenIntrospectionRequest{
					Client: &corev1.Client{
						ClientId: "s6BhdRkqt3",
					},
					Token: "cwE.HcbVtkyQCyCUfjxYvjHNODfTbVpSlmyo",
				},
			},
			prepare: func(clients *storagemock.MockClientReader, _ *storagemock.MockAuthorizationRequestReader, at *tokenmock.MockAccessTokenGenerator, _ *storagemock.MockSession, tokens *storagemock.MockToken) {
				clients.EXPECT().Get(gomock.Any(), "s6BhdRkqt3").Return(&corev1.Client{
					GrantTypes: []string{oidc.GrantTypeClientCredentials},
				}, nil)
				tokens.EXPECT().GetByValue(gomock.Any(), "cwE.HcbVtkyQCyCUfjxYvjHNODfTbVpSlmyo").Return(nil, storage.ErrNotFound)
			},
			wantErr: true,
			want: &corev1.TokenIntrospectionResponse{
				Token: &corev1.Token{
					Value:  "cwE.HcbVtkyQCyCUfjxYvjHNODfTbVpSlmyo",
					Status: corev1.TokenStatus_TOKEN_STATUS_INVALID,
				},
			},
		},
		{
			name: "token storage error",
			args: args{
				ctx: context.Background(),
				req: &corev1.TokenIntrospectionRequest{
					Client: &corev1.Client{
						ClientId: "s6BhdRkqt3",
					},
					Token: "cwE.HcbVtkyQCyCUfjxYvjHNODfTbVpSlmyo",
				},
			},
			prepare: func(clients *storagemock.MockClientReader, _ *storagemock.MockAuthorizationRequestReader, at *tokenmock.MockAccessTokenGenerator, _ *storagemock.MockSession, tokens *storagemock.MockToken) {
				clients.EXPECT().Get(gomock.Any(), "s6BhdRkqt3").Return(&corev1.Client{
					GrantTypes: []string{oidc.GrantTypeClientCredentials},
				}, nil)
				tokens.EXPECT().GetByValue(gomock.Any(), "cwE.HcbVtkyQCyCUfjxYvjHNODfTbVpSlmyo").Return(nil, fmt.Errorf("foo"))
			},
			wantErr: true,
			want: &corev1.TokenIntrospectionResponse{
				Token: &corev1.Token{
					Value:  "cwE.HcbVtkyQCyCUfjxYvjHNODfTbVpSlmyo",
					Status: corev1.TokenStatus_TOKEN_STATUS_INVALID,
				},
			},
		},
		// ---------------------------------------------------------------------
		{
			name: "valid",
			args: args{
				ctx: context.Background(),
				req: &corev1.TokenIntrospectionRequest{
					Client: &corev1.Client{
						ClientId: "s6BhdRkqt3",
					},
					Token: "cwE.HcbVtkyQCyCUfjxYvjHNODfTbVpSlmyo",
				},
			},
			prepare: func(clients *storagemock.MockClientReader, _ *storagemock.MockAuthorizationRequestReader, at *tokenmock.MockAccessTokenGenerator, _ *storagemock.MockSession, tokens *storagemock.MockToken) {
				clients.EXPECT().Get(gomock.Any(), "s6BhdRkqt3").Return(&corev1.Client{}, nil)
				tokens.EXPECT().GetByValue(gomock.Any(), "cwE.HcbVtkyQCyCUfjxYvjHNODfTbVpSlmyo").Return(&corev1.Token{
					Status:  corev1.TokenStatus_TOKEN_STATUS_ACTIVE,
					TokenId: "123456789",
					Value:   "cwE.HcbVtkyQCyCUfjxYvjHNODfTbVpSlmyo",
				}, nil)
			},
			wantErr: false,
			want: &corev1.TokenIntrospectionResponse{
				Token: &corev1.Token{
					Value:  "cwE.HcbVtkyQCyCUfjxYvjHNODfTbVpSlmyo",
					Status: corev1.TokenStatus_TOKEN_STATUS_ACTIVE,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			// Arm mocks
			clients := storagemock.NewMockClientReader(ctrl)
			authorizationRequests := storagemock.NewMockAuthorizationRequestReader(ctrl)
			accessTokens := tokenmock.NewMockAccessTokenGenerator(ctrl)
			idTokens := tokenmock.NewMockIDTokenGenerator(ctrl)
			sessions := storagemock.NewMockSession(ctrl)
			tokens := storagemock.NewMockToken(ctrl)

			// Prepare them
			if tt.prepare != nil {
				tt.prepare(clients, authorizationRequests, accessTokens, sessions, tokens)
			}

			// Instanciate service
			underTest := New(accessTokens, idTokens, clients, authorizationRequests, sessions, tokens)

			got, err := underTest.Introspect(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("service.Introspect() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if diff := cmp.Diff(got, tt.want, cmpOpts...); diff != "" {
				t.Errorf("service.Introspect() res = %s", diff)
			}
		})
	}
}