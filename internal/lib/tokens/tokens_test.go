package tokens_test

import (
	"testing"

	"github.com/Onnywrite/grpc-auth/internal/lib/tokens"
	"github.com/Onnywrite/grpc-auth/internal/models"
	"github.com/stretchr/testify/assert"
)

const (
	secret = "123aboba321"
)

func TestRefresh(t *testing.T) {
	tests := []struct {
		name string
		tkn  *models.RefreshToken
		exp  string
		err  error
	}{
		{
			name: "OK",
			tkn: &models.RefreshToken{
				SessionUUID: "23897ac4-6a95-4d9a-a7f1-264e3ce10084",
				Rotation:    0,
				Exp:         123456789,
			},
			exp: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjEyMzQ1Njc4OSwicm90YXRpb24iOjAsInNlc3Npb25fdXVpZCI6IjIzODk3YWM0LTZhOTUtNGQ5YS1hN2YxLTI2NGUzY2UxMDA4NCJ9.lGG9c8aR3URylFAs05uPK59wqMg2vDLRdgsGrSJ5bVY",
			err: nil,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(tt *testing.T) {
			exp, err := tokens.Refresh(tc.tkn, secret)
			assert.ErrorIs(tt, err, tc.err)
			assert.Equal(tt, tc.exp, exp)
		})
	}
}

func TestParseRefresh(t *testing.T) {
	tests := []struct {
		name   string
		tkn    string
		exp    *models.RefreshToken
		expErr bool
		hasErr string
	}{
		{
			name: "OK",
			tkn:  "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjM0OTYxMTQ0NDAwLCJyb3RhdGlvbiI6MCwic2Vzc2lvbl91dWlkIjoiMzQzZWJlMTMtNzVhYy00ODkwLWJjYTAtYjdkMTFjZDk2MWM3In0.sd5BHEafPaVFefTolDZge2Nj6X5dQrP0W-QfCwb_2AE",
			exp: &models.RefreshToken{
				SessionUUID: "343ebe13-75ac-4890-bca0-b7d11cd961c7",
				Rotation:    0,
				Exp:         34961144400,
			},
			expErr: false,
		},
		{
			name: "Expired",
			tkn:  "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjEyMzQ1Njc4OSwicm90YXRpb24iOjAsInNlc3Npb25fdXVpZCI6IjIzODk3YWM0LTZhOTUtNGQ5YS1hN2YxLTI2NGUzY2UxMDA4NCJ9.lGG9c8aR3URylFAs05uPK59wqMg2vDLRdgsGrSJ5bVY",
			exp: &models.RefreshToken{
				SessionUUID: "23897ac4-6a95-4d9a-a7f1-264e3ce10084",
				Rotation:    0,
				Exp:         123456789,
			},
			expErr: true,
			hasErr: tokens.ErrTokenExpired,
		},
		{
			name:   "Invalid exp",
			tkn:    "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOiIxMjM0NTY3ODkiLCJyb3RhdGlvbiI6MCwic2Vzc2lvbl91dWlkIjoiMjM4OTdhYzQtNmE5NS00ZDlhLWE3ZjEtMjY0ZTNjZTEwMDg0In0.X1SnQOKuEMcVUO45MbJZ3B1IZFCmZUJTGPqkXb-XiLU",
			exp:    nil,
			expErr: true,
			hasErr: tokens.ErrInvalidData,
		},
		{
			name:   "Invalid UUID",
			tkn:    "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjEyMzQ1Njc4OSwicm90YXRpb24iOjAsInNlc3Npb25fdXVpZCI6MjM4OTd9.son34AxbDXgeEM9vb3-r693WmSQLlYzffLfri4fGt_0",
			exp:    nil,
			expErr: true,
			hasErr: tokens.ErrInvalidData,
		},
		{
			name:   "Invalid rotation",
			tkn:    "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjEyMzQ1Njc4OSwicm90YXRpb24iOlsicGVhY2UiXSwic2Vzc2lvbl91dWlkIjoiMjM4OTdhYzQtNmE5NS00ZDlhLWE3ZjEtMjY0ZTNjZTEwMDg0In0.X_Ofb7zxZgq_dzZfNKI4sFfIidHF-uBQ-E-sy5n9qIQ",
			exp:    nil,
			expErr: true,
			hasErr: tokens.ErrInvalidData,
		},
		// {
		// 	name: "Unexpected signing method",
		// 	tkn:  "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjEyMzQ1Njc4OSwicm90YXRpb24iOjEsInNlc3Npb25fdXVpZCI6IjIzODk3YWM0LTZhOTUtNGQ5YS1hN2YxLTI2NGUzY2UxMDA4NCJ9.nSc2T0rVY-tbR6Ozpy6s5k5bMtk5rhCUwiDncWMRNWVFi_395KQP20Nknqc_A3ZGK7dJRPuVc-0wfLtJaYN2FCTpknVPM8KM4JwSPg0OKGKPEdJdTVEvcJOgn7FBtAkUrLgqPizo89AzW7kgqy7W56T9wjc33uO2GPBgZpBeGs6ipEGhr6tkKN3RB9YH6FKYi2AXan3h-s7q-K2NByKRiPMin6Y9aYuHBTLoeiI-G12nnv3SDONFHFO2VFgp1F6rZnRlebEAuRAFwo_0-3YgjSUeAPGZL-O34sePcnkPDxuVnn7M4IGAY61Dc9dnHTBS9FIsZoeceGVvwpxnpD_2fxhVYIBpQGUK0cTlwui1O-XIr9KCe1EmuTwAyc4kSmRqxFHzvMhjacXnnAjNdoxin15vZhyjv31B_XeJWTaK9My6CtT05k40pFI0IcK2BMfsta3c0o-uAIPz4s5wODO26ajHX8t-x9ccxt84K7WLN7_04gjG3_5Ir2QOuJnW3fuI9_pxcx5Q3hZAk-jxyCNTrwQsHeg6nfzltIOLrKjq8pQgoM0r4II9sMC7Cb-H-OuBSpZuVUKK2j0yH3EzsYGbQQXPt8mk1oFxPmupJ6njrVdOmPFTiD44adSgHQxPnJhUio_S-LU5ZbjwuXt8yw6bHgk80vDDqy9uFx_TZCw8qkU",
		// 	exp:  nil,
		// 	err:  tokens.ErrUnexpectedSigningMethod,
		// },
	}

	t.Setenv("TOKEN_SECRET", "123aboba321")

	for _, tc := range tests {
		t.Run(tc.name, func(tt *testing.T) {
			exp, err := tokens.ParseRefresh(tc.tkn, secret)
			if tc.expErr {
				assert.True(tt, err.Has(tc.hasErr))
			}
			assert.EqualValues(tt, tc.exp, exp)
		})
	}
}

func TestAccess(t *testing.T) {
	tests := []struct {
		name   string
		tkn    *models.AccessToken
		exp    string
		expErr bool
		hasErr string
	}{
		{
			name: "OK",
			tkn: &models.AccessToken{
				Id:        660,
				ServiceId: 1220,
				Roles:     []string{"role1", "role2"},
				Exp:       987654321,
			},
			exp:    "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjk4NzY1NDMyMSwiaWQiOjY2MCwicm9sZXMiOlsicm9sZTEiLCJyb2xlMiJdLCJzZXJ2aWNlX2lkIjoxMjIwfQ.eGDFeVAaR70z76_-9gy21VDtp8hIX0woRKF4wlC4BDQ",
			expErr: false,
		},
	}

	t.Setenv("TOKEN_SECRET", "123aboba321")

	for _, tc := range tests {
		t.Run(tc.name, func(tt *testing.T) {
			exp, err := tokens.Access(tc.tkn, secret)
			if tc.expErr {
				assert.True(tt, err.Has(tc.hasErr))
			}
			assert.Equal(tt, tc.exp, exp)
		})
	}
}

func TestParseAccess(t *testing.T) {
	tests := []struct {
		name   string
		tkn    string
		exp    *models.AccessToken
		expErr bool
		hasErr string
	}{
		{
			name: "OK",
			tkn:  "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjM0OTYxMTQ0NDAwLCJpZCI6NjYwLCJyb2xlcyI6WyJyb2xlMSIsInJvbGUyIl0sInNlcnZpY2VfaWQiOjEyMjB9.GlReJRRY6xheKmgMEH9aUGotTueovkJI0paziAjgDPQ",
			exp: &models.AccessToken{
				Id:        660,
				ServiceId: 1220,
				Roles:     []string{"role1", "role2"},
				Exp:       34961144400,
			},
			expErr: false,
		},
		{
			name: "Expired",
			tkn:  "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjk4NzY1NDMyMSwiaWQiOjY2MCwicm9sZXMiOlsicm9sZTEiLCJyb2xlMiJdLCJzZXJ2aWNlX2lkIjoxMjIwfQ.eGDFeVAaR70z76_-9gy21VDtp8hIX0woRKF4wlC4BDQ",
			exp: &models.AccessToken{
				Id:        660,
				ServiceId: 1220,
				Roles:     []string{"role1", "role2"},
				Exp:       987654321,
			},
			expErr: true,
			hasErr: tokens.ErrTokenExpired,
		},
		{
			name:   "Invalid exp",
			tkn:    "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOiIzNDk2MTE0NDQwMCIsImlkIjo2NjAsInJvbGVzIjpbInJvbGUxIiwicm9sZTIiXSwic2VydmljZV9pZCI6MTIyMH0.xUwKLt7yvKkFxweH1SZIJD_wb7VZtNj8PuPWxJNyKqM",
			exp:    nil,
			expErr: true,
			hasErr: tokens.ErrInvalidData,
		},
		{
			name:   "Invalid ID",
			tkn:    "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjM0OTYxMTQ0NDAwLCJpZCI6W10sInJvbGVzIjpbInJvbGUxIiwicm9sZTIiXSwic2VydmljZV9pZCI6MTIyMH0.iLloFr5PLdLTeIlRdQxVGKcdTlTSKoUDUbYOMFYVUDs",
			exp:    nil,
			expErr: true,
			hasErr: tokens.ErrInvalidData,
		},
		{
			name:   "Invalid roles",
			tkn:    "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjM0OTYxMTQ0NDAwLCJpZCI6NjYwLCJyb2xlcyI6NzcwLCJzZXJ2aWNlX2lkIjoxMjIwfQ.pWL8iMDQZcggtHvrzZEalIK94G_Aohd9aJFJsZ4_Mdw",
			exp:    nil,
			expErr: true,
			hasErr: tokens.ErrInvalidData,
		},
		{
			name:   "Invalid service id",
			tkn:    "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjM0OTYxMTQ0NDAwLCJpZCI6NjYwLCJyb2xlcyI6WyJyb2xlMSIsInJvbGUyIl0sInNlcnZpY2VfaWQiOnsiaWQiOjEyMjB9fQ.IkPrkIbgy_dc46H50ZAueiT0PlI_BW96eokkUl1XpFQ",
			exp:    nil,
			expErr: true,
			hasErr: tokens.ErrInvalidData,
		},
		// {
		// 	name: "Unexpected signing method",
		// 	tkn:  "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjEyMzQ1Njc4OSwicm90YXRpb24iOjEsInNlc3Npb25fdXVpZCI6IjIzODk3YWM0LTZhOTUtNGQ5YS1hN2YxLTI2NGUzY2UxMDA4NCJ9.nSc2T0rVY-tbR6Ozpy6s5k5bMtk5rhCUwiDncWMRNWVFi_395KQP20Nknqc_A3ZGK7dJRPuVc-0wfLtJaYN2FCTpknVPM8KM4JwSPg0OKGKPEdJdTVEvcJOgn7FBtAkUrLgqPizo89AzW7kgqy7W56T9wjc33uO2GPBgZpBeGs6ipEGhr6tkKN3RB9YH6FKYi2AXan3h-s7q-K2NByKRiPMin6Y9aYuHBTLoeiI-G12nnv3SDONFHFO2VFgp1F6rZnRlebEAuRAFwo_0-3YgjSUeAPGZL-O34sePcnkPDxuVnn7M4IGAY61Dc9dnHTBS9FIsZoeceGVvwpxnpD_2fxhVYIBpQGUK0cTlwui1O-XIr9KCe1EmuTwAyc4kSmRqxFHzvMhjacXnnAjNdoxin15vZhyjv31B_XeJWTaK9My6CtT05k40pFI0IcK2BMfsta3c0o-uAIPz4s5wODO26ajHX8t-x9ccxt84K7WLN7_04gjG3_5Ir2QOuJnW3fuI9_pxcx5Q3hZAk-jxyCNTrwQsHeg6nfzltIOLrKjq8pQgoM0r4II9sMC7Cb-H-OuBSpZuVUKK2j0yH3EzsYGbQQXPt8mk1oFxPmupJ6njrVdOmPFTiD44adSgHQxPnJhUio_S-LU5ZbjwuXt8yw6bHgk80vDDqy9uFx_TZCw8qkU",
		// 	exp:  nil,
		// 	err:  tokens.ErrUnexpectedSigningMethod,
		// },
	}

	t.Setenv("TOKEN_SECRET", "123aboba321")

	for _, tc := range tests {
		t.Run(tc.name, func(tt *testing.T) {
			exp, err := tokens.ParseAccess(tc.tkn, secret)
			if tc.expErr {
				assert.True(tt, err.Has(tc.hasErr))
			}
			assert.EqualValues(tt, tc.exp, exp)
		})
	}
}
