package config_test

/*
func TestMustLoadByPath(t *testing.T) {
	tests := []struct {
		name     string
		envs     map[string]string
		expected config.Config
		panic    bool
	}{
		{
			name: "everything's okay",
			envs: map[string]string{},
			expected: config.Config{
				Environment: "local",
				Conn:        "test_conn",
				TokenTTL:    10 * time.Hour,
				GRPC: config.GRPCConfig{
					Port:    5050,
					Timeout: 10 * time.Hour,
				},
			},
			panic: false,
		},

		{
			name: "1 env inside config",
			envs: map[string]string{"CONN": "new_test_conn"},
			expected: config.Config{
				Environment: "local",
				Conn:        "new_test_conn",
				TokenTTL:    10 * time.Hour,
				GRPC: config.GRPCConfig{
					Port:    5050,
					Timeout: 10 * time.Hour,
				},
			},
			panic: false,
		},

		{
			name: "1 env inside grpc config",
			envs: map[string]string{"GRPC_TIMEOUT": "30m"},
			expected: config.Config{
				Environment: "local",
				Conn:        "test_conn",
				TokenTTL:    10 * time.Hour,
				GRPC: config.GRPCConfig{
					Port:    5050,
					Timeout: 30 * time.Minute,
				},
			},
			panic: false,
		},

		{
			name: "2 envs inside both config and grpc config",
			envs: map[string]string{"CONN": "new_test_conn", "GRPC_TIMEOUT": "30m"},
			expected: config.Config{
				Environment: "local",
				Conn:        "new_test_conn",
				TokenTTL:    10 * time.Hour,
				GRPC: config.GRPCConfig{
					Port:    5050,
					Timeout: 30 * time.Minute,
				},
			},
			panic: false,
		},

		{
			name: "all envs",
			envs: map[string]string{
				"ENV":             "dev",
				"CONN":            "brand_new_test_conn",
				"MIGRATIONS_PATH": "./other_migrations",
				"TOKEN_TTL":       "2m",
				"GRPC_PORT":       "5050",
				"GRPC_TIMEOUT":    "10s"},
			expected: config.Config{
				Environment: "dev",
				Conn:        "brand_new_test_conn",
				TokenTTL:    2 * time.Minute,
				GRPC: config.GRPCConfig{
					Port:    5050,
					Timeout: 10 * time.Second,
				},
			},
			panic: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for env, val := range tt.envs {
				t.Setenv(env, val)
			}

			for env := range tt.envs {
				t.Logf("%s=%s\n", env, os.Getenv(env))
			}

			cfg := config.MustLoadByPath("./config_test.yaml")

			panic := recover()
			if (panic == nil && tt.panic) || (panic != nil && !tt.panic) {
				t.Errorf("unexpected panic")
			}

			v := reflect.ValueOf(*cfg)
			ve := reflect.ValueOf(tt.expected)
			tp := reflect.TypeOf(*cfg)
			for i := range v.NumField() {
				f := v.Field(i)
				fe := ve.Field(i)
				if !f.Equal(fe) {
					t.Errorf("unexpected %s: %s",
						tp.Field(i).Name,
						f.String(),
					)
				} else {
					t.Logf("expected %s: %s",
						tp.Field(i).Name,
						f.String(),
					)
				}
			}
		})
	}
}
*/
