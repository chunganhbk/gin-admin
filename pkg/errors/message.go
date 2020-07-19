package errors

var MsgFlags = map[int]string{
	SUCCESS:                        "ok",
	INVALID_PARAMS:                 "Request parameter error - %s",
	ERROR_AUTH_CHECK_TOKEN_FAIL:    "Token authentication failed",
	ERROR_AUTH_CHECK_TOKEN_TIMEOUT: "Token time out",
	ERROR_AUTH_TOKEN:               "Token build failed",
	ERROR_AUTH:                     "Token error",
	ERROR:                          "fail",
	ERROR_INTERNAL_SERVER:          "Server error",
	ERROR_EXIST_EMAIL:              "The Email Address entered already exists in the system",
	ERROR_BAD_REQUEST:              "Request error",
	ERROR_INVALID_PARENT:           "Invalid parent node",
	ERROR_ALLOW_DELETE_WITH_CHILD:  "Contains children, cannot be deleted",
	ERROR_NOT_ALLOW_DELETE:         "Resources are not allowed to be deleted",
	ERROR_INVALID_OLD_PASS:         "Old password is incorrect",
	ERROR_NOT_FOUND:                "Resource does not exist",
	ERROR_PASSWORD_REQUIRED:        "Password is required",
	ERROR_EXIST_MENU_NAME:          "Menu name already exists",
	ERROR_USER_DISABLED:            "User is disabled, please contact administrator",
	ERROR_NO_PERRMISSION:           "No access",
	ERROR_METHOD_NOT_ALLOW:         "Method is not allowed",
	ERROR_TOO_MANY_REQUEST:         "Requests are too frequent",
	ERROR_LOGIN_FAILED:             "Email or password is invalid",
	ERROR_EXIST_ROLE:               "Role name already exists",
	ERROR_NOT_EXIST_USER:           "Account is invalid",
	ERROR_EXIST_ROLE_USER:          "The role has been given to the user and is not allowed to be deleted",
	ERROR_NOT_EXIST_ROLE:           "Role user is disabled, please contact administrator",
	ERROR_TOKEN_EXPIRED:            "Token is expired",
	ERROR_TOKEN_INVALID:            "Token not active yet",
	ERROR_TOKEN_MALFORMALED:        "That's not even a token",
}

// GetMsg get error information based on Code
func GetMsg(code int) string {
	msg, ok := MsgFlags[code]
	if ok {
		return msg
	}
	return MsgFlags[ERROR]
}
