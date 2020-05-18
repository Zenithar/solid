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
	"time"

	corev1 "go.zenithar.org/solid/api/gen/go/oidc/core/v1"
	"go.zenithar.org/solid/api/oidc"
	"go.zenithar.org/solid/pkg/rfcerrors"
	"go.zenithar.org/solid/pkg/storage"
	storagemock "go.zenithar.org/solid/pkg/storage/mock"
	tokenmock "go.zenithar.org/solid/pkg/token/mock"

	"github.com/golang/mock/gomock"
	"github.com/google/go-cmp/cmp"
)

func Test_service_refreshToken(t *testing.T) {
	type args struct {
		ctx    context.Context
		client *corev1.Client
		req    *corev1.TokenRequest
	}
	tests := []struct {
		name    string
		args    args
		prepare func(*storagemock.MockToken, *tokenmock.MockAccessTokenGenerator)
		want    *corev1.TokenResponse
		wantErr bool
	}{
		{
			name: "nil client",
			args: args{
				ctx: context.Background(),
				req: &corev1.TokenRequest{
					GrantType: oidc.GrantTypeClientCredentials,
					Grant: &corev1.TokenRequest_ClientCredentials{
						ClientCredentials: &corev1.GrantClientCredentials{},
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
				client: &corev1.Client{},
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
				client: &corev1.Client{},
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
			name: "empty refresh_token",
			args: args{
				ctx:    context.Background(),
				client: &corev1.Client{},
				req: &corev1.TokenRequest{
					Client: &corev1.Client{
						ClientId: "s6BhdRkqt3",
					},
					GrantType: oidc.GrantTypeRefreshToken,
					Grant: &corev1.TokenRequest_RefreshToken{
						RefreshToken: &corev1.GrantRefreshToken{
							RefreshToken: "",
						},
					},
				},
			},
			wantErr: true,
			want: &corev1.TokenResponse{
				Error: rfcerrors.InvalidRequest(""),
			},
		},
		{
			name: "client not support grant_type",
			args: args{
				ctx: context.Background(),
				client: &corev1.Client{
					GrantTypes: []string{oidc.GrantTypeAuthorizationCode},
				},
				req: &corev1.TokenRequest{
					Client: &corev1.Client{
						ClientId: "s6BhdRkqt3",
					},
					GrantType: oidc.GrantTypeRefreshToken,
					Grant: &corev1.TokenRequest_RefreshToken{
						RefreshToken: &corev1.GrantRefreshToken{
							RefreshToken: "LHT.djeMMoErRAsLuXLlDYZDGdodfVLOduDi",
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
			name: "refresh token not found",
			args: args{
				ctx: context.Background(),
				client: &corev1.Client{
					GrantTypes: []string{oidc.GrantTypeRefreshToken},
				},
				req: &corev1.TokenRequest{
					Client: &corev1.Client{
						ClientId: "s6BhdRkqt3",
					},
					GrantType: oidc.GrantTypeRefreshToken,
					Grant: &corev1.TokenRequest_RefreshToken{
						RefreshToken: &corev1.GrantRefreshToken{
							RefreshToken: "LHT.djeMMoErRAsLuXLlDYZDGdodfVLOduDi",
						},
					},
				},
			},
			prepare: func(tokens *storagemock.MockToken, _ *tokenmock.MockAccessTokenGenerator) {
				tokens.EXPECT().GetByValue(gomock.Any(), "LHT.djeMMoErRAsLuXLlDYZDGdodfVLOduDi").Return(nil, storage.ErrNotFound)
			},
			wantErr: true,
			want: &corev1.TokenResponse{
				Error: rfcerrors.InvalidRequest(""),
			},
		},
		{
			name: "refresh token storage error",
			args: args{
				ctx: context.Background(),
				client: &corev1.Client{
					GrantTypes: []string{oidc.GrantTypeRefreshToken},
				},
				req: &corev1.TokenRequest{
					Client: &corev1.Client{
						ClientId: "s6BhdRkqt3",
					},
					GrantType: oidc.GrantTypeRefreshToken,
					Grant: &corev1.TokenRequest_RefreshToken{
						RefreshToken: &corev1.GrantRefreshToken{
							RefreshToken: "LHT.djeMMoErRAsLuXLlDYZDGdodfVLOduDi",
						},
					},
				},
			},
			prepare: func(tokens *storagemock.MockToken, _ *tokenmock.MockAccessTokenGenerator) {
				tokens.EXPECT().GetByValue(gomock.Any(), "LHT.djeMMoErRAsLuXLlDYZDGdodfVLOduDi").Return(nil, fmt.Errorf("foo"))
			},
			wantErr: true,
			want: &corev1.TokenResponse{
				Error: rfcerrors.ServerError(""),
			},
		},
		{
			name: "refresh token is not active",
			args: args{
				ctx: context.Background(),
				client: &corev1.Client{
					GrantTypes: []string{oidc.GrantTypeRefreshToken},
				},
				req: &corev1.TokenRequest{
					Client: &corev1.Client{
						ClientId: "s6BhdRkqt3",
					},
					GrantType: oidc.GrantTypeRefreshToken,
					Grant: &corev1.TokenRequest_RefreshToken{
						RefreshToken: &corev1.GrantRefreshToken{
							RefreshToken: "LHT.djeMMoErRAsLuXLlDYZDGdodfVLOduDi",
						},
					},
				},
			},
			prepare: func(tokens *storagemock.MockToken, _ *tokenmock.MockAccessTokenGenerator) {
				tokens.EXPECT().GetByValue(gomock.Any(), "LHT.djeMMoErRAsLuXLlDYZDGdodfVLOduDi").Return(&corev1.Token{
					Value:     "LHT.djeMMoErRAsLuXLlDYZDGdodfVLOduDi",
					TokenId:   "0123456789",
					Status:    corev1.TokenStatus_TOKEN_STATUS_REVOKED,
					TokenType: corev1.TokenType_TOKEN_TYPE_ACCESS_TOKEN,
				}, nil)
			},
			wantErr: true,
			want: &corev1.TokenResponse{
				Error: rfcerrors.InvalidRequest(""),
			},
		},
		{
			name: "refresh token is not a refresh_token",
			args: args{
				ctx: context.Background(),
				client: &corev1.Client{
					GrantTypes: []string{oidc.GrantTypeRefreshToken},
				},
				req: &corev1.TokenRequest{
					Client: &corev1.Client{
						ClientId: "s6BhdRkqt3",
					},
					GrantType: oidc.GrantTypeRefreshToken,
					Grant: &corev1.TokenRequest_RefreshToken{
						RefreshToken: &corev1.GrantRefreshToken{
							RefreshToken: "LHT.djeMMoErRAsLuXLlDYZDGdodfVLOduDi",
						},
					},
				},
			},
			prepare: func(tokens *storagemock.MockToken, _ *tokenmock.MockAccessTokenGenerator) {
				tokens.EXPECT().GetByValue(gomock.Any(), "LHT.djeMMoErRAsLuXLlDYZDGdodfVLOduDi").Return(&corev1.Token{
					Value:     "LHT.djeMMoErRAsLuXLlDYZDGdodfVLOduDi",
					TokenId:   "0123456789",
					Status:    corev1.TokenStatus_TOKEN_STATUS_ACTIVE,
					TokenType: corev1.TokenType_TOKEN_TYPE_ACCESS_TOKEN,
				}, nil)
			},
			wantErr: true,
			want: &corev1.TokenResponse{
				Error: rfcerrors.InvalidRequest(""),
			},
		},
		{
			name: "refresh token doesn't have metadata",
			args: args{
				ctx: context.Background(),
				client: &corev1.Client{
					GrantTypes: []string{oidc.GrantTypeRefreshToken},
				},
				req: &corev1.TokenRequest{
					Client: &corev1.Client{
						ClientId: "s6BhdRkqt3",
					},
					GrantType: oidc.GrantTypeRefreshToken,
					Grant: &corev1.TokenRequest_RefreshToken{
						RefreshToken: &corev1.GrantRefreshToken{
							RefreshToken: "LHT.djeMMoErRAsLuXLlDYZDGdodfVLOduDi",
						},
					},
				},
			},
			prepare: func(tokens *storagemock.MockToken, _ *tokenmock.MockAccessTokenGenerator) {
				tokens.EXPECT().GetByValue(gomock.Any(), "LHT.djeMMoErRAsLuXLlDYZDGdodfVLOduDi").Return(&corev1.Token{
					Value:     "LHT.djeMMoErRAsLuXLlDYZDGdodfVLOduDi",
					TokenId:   "0123456789",
					Status:    corev1.TokenStatus_TOKEN_STATUS_ACTIVE,
					TokenType: corev1.TokenType_TOKEN_TYPE_REFRESH_TOKEN,
				}, nil)
			},
			wantErr: true,
			want: &corev1.TokenResponse{
				Error: rfcerrors.ServerError(""),
			},
		},
		{
			name: "refresh token expired",
			args: args{
				ctx: context.Background(),
				client: &corev1.Client{
					GrantTypes: []string{oidc.GrantTypeRefreshToken},
				},
				req: &corev1.TokenRequest{
					Client: &corev1.Client{
						ClientId: "s6BhdRkqt3",
					},
					GrantType: oidc.GrantTypeRefreshToken,
					Grant: &corev1.TokenRequest_RefreshToken{
						RefreshToken: &corev1.GrantRefreshToken{
							RefreshToken: "LHT.djeMMoErRAsLuXLlDYZDGdodfVLOduDi",
						},
					},
				},
			},
			prepare: func(tokens *storagemock.MockToken, _ *tokenmock.MockAccessTokenGenerator) {
				timeFunc = func() time.Time { return time.Unix(100, 0) }
				tokens.EXPECT().GetByValue(gomock.Any(), "LHT.djeMMoErRAsLuXLlDYZDGdodfVLOduDi").Return(&corev1.Token{
					Value:     "LHT.djeMMoErRAsLuXLlDYZDGdodfVLOduDi",
					TokenId:   "0123456789",
					Status:    corev1.TokenStatus_TOKEN_STATUS_ACTIVE,
					TokenType: corev1.TokenType_TOKEN_TYPE_REFRESH_TOKEN,
					Metadata: &corev1.TokenMeta{
						Audience:  "mDuGcLjmamjNpLmYZMLIshFcXUDCNDcH",
						Scope:     "openid profile email offline_access",
						IssuedAt:  1,
						ExpiresAt: 2,
					},
				}, nil)
			},
			wantErr: true,
			want: &corev1.TokenResponse{
				Error: rfcerrors.InvalidRequest(""),
			},
		},
		{
			name: "refresh token / client_id mismatch",
			args: args{
				ctx: context.Background(),
				client: &corev1.Client{
					GrantTypes: []string{oidc.GrantTypeRefreshToken},
				},
				req: &corev1.TokenRequest{
					Client: &corev1.Client{
						ClientId: "s6BhdRkqt3",
					},
					GrantType: oidc.GrantTypeRefreshToken,
					Grant: &corev1.TokenRequest_RefreshToken{
						RefreshToken: &corev1.GrantRefreshToken{
							RefreshToken: "LHT.djeMMoErRAsLuXLlDYZDGdodfVLOduDi",
						},
					},
				},
			},
			prepare: func(tokens *storagemock.MockToken, _ *tokenmock.MockAccessTokenGenerator) {
				timeFunc = func() time.Time { return time.Unix(1, 0) }
				tokens.EXPECT().GetByValue(gomock.Any(), "LHT.djeMMoErRAsLuXLlDYZDGdodfVLOduDi").Return(&corev1.Token{
					Value:     "LHT.djeMMoErRAsLuXLlDYZDGdodfVLOduDi",
					TokenId:   "0123456789",
					Status:    corev1.TokenStatus_TOKEN_STATUS_ACTIVE,
					TokenType: corev1.TokenType_TOKEN_TYPE_REFRESH_TOKEN,
					Metadata: &corev1.TokenMeta{
						ClientId:  "123458",
						Audience:  "mDuGcLjmamjNpLmYZMLIshFcXUDCNDcH",
						Scope:     "openid profile email offline_access",
						IssuedAt:  1,
						ExpiresAt: 604801,
					},
				}, nil)
			},
			wantErr: true,
			want: &corev1.TokenResponse{
				Error: rfcerrors.InvalidRequest(""),
			},
		},
		// ---------------------------------------------------------------------
		{
			name: "error during accessToken generation",
			args: args{
				ctx: context.Background(),
				client: &corev1.Client{
					GrantTypes: []string{oidc.GrantTypeRefreshToken},
				},
				req: &corev1.TokenRequest{
					Client: &corev1.Client{
						ClientId: "s6BhdRkqt3",
					},
					GrantType: oidc.GrantTypeRefreshToken,
					Grant: &corev1.TokenRequest_RefreshToken{
						RefreshToken: &corev1.GrantRefreshToken{
							RefreshToken: "LHT.djeMMoErRAsLuXLlDYZDGdodfVLOduDi",
						},
					},
				},
			},
			prepare: func(tokens *storagemock.MockToken, at *tokenmock.MockAccessTokenGenerator) {
				timeFunc = func() time.Time { return time.Unix(1, 0) }
				tokens.EXPECT().GetByValue(gomock.Any(), "LHT.djeMMoErRAsLuXLlDYZDGdodfVLOduDi").Return(&corev1.Token{
					Value:     "LHT.djeMMoErRAsLuXLlDYZDGdodfVLOduDi",
					TokenId:   "0123456789",
					TokenType: corev1.TokenType_TOKEN_TYPE_REFRESH_TOKEN,
					Status:    corev1.TokenStatus_TOKEN_STATUS_ACTIVE,
					Metadata: &corev1.TokenMeta{
						Audience:  "mDuGcLjmamjNpLmYZMLIshFcXUDCNDcH",
						Scope:     "openid profile email offline_access",
						IssuedAt:  1,
						ExpiresAt: 604801,
					},
				}, nil)
				at.EXPECT().Generate(gomock.Any(), gomock.Any(), gomock.Any()).Return("", fmt.Errorf("foo"))
			},
			wantErr: true,
			want: &corev1.TokenResponse{
				Error: rfcerrors.ServerError(""),
			},
		},
		{
			name: "empty access token value",
			args: args{
				ctx: context.Background(),
				client: &corev1.Client{
					GrantTypes: []string{oidc.GrantTypeRefreshToken},
				},
				req: &corev1.TokenRequest{
					Client: &corev1.Client{
						ClientId: "s6BhdRkqt3",
					},
					GrantType: oidc.GrantTypeRefreshToken,
					Grant: &corev1.TokenRequest_RefreshToken{
						RefreshToken: &corev1.GrantRefreshToken{
							RefreshToken: "LHT.djeMMoErRAsLuXLlDYZDGdodfVLOduDi",
						},
					},
				},
			},
			prepare: func(tokens *storagemock.MockToken, at *tokenmock.MockAccessTokenGenerator) {
				timeFunc = func() time.Time { return time.Unix(1, 0) }
				tokens.EXPECT().GetByValue(gomock.Any(), "LHT.djeMMoErRAsLuXLlDYZDGdodfVLOduDi").Return(&corev1.Token{
					Value:     "LHT.djeMMoErRAsLuXLlDYZDGdodfVLOduDi",
					TokenId:   "0123456789",
					TokenType: corev1.TokenType_TOKEN_TYPE_REFRESH_TOKEN,
					Status:    corev1.TokenStatus_TOKEN_STATUS_ACTIVE,
					Metadata: &corev1.TokenMeta{
						Audience:  "mDuGcLjmamjNpLmYZMLIshFcXUDCNDcH",
						Scope:     "openid profile email offline_access",
						IssuedAt:  1,
						ExpiresAt: 604801,
					},
				}, nil)
				at.EXPECT().Generate(gomock.Any(), gomock.Any(), gomock.Any()).Return("", nil)
			},
			wantErr: true,
			want: &corev1.TokenResponse{
				Error: rfcerrors.ServerError(""),
			},
		},
		{
			name: "token storage error",
			args: args{
				ctx: context.Background(),
				client: &corev1.Client{
					GrantTypes: []string{oidc.GrantTypeRefreshToken},
				},
				req: &corev1.TokenRequest{
					Client: &corev1.Client{
						ClientId: "s6BhdRkqt3",
					},
					GrantType: oidc.GrantTypeRefreshToken,
					Grant: &corev1.TokenRequest_RefreshToken{
						RefreshToken: &corev1.GrantRefreshToken{
							RefreshToken: "LHT.djeMMoErRAsLuXLlDYZDGdodfVLOduDi",
						},
					},
				},
			},
			prepare: func(tokens *storagemock.MockToken, at *tokenmock.MockAccessTokenGenerator) {
				timeFunc = func() time.Time { return time.Unix(1, 0) }
				tokens.EXPECT().GetByValue(gomock.Any(), "LHT.djeMMoErRAsLuXLlDYZDGdodfVLOduDi").Return(&corev1.Token{
					Value:     "LHT.djeMMoErRAsLuXLlDYZDGdodfVLOduDi",
					TokenId:   "0123456789",
					TokenType: corev1.TokenType_TOKEN_TYPE_REFRESH_TOKEN,
					Status:    corev1.TokenStatus_TOKEN_STATUS_ACTIVE,
					Metadata: &corev1.TokenMeta{
						Audience:  "mDuGcLjmamjNpLmYZMLIshFcXUDCNDcH",
						Scope:     "openid profile email offline_access",
						IssuedAt:  1,
						ExpiresAt: 604801,
					},
				}, nil)
				at.EXPECT().Generate(gomock.Any(), gomock.Any(), gomock.Any()).Return("xtU.GvmXVrPVNqSnHjpZbEarIqOPAlfXfQpM", nil)
				tokens.EXPECT().Create(gomock.Any(), gomock.Any()).Return(fmt.Errorf("foo"))
			},
			wantErr: true,
			want: &corev1.TokenResponse{
				Error: rfcerrors.ServerError(""),
			},
		},
		{
			name: "rt generation error",
			args: args{
				ctx: context.Background(),
				client: &corev1.Client{
					GrantTypes: []string{oidc.GrantTypeRefreshToken},
				},
				req: &corev1.TokenRequest{
					Client: &corev1.Client{
						ClientId: "s6BhdRkqt3",
					},
					GrantType: oidc.GrantTypeRefreshToken,
					Grant: &corev1.TokenRequest_RefreshToken{
						RefreshToken: &corev1.GrantRefreshToken{
							RefreshToken: "LHT.djeMMoErRAsLuXLlDYZDGdodfVLOduDi",
						},
					},
				},
			},
			prepare: func(tokens *storagemock.MockToken, at *tokenmock.MockAccessTokenGenerator) {
				timeFunc = func() time.Time { return time.Unix(1, 0) }
				tokens.EXPECT().GetByValue(gomock.Any(), "LHT.djeMMoErRAsLuXLlDYZDGdodfVLOduDi").Return(&corev1.Token{
					Value:     "LHT.djeMMoErRAsLuXLlDYZDGdodfVLOduDi",
					TokenId:   "0123456789",
					TokenType: corev1.TokenType_TOKEN_TYPE_REFRESH_TOKEN,
					Status:    corev1.TokenStatus_TOKEN_STATUS_ACTIVE,
					Metadata: &corev1.TokenMeta{
						Audience:  "mDuGcLjmamjNpLmYZMLIshFcXUDCNDcH",
						Scope:     "openid profile email offline_access",
						IssuedAt:  1,
						ExpiresAt: 2,
					},
				}, nil)
				atGen := at.EXPECT().Generate(gomock.Any(), gomock.Any(), gomock.Any()).Return("xtU.GvmXVrPVNqSnHjpZbEarIqOPAlfXfQpM", nil)
				atSave := tokens.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil)
				at.EXPECT().Generate(gomock.Any(), gomock.Any(), gomock.Any()).Return("JHP.HscxBIrTOYZWgupVlrABwkdbhtqVFrmr", nil).After(atGen)
				tokens.EXPECT().Create(gomock.Any(), gomock.Any()).Return(fmt.Errorf("foo")).After(atSave)
			},
			wantErr: true,
			want: &corev1.TokenResponse{
				Error: rfcerrors.ServerError(""),
			},
		},
		{
			name: "rt revocation error",
			args: args{
				ctx: context.Background(),
				client: &corev1.Client{
					GrantTypes: []string{oidc.GrantTypeRefreshToken},
				},
				req: &corev1.TokenRequest{
					Client: &corev1.Client{
						ClientId: "s6BhdRkqt3",
					},
					GrantType: oidc.GrantTypeRefreshToken,
					Grant: &corev1.TokenRequest_RefreshToken{
						RefreshToken: &corev1.GrantRefreshToken{
							RefreshToken: "LHT.djeMMoErRAsLuXLlDYZDGdodfVLOduDi",
						},
					},
				},
			},
			prepare: func(tokens *storagemock.MockToken, at *tokenmock.MockAccessTokenGenerator) {
				timeFunc = func() time.Time { return time.Unix(1, 0) }
				tokens.EXPECT().GetByValue(gomock.Any(), "LHT.djeMMoErRAsLuXLlDYZDGdodfVLOduDi").Return(&corev1.Token{
					Value:     "LHT.djeMMoErRAsLuXLlDYZDGdodfVLOduDi",
					TokenId:   "0123456789",
					TokenType: corev1.TokenType_TOKEN_TYPE_REFRESH_TOKEN,
					Status:    corev1.TokenStatus_TOKEN_STATUS_ACTIVE,
					Metadata: &corev1.TokenMeta{
						Audience:  "mDuGcLjmamjNpLmYZMLIshFcXUDCNDcH",
						Scope:     "openid profile email offline_access",
						IssuedAt:  1,
						ExpiresAt: 2,
					},
				}, nil)
				atGen := at.EXPECT().Generate(gomock.Any(), gomock.Any(), gomock.Any()).Return("xtU.GvmXVrPVNqSnHjpZbEarIqOPAlfXfQpM", nil)
				atSave := tokens.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil)
				at.EXPECT().Generate(gomock.Any(), gomock.Any(), gomock.Any()).Return("JHP.HscxBIrTOYZWgupVlrABwkdbhtqVFrmr", nil).After(atGen)
				tokens.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil).After(atSave)
				tokens.EXPECT().Revoke(gomock.Any(), "0123456789").Return(fmt.Errorf("foo"))
			},
			wantErr: true,
			want: &corev1.TokenResponse{
				Error: rfcerrors.ServerError(""),
			},
		},
		// ---------------------------------------------------------------------
		{
			name: "valid",
			args: args{
				ctx: context.Background(),
				client: &corev1.Client{
					GrantTypes: []string{oidc.GrantTypeRefreshToken},
				},
				req: &corev1.TokenRequest{
					Client: &corev1.Client{
						ClientId: "s6BhdRkqt3",
					},
					GrantType: oidc.GrantTypeRefreshToken,
					Grant: &corev1.TokenRequest_RefreshToken{
						RefreshToken: &corev1.GrantRefreshToken{
							RefreshToken: "LHT.djeMMoErRAsLuXLlDYZDGdodfVLOduDi",
						},
					},
				},
			},
			prepare: func(tokens *storagemock.MockToken, at *tokenmock.MockAccessTokenGenerator) {
				timeFunc = func() time.Time { return time.Unix(1, 0) }
				tokens.EXPECT().GetByValue(gomock.Any(), "LHT.djeMMoErRAsLuXLlDYZDGdodfVLOduDi").Return(&corev1.Token{
					Value:     "LHT.djeMMoErRAsLuXLlDYZDGdodfVLOduDi",
					TokenId:   "0123456789",
					TokenType: corev1.TokenType_TOKEN_TYPE_REFRESH_TOKEN,
					Status:    corev1.TokenStatus_TOKEN_STATUS_ACTIVE,
					Metadata: &corev1.TokenMeta{
						Audience:  "mDuGcLjmamjNpLmYZMLIshFcXUDCNDcH",
						Scope:     "openid profile email offline_access",
						IssuedAt:  1,
						ExpiresAt: 604801,
					},
				}, nil)
				at.EXPECT().Generate(gomock.Any(), gomock.Any(), gomock.Any()).Return("xtU.GvmXVrPVNqSnHjpZbEarIqOPAlfXfQpM", nil)
				tokens.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil)
			},
			wantErr: false,
			want: &corev1.TokenResponse{
				AccessToken: &corev1.Token{
					Value:     "xtU.GvmXVrPVNqSnHjpZbEarIqOPAlfXfQpM",
					TokenId:   "0123456789",
					TokenType: corev1.TokenType_TOKEN_TYPE_ACCESS_TOKEN,
					Status:    corev1.TokenStatus_TOKEN_STATUS_ACTIVE,
					Metadata: &corev1.TokenMeta{
						Audience:  "mDuGcLjmamjNpLmYZMLIshFcXUDCNDcH",
						Scope:     "openid profile email offline_access",
						IssuedAt:  1,
						ExpiresAt: 3601,
					},
				},
			},
		},
		{
			name: "valid with new rt",
			args: args{
				ctx: context.Background(),
				client: &corev1.Client{
					GrantTypes: []string{oidc.GrantTypeRefreshToken},
				},
				req: &corev1.TokenRequest{
					Client: &corev1.Client{
						ClientId: "s6BhdRkqt3",
					},
					GrantType: oidc.GrantTypeRefreshToken,
					Grant: &corev1.TokenRequest_RefreshToken{
						RefreshToken: &corev1.GrantRefreshToken{
							RefreshToken: "LHT.djeMMoErRAsLuXLlDYZDGdodfVLOduDi",
						},
					},
				},
			},
			prepare: func(tokens *storagemock.MockToken, at *tokenmock.MockAccessTokenGenerator) {
				timeFunc = func() time.Time { return time.Unix(1, 0) }
				tokens.EXPECT().GetByValue(gomock.Any(), "LHT.djeMMoErRAsLuXLlDYZDGdodfVLOduDi").Return(&corev1.Token{
					Value:     "LHT.djeMMoErRAsLuXLlDYZDGdodfVLOduDi",
					TokenId:   "0123456789",
					TokenType: corev1.TokenType_TOKEN_TYPE_REFRESH_TOKEN,
					Status:    corev1.TokenStatus_TOKEN_STATUS_ACTIVE,
					Metadata: &corev1.TokenMeta{
						Audience:  "mDuGcLjmamjNpLmYZMLIshFcXUDCNDcH",
						Scope:     "openid profile email offline_access",
						IssuedAt:  1,
						ExpiresAt: 2,
					},
				}, nil)
				atGen := at.EXPECT().Generate(gomock.Any(), gomock.Any(), gomock.Any()).Return("xtU.GvmXVrPVNqSnHjpZbEarIqOPAlfXfQpM", nil)
				atSave := tokens.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil)
				at.EXPECT().Generate(gomock.Any(), gomock.Any(), gomock.Any()).Return("JHP.HscxBIrTOYZWgupVlrABwkdbhtqVFrmr", nil).After(atGen)
				tokens.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil).After(atSave)
				tokens.EXPECT().Revoke(gomock.Any(), "0123456789").Return(nil)
			},
			wantErr: false,
			want: &corev1.TokenResponse{
				AccessToken: &corev1.Token{
					Value:     "xtU.GvmXVrPVNqSnHjpZbEarIqOPAlfXfQpM",
					TokenId:   "0123456789",
					TokenType: corev1.TokenType_TOKEN_TYPE_ACCESS_TOKEN,
					Status:    corev1.TokenStatus_TOKEN_STATUS_ACTIVE,
					Metadata: &corev1.TokenMeta{
						Audience:  "mDuGcLjmamjNpLmYZMLIshFcXUDCNDcH",
						Scope:     "openid profile email offline_access",
						IssuedAt:  1,
						ExpiresAt: 3601,
					},
				},
				RefreshToken: &corev1.Token{
					Value:     "JHP.HscxBIrTOYZWgupVlrABwkdbhtqVFrmr",
					TokenId:   "0123456789",
					TokenType: corev1.TokenType_TOKEN_TYPE_REFRESH_TOKEN,
					Status:    corev1.TokenStatus_TOKEN_STATUS_ACTIVE,
					Metadata: &corev1.TokenMeta{
						Audience:  "mDuGcLjmamjNpLmYZMLIshFcXUDCNDcH",
						Scope:     "openid profile email offline_access",
						IssuedAt:  1,
						ExpiresAt: 604801,
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			// Arm mocks
			sessions := storagemock.NewMockSession(ctrl)
			accessTokens := tokenmock.NewMockAccessTokenGenerator(ctrl)
			idTokens := tokenmock.NewMockIDTokenGenerator(ctrl)
			tokens := storagemock.NewMockToken(ctrl)

			// Prepare them
			if tt.prepare != nil {
				tt.prepare(tokens, accessTokens)
			}

			s := &service{
				sessions:             sessions,
				tokens:               tokens,
				accessTokenGenerator: accessTokens,
				idTokenGenerator:     idTokens,
			}
			got, err := s.refreshToken(tt.args.ctx, tt.args.client, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("service.refreshToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if diff := cmp.Diff(got, tt.want, cmpOpts...); diff != "" {
				t.Errorf("service.refreshToken() res = %s", diff)
			}
		})
	}
}
