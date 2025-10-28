package storage

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMemStorage_Get(t *testing.T) {
	type fields struct {
		links map[string]string
	}
	type args struct {
		key string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
		want1  bool
	}{
		{
			name: "Get link successfully",
			fields: fields{
				links: map[string]string{
					"link1": "value1",
				},
			},
			args: args{
				key: "link1",
			},
			want:  "value1",
			want1: true,
		},
		{
			name: "Get link not found",
			fields: fields{
				links: map[string]string{
					"link2": "value1",
				},
			},
			args: args{
				key: "link1",
			},
			want:  "",
			want1: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &MemStorage{
				links: tt.fields.links,
			}
			got, got1 := m.Get(tt.args.key)
			if got != tt.want {
				t.Errorf("Get() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("Get() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestMemStorage_Set(t *testing.T) {
	type fields struct {
		links map[string]string
	}
	type args struct {
		key   string
		value string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "Set link successfully",
			fields: fields{
				links: map[string]string{
					"link1": "value1",
				},
			},
			args: args{
				key:   "link1",
				value: "value2",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &MemStorage{
				links: tt.fields.links,
			}
			m.Set(tt.args.key, tt.args.value)
			assert.Equal(t, tt.args.value, m.links[tt.args.key])
		})
	}
}

func TestNewMemStorage(t *testing.T) {
	tests := []struct {
		name string
		want *MemStorage
	}{
		{
			name: "success",
			want: &MemStorage{
				links: make(map[string]string),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewMemStorage(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewMemStorage() = %v, want %v", got, tt.want)
			}
		})
	}
}
