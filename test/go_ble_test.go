package test

import (
	"testing"

	"tinygo.org/x/bluetooth"
)

var adapter = bluetooth.DefaultAdapter

func Test_bluetooth(t *testing.T) {
	if err := adapter.Enable(); err != nil {
		t.Fatal(err)
	}

	err1 := adapter.Scan(func(adapter *bluetooth.Adapter, device bluetooth.ScanResult) {
		t.Log("found device:", device.Address.String(), device.RSSI, device.LocalName())
	})
	if err1 != nil {
		t.Fatal(err1)
	}
}
