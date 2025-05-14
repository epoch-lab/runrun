package protocol

import "testing"

func TestLogin(t *testing.T) {
	info := GenerateFakeClient()
	user, err := Login("your_phone", "your_passwd", info)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("Login successful: %+v", user)
}

var mockClientInfo = ClientInfo{
	AppVersion:  "1.8.2",
	Brand:       "Xiaomi",
	DeviceToken: "",
	DeviceType:  "Xiaomi_2201123C",
	MobileType:  "android",
	SysVersion:  "13.0",
}

func TestLoginFixed(t *testing.T) {
	user, err := Login("your_phone", "your_passwd", mockClientInfo)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("Login successful: %+v", user)
}
