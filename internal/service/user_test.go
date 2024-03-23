package service

import (
	"context"
	"fmt"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"

	"go-web-template/internal/model"
	"go-web-template/internal/repository"
)

func Test_user_GetByName(t *testing.T) {
	type fields struct {
		useRepo repository.User
	}
	type args struct {
		ctx  context.Context
		name string
	}
	type behavior func(fields, args)
	tests := []struct {
		name     string
		fields   fields
		args     args
		behavior behavior
		want     model.User
		hasErr   assert.ErrorAssertionFunc
	}{
		{
			name: "when user profile is found then return user",
			fields: fields{
				useRepo: repository.NewMockUser(t),
			},
			args: args{
				ctx:  context.Background(),
				name: "Gopher",
			},
			behavior: func(fields fields, args args) {
				var (
					userRepo = fields.useRepo.(*repository.MockUser)
				)

				userRepo.On("UserProfile", args.ctx, args.name).
					Return(fmt.Sprintf("hi %s", args.name), nil).
					Times(1)

			},
			want: model.User{
				Username: "Gopher",
				Profile:  "hi Gopher",
			},
			hasErr: assert.NoError,
		},
		{
			name: "when user profile is not found then return error",
			fields: fields{
				useRepo: repository.NewMockUser(t),
			},
			args: args{
				ctx:  context.Background(),
				name: "Gopher",
			},
			behavior: func(fields fields, args args) {
				var (
					userRepo = fields.useRepo.(*repository.MockUser)
				)

				userRepo.On("UserProfile", args.ctx, args.name).
					Return("", fmt.Errorf("not found")).
					Times(1)
			},
			want:   model.User{},
			hasErr: assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// mock
			tt.behavior(tt.fields, tt.args)

			u := &user{
				useRepo: tt.fields.useRepo,
			}

			got, err := u.GetByName(tt.args.ctx, tt.args.name)

			// assert
			if tt.hasErr(t, err, "user.GetByName() error = %v, wantErr %v") {
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetByName() got = %v, want %v", got, tt.want)
				return
			}
		})
	}
}
