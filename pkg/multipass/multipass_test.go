package multipass

import (
	"testing"
)

func TestParseIP(t *testing.T) {
	tests := []struct {
		name    string
		output  string
		want    string
		wantErr bool
	}{
		{
			name: "Valid output",
			output: `Name:           eks-vm
State:          Running
IPv4:           192.168.64.2
Release:        Ubuntu 22.04 LTS`,
			want:    "192.168.64.2",
			wantErr: false,
		},
		{
			name: "Missing IP",
			output: `Name:           eks-vm
State:          Stopped`,
			want:    "",
			wantErr: true,
		},
		{
			name: "Garbage data",
			output: `Some random text
that does not make sense`,
			want:    "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseIP(tt.output)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseIP() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("parseIP() = %v, want %v", got, tt.want)
			}
		})
	}
}
