package config

import (
	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/snap"
)

var SnapClient snap.Client

func InitMidtrans() {
	// Server Key Sandbox
	midtrans.ServerKey = os.Getenv("MIDTRANS_SERVER_KEY")
	midtrans.Environment = midtrans.Sandbox

	SnapClient.New(midtrans.ServerKey, midtrans.Sandbox)
}
