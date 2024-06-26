package driver

import (
	"go.einride.tech/can/pkg/candevice"
	"testing"
)

func TestCan(t *testing.T) {
	// Error handling omitted to keep example simple
	d, _ := candevice.New("can0")
	d.SetBitrate(250000)
	d.SetUp()
	defer d.SetDown()
}
