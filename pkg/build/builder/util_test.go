package builder

import (
	"fmt"
	"path/filepath"
	"strings"
	"testing"
)

func TestMergeEnv(t *testing.T) {
	tests := []struct {
		oldEnv   []string
		newEnv   []string
		expected []string
	}{
		{
			oldEnv:   []string{"one=1", "two=2"},
			newEnv:   []string{"three=3", "four=4"},
			expected: []string{"one=1", "two=2", "three=3", "four=4"},
		},
		{
			oldEnv:   []string{"one=1", "two=2", "four=4"},
			newEnv:   []string{"three=3", "four=4=5=6"},
			expected: []string{"one=1", "two=2", "three=3", "four=4=5=6"},
		},
		{
			oldEnv:   []string{"one=1", "two=2", "three=3"},
			newEnv:   []string{"two=002", "four=4"},
			expected: []string{"one=1", "two=002", "three=3", "four=4"},
		},
		{
			oldEnv:   []string{"one=1", "=2"},
			newEnv:   []string{"=3", "two=2"},
			expected: []string{"one=1", "=3", "two=2"},
		},
		{
			oldEnv:   []string{"one=1", "two"},
			newEnv:   []string{"two=2", "three=3"},
			expected: []string{"one=1", "two=2", "three=3"},
		},
	}
	for _, tc := range tests {
		result := MergeEnv(tc.oldEnv, tc.newEnv)
		toCheck := map[string]struct{}{}
		for _, e := range tc.expected {
			toCheck[e] = struct{}{}
		}
		for _, e := range result {
			if _, exists := toCheck[e]; !exists {
				t.Errorf("old = %s, new = %s: %s not expected in result",
					strings.Join(tc.oldEnv, ","), strings.Join(tc.newEnv, ","), e)
				continue
			}
			delete(toCheck, e)
		}
		if len(toCheck) > 0 {
			t.Errorf("old = %s, new = %s: did not get expected values in result: %#v",
				strings.Join(tc.oldEnv, ","), strings.Join(tc.newEnv, ","), toCheck)
		}
	}
}
func TestNameForBuildVolume(t *testing.T) {
	type args struct {
		objName string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Secret One",
			args: args{objName: "secret-one"},
			want: fmt.Sprintf("secret-one-%s", buildVolumeSuffix),
		},
		{
			name: "ConfigMap One",
			args: args{objName: "configmap-one"},
			want: fmt.Sprintf("configmap-one-%s", buildVolumeSuffix),
		},
		{
			name: "Greater than 47 characters",
			args: args{objName: "build-volume-larger-than-47-characters-but-less-than-63"},
			want: fmt.Sprintf("build-volume-larger-than-47-characte-8c2b6813-%s", buildVolumeSuffix),
		},
		{
			name: "Should convert to lowercase",
			args: args{objName: "Secret-One"},
			want: fmt.Sprintf("secret-one-%s", buildVolumeSuffix),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NameForBuildVolume(tt.args.objName); got != tt.want {
				t.Errorf("NameForBuildVolume() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPathForBuildVolume(t *testing.T) {
	type args struct {
		objName string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Secret One",
			args: args{"secret-one"},
			want: filepath.Join(buildVolumeMountPath, fmt.Sprintf("secret-one-%s", buildVolumeSuffix)),
		},
		{
			name: "ConfigMap One",
			args: args{"configmap-one"},
			want: filepath.Join(buildVolumeMountPath, fmt.Sprintf("configmap-one-%s", buildVolumeSuffix)),
		},
		{
			name: "Greater than 47 characters",
			args: args{objName: "build-volume-larger-than-47-characters-but-less-than-63"},
			want: filepath.Join(buildVolumeMountPath, fmt.Sprintf("build-volume-larger-than-47-characte-8c2b6813-%s", buildVolumeSuffix)),
		},
		{
			name: "Should convert to lowercase",
			args: args{"Secret-One"},
			want: filepath.Join(buildVolumeMountPath, fmt.Sprintf("secret-one-%s", buildVolumeSuffix)),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := PathForBuildVolume(tt.args.objName); got != tt.want {
				t.Errorf("PathForBuildVolume() = %v, want %v", got, tt.want)
			}
		})
	}
}
