package handlers

// func TestPostMetrics(t *testing.T) {
// 	type args struct {
// 		typeMetric string
// 		name       string
// 		value      string
// 	}
// 	type want struct {
// 		code        int
// 		contentType string
// 	}
// 	tests := []struct {
// 		name string
// 		args args
// 		want want
// 	}{
// 		{
// 			name: "positive test #1",
// 			args: args{
// 				typeMetric: "gauge",
// 				name:       "test",
// 				value:      "123.5324523",
// 			},
// 			want: want{
// 				code:        200,
// 				contentType: "text/plain",
// 			},
// 		},
// 		{
// 			name: "negative test #2",
// 			args: args{
// 				typeMetric: "asdas",
// 				name:       "test",
// 				value:      "123.5324523",
// 			},
// 			want: want{
// 				code:        400,
// 				contentType: "text/plain",
// 			},
// 		},
// 		{
// 			name: "negative test #3",
// 			args: args{
// 				typeMetric: "gauge",
// 				name:       "",
// 				value:      "123.5324523",
// 			},
// 			want: want{
// 				code:        404,
// 				contentType: "text/plain",
// 			},
// 		},
// 		{
// 			name: "negative test #4",
// 			args: args{
// 				typeMetric: "gauge",
// 				name:       "test",
// 				value:      "asdasda",
// 			},
// 			want: want{
// 				code:        400,
// 				contentType: "text/plain",
// 			},
// 		},
// 		{
// 			name: "negative test #5",
// 			args: args{
// 				typeMetric: "counter",
// 				name:       "test",
// 				value:      "1234123.23123",
// 			},
// 			want: want{
// 				code:        400,
// 				contentType: "text/plain",
// 			},
// 		},
// 		{
// 			name: "positive test #6",
// 			args: args{
// 				typeMetric: "counter",
// 				name:       "test",
// 				value:      "1234123",
// 			},
// 			want: want{
// 				code:        200,
// 				contentType: "text/plain",
// 			},
// 		},
// 	}
// 	db := storage.NewMemStorage()
// 	for _, test := range tests {
// 		t.Run(test.name, func(t *testing.T) {
// 			paramURL := "/update/" + test.args.typeMetric + "/" + test.args.name + "/" + test.args.value
// 			request := httptest.NewRequest(http.MethodPost, paramURL, nil)
// 			// создаём новый Recorder
// 			w := httptest.NewRecorder()
// 			httpFunc := PostMetrics(db)
// 			httpFunc(w, r)
// 			res := w.Result()
// 			// проверяем код ответа
// 			assert.Equal(t, test.want.code, res.StatusCode)
// 			// получаем и проверяем тело запроса
// 			defer res.Body.Close()
// 			_, err := io.ReadAll(res.Body)
// 			require.NoError(t, err)
// 			assert.Equal(t, test.want.contentType, res.Header.Get("Content-Type"))
// 		})
// 	}
// }

// func TestGetAll(t *testing.T) {
// 	type args struct {
// 		db *storage.MemStorage
// 	}
// 	tests := []struct {
// 		name string
// 		args args
// 		want http.HandlerFunc
// 	}{
// 		// TODO: Add test cases.
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			if got := GetAll(tt.args.db); !reflect.DeepEqual(got, tt.want) {
// 				t.Errorf("GetAll() = %v, want %v", got, tt.want)
// 			}
// 		})
// 	}
// }

// func TestGetValue(t *testing.T) {
// 	type args struct {
// 		db *storage.MemStorage
// 	}
// 	tests := []struct {
// 		name string
// 		args args
// 		want http.HandlerFunc
// 	}{
// 		// TODO: Add test cases.
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			if got := GetValue(tt.args.db); !reflect.DeepEqual(got, tt.want) {
// 				t.Errorf("GetValue() = %v, want %v", got, tt.want)
// 			}
// 		})
// 	}
// }
