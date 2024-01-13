package dal4teamus

import (
	"github.com/dal-go/dalgo/dal"
	"reflect"
	"testing"
)

func TestNewTeamKey(t *testing.T) {
	type args struct {
		id string
	}
	t.Run("should_panic_on_empty_id", func(t *testing.T) {
		defer func() {
			if recover() == nil {
				t.Fatal("no expected panic")
			}
		}()
		NewTeamKey("")
	})
	t.Run("should_panic_on_long_id", func(t *testing.T) {
		defer func() {
			if recover() == nil {
				t.Fatal("no expected panic")
			}
		}()
		NewTeamKey("0123456789012345678901234567891")
	})
	tests := []struct {
		name string
		args args
		want *dal.Key
	}{
		{
			name: "should_pass",
			args: args{id: "TestTeam"},
			want: dal.NewKeyWithID(TeamsCollection, "TestTeam"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewTeamKey(tt.args.id); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewTeamKey() = %v, want %v", got, tt.want)
			}
		})
	}
}
