package cmds4listusbot

import "testing"

func Test_cleanListItemTitle(t *testing.T) {
	type args struct {
		title string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "empty",
			args: args{title: ""},
			want: "",
		},
		{
			name: "spaces",
			args: args{title: "  "},
			want: "",
		},
		{
			name: "same",
			args: args{title: "same"},
			want: "same",
		},
		{
			name: "same with spaces",
			args: args{title: " same with spaces "},
			want: "same with spaces",
		},
		{
			name: "same with tabs",
			args: args{title: "\tsame with spaces\t"},
			want: "same with spaces",
		},
		{
			name: "same with dashes",
			args: args{title: "-same with spaces"},
			want: "same with spaces",
		},
		{
			name: "same with tabbed and spaced dashes",
			args: args{title: "\t- same with spaces"},
			want: "same with spaces",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := cleanListItemTitle(tt.args.title); got != tt.want {
				t.Errorf("cleanListItemTitle() = %v, want %v", got, tt.want)
			}
		})
	}
}
