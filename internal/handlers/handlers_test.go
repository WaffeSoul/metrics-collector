package handlers

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/WaffeSoul/metrics-collector/internal/storage"
)

func TestPostMetricsJSON(t *testing.T) {
	type args struct {
		data        []byte
		contentType string
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
				data:        []byte("{\"id\": \"test\",\"type\": \"gauge\", \"value\": 123.5324523}"),
				contentType: "application/json",
			},
			want: want{
				code:        200,
				contentType: "text/plain",
			},
		},
		{
			name: "negative test #2",
			args: args{
				data:        []byte("{\"id\": \"test\",\"type\": \"asdas\", \"value\": 123.5324523}"),
				contentType: "application/json",
			},
			want: want{
				code:        400,
				contentType: "text/plain",
			},
		},
		{
			name: "negative test #3",
			args: args{
				data:        []byte("{\"type\": \"gauge\", \"value\": 123.5324523}"),
				contentType: "application/json",
			},
			want: want{
				code:        404,
				contentType: "text/plain",
			},
		},
		{
			name: "negative test #4",
			args: args{
				data:        []byte("{\"id\": \"test\",\"type\": \"gauge\", \"value\": \"asdasda\"}"),
				contentType: "application/json",
			},
			want: want{
				code:        400,
				contentType: "text/plain",
			},
		},
		{
			name: "negative test #5",
			args: args{
				data:        []byte("{\"id\": \"test\",\"type\": \"counter\", \"delta\": 123.5324523}"),
				contentType: "application/json",
			},
			want: want{
				code:        400,
				contentType: "text/plain",
			},
		},
		{
			name: "positive test #6",
			args: args{
				data:        []byte("{\"id\": \"test\",\"type\": \"counter\", \"delta\": 1234123}"),
				contentType: "application/json",
			},
			want: want{
				code:        200,
				contentType: "text/plain",
			},
		},
	}
	db, _ := storage.New("mem", 10, "test.txt", "")
	r := chi.NewRouter()
	r.Route("/", func(r chi.Router) {
		r.Post("/update/", PostMetricJSON(db))
	})

	ts := httptest.NewServer(r)
	defer ts.Close()
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			paramURL := "/update/"
			w, err := http.Post(ts.URL+paramURL, "application/json", bytes.NewBuffer(test.args.data))
			require.NoError(t, err)
			assert.Equal(t, test.want.code, w.StatusCode)
			defer w.Body.Close()
			_, err = io.ReadAll(w.Body)
			require.NoError(t, err)
			assert.Equal(t, test.want.contentType, w.Header.Get("Content-Type"))
		})
	}
}

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
	db, _ := storage.New("mem", 10, "test.txt", "")
	r := chi.NewRouter()
	r.Route("/", func(r chi.Router) {
		r.Post("/update/{type}/{name}/{value}", PostMetric(db))
	})

	ts := httptest.NewServer(r)
	defer ts.Close()
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			paramURL := "/update/" + test.args.typeMetric + "/" + test.args.name + "/" + test.args.value
			request, err := http.NewRequest(http.MethodPost, ts.URL+paramURL, nil)
			require.NoError(t, err)
			w, err := http.DefaultClient.Do(request)
			require.NoError(t, err)
			assert.Equal(t, test.want.code, w.StatusCode)
			defer w.Body.Close()
			_, err = io.ReadAll(w.Body)
			require.NoError(t, err)
			assert.Equal(t, test.want.contentType, w.Header.Get("Content-Type"))
		})
	}
}

func TestGetValueJSON(t *testing.T) {
	type args struct {
		data        []byte
		contentType string
	}
	type want struct {
		code        int
		contentType string
		body        []byte
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "positive test #1",
			args: args{
				data:        []byte("{\"id\":\"test\",\"type\":\"gauge\"}"),
				contentType: "application/json",
			},
			want: want{
				code:        200,
				contentType: "application/json",
				body:        []byte("{\"id\":\"test\",\"type\":\"gauge\",\"value\":123.5324523}"),
			},
		},
		{
			name: "positive test #2",
			args: args{
				data:        []byte("{\"id\":\"test\",\"type\":\"counter\"}"),
				contentType: "application/json",
			},
			want: want{
				code:        200,
				contentType: "application/json",
				body:        []byte("{\"id\":\"test\",\"type\":\"counter\",\"delta\":123}"),
			},
		},
		{
			name: "negative test #3",
			args: args{
				data:        []byte("{\"id\":\"test2\",\"type\":\"gauge\"}"),
				contentType: "application/json",
			},
			want: want{
				code:        404,
				contentType: "application/json",
				body:        []byte(""),
			},
		},
		{
			name: "negative test #4",
			args: args{
				data:        []byte("{\"id\":\"test2\",\"type\":\"counter\"}"),
				contentType: "application/json",
			},
			want: want{
				code:        404,
				contentType: "application/json",
				body:        []byte(""),
			},
		},
		{
			name: "negative test #5",
			args: args{
				data:        []byte("{\"id\":\"test\",\"type\":\"gaugeasd\"}"),
				contentType: "application/json",
			},
			want: want{
				code:        400,
				contentType: "application/json",
				body:        []byte(""),
			},
		},
		{
			name: "negative test #6",
			args: args{
				data:        []byte("{\"type\":\"counter\"}"),
				contentType: "application/json",
			},
			want: want{
				code:        404,
				contentType: "application/json",
				body:        []byte(""),
			},
		},
		{
			name: "negative test #7",
			args: args{
				data:        []byte("{\"id\":\"test\"}"),
				contentType: "application/json",
			},
			want: want{
				code:        400,
				contentType: "application/json",
				body:        []byte(""),
			},
		},
		{
			name: "negative test #8",
			args: args{
				data:        []byte(""),
				contentType: "application/json",
			},
			want: want{
				code:        400,
				contentType: "application/json",
				body:        []byte(""),
			},
		},
	}
	db, _ := storage.New("mem", 10, "test.txt", "")
	db.DB.Add("gauge", "test", "123.5324523")
	db.DB.Add("counter", "test", "123")
	r := chi.NewRouter()
	r.Post("/value/", GetValueJSON(db))
	ts := httptest.NewServer(r)
	defer ts.Close()
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			paramURL := "/value/"
			w, err := http.Post(ts.URL+paramURL, "application/json", bytes.NewBuffer(test.args.data))
			require.NoError(t, err)
			assert.Equal(t, test.want.code, w.StatusCode)
			defer w.Body.Close()
			body, err := io.ReadAll(w.Body)
			require.NoError(t, err)
			assert.Equal(t, test.want.contentType, w.Header.Get("Content-Type"))
			assert.Equal(t, test.want.body, body)
		})
	}
}

func TestGetValue(t *testing.T) {
	type args struct {
		typeMetric string
		name       string
	}
	type want struct {
		code        int
		contentType string
		body        string
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
			},
			want: want{
				code:        200,
				contentType: "text/plain",
				body:        "123.5324523",
			},
		},
		{
			name: "positive test #2",
			args: args{
				typeMetric: "counter",
				name:       "test",
			},
			want: want{
				code:        200,
				contentType: "text/plain",
				body:        "123",
			},
		},
		{
			name: "negative test #3",
			args: args{
				typeMetric: "gauge",
				name:       "test2",
			},
			want: want{
				code:        404,
				contentType: "text/plain",
				body:        "",
			},
		},
		{
			name: "negative test #4",
			args: args{
				typeMetric: "counter",
				name:       "test2",
			},
			want: want{
				code:        404,
				contentType: "text/plain",
				body:        "",
			},
		},
		{
			name: "negative test #5",
			args: args{
				typeMetric: "gaugeasd",
				name:       "test",
			},
			want: want{
				code:        400,
				contentType: "text/plain",
				body:        "",
			},
		},
		{
			name: "negative test #6",
			args: args{
				typeMetric: "counter",
				name:       "",
			},
			want: want{
				code:        404,
				contentType: "text/plain; charset=utf-8",
				body:        "404 page not found\n",
			},
		},
		{
			name: "negative test #7",
			args: args{
				typeMetric: "",
				name:       "test",
			},
			want: want{
				code:        400,
				contentType: "text/plain",
				body:        "",
			},
		},
		{
			name: "negative test #8",
			args: args{
				typeMetric: "",
				name:       "",
			},
			want: want{
				code:        404,
				contentType: "text/plain; charset=utf-8",
				body:        "404 page not found\n",
			},
		},
	}
	db, _ := storage.New("mem", 10, "test.txt", "")
	db.DB.Add("gauge", "test", "123.5324523")
	db.DB.Add("counter", "test", "123")
	r := chi.NewRouter()
	r.Get("/value/{type}/{name}", GetValue(db))
	ts := httptest.NewServer(r)
	defer ts.Close()
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			paramURL := "/value/" + test.args.typeMetric + "/" + test.args.name
			request, err := http.NewRequest(http.MethodGet, ts.URL+paramURL, nil)
			require.NoError(t, err)
			w, err := http.DefaultClient.Do(request)
			require.NoError(t, err)
			assert.Equal(t, test.want.code, w.StatusCode)
			defer w.Body.Close()
			body, err := io.ReadAll(w.Body)
			require.NoError(t, err)
			assert.Equal(t, test.want.contentType, w.Header.Get("Content-Type"))
			assert.Equal(t, test.want.body, string(body))
		})
	}
}

func TestGetAll(t *testing.T) {
	type want struct {
		code        int
		contentType string
		body        string
	}
	tests := []struct {
		name string
		want want
	}{
		{
			name: "positive test #1",
			want: want{
				code:        200,
				contentType: "text/plain",
				body:        "test: 123\ntest: 123.5324523\n",
			},
		},
	}
	db, _ := storage.New("mem", 10, "test.txt", "")
	db.DB.Add("gauge", "test", "123.5324523")
	db.DB.Add("counter", "test", "123")
	r := chi.NewRouter()
	r.Get("/", GetAll(db))
	ts := httptest.NewServer(r)
	defer ts.Close()
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			request, err := http.NewRequest(http.MethodGet, ts.URL, nil)
			require.NoError(t, err)
			w, err := http.DefaultClient.Do(request)
			require.NoError(t, err)
			assert.Equal(t, test.want.code, w.StatusCode)
			defer w.Body.Close()
			body, err := io.ReadAll(w.Body)
			require.NoError(t, err)
			assert.Equal(t, test.want.contentType, w.Header.Get("Content-Type"))
			assert.Equal(t, test.want.body, string(body))
		})
	}
}
