package pathz

import (
	"path"
	"path/filepath"
	"runtime"
	"testing"
)

func TestRootDir(t *testing.T) {
	t.Parallel()
	_, b, _, _ := runtime.Caller(0)
	d := path.Join(path.Dir(b))
	root := filepath.Dir(filepath.Dir(d))
	tests := []struct {
		name string
		want string
	}{
		{
			name: "root",
			want: root,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := RootDir(); got != tt.want {
				t.Errorf("RootDir() = %v, want %v", got, tt.want)
			}
		})
	}
}
