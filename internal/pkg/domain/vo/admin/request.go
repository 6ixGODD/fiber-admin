package admin

type (
	InsertNoticeRequest struct {
		Title      *string `json:"title" validate:"required,max=100,min=1"`
		Content    *string `json:"content" validate:"required,max=10000,min=1"`
		NoticeType *string `json:"notice_type" validate:"required,noticeType"`
	}

	UpdateNoticeRequest struct {
		NoticeID   *string `json:"notice_id" validate:"required,mongodb"`
		Title      *string `json:"title" validate:"omitnil,max=100,min=1"`
		Content    *string `json:"content" validate:"omitnil,max=10000,min=1"`
		NoticeType *string `json:"notice_type" validate:"omitnil,noticeType"`
	}

	DeleteNoticeRequest struct {
		NoticeID *string `query:"noticeID" validate:"required,mongodb"`
	}

	InsertUserRequest struct {
		Username     *string `json:"username" validate:"required,min=3,max=20"`
		Email        *string `json:"email" validate:"required,email,max=100"`
		Password     *string `json:"password" validate:"required,min=8,max=20"`
		Organization *string `json:"organization" validate:"required,max=100"`
	}

	GetUserRequest struct {
		UserID *string `query:"userID" validate:"required,mongodb"`
	}

	GetUserListRequest struct {
		Page               *int64  `query:"page" validate:"required,numeric,min=1"`
		PageSize           *int64  `query:"pageSize" validate:"required,numeric,min=1,max=100"`
		Desc               *bool   `query:"desc" validate:"required"`
		Role               *string `query:"role" validate:"omitnil,userRole"`
		LastLoginStartTime *string `query:"lastLoginStartTime" validate:"omitnil,rfc3339,earlierThan=LastLoginEndTime"`
		LastLoginEndTime   *string `query:"lastLoginEndTime" validate:"omitnil,rfc3339"`
		CreateStartTime    *string `query:"createStartTime" validate:"omitnil,rfc3339,earlierThan=CreateEndTime"`
		CreateEndTime      *string `query:"createEndTime" validate:"omitnil,rfc3339"`
		Query              *string `query:"query" validate:"omitnil,max=100"`
	}

	UpdateUserRequest struct {
		UserID       *string `json:"user_id" validate:"required"`
		Username     *string `json:"username" validate:"omitnil,min=3,max=20"`
		Email        *string `json:"email" validate:"omitnil,email"`
		Organization *string `json:"organization" validate:"omitnil,max=100"`
	}

	DeleteUserRequest struct {
		UserID *string `query:"userID" validate:"required,mongodb"`
	}

	ChangeUserPasswordRequest struct {
		UserID      *string `json:"user_id" validate:"required,mongodb"`
		NewPassword *string `json:"new_password" validate:"required,min=8,max=20"`
	}

	InsertDocumentationRequest struct {
		Title   *string `json:"title" validate:"required,max=100,min=1"`
		Content *string `json:"content" validate:"required,max=10000,min=1"`
	}

	UpdateDocumentationRequest struct {
		DocumentationID *string `json:"documentation_id" validate:"required,mongodb"`
		Title           *string `json:"title" validate:"omitnil,max=100,min=1"`
		Content         *string `json:"content" validate:"omitnil,max=10000,min=1"`
	}

	DeleteDocumentationRequest struct {
		DocumentationID *string `query:"documentationID" validate:"required,mongodb"`
	}

	GetLoginLogListRequest struct {
		Page            *int64  `query:"page" validate:"required,numeric,min=1"`
		PageSize        *int64  `query:"pageSize" validate:"required,numeric,min=1,max=100"`
		Desc            *bool   `query:"desc" validate:"required"`
		Query           *string `query:"query" validate:"omitnil,max=100"`
		CreateStartTime *string `query:"createStartTime" validate:"omitnil,rfc3339,earlierThan=CreateEndTime"`
		CreateEndTime   *string `query:"createEndTime" validate:"omitnil,rfc3339"`
	}

	GetOperationLogListRequest struct {
		Page            *int64  `query:"page" validate:"required,numeric,min=1"`
		PageSize        *int64  `query:"pageSize" validate:"required,numeric,min=1,max=100"`
		Desc            *bool   `query:"desc" validate:"required"`
		Query           *string `query:"query" validate:"omitnil,max=100"`
		Operation       *string `query:"operation" validate:"omitnil,operationType"`
		EntityType      *string `query:"entityType" validate:"omitnil,entityType"`
		Status          *string `query:"status" validate:"omitnil,operationStatus"`
		CreateStartTime *string `query:"createStartTime" validate:"omitnil,rfc3339,earlierThan=CreateEndTime"`
		CreateEndTime   *string `query:"createEndTime" validate:"omitnil,rfc3339"`
	}
)
