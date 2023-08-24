package service

import (
	"context"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"ozon/domain"
	mock_service "ozon/internal/service/mocks"
	"testing"
)

func TestGet(t *testing.T) {
	type mockBehavior func(r *mock_service.MockRepository, shortUrl string)
	tests := []struct {
		name string
		in   string
		mockBehavior
		want    string
		wantErr error
	}{
		{
			name: "ok",
			in:   "9Txmx1CPWV",
			mockBehavior: func(r *mock_service.MockRepository, shortUrl string) {
				r.EXPECT().Get(context.Background(), shortUrl).Return("google.com", nil)
			},
			want:    "google.com",
			wantErr: nil,
		},
		{
			name:         "Bad request",
			in:           "9Txmx1CPWVT",
			mockBehavior: func(r *mock_service.MockRepository, shortUrl string) {},
			want:         "",
			wantErr:      domain.ErrInvalidArgument,
		},
		{
			name:         "Bad request",
			in:           "",
			mockBehavior: func(r *mock_service.MockRepository, shortUrl string) {},
			want:         "",
			wantErr:      domain.ErrInvalidArgument,
		},
		{
			name: "Not found",
			in:   "9Txmj1CPWV",
			mockBehavior: func(r *mock_service.MockRepository, shortUrl string) {
				r.EXPECT().Get(context.Background(), shortUrl).Return("", domain.ErrNotFound)
			},
			want:    "",
			wantErr: domain.ErrNotFound,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			s := mock_service.NewMockRepository(ctrl)
			tt.mockBehavior(s, tt.in)
			service := New(s)
			response, err := service.Get(context.Background(), tt.in)
			assert.Equal(t, tt.wantErr, err)
			assert.Equal(t, tt.want, response)
		})
	}
}
func TestCreate(t *testing.T) {
	type mockBehavior func(r *mock_service.MockRepository, shortUrl, originalUrl string)
	type args struct {
		ShortUrl    string
		OriginalUrl string
	}
	tests := []struct {
		name string
		args
		mockBehavior
		want    string
		wantErr error
	}{
		{
			name: "ok",
			args: args{OriginalUrl: "google.com", ShortUrl: "9Txmx1CPWV"},
			mockBehavior: func(r *mock_service.MockRepository, shortUrl, originalUrl string) {
				r.EXPECT().Create(context.Background(), shortUrl, originalUrl).Return(nil)
			},
			want:    "9Txmx1CPWV",
			wantErr: nil,
		},
		{
			name:         "invalid format",
			args:         args{OriginalUrl: "googlecom", ShortUrl: ""},
			mockBehavior: func(r *mock_service.MockRepository, shortUrl, originalUrl string) {},
			want:         "",
			wantErr:      domain.ErrInvalidArgument,
		},
		{
			name:         "invalid format",
			args:         args{OriginalUrl: "", ShortUrl: ""},
			mockBehavior: func(r *mock_service.MockRepository, shortUrl, originalUrl string) {},
			want:         "",
			wantErr:      domain.ErrInvalidArgument,
		},
		{
			name:         "Bad request",
			args:         args{OriginalUrl: "", ShortUrl: "9Txmx1CPWV"},
			mockBehavior: func(r *mock_service.MockRepository, shortUrl, originalUrl string) {},
			want:         "",
			wantErr:      domain.ErrInvalidArgument,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			s := mock_service.NewMockRepository(ctrl)
			tt.mockBehavior(s, tt.args.ShortUrl, tt.args.OriginalUrl)
			service := New(s)
			response, err := service.Create(context.Background(), tt.args.OriginalUrl)
			assert.Equal(t, tt.wantErr, err)
			assert.Equal(t, tt.want, response)
		})
	}
}
func TestHash(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want string
	}{
		{
			name: "ok",
			in:   "google.com",
			want: "9Txmx1CPWV",
		},
		{
			name: "ok",
			in:   "sad.asd",
			want: "H4iUYyizdo",
		},
		{
			name: "ok",
			in:   "asd.asd",
			want: "dvFbab4Wnj",
		},
		{
			name: "ok",
			in:   "qwee.qweqwe",
			want: "vfS0QpXBCN",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			s := mock_service.NewMockRepository(ctrl)
			service := New(s)
			response := service.Hash(tt.in)
			assert.Equal(t, tt.want, response)
		})
	}
}
