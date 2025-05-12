package internal

import "testing"

func TestHashPassword(t *testing.T) {
	type args struct {
		plainTextPassword string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{"Base case", args{""}, "", true},
		{"Simple password", args{"password"}, "password", false},
		{"Simple password with special characters", args{"p@ssw0rd"}, "p@ssw0rd", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := HashPassword(tt.args.plainTextPassword)
			if (err != nil) != tt.wantErr {
				t.Errorf("HashPassword() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if err != nil && got != tt.want {
				t.Errorf("HashPassword() = %v, want %v", got, tt.want)
				return // this is a test case where we expect an error, so we don't need to check the hash
			}

			if tt.args.plainTextPassword != "" && !VerifyPassword(got, tt.want) {
				t.Errorf("HashPassword() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestVerifyPassword(t *testing.T) {
	type args struct {
		hashedPassword string
		password       string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{"Base case", args{"", ""}, false},
		{"Basic password", args{"password", "password"}, true},
		{"Basic password with special characters", args{"p@ssw0rd", "p@ssw0rd"}, true},
		{"Incorrect password", args{"p@ssw0rd", "wrongpassword"}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Hash the password for the test case
			// We hash the password here to simulate the case where we have a hashed password.
			// However, when we do our base case with an empty input, this is invalid.
			// But we don't want to stop the test case since we're testing the VerifyPassword function.
			hashedPassword, _ := HashPassword(tt.args.hashedPassword)
			// Verify the password
			if got := VerifyPassword(hashedPassword, tt.args.password); got != tt.want {
				t.Errorf("VerifyPassword() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEncodeSHA256(t *testing.T) {
	type args struct {
		data []byte
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"Base case", args{[]byte("")}, ""},
		{"Simple string", args{[]byte("hello world")}, "b94d27b9934d3e08a52e52d7da7dabfac484efe37a5380ee9088f7ace2efcde9"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := EncodeSHA256(tt.args.data); got != tt.want {
				t.Errorf("EncodeSHA256() = %v, want %v", got, tt.want)
			}
		})
	}
}
