package storage

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMemStorage_SetGauge(t *testing.T) {

	ms := MemStorage{
		metrics: metrics{
			gauge:   make(map[string]float64),
			counter: make(map[string]int64),
		},
	}

	type args struct {
		name string
		val  float64
	}

	tests := []struct {
		name string
		args args
	}{
		{
			name: "add pos val",
			args: args{name: "testgg1", val: 1.234},
		},
		{
			name: "add neg val",
			args: args{name: "testgg2", val: -1.234},
		},
		{
			name: "add zero val",
			args: args{name: "testgg3", val: 0.0},
		},
		{
			name: "add small val",
			args: args{name: "testgg4", val: 0.00000000000000000013},
		},
		{
			name: "add big val",
			args: args{name: "testgg4", val: 32453453243534534.99999999999999999999999993},
		},
		{
			name: "new val",
			args: args{name: "testgg2", val: 4123.3434},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ms.SetGauge(tt.args.name, tt.args.val)
			assert.Equal(t, tt.args.val, ms.metrics.gauge[tt.args.name])
		})
	}
}

func TestMemStorage_SetCounter(t *testing.T) {

	ms := MemStorage{
		metrics: metrics{
			gauge:   make(map[string]float64),
			counter: make(map[string]int64),
		},
	}

	type args struct {
		name string
		val  int64
	}

	tests := []struct {
		name string
		args args
	}{
		{
			name: "add pos val",
			args: args{name: "testc", val: 13335},
		},
		{
			name: "add neg val",
			args: args{name: "testc", val: -145},
		},
		{
			name: "add zero val",
			args: args{name: "testc", val: 0},
		},
	}

	var cnt int64

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ms.SetCounter(tt.args.name, tt.args.val)
			cnt += tt.args.val
			assert.Equal(t, cnt, ms.metrics.counter[tt.args.name])
		})
	}
}

func TestMemStorage_GetGauge(t *testing.T) {
	type args struct {
		name string
	}
	tests := []struct {
		name     string
		ms       MemStorage
		args     args
		want     float64
		wantBool bool
	}{
		{
			name: "test get ok",
			ms: MemStorage{
				metrics: metrics{
					gauge:   map[string]float64{"tstgg": 1.003},
					counter: make(map[string]int64),
				},
			},
			args:     args{name: "tstgg"},
			want:     1.003,
			wantBool: true,
		},
		{
			name: "test get NotOK",
			ms: MemStorage{
				metrics: metrics{
					gauge:   map[string]float64{"tstggErr": 1.003},
					counter: make(map[string]int64),
				},
			},
			args:     args{name: "tstgg"},
			want:     1.003,
			wantBool: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok := tt.ms.GetGauge(tt.args.name)

			assert.Equal(t, tt.wantBool, ok)

			if ok {
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func TestMemStorage_GetCounter(t *testing.T) {
	type args struct {
		name string
	}
	tests := []struct {
		name     string
		ms       MemStorage
		args     args
		want     int64
		wantBool bool
	}{
		{
			name: "test get ok",
			ms: MemStorage{
				metrics: metrics{
					gauge:   make(map[string]float64),
					counter: map[string]int64{"tstc": 103},
				},
			},
			args:     args{name: "tstc"},
			want:     103,
			wantBool: true,
		},
		{
			name: "test get NotOK",
			ms: MemStorage{
				metrics: metrics{
					gauge:   make(map[string]float64),
					counter: map[string]int64{"tstErr": 103},
				},
			},
			args:     args{name: "tstgg"},
			want:     103,
			wantBool: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok := tt.ms.GetCounter(tt.args.name)

			assert.Equal(t, tt.wantBool, ok)

			if ok {
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func TestMemStorage_GetGaugeList(t *testing.T) {
	tests := []struct {
		name string
		ms   MemStorage
		want []string
	}{
		{
			name: "test get empty",
			ms: MemStorage{
				metrics: metrics{
					gauge:   map[string]float64{},
					counter: make(map[string]int64),
				},
			},
			want: []string{},
		},
		{
			name: "test get list",
			ms: MemStorage{
				metrics: metrics{
					gauge:   map[string]float64{"tstgg1": 1.003, "tstgg2": 1.004, "tstgg3": 5.0},
					counter: make(map[string]int64),
				},
			},
			want: []string{
				fmt.Sprintf("Name: %s,\tValue: %v", "tstgg1", 1.003),
				fmt.Sprintf("Name: %s,\tValue: %v", "tstgg2", 1.004),
				fmt.Sprintf("Name: %s,\tValue: %v", "tstgg3", 5.0),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.ms.GetGaugeList()
			assert.Equal(t, len(tt.want), len(got))
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestMemStorage_GetCounterList(t *testing.T) {
	tests := []struct {
		name string
		ms   MemStorage
		want []string
	}{
		{
			name: "test get empty",
			ms: MemStorage{
				metrics: metrics{
					gauge:   make(map[string]float64),
					counter: map[string]int64{},
				},
			},
			want: []string{},
		},
		{
			name: "test get list",
			ms: MemStorage{
				metrics: metrics{
					gauge:   make(map[string]float64),
					counter: map[string]int64{"tstgg1": 1, "tstgg2": 4, "tstgg3": 5},
				},
			},
			want: []string{
				fmt.Sprintf("Name: %s,\tValue: %v", "tstgg1", 1),
				fmt.Sprintf("Name: %s,\tValue: %v", "tstgg2", 4),
				fmt.Sprintf("Name: %s,\tValue: %v", "tstgg3", 5),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.ms.GetCounterList()
			assert.Equal(t, len(tt.want), len(got))
			assert.Equal(t, tt.want, got)
		})
	}
}
