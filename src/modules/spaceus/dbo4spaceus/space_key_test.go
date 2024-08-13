package dbo4spaceus

import (
	"github.com/dal-go/dalgo/dal"
	"reflect"
	"testing"
)

func TestNewSpaceKey(t *testing.T) {
	type args struct {
		id string
	}
	t.Run("should_panic_on_empty_id", func(t *testing.T) {
		defer func() {
			if recover() == nil {
				t.Fatal("no expected panic")
			}
		}()
		NewSpaceKey("")
	})
	t.Run("should_panic_on_long_id", func(t *testing.T) {
		defer func() {
			if recover() == nil {
				t.Fatal("no expected panic")
			}
		}()
		NewSpaceKey("0123456789012345678901234567891")
	})
	tests := []struct {
		name string
		args args
		want *dal.Key
	}{
		{
			name: "should_pass",
			args: args{id: "TestSpace"},
			want: dal.NewKeyWithID(SpacesCollection, "TestSpace"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewSpaceKey(tt.args.id); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewSpaceKey() = %v, want %v", got, tt.want)
			}
		})
	}
}
