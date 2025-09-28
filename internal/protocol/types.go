package protocol

// Uppercase
// What if you guys really want to use it...
type ClientInfo struct {
	AppVersion, Brand, DeviceToken, DeviceType, MobileType, SysVersion string
}
type Oauth struct {
	RefreshToken string
	Token        string
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
	OauthToken       Oauth
}

// I bet case-insensitive match will work properly
type response[T any] struct {
	Code     int32
	Msg      string
	Response T
}

type schoolBound struct {
	SiteName    string
	SiteBound   string
	BoundCenter string
}
type RunStandard struct {
	StandardID     int64
	SchoolID       int64
	BoyOnceTimeMin int64
	BoyOnceTimeMax int64
	SemesterYear   string
}
type location struct {
	ID       int32
	Location string
	Edge     []int32
}

// AuthRequest defines the structure for the authentication endpoint request body.
type AuthRequest struct {
	Account  string `json:"account"`
	Password string `json:"password"`
}
