// 这份代码是用来和unirun服务器交互用的
// 改编自msojocs/AutoRun
package internal

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"net/url"
	"strings"
	"time"
)

func stringToMD5(input string) string {
	hasher := md5.New()
	hasher.Write([]byte(input))
	return hex.EncodeToString(hasher.Sum(nil))
}
func sign(param *map[string]string, body string) string {
	const APPKEY = "389885588s0648fa"
	const APPSECRET = "56E39A1658455588885690425C0FD16055A21676"
	str := ""
	if param != nil {
		for k, v := range *param {
			str += k
			str += v
		}
	}
	str += APPKEY
	str += APPSECRET
	// why the fuck the replacing action is needed
	replaced := false
	str = strings.Map(func(r rune) rune {
		switch r {
		case ' ', '~', '!', '(', ')', '\'':
			replaced = true
			return -1
		default:
			return r
		}
	}, str)
	if replaced {
		str = url.QueryEscape(str)
	}
	ret := strings.ToUpper(stringToMD5(str))
	if replaced {
		ret += "encodeutf8"
	}
	return ret
}

// Uppercase
// What if you guys really want to use it...
type ClientInfo struct {
	AppVersion, Brand, DeviceToken, DeviceType, MobileType, SysVersion string
}

// this function is used to generate random phone model information
// note that appVersion is always 1.8.2, which could lead to suspicion!!!
func generateFakeClient() ClientInfo {
	brands := []string{
		"Xiaomi", "HUAWEI", "HONOR", "OPPO", "vivo", "OnePlus", "Samsung",
	}

	// Common Android versions
	androidVersions := []string{
		"11.0", "12.0", "13.0", "14.0",
	}

	// Generate random device token (16 characters)
	const letterBytes = "abcdef0123456789"
	deviceToken := make([]byte, 16)
	for i := range deviceToken {
		deviceToken[i] = letterBytes[rand.Intn(len(letterBytes))]
	}

	// Random model numbers
	modelNumbers := []string{"2201123C", "2207122C", "22081212C", "23046PNC9C"}

	brand := brands[rand.Intn(len(brands))]
	model := modelNumbers[rand.Intn(len(modelNumbers))]

	return ClientInfo{
		AppVersion:  "1.8.2",                                          // Fixed version
		Brand:       brand,                                            // Random brand
		DeviceToken: string(deviceToken),                              // Random device token
		DeviceType:  fmt.Sprintf("%s_%s", brand, model),               // Brand_Model format
		MobileType:  "android",                                        // Fixed as android
		SysVersion:  androidVersions[rand.Intn(len(androidVersions))], // Random Android version
	}
}

// I bet case-insensitive match will work properly
// how could this be this long??
type UserInfo struct {
	UserID           int64
	StudentID        int64
	RegisterCode     string
	StudentName      string
	Gender           string
	SchoolID         int64
	SchoolName       string
	ClassID          int64
	StudentClass     int32
	ClassName        string
	StartSchool      int32
	CollegeCode      string
	CollegeName      string
	MajorCode        string
	MajorName        string
	NationCode       string
	Birthday         string
	IDCardNo         string
	AddrDetail       string
	StudentSource    string
	UserVerifyStatus string //hey how could this be string??
	OauthToken       struct {
		RefreshToken string
		Token        string
	}
}

// I bet case-insensitive match will work properly
type response[T any] struct {
	Code     int32
	Msg      string
	Response T
}

const appKey = "389885588s0648fa"
const host = "https://run-lb.tanmasports.com/"

func Login(phone, password string, info ClientInfo) (UserInfo, error) {
	const API = host + "api/run/login"
	// Convert body to JSON string
	body, err := json.Marshal(map[string]string{
		"appVersion":  info.AppVersion,
		"brand":       info.Brand,
		"deviceToken": info.DeviceToken,
		"deviceType":  info.DeviceType,
		"mobileType":  info.MobileType,
		"password":    password,
		"sysVersion":  info.SysVersion,
		"userPhone":   phone,
	})
	if err != nil {
		return UserInfo{}, fmt.Errorf("marshal body failed: %w", err)
	}
	sign := sign(nil, string(body))
	req, _ := http.NewRequest(http.MethodPost, API, bytes.NewBuffer(body))
	req.Header.Set("sign", sign)
	req.Header.Set("appkey", appKey)
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")
	resp, err := (&http.Client{Timeout: 10 * time.Second}).Do(req)
	if err != nil {
		return UserInfo{}, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

}
