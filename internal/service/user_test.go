package service

import (
	"context"
	"testing"

	"go-web-template/internal/config"
	"go-web-template/internal/repository"
)

func Test_user_VerifyAdmin(t *testing.T) {
	type fields struct {
		cfg     *config.Config
		useRepo repository.User
	}
	type args struct {
		ctx  context.Context
		name string
		pass string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantOk  bool
		wantErr bool
	}{
		{
			name: "empty name or pass",
			fields: fields{
				cfg: &config.Config{},
			},
			args: args{
				ctx:  context.Background(),
				name: "",
				pass: "",
			},
			wantOk:  false,
			wantErr: true,
		},
		{
			name: "admin verified",
			fields: fields{
				cfg: &config.Config{
					App: config.AppConfig{
						AdminUser: "admin",
						AdminPass: "admin123",
					},
				},
			},
			args: args{
				ctx:  context.Background(),
				name: "admin",
				pass: "admin123",
			},
			wantOk:  true,
			wantErr: false,
		},
		{
			name: "admin verify failed",
			fields: fields{
				cfg: &config.Config{
					App: config.AppConfig{
						AdminUser: "admin",
						AdminPass: "admin123",
					},
				},
			},
			args: args{
				ctx:  context.Background(),
				name: "admin",
				pass: "admin",
			},
			wantOk:  false,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &user{
				cfg:     tt.fields.cfg,
				useRepo: tt.fields.useRepo,
			}
			gotOk, err := u.VerifyAdmin(tt.args.ctx, tt.args.name, tt.args.pass)
			if (err != nil) != tt.wantErr {
				t.Errorf("VerifyAdmin() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotOk != tt.wantOk {
				t.Errorf("VerifyAdmin() gotOk = %v, want %v", gotOk, tt.wantOk)
			}
		})
	}
}
