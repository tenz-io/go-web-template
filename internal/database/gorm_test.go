package database

import "testing"

func Test_hashPasswordWithSalt(t *testing.T) {
	type args struct {
		password string
		salt     string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "valid password and salt",
			args: args{
				password: intiAdminPassword,
				salt:     initAdminSalt,
			},
			want: "j4rSB8DxWtMzguqCMysQZchdJV8983/u3nnbqWqFJDE=", // Expected hash value
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := hashPasswordWithSalt(tt.args.password, tt.args.salt)
			if (err != nil) != tt.wantErr {
				t.Errorf("hashPasswordWithSalt() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("hashPasswordWithSalt() got = %v, want %v", got, tt.want)
			}
		})
	}
}
