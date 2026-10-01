package config

import (
	"reflect"
	"testing"
)

func TestAdminOrigins(t *testing.T) {
	for _, tc := range []struct {
		name, env, list, app string
		want                 []string
		invalid              bool
	}{
		{name: "dev default", want: []string{"http://localhost:5174", "http://127.0.0.1:5174"}},
		{name: "trim list", list: " http://localhost:5174, http://127.0.0.1:5174 , ", want: []string{"http://localhost:5174", "http://127.0.0.1:5174"}},
		{name: "explicit overrides legacy", list: "http://localhost:5174", app: "http://localhost:5173", want: []string{"http://localhost:5174"}},
		{name: "production empty", env: "production", want: []string{}},
		{name: "production configured", env: "production", list: "https://admin.example.invalid", want: []string{"https://admin.example.invalid"}},
		{name: "production same origin", env: "production", app: "https://site.example.invalid", want: []string{"https://site.example.invalid"}},
		{name: "wildcard", list: "*", invalid: true},
		{name: "path", list: "https://site.example.invalid/path", invalid: true},
		{name: "insecure production", env: "production", list: "http://localhost:5174", invalid: true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			got, err := adminOrigins(tc.env, tc.list, tc.app)
			if (err != nil) != tc.invalid {
				t.Fatalf("unexpected error: %v", err)
			}
			if !tc.invalid && !reflect.DeepEqual(got, tc.want) {
				t.Fatalf("got %v, want %v", got, tc.want)
			}
		})
	}
}
