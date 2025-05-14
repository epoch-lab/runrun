// Package protocol 提供与 unirun 服务器交互的相关函数。
// 改编自 msojocs/AutoRun
package protocol

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	c "runrun/common"
	"time"
)

const appKey = "389885588s0648fa"
const host = "https://run-lb.tanmasports.com/"

// Login 使用手机号、密码和客户端信息进行登录。
// 成功返回 UserInfo，失败返回错误信息。
// 考虑使用 GenerateFakeClient()随机生成一个可被此函数接受的info 这个info是怎么构造的依你而定
// 返回的UserInfo可以直接用到后续的Submit函数 其OauthToken成员则可以直接用到后续的Get系列函数中
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
		return nil, fmt.Errorf("%s: marshal body failed: %w", c.CurrentFunctionName(), err)
	}
	req, _ := http.NewRequest(http.MethodPost, API, bytes.NewBuffer(body))
	req.Header.Set("sign", sign(nil, string(body)))
	req.Header.Set("appkey", appKey)
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")
	resp, err := (&http.Client{
		Timeout:   10 * time.Second,
		Transport: &http.Transport{TLSClientConfig: &tls.Config{MinVersion: tls.VersionTLS12}},
	}).Do(req)
	if err != nil {
		return nil, fmt.Errorf("%s: request failed: %w", c.CurrentFunctionName(), err)
	}
	defer resp.Body.Close()
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("%s: read response body failed: %w", c.CurrentFunctionName(), err)
	}
	var respBody response[UserInfo]
	if err := json.Unmarshal(bodyBytes, &respBody); err != nil {
		return nil, fmt.Errorf("%s: parse response failed: %w", c.CurrentFunctionName(), err)
	}
	if respBody.Code == 10000 {
		return &respBody.Response, nil
	} else {
		return nil, fmt.Errorf("%s: login failed: %s", c.CurrentFunctionName(), respBody.Msg)
	}
}

// GetUserInfo 根据 Oauth token 获取用户信息。
// 成功返回 UserInfo，失败返回错误信息。
func GetUserInfo(token Oauth) (*UserInfo, error) {
	const API = host + "v1/auth/query/token"
	req, _ := http.NewRequest(http.MethodGet, API, nil)
	req.Header.Set("sign", sign(nil, ""))
	req.Header.Set("token", token.Token)
	req.Header.Set("appkey", appKey)
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")
	req.Header.Set("User-Agent", "okhttp/3.12.0")
	resp, err := (&http.Client{
		Timeout:   10 * time.Second,
		Transport: &http.Transport{TLSClientConfig: &tls.Config{MinVersion: tls.VersionTLS12}},
	}).Do(req)
	if err != nil {
		return nil, fmt.Errorf("%s: request failed: %w", c.CurrentFunctionName(), err)
	}
	defer resp.Body.Close()
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("%s: read response body failed: %w", c.CurrentFunctionName(), err)
	}
	var respBody response[UserInfo]
	if err := json.Unmarshal(bodyBytes, &respBody); err != nil {
		return nil, fmt.Errorf("%s: parse response failed: %w", c.CurrentFunctionName(), err)
	}
	return &respBody.Response, nil
}

// getSchoolBound 获取指定学校的围栏信息。
// 成功返回 schoolBound 切片，失败返回错误信息。
func getSchoolBound(token Oauth, schoolID int64) ([]schoolBound, error) {
	API := fmt.Sprintf("%sv1/unirun/querySchoolBound?schoolId=%d", host, schoolID)
	req, _ := http.NewRequest(http.MethodGet, API, nil)
	req.Header.Set("sign", sign(&map[string]string{
		"schoolId": fmt.Sprintf("%d", schoolID),
	}, ""))
	req.Header.Set("token", token.Token)
	req.Header.Set("appkey", appKey)
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")
	resp, err := (&http.Client{
		Timeout:   10 * time.Second,
		Transport: &http.Transport{TLSClientConfig: &tls.Config{MinVersion: tls.VersionTLS12}},
	}).Do(req)
	if err != nil {
		return nil, fmt.Errorf("%s: request failed: %w", c.CurrentFunctionName(), err)
	}
	defer resp.Body.Close()
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("%s: read response body failed: %w", c.CurrentFunctionName(), err)
	}
	var respBody response[[]schoolBound]
	if err := json.Unmarshal(bodyBytes, &respBody); err != nil {
		return nil, fmt.Errorf("%s: parse response failed: %w", c.CurrentFunctionName(), err)
	}
	return respBody.Response, nil
}

// GetRunStandard 获取指定学校的跑步标准。
// 成功返回 RunStandard 指针，失败返回错误信息。
func GetRunStandard(token Oauth, schoolID int64) (*RunStandard, error) {
	API := fmt.Sprintf("%sv1/unirun/query/runStandard?schoolId=%d", host, schoolID)
	req, _ := http.NewRequest(http.MethodGet, API, nil)
	req.Header.Set("sign", sign(&map[string]string{
		"schoolId": fmt.Sprintf("%d", schoolID),
	}, ""))
	req.Header.Set("token", token.Token)
	req.Header.Set("appkey", appKey)
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")
	resp, err := (&http.Client{
		Timeout:   10 * time.Second,
		Transport: &http.Transport{TLSClientConfig: &tls.Config{MinVersion: tls.VersionTLS12}},
	}).Do(req)
	if err != nil {
		return nil, fmt.Errorf("%s: request failed: %w", c.CurrentFunctionName(), err)
	}
	defer resp.Body.Close()
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("%s: read response body failed: %w", c.CurrentFunctionName(), err)
	}
	var respBody response[RunStandard]
	if err := json.Unmarshal(bodyBytes, &respBody); err != nil {
		return nil, fmt.Errorf("%s: parse response failed: %w", c.CurrentFunctionName(), err)
	}
	return &respBody.Response, nil
}

// Submit 提交一次跑步记录。
// user：用户信息 由Login得来 自行注意OauthToken不可持久化保存
// client：客户端信息 用GenerateFakeClient()生成或者怎么样产生
// duration：用时（分钟）
// distance：距离（米）
// 成功返回 nil，失败返回错误信息。
func Submit(user UserInfo, client ClientInfo, duration int32, distance int64) error {
	const API = host + "v1/unirun/save/run/record/new"
	runstd, err := GetRunStandard(user.OauthToken, user.SchoolID)
	if err != nil {
		return fmt.Errorf("%s: %s", c.CurrentFunctionName(), err.Error())
	}
	yearSemester := runstd.SemesterYear
	siteBound, err := getSchoolBound(user.OauthToken, user.SchoolID)
	if err != nil {
		return fmt.Errorf("%s: %s", c.CurrentFunctionName(), err.Error())
	}
	body, err := json.Marshal(map[string]string{
		"againRunStatus":     "0",
		"againRunTime":       "0",
		"userId":             fmt.Sprintf("%d", user.UserID),
		"appVersion":         client.AppVersion,
		"brand":              client.Brand,
		"mobileType":         client.MobileType,
		"sysVersion":         client.SysVersion,
		"runDistance":        fmt.Sprintf("%d", distance),
		"runTime":            fmt.Sprintf("%d", duration),
		"yearSemester":       yearSemester,
		"realityTrackPoints": siteBound[0].SiteBound + "--",
		"recordDate":         time.Now().Format("2006-01-02"),
		"trackPoints":        genTrack(distance),
		"distanceTimeStatus": "1",
		"innerSchool":        "1",
		"vocalStatus":        "1",
	})
	if err != nil {
		return fmt.Errorf("%s: marshal body failed: %w", c.CurrentFunctionName(), err)
	}
	req, _ := http.NewRequest(http.MethodPost, API, bytes.NewBuffer(body))
	req.Header.Set("sign", sign(nil, string(body)))
	req.Header.Set("token", user.OauthToken.Token)
	req.Header.Set("appkey", appKey)
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")
	resp, err := (&http.Client{
		Timeout:   10 * time.Second,
		Transport: &http.Transport{TLSClientConfig: &tls.Config{MinVersion: tls.VersionTLS12}},
	}).Do(req)
	if err != nil {
		return fmt.Errorf("%s: request failed: %w", c.CurrentFunctionName(), err)
	}
	defer resp.Body.Close()
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("%s: read response body failed: %w", c.CurrentFunctionName(), err)
	}
	var respBody response[struct{}]
	if err := json.Unmarshal(bodyBytes, &respBody); err != nil {
		return fmt.Errorf("%s: parse response failed: %w", c.CurrentFunctionName(), err)
	}
	if respBody.Code == 10000 {
		return nil
	} else {
		return fmt.Errorf("%s: submit failed: %s", c.CurrentFunctionName(), respBody.Msg)
	}
}
