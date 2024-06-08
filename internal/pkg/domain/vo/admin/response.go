package admin

type (
	GetUserResponse struct {
		UserID       string `json:"user_id"`
		Username     string `json:"username"`
		Email        string `json:"email"`
		Role         string `json:"role"`
		Organization string `json:"organization"`
		LastLogin    string `json:"last_login"`
		CreatedAt    string `json:"created_at"`
		UpdatedAt    string `json:"updated_at"`
	}

	GetUserListResponse struct {
		Total    int64              `json:"total"`
		UserList []*GetUserResponse `json:"user_list"`
	}

	GetLoginLogResponse struct {
		LoginLogID string `json:"login_log_id"`
		UserID     string `json:"user_id"`
		Username   string `json:"username"`
		Email      string `json:"email"`
		IPAddress  string `json:"ip_address"`
		UserAgent  string `json:"user_agent"`
		CreatedAt  string `json:"created_at"`
	}

	GetLoginLogListResponse struct {
		Total        int64                  `json:"total"`
		LoginLogList []*GetLoginLogResponse `json:"login_log_list"`
	}

	GetOperationLogResponse struct {
		OperationLogID string `json:"operation_log_id"`
		UserID         string `json:"user_id"`
		Username       string `json:"username"`
		Email          string `json:"email"`
		IPAddress      string `json:"ip_address"`
		UserAgent      string `json:"user_agent"`
		Operation      string `json:"operation"`
		EntityID       string `json:"entity_id"`
		EntityType     string `json:"entity_type"`
		Description    string `json:"description"`
		Status         string `json:"status"`
		CreatedAt      string `json:"created_at"`
	}

	GetOperationLogListResponse struct {
		Total            int64                      `json:"total"`
		OperationLogList []*GetOperationLogResponse `json:"operation_log_list"`
	}

	GetErrorLogResponse struct {
		ErrorLogID     string `json:"error_log_id"`
		UserID         string `json:"user_id"`
		Username       string `json:"username"`
		IPAddress      string `json:"ip_address"`
		UserAgent      string `json:"user_agent"`
		RequestURL     string `json:"request_uri"`
		RequestMethod  string `json:"request_method"`
		RequestPayload string `json:"request_payload"`
		ErrorCode      string `json:"error_code"`
		ErrorMsg       string `json:"error_msg"`
		Stack          string `json:"stack"`
		CreatedAt      string `json:"created_at"`
	}

	GetErrorLogListResponse struct {
		Total        int64                  `json:"total"`
		ErrorLogList []*GetErrorLogResponse `json:"error_log_list"`
	}
)
