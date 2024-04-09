package handlers

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/WaffeSoul/metrics-collector/internal/storage"
)

func TestPostMetrics(t *testing.T) {
	type args struct {
		typeMetric string
		name       string
		value      string
	}
	type want struct {
		code        int
		contentType string
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "positive test #1",
			args: args{
				typeMetric: "gauge",
				name:       "test",
				value:      "123.5324523",
			},
			want: want{
				code:        200,
				contentType: "text/plain",
			},
		},
		{
			name: "negative test #2",
			args: args{
				typeMetric: "asdas",
				name:       "test",
				value:      "123.5324523",
			},
			want: want{
				code:        400,
				contentType: "text/plain",
			},
		},
		{
			name: "negative test #3",
			args: args{
				typeMetric: "gauge",
				name:       "",
				value:      "123.5324523",
			},
			want: want{
				code:        404,
				contentType: "text/plain",
			},
		},
		{
			name: "negative test #4",
			args: args{
				typeMetric: "gauge",
				name:       "test",
				value:      "asdasda",
			},
			want: want{
				code:        400,
				contentType: "text/plain",
			},
		},
		{
			name: "negative test #5",
			args: args{
				typeMetric: "counter",
				name:       "test",
				value:      "1234123.23123",
			},
			want: want{
				code:        400,
				contentType: "text/plain",
			},
		},
		{
			name: "positive test #6",
			args: args{
				typeMetric: "counter",
				name:       "test",
				value:      "1234123",
			},
			want: want{
				code:        200,
				contentType: "text/plain",
			},
		},
	}
	storage.StorageGause = storage.Init()
	storage.StorageConter = storage.Init()
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			paramUrl := "/update/" + test.args.typeMetric + "/" + test.args.name + "/" + test.args.value
			request := httptest.NewRequest(http.MethodPost, paramUrl, nil)
			// создаём новый Recorder
			w := httptest.NewRecorder()
			PostMetrics(w, request)
			res := w.Result()
			// проверяем код ответа
			assert.Equal(t, test.want.code, res.StatusCode)
			// получаем и проверяем тело запроса
			defer res.Body.Close()
			_, err := io.ReadAll(res.Body)
			require.NoError(t, err)
			assert.Equal(t, test.want.contentType, res.Header.Get("Content-Type"))
		})
	}
}
