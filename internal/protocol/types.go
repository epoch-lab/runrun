package protocol

// Uppercase
// What if you guys really want to use it...
type ClientInfo struct {
	AppVersion, Brand, DeviceToken, DeviceType, MobileType, SysVersion string
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
