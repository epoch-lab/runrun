package protocol

import (
	"crypto/md5"
	"encoding/hex"
	"net/url"
	"strings"
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
	str += body
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
