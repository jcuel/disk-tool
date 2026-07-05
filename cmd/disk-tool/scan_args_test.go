package main

import "testing"

func TestParseScanArgs(t *testing.T) {
	tests := []struct {
		name    string
		args    []string
		want    scanFlags
		wantErr bool
	}{
		{
			name: "flags after path",
			args: []string{`C:\Users`, "--json"},
			want: scanFlags{json: true, path: `C:\Users`},
		},
		{
			name: "flags before path",
			args: []string{"--json", `C:\Users`},
			want: scanFlags{json: true, path: `C:\Users`},
		},
		{
			name: "json and full after path",
			args: []string{".", "--json", "--full"},
			want: scanFlags{json: true, full: true, path: "."},
		},
		{
			name: "json and full before path",
			args: []string{"--full", "--json", "/tmp"},
			want: scanFlags{json: true, full: true, path: "/tmp"},
		},
		{
			name:    "missing path",
			args:    []string{"--json"},
			wantErr: true,
		},
		{
			name:    "unknown flag",
			args:    []string{".", "--verbose"},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseScanArgs(tt.args)
			if tt.wantErr {
				if err == nil {
					t.Fatal("expected error")
				}
				return
			}
			if err != nil {
				t.Fatal(err)
			}
			if got != tt.want {
				t.Fatalf("got %+v, want %+v", got, tt.want)
			}
		})
	}
}
