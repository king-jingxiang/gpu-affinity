// +build linux

package ibverbs

import (
	"C"
)
import (
	"reflect"
	"testing"
)

func TestIbvGetDeviceList(t *testing.T) {
	tests := []struct {
		name    string
		want    []IbvDevice
		wantErr bool
	}{
		// TODO: Add test cases.
		"mlx5_0",
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := IbvGetDeviceList()
			if (err != nil) != tt.wantErr {
				t.Errorf("IbvGetDeviceList() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got.name, tt.want) {
				t.Errorf("IbvGetDeviceList() = %v, want %v", got, tt.want)
			}
		})
	}
}
