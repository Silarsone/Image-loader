package service

import (
	"context"
	"fmt"
	"github.com/Silarsone/image-loader/internal/mock"
	"github.com/Silarsone/image-loader/internal/model"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"io"
	"strings"
	"testing"
)

func TestController_GetImageObjects(t *testing.T) {
	mockCtrl := gomock.NewController(t)

	mockRepo := mock.NewMockrepository(mockCtrl)
	mockImageRepo := mock.NewMockimageRepository(mockCtrl)
	mockMinio := mock.NewMockfileStorage(mockCtrl)

	tt := []struct {
		name           string
		tgID           int64
		userID         int
		images         []model.Image
		reader         []io.Reader
		err            error
		expectedReader []io.Reader
		expectedErr    error
		mockF          func(tgID int64, userID int, images []model.Image, objects []io.Reader, err error)
	}{
		{
			name:           "success",
			tgID:           1,
			userID:         2,
			reader:         []io.Reader{strings.NewReader("lalalalala")},
			images:         []model.Image{{UserID: 2}},
			err:            nil,
			expectedReader: []io.Reader{strings.NewReader("lalalalala")},
			expectedErr:    nil,
			mockF: func(tgID int64, userID int, images []model.Image, objects []io.Reader, err error) {
				mockRepo.EXPECT().
					CheckTgAuth(gomock.AssignableToTypeOf(context.Background()), tgID).
					Return(userID, nil)

				mockImageRepo.EXPECT().
					GetImages(gomock.AssignableToTypeOf(context.Background()), userID).
					Return(images, nil)

				mockMinio.EXPECT().
					GetObjects(gomock.AssignableToTypeOf(context.Background()), images).
					Return(objects, nil)
			},
		},
		{
			name:           "failed to check tg auth",
			tgID:           2,
			userID:         0,
			reader:         nil,
			images:         nil,
			err:            fmt.Errorf("failed tg auth"),
			expectedReader: nil,
			expectedErr:    fmt.Errorf("failed tg auth"),
			mockF: func(tgID int64, userID int, images []model.Image, objects []io.Reader, err error) {
				mockRepo.EXPECT().
					CheckTgAuth(gomock.AssignableToTypeOf(context.Background()), tgID).
					Return(userID, err)
			},
		},
		{
			name:           "failed to get images",
			tgID:           3,
			userID:         0,
			reader:         nil,
			images:         nil,
			err:            fmt.Errorf("failed to get images"),
			expectedReader: nil,
			expectedErr:    fmt.Errorf("failed to get images"),
			mockF: func(tgID int64, userID int, images []model.Image, objects []io.Reader, err error) {
				mockRepo.EXPECT().
					CheckTgAuth(gomock.AssignableToTypeOf(context.Background()), tgID).
					Return(userID, nil)

				mockImageRepo.EXPECT().
					GetImages(gomock.AssignableToTypeOf(context.Background()), userID).
					Return(images, err)
			},
		},
		{
			name:           "failed to get objects",
			tgID:           1,
			userID:         0,
			reader:         nil,
			images:         []model.Image{{UserID: 2}},
			err:            fmt.Errorf("failed to get objects"),
			expectedReader: nil,
			expectedErr:    fmt.Errorf("failed to get objects"),
			mockF: func(tgID int64, userID int, images []model.Image, objects []io.Reader, err error) {
				mockRepo.EXPECT().
					CheckTgAuth(gomock.AssignableToTypeOf(context.Background()), tgID).
					Return(userID, nil)

				mockImageRepo.EXPECT().
					GetImages(gomock.AssignableToTypeOf(context.Background()), userID).
					Return(images, nil)

				mockMinio.EXPECT().
					GetObjects(gomock.AssignableToTypeOf(context.Background()), images).
					Return(nil, err)
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockF(tc.tgID, tc.userID, tc.images, tc.reader, tc.err)

			controller := NewController(mockRepo, mockImageRepo, nil, mockMinio)

			res, err := controller.GetImageObjects(context.Background(), tc.tgID)
			assert.EqualValues(t, tc.expectedErr, err)
			assert.EqualValues(t, tc.expectedReader, res)
		})
	}
}

func BenchmarkForABenchmark(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ForABenchmark()
	}
}
