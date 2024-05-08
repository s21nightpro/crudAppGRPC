package server

import (
	"context"
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/s21nightpro/crudAppGRPC/cmd/api"
	"github.com/s21nightpro/crudAppGRPC/internal/server/mocks"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

//func TestNewServer(t *testing.T) {
//	type args struct {
//		cache Cache
//		db    DB
//	}
//	tests := []struct {
//		name string
//		args args
//		want *Server
//	}{
//		// TODO: Add test cases.
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			if got := NewServer(tt.args.cache, tt.args.db); !reflect.DeepEqual(got, tt.want) {
//				t.Errorf("NewServer() = %v, want %v", got, tt.want)
//			}
//		})
//	}
//}

func TestServer_CreateUser(t *testing.T) {
	type fields struct {
		UnimplementedUserServiceServer api.UnimplementedUserServiceServer
		cache                          *mocks.MockCache
		db                             *mocks.MockDB
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCache := mocks.NewMockCache(ctrl)
	mockDB := mocks.NewMockDB(ctrl)

	f := fields{
		cache: mockCache,
		db:    mockDB,
	}

	type args struct {
		ctx context.Context
		req *api.CreateUserRequest
	}
	tests := []struct {
		name        string
		fields      fields
		args        args
		want        *api.User
		prepareFunc func(f *fields)
		wantErr     bool
	}{
		{
			name:   "easy case",
			fields: f,
			args: args{context.Background(), &api.CreateUserRequest{
				Name:  "Aboba",
				Email: "loh@gmail.com",
			}},
			want: &api.User{
				Name:  "Aboba",
				Email: "loh@gmail.com",
			},
			prepareFunc: func(f *fields) {
				f.db.EXPECT().Exec(gomock.Any(), gomock.Any()).Return(nil, nil)
			},
			wantErr: false,
		},
		{
			name: "error case",

			fields: f,
			args:   args{context.Background(), &api.CreateUserRequest{}},
			want:   &api.User{},
			prepareFunc: func(f *fields) {
				f.db.EXPECT().Exec(gomock.Any(), gomock.Any()).Return(nil, errors.New("error"))
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.prepareFunc(&tt.fields)

			s := NewServer(tt.fields.cache, tt.fields.db)
			got, err := s.CreateUser(tt.args.ctx, tt.args.req)

			if (err != nil) != tt.wantErr {
				t.Errorf("CreateUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				assert.Equal(t, got.Name, tt.want.Name)
				assert.Equal(t, got.Email, tt.want.Email)
				assert.NotEmpty(t, got.Id)
			}
		})
	}
}

func TestServer_DeleteUser(t *testing.T) {
	type fields struct {
		UnimplementedUserServiceServer api.UnimplementedUserServiceServer
		cache                          *mocks.MockCache
		db                             *mocks.MockDB
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCache := mocks.NewMockCache(ctrl)
	mockDB := mocks.NewMockDB(ctrl)

	f := fields{
		cache: mockCache,
		db:    mockDB,
	}

	type args struct {
		ctx context.Context
		req *api.DeleteUserRequest
	}
	tests := []struct {
		name        string
		fields      fields
		args        args
		wantErr     bool
		prepareFunc func(f *fields)
	}{
		{
			name:   "successful delete",
			fields: f,
			args: args{context.Background(), &api.DeleteUserRequest{
				Id: "someUserId",
			}},
			wantErr: false,
			prepareFunc: func(f *fields) {
				f.db.EXPECT().Exec(gomock.Any(), gomock.Any()).Return(nil, nil)
			},
		},
		{
			name:    "delete error",
			fields:  f,
			args:    args{context.Background(), &api.DeleteUserRequest{}},
			wantErr: true,
			prepareFunc: func(f *fields) {
				f.db.EXPECT().Exec(gomock.Any(), gomock.Any()).Return(nil, errors.New("error"))
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.prepareFunc(&tt.fields)

			s := NewServer(tt.fields.cache, tt.fields.db)
			_, err := s.DeleteUser(tt.args.ctx, tt.args.req) // Добавлено второе значение для ошибки

			if (err != nil) != tt.wantErr {
				t.Errorf("DeleteUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestServer_GetUser(t *testing.T) {
	type fields struct {
		UnimplementedUserServiceServer api.UnimplementedUserServiceServer
		cache                          Cache
		db                             DB
	}
	type args struct {
		ctx context.Context
		req *api.GetUserRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *api.User
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Server{
				UnimplementedUserServiceServer: tt.fields.UnimplementedUserServiceServer,
				cache:                          tt.fields.cache,
				db:                             tt.fields.db,
			}
			got, err := s.GetUser(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetUser() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestServer_UpdateUser(t *testing.T) {
	type fields struct {
		UnimplementedUserServiceServer api.UnimplementedUserServiceServer
		cache                          Cache
		db                             DB
	}
	type args struct {
		ctx context.Context
		req *api.UpdateUserRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *api.User
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Server{
				UnimplementedUserServiceServer: tt.fields.UnimplementedUserServiceServer,
				cache:                          tt.fields.cache,
				db:                             tt.fields.db,
			}
			got, err := s.UpdateUser(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("UpdateUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("UpdateUser() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestServer_userExists(t *testing.T) {
	type fields struct {
		UnimplementedUserServiceServer api.UnimplementedUserServiceServer
		cache                          Cache
		db                             DB
	}
	type args struct {
		id string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    bool
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Server{
				UnimplementedUserServiceServer: tt.fields.UnimplementedUserServiceServer,
				cache:                          tt.fields.cache,
				db:                             tt.fields.db,
			}
			got, err := s.userExists(tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("userExists() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("userExists() got = %v, want %v", got, tt.want)
			}
		})
	}
}
