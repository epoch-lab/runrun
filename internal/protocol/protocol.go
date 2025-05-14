// 这份代码是用来和unirun服务器交互用的
// 改编自msojocs/AutoRun
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
