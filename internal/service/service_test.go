package service

import (
	"avito-tech-backend/domain"
	"avito-tech-backend/internal/service/mocks"
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetHistory(t *testing.T) {
	type mockBehavior func(r *mocks.Repository, timeBegin, timeEnd int64, userId int)
	type args struct {
		TimeBegin int64
		TimeEnd   int64
		UserId    int
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
			args: args{TimeBegin: 1693063272, TimeEnd: 1693064272, UserId: 1},
			mockBehavior: func(r *mocks.Repository, timeBegin, timeEnd int64, userId int) {
				r.EXPECT().GetHistory(context.Background(), timeBegin, timeEnd, userId).Return("file.csv", nil)
			},
			want:    "file.csv",
			wantErr: nil,
		},
		{
			name:         "Invalid argument",
			args:         args{TimeBegin: 1693064272, TimeEnd: 1693063272, UserId: 1},
			mockBehavior: func(r *mocks.Repository, timeBegin, timeEnd int64, userId int) {},
			want:         "",
			wantErr:      domain.ErrInvalidArgument,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := mocks.NewRepository(t)
			tt.mockBehavior(mockRepo, tt.args.TimeBegin, tt.args.TimeEnd, tt.args.UserId)
			service := New(mockRepo)
			response, responseErr := service.GetHistory(tt.args.TimeBegin, tt.args.TimeEnd, tt.args.UserId)
			assert.Equal(t, tt.want, response)
			assert.Equal(t, tt.wantErr, responseErr)
		})
	}
}
