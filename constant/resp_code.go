package constant

const (
	// user related response
	UserAlreadyExist   = 1001
	UserAddSuccess     = 1002
	UserAuthSuccess    = 1003
	UserAuthError      = 1004
	UserAuthTimeout    = 1005
	UserSignoutSuccess = 1006

	// JWT related response
	JwtGenerationError = 2001
	JwtMissingError    = 2002
	JwtParseError      = 2003

	//Bucket related response
	BucketAlreadyExist  = 3001
	BucketAddSuccess    = 3002
	BucketNotExist      = 3003
	BucketDeleteSuccess = 3004
	BucketUpdateSuccess = 3005
	BucketGetSuccess    = 3006

	// Photo related response
	PhotoAlreadyExist  = 4001
	PhotoAddInProcess = 4002
	PhotoUploadSuccess = 4003
	PhotoUploadError   = 4004
	PhotoNotExist      = 4005
	PhotoDeleteSuccess = 4006
	PhotoUpdateSuccess = 4007
	PhotoGetSuccess    = 4008

	// Internal server response
	InternalServerError = 5001
	PaginationSuccess    = 6001
	InvalidParams        = 7001
)

var Message map[int]string

func init() {
	Message = make(map[int]string)
	Message[InvalidParams] = "Invalid parameters."
	Message[UserAlreadyExist] = "User already exists."
	Message[UserAddSuccess] = "Add user success."
	Message[UserAuthSuccess] = "User authentication success."
	Message[UserAuthError] = "User authentication fail."
	Message[UserAuthTimeout] = "User authentication timeout."
	Message[UserSignoutSuccess] = "User sign out success."
	Message[JwtGenerationError] = "JWT generation fail."
	Message[JwtMissingError] = "JWT is missing."
	Message[InternalServerError] = "Internal server error."
	Message[BucketAlreadyExist] = "Bucket already exists."
	Message[BucketAddSuccess] = "Add bucket success."
	Message[BucketNotExist] = "Bucket does not exist."
	Message[BucketDeleteSuccess] = "Bucket delete success."
	Message[BucketUpdateSuccess] = "Bucket update success."
	Message[BucketGetSuccess] = "Bucket get success."
	Message[PhotoAlreadyExist] = "Photo already exists."
	Message[PhotoAddInProcess] = "Adding photo is in process."
	Message[PhotoUploadSuccess] = "Photo upload success."
	Message[PhotoUploadError] = "Photo upload error."
	Message[PhotoNotExist] = "Photo does not exist."
	Message[PhotoDeleteSuccess] = "Photo delete success."
	Message[PhotoGetSuccess] = "Photo get success."
}

// GetMessage func to get response description according to the code
func GetMessage(code int) string {
	msg, ok := Message[code]
	if ok {
		return msg
	}
	return ""
}
