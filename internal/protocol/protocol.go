// 这份代码是用来和unirun服务器交互用的
// 改编自msojocs/AutoRun
package protocol

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const appKey = "389885588s0648fa"
const host = "https://run-lb.tanmasports.com/"

func Login(phone, password string, info ClientInfo) (*UserInfo, error) {
	const API = host + "v1/auth/login/password"
	// Convert body to JSON string
	body, err := json.Marshal(map[string]string{
		"appVersion":  info.AppVersion,
		"brand":       info.Brand,
		"deviceToken": info.DeviceToken,
		"deviceType":  info.DeviceType,
		"mobileType":  info.MobileType,
		"password":    stringToMD5(password),
		"sysVersion":  info.SysVersion,
		"userPhone":   phone,
	})
	if err != nil {
		return nil, fmt.Errorf("marshal body failed: %w", err)
	}
	sign := sign(nil, string(body))
	req, _ := http.NewRequest(http.MethodPost, API, bytes.NewBuffer(body))
	req.Header.Set("sign", sign)
	req.Header.Set("appkey", appKey)
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")
	resp, err := (&http.Client{Timeout: 10 * time.Second}).Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response body failed: %w", err)
	}
	var respBody response[UserInfo]
	if err := json.Unmarshal(bodyBytes, &respBody); err != nil {
		return nil, fmt.Errorf("parse response failed: %w", err)
	}
	if respBody.Code == 10000 {
		return &respBody.Response, nil
	} else {
		return nil, fmt.Errorf("login failed: %s", respBody.Msg)
	}
}
