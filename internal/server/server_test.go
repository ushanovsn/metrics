package server

import (
	"github.com/stretchr/testify/assert"
	"github.com/ushanovsn/metrics/internal/storage"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestServerMux(t *testing.T) {
	type want struct {
		code    int
		mType   string
		mName   string
		mValueF float64
		mValueI int64
		isOkVal bool
	}

	tests := []struct {
		name string
		args string
		want want
	}{
		{
			name: "test OK gauge (200)",
			args: "/update/gauge/testGauge1/12.0",
			want: want{
				code:    200,
				mType:   "gauge",
				mName:   "testGauge1",
				mValueF: 12.0,
				isOkVal: true,
			},
		},
		{
			name: "test OK counter (200)",
			args: "/update/counter/testCounter1/13",
			want: want{
				code:    200,
				mType:   "counter",
				mName:   "testCounter1",
				mValueI: 13,
				isOkVal: true,
			},
		},
		{
			name: "test OK gauge (200) add",
			args: "/update/gauge/testGauge1/13.0",
			want: want{
				code:    200,
				mType:   "gauge",
				mName:   "testGauge1",
				mValueF: 13.0,
				isOkVal: true,
			},
		},
		{
			name: "test OK counter (200) add",
			args: "/update/counter/testCounter1/12",
			want: want{
				code:    200,
				mType:   "counter",
				mName:   "testCounter1",
				mValueI: 25,
				isOkVal: true,
			},
		},
		{
			name: "test BadReq gauge (400)",
			args: "/update/gauge1/testGauge2/12.0",
			want: want{
				code:    400,
				mType:   "gauge",
				mName:   "testGauge2",
				mValueF: 0.0,
				isOkVal: false,
			},
		},
		{
			name: "test BadReq counter (400)",
			args: "/update/1counter/testCounter2/13",
			want: want{
				code:    400,
				mType:   "counter",
				mName:   "testCounter2",
				mValueI: 0,
				isOkVal: false,
			},
		},
		{
			name: "test NoFound gauge (404)",
			args: "/update/gauge/12.0",
			want: want{
				code:    404,
				mType:   "gauge",
				mName:   "",
				mValueF: 0.0,
				isOkVal: false,
			},
		},
		{
			name: "test NoFound counter (404)",
			args: "/update/counter/13",
			want: want{
				code:    404,
				mType:   "counter",
				mName:   "",
				mValueI: 0,
				isOkVal: false,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			//request for handler
			request := httptest.NewRequest(http.MethodPost, tt.args, nil)
			w := httptest.NewRecorder()

			b := ServerMux()
			b.ServeHTTP(w, request)

			res := w.Result()

			assert.Equal(t, tt.want.code, res.StatusCode)

			switch tt.want.mType {
			case "gauge":
				val, okRes := storage.Metr.GetGauge(tt.want.mName)
				assert.Equal(t, tt.want.mValueF, val)
				assert.Equal(t, tt.want.isOkVal, okRes)
			case "counter":
				val, okRes := storage.Metr.GetCounter(tt.want.mName)
				assert.Equal(t, tt.want.mValueI, val)
				assert.Equal(t, tt.want.isOkVal, okRes)
			}

			res.Body.Close()
		})
	}
}
