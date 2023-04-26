package test

import (
	"context"
	"os"
	"os/signal"
	"testing"

	"github.com/hootrhino/rulex/glogger"
	"github.com/muka/go-bluetooth/api"
	"github.com/muka/go-bluetooth/api/beacon"
	"github.com/muka/go-bluetooth/bluez/profile/adapter"
	"github.com/muka/go-bluetooth/bluez/profile/device"
	eddystone "github.com/suapapa/go_eddystone"
)

func Test_BLE_Server(t *testing.T) {
	Run("AABBCCDD", false)
}

//shows how to watch for new devices and list them

func Run(adapterID string, onlyBeacon bool) error {

	//clean up connection on exit
	defer api.Exit()

	a, err := adapter.GetAdapter(adapterID)
	if err != nil {
		return err
	}

	glogger.GLogger.Debug("Flush cached devices")
	err = a.FlushDevices()
	if err != nil {
		return err
	}

	glogger.GLogger.Debug("Start discovery")
	discovery, cancel, err := api.Discover(a, nil)
	if err != nil {
		return err
	}
	defer cancel()

	go func() {

		for ev := range discovery {

			if ev.Type == adapter.DeviceRemoved {
				continue
			}

			dev, err := device.NewDevice1(ev.Path)
			if err != nil {
				glogger.GLogger.Errorf("%s: %s", ev.Path, err)
				continue
			}

			if dev == nil {
				glogger.GLogger.Errorf("%s: not found", ev.Path)
				continue
			}

			glogger.GLogger.Infof("name=%s addr=%s rssi=%d", dev.Properties.Name, dev.Properties.Address, dev.Properties.RSSI)

			go func(ev *adapter.DeviceDiscovered) {
				err = handleBeacon(dev)
				if err != nil {
					glogger.GLogger.Errorf("%s: %s", ev.Path, err)
				}
			}(ev)
		}

	}()

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt) // get notified of all OS signals

	sig := <-ch
	glogger.GLogger.Infof("Received signal [%v]; shutting down...\n", sig)

	return nil
}

func handleBeacon(dev *device.Device1) error {

	b, err := beacon.NewBeacon(dev)
	if err != nil {
		return err
	}

	beaconUpdated, err := b.WatchDeviceChanges(context.Background())
	if err != nil {
		return err
	}

	isBeacon := <-beaconUpdated
	if !isBeacon {
		return nil
	}

	name := b.Device.Properties.Alias
	if name == "" {
		name = b.Device.Properties.Name
	}

	glogger.GLogger.Debugf("Found beacon %s %s", b.Type, name)

	if b.IsEddystone() {
		ed := b.GetEddystone()
		switch ed.Frame {
		case eddystone.UID:
			glogger.GLogger.Debugf(
				"Eddystone UID %s instance %s (%ddbi)",
				ed.UID,
				ed.InstanceUID,
				ed.CalibratedTxPower,
			)
		case eddystone.TLM:
			glogger.GLogger.Debugf(
				"Eddystone TLM temp:%.0f batt:%d last reboot:%d advertising pdu:%d (%ddbi)",
				ed.TLMTemperature,
				ed.TLMBatteryVoltage,
				ed.TLMLastRebootedTime,
				ed.TLMAdvertisingPDU,
				ed.CalibratedTxPower,
			)
		case eddystone.URL:
			glogger.GLogger.Debugf(
				"Eddystone URL %s (%ddbi)",
				ed.URL,
				ed.CalibratedTxPower,
			)
		}

	}
	if b.IsIBeacon() {
		ibeacon := b.GetIBeacon()
		glogger.GLogger.Debugf(
			"IBeacon %s (%ddbi) (major=%d minor=%d)",
			ibeacon.ProximityUUID,
			ibeacon.MeasuredPower,
			ibeacon.Major,
			ibeacon.Minor,
		)
	}

	return nil
}
