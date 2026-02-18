package storage

import (
	"context"
	"reflect"
	"testing"

	"github.com/Foga2H/ya-go-url-shortener/internal/repository"
	"github.com/stretchr/testify/assert"
)

func TestMemStorage_Get(t *testing.T) {
	type fields struct {
		links map[string]Link
	}
	type args struct {
		key string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   repository.UserLink
		want1  bool
	}{
		{
			name: "Get link successfully",
			fields: fields{
				links: map[string]Link{
					"link1": {OriginalURL: "value1", UserID: "u1"},
				},
			},
			args: args{
				key: "link1",
			},
			want: repository.UserLink{
				ShortURL:    "link1",
				OriginalURL: "value1",
			},
			want1: true,
		},
		{
			name: "Get link not found",
			fields: fields{
				links: map[string]Link{
					"link2": {OriginalURL: "value1", UserID: "u1"},
				},
			},
			args: args{
				key: "link1",
			},
			want:  repository.UserLink{},
			want1: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &MemStorage{
				links: tt.fields.links,
			}
			got, got1 := m.Get(context.Background(), tt.args.key)
			if !reflect.DeepEqual(got, tt.want) {
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
		links map[string]Link
	}
	type args struct {
		userID string
		key    string
		value  string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "Set link successfully",
			fields: fields{
				links: map[string]Link{
					"link1": {OriginalURL: "value1", UserID: "u1"},
				},
			},
			args: args{
				userID: "u2",
				key:    "link1",
				value:  "value2",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &MemStorage{
				links: tt.fields.links,
			}
			gotKey, err := m.Set(context.Background(), tt.args.userID, tt.args.key, tt.args.value)
			assert.NoError(t, err)
			assert.Equal(t, tt.args.key, gotKey)
			assert.Equal(t, tt.args.value, m.links[tt.args.key].OriginalURL)
			assert.Equal(t, tt.args.userID, m.links[tt.args.key].UserID)
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
				links: make(map[string]Link),
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
