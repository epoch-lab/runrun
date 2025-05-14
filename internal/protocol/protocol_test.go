package protocol

import (
	"fmt"
	"testing"
)

const phone = "your_phone"
const passwd = "your_passwd"

func TestLogin(t *testing.T) {
	info := GenerateFakeClient()
	user, err := Login(phone, passwd, info)
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
	user, err := Login(phone, passwd, mockClientInfo)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("Login successful: %+v", user)
}
func TestGenerateTrack(t *testing.T) {
	str := genTrack(5024)
	fmt.Println(str)
}
func TestGetUserInfo(t *testing.T) {
	user, err := Login(phone, passwd, mockClientInfo)
	if err != nil {
		t.Fatal(err)
	}
	info, err := GetUserInfo(user.OauthToken)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("%v", info)
}
func TestGetSchoolBound(t *testing.T) {
	user, err := Login(phone, passwd, mockClientInfo)
	if err != nil {
		t.Fatal(err)
	}
	info, err := getSchoolBound(user.OauthToken, user.SchoolID)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("%v", info)
}
func TestGetRunStandard(t *testing.T) {
	user, err := Login(phone, passwd, mockClientInfo)
	if err != nil {
		t.Fatal(err)
	}
	info, err := GetRunStandard(user.OauthToken, user.SchoolID)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("%v", info)
}
func TestSubmitFixed(t *testing.T) {
	user, err := Login(phone, passwd, mockClientInfo)
	if err != nil {
		t.Fatal(err)
	}
	if err := Submit(*user, mockClientInfo, 57, 5120); err != nil {
		t.Fatal(err)
	}
}
