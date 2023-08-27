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
	tests := []struct {
		name        string
		inTimeBegin int64
		inTimeEnd   int64
		inUserId    int
		mockBehavior
		want    string
		wantErr error
	}{
		{
			name:        "ok",
			inTimeBegin: 1693063272,
			inTimeEnd:   1693064272,
			inUserId:    1,
			mockBehavior: func(r *mocks.Repository, timeBegin, timeEnd int64, userId int) {
				r.EXPECT().GetHistory(context.Background(), timeBegin, timeEnd, userId).Return("file.csv", nil)
			},
			want:    "file.csv",
			wantErr: nil,
		},
		{
			name:         "Invalid argument",
			inTimeBegin:  1693064272,
			inTimeEnd:    1693063272,
			inUserId:     1,
			mockBehavior: func(r *mocks.Repository, timeBegin, timeEnd int64, userId int) {},
			want:         "",
			wantErr:      domain.ErrInvalidArgument,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := mocks.NewRepository(t)
			service := New(r)
			response, err := service.GetHistory(tt.inTimeBegin, tt.inTimeEnd, tt.inUserId)
			assert.Equal(t, tt.wantErr, err)
			assert.Equal(t, tt.want, response)
		})
	}
}
