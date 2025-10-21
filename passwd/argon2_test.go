// Original Licenses under Apache-2.0 by the openlane https://github.com/theopenlane
// SPDX-License-Identifier: Apache-2.0

package passwd

import (
	"testing"
)

func TestDefaultArgon2Config(t *testing.T) {
	config := DefaultArgon2Config()

	if config.Time != 1 {
		t.Errorf("DefaultArgon2Config().Time = %v, want %v", config.Time, 1)
	}

	if config.Memory != 64*1024 {
		t.Errorf("DefaultArgon2Config().Memory = %v, want %v", config.Memory, 64*1024)
	}

	if config.Threads == 0 {
		t.Error("DefaultArgon2Config().Threads = 0, want > 0")
	}

	if config.KeyLen != 32 {
		t.Errorf("DefaultArgon2Config().KeyLen = %v, want %v", config.KeyLen, 32)
	}

	if config.SaltLen != 16 {
		t.Errorf("DefaultArgon2Config().SaltLen = %v, want %v", config.SaltLen, 16)
	}
}

func TestCreateDerivedKey(t *testing.T) {
	tests := []struct {
		name     string
		password string
		wantErr  bool
	}{
		{
			name:     "valid password",
			password: "test-password-123!",
			wantErr:  false,
		},
		{
			name:     "empty password",
			password: "",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dk, err := CreateDerivedKey(tt.password)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateDerivedKey() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && !IsDerivedKey(dk) {
				t.Errorf("CreateDerivedKey() = %v, is not a valid derived key", dk)
			}
		})
	}
}

func TestCreateDerivedKeyWithConfig(t *testing.T) {
	tests := []struct {
		name     string
		password string
		config   Argon2Config
		wantErr  bool
	}{
		{
			name:     "valid password with default config",
			password: "test-password-123!",
			config:   DefaultArgon2Config(),
			wantErr:  false,
		},
		{
			name:     "valid password with custom config",
			password: "test-password-123!",
			config: Argon2Config{
				Time:    2,
				Memory:  128 * 1024,
				Threads: 4,
				KeyLen:  32,
				SaltLen: 16,
			},
			wantErr: false,
		},
		{
			name:     "empty password",
			password: "",
			config:   DefaultArgon2Config(),
			wantErr:  true,
		},
		{
			name:     "zero memory",
			password: "test-password-123!",
			config: Argon2Config{
				Time:    1,
				Memory:  0,
				Threads: 1,
				KeyLen:  32,
				SaltLen: 16,
			},
			wantErr: true,
		},
		{
			name:     "zero threads",
			password: "test-password-123!",
			config: Argon2Config{
				Time:    1,
				Memory:  64 * 1024,
				Threads: 0,
				KeyLen:  32,
				SaltLen: 16,
			},
			wantErr: true,
		},
		{
			name:     "zero key length",
			password: "test-password-123!",
			config: Argon2Config{
				Time:    1,
				Memory:  64 * 1024,
				Threads: 1,
				KeyLen:  0,
				SaltLen: 16,
			},
			wantErr: true,
		},
		{
			name:     "zero salt length",
			password: "test-password-123!",
			config: Argon2Config{
				Time:    1,
				Memory:  64 * 1024,
				Threads: 1,
				KeyLen:  32,
				SaltLen: 0,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dk, err := CreateDerivedKeyWithConfig(tt.password, tt.config)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateDerivedKeyWithConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if !IsDerivedKey(dk) {
					t.Errorf("CreateDerivedKeyWithConfig() = %v, is not a valid derived key", dk)
				}

				// Verify the configuration was applied
				config, err := GetDerivedKeyConfig(dk)
				if err != nil {
					t.Errorf("GetDerivedKeyConfig() error = %v", err)
					return
				}

				if config.Time != tt.config.Time {
					t.Errorf("config.Time = %v, want %v", config.Time, tt.config.Time)
				}

				if config.Memory != tt.config.Memory {
					t.Errorf("config.Memory = %v, want %v", config.Memory, tt.config.Memory)
				}

				if config.Threads != tt.config.Threads {
					t.Errorf("config.Threads = %v, want %v", config.Threads, tt.config.Threads)
				}
			}
		})
	}
}

func TestVerifyDerivedKey(t *testing.T) {
	password := "test-password-123!"

	dk, err := CreateDerivedKey(password)
	if err != nil {
		t.Fatalf("CreateDerivedKey() error = %v", err)
	}

	tests := []struct {
		name     string
		dk       string
		password string
		want     bool
		wantErr  bool
	}{
		{
			name:     "valid password",
			dk:       dk,
			password: password,
			want:     true,
			wantErr:  false,
		},
		{
			name:     "invalid password",
			dk:       dk,
			password: "wrong-password",
			want:     false,
			wantErr:  false,
		},
		{
			name:     "empty derived key",
			dk:       "",
			password: password,
			want:     false,
			wantErr:  true,
		},
		{
			name:     "empty password",
			dk:       dk,
			password: "",
			want:     false,
			wantErr:  true,
		},
		{
			name:     "invalid derived key format",
			dk:       "invalid-format",
			password: password,
			want:     false,
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := VerifyDerivedKey(tt.dk, tt.password)
			if (err != nil) != tt.wantErr {
				t.Errorf("VerifyDerivedKey() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if got != tt.want {
				t.Errorf("VerifyDerivedKey() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParseDerivedKey(t *testing.T) {
	password := "test-password-123!"

	dk, err := CreateDerivedKey(password)
	if err != nil {
		t.Fatalf("CreateDerivedKey() error = %v", err)
	}

	tests := []struct {
		name    string
		dk      string
		wantErr bool
	}{
		{
			name:    "valid derived key",
			dk:      dk,
			wantErr: false,
		},
		{
			name:    "empty derived key",
			dk:      "",
			wantErr: true,
		},
		{
			name:    "invalid format",
			dk:      "invalid-format",
			wantErr: true,
		},
		{
			name:    "invalid algorithm",
			dk:      "$argon2d$v=19$m=65536,t=1,p=2$salt$key",
			wantErr: true,
		},
		{
			name:    "invalid version",
			dk:      "$argon2id$v=18$m=65536,t=1,p=2$salt$key",
			wantErr: true,
		},
		{
			name:    "invalid memory format",
			dk:      "$argon2id$v=19$m=invalid,t=1,p=2$salt$key",
			wantErr: true,
		},
		{
			name:    "invalid time format",
			dk:      "$argon2id$v=19$m=65536,t=invalid,p=2$salt$key",
			wantErr: true,
		},
		{
			name:    "invalid threads format",
			dk:      "$argon2id$v=19$m=65536,t=1,p=invalid$salt$key",
			wantErr: true,
		},
		{
			name:    "invalid salt format",
			dk:      "$argon2id$v=19$m=65536,t=1,p=2$invalid-salt$key",
			wantErr: true,
		},
		{
			name:    "invalid key format",
			dk:      "$argon2id$v=19$m=65536,t=1,p=2$salt$invalid-key",
			wantErr: true,
		},
		{
			name:    "missing parts",
			dk:      "$argon2id$v=19$m=65536,t=1,p=2$salt",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, _, _, _, _, err := ParseDerivedKey(tt.dk)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseDerivedKey() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestIsDerivedKey(t *testing.T) {
	password := "test-password-123!"

	dk, err := CreateDerivedKey(password)
	if err != nil {
		t.Fatalf("CreateDerivedKey() error = %v", err)
	}

	tests := []struct {
		name string
		dk   string
		want bool
	}{
		{
			name: "valid derived key",
			dk:   dk,
			want: true,
		},
		{
			name: "empty string",
			dk:   "",
			want: false,
		},
		{
			name: "invalid format",
			dk:   "invalid-format",
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsDerivedKey(tt.dk); got != tt.want {
				t.Errorf("IsDerivedKey() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetDerivedKeyConfig(t *testing.T) {
	password := "test-password-123!"
	config := Argon2Config{
		Time:    2,
		Memory:  128 * 1024,
		Threads: 4,
		KeyLen:  32,
		SaltLen: 16,
	}

	dk, err := CreateDerivedKeyWithConfig(password, config)
	if err != nil {
		t.Fatalf("CreateDerivedKeyWithConfig() error = %v", err)
	}

	tests := []struct {
		name    string
		dk      string
		want    Argon2Config
		wantErr bool
	}{
		{
			name:    "valid derived key",
			dk:      dk,
			want:    config,
			wantErr: false,
		},
		{
			name:    "empty derived key",
			dk:      "",
			want:    Argon2Config{},
			wantErr: true,
		},
		{
			name:    "invalid format",
			dk:      "invalid-format",
			want:    Argon2Config{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetDerivedKeyConfig(tt.dk)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetDerivedKeyConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if got.Time != tt.want.Time {
					t.Errorf("GetDerivedKeyConfig().Time = %v, want %v", got.Time, tt.want.Time)
				}

				if got.Memory != tt.want.Memory {
					t.Errorf("GetDerivedKeyConfig().Memory = %v, want %v", got.Memory, tt.want.Memory)
				}

				if got.Threads != tt.want.Threads {
					t.Errorf("GetDerivedKeyConfig().Threads = %v, want %v", got.Threads, tt.want.Threads)
				}

				if got.KeyLen != tt.want.KeyLen {
					t.Errorf("GetDerivedKeyConfig().KeyLen = %v, want %v", got.KeyLen, tt.want.KeyLen)
				}

				if got.SaltLen != tt.want.SaltLen {
					t.Errorf("GetDerivedKeyConfig().SaltLen = %v, want %v", got.SaltLen, tt.want.SaltLen)
				}
			}
		})
	}
}
