package nils

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestString(t *testing.T) {
	in := "hi"
	dv := []string{"default value", "default value 2"}
	type args struct {
		s            *string
		defaultValue []string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"nil without default", args{nil, nil}, ""},
		{"nil with default", args{nil, dv}, "default value"},
		{"hi without default", args{&in, nil}, "hi"},
		{"hi with default", args{&in, dv}, "hi"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ss := String(tt.args.s, tt.args.defaultValue...)
			require.Equal(t, ss, tt.want)
		})
	}
}
