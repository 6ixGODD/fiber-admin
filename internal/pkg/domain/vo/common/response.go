package common

type (
	LoginResponse struct {
		AccessToken  string  `json:"access_token"`
		RefreshToken string  `json:"refresh_token"`
		ExpiresIn    float64 `json:"expires_in"`
		Meta         struct {
			UserID   string `json:"user_id"`
			Username string `json:"username"`
			Email    string `json:"email"`
			Role     string `json:"role"`
		} `json:"meta"`
	}

	RefreshTokenResponse struct {
		AccessToken  string  `json:"access_token"`
		RefreshToken string  `json:"refresh_token"`
		ExpiresIn    float64 `json:"expires_in"`
		Meta         struct {
			UserID   string `json:"user_id"`
			Username string `json:"username"`
			Email    string `json:"email"`
			Role     string `json:"role"`
		} `json:"meta"`
	}

	GetNoticeResponse struct {
		NoticeID   string `json:"notice_id"`
		Title      string `json:"title"`
		Content    string `json:"content"`
		NoticeType string `json:"type"`
		CreatedAt  string `json:"created_at"`
		UpdatedAt  string `json:"updated_at"`
	}

	NoticeSummary struct {
		NoticeID   string `json:"notice_id"`
		Title      string `json:"title"`
		NoticeType string `json:"type"`
		CreatedAt  string `json:"created_at"`
	}

	GetNoticeListResponse struct {
		Total             int64            `json:"total"`
		NoticeSummaryList []*NoticeSummary `json:"notice_summary_list"`
	}

	GetDocumentationResponse struct {
		DocumentID string `json:"document_id"`
		Title      string `json:"title"`
		Content    string `json:"content"`
		CreatedAt  string `json:"created_at"`
		UpdatedAt  string `json:"updated_at"`
	}

	DocumentationSummary struct {
		DocumentID string `json:"document_id"`
		Title      string `json:"title"`
		CreatedAt  string `json:"created_at"`
	}

	GetDocumentationListResponse struct {
		Total                    int64                   `json:"total"`
		DocumentationSummaryList []*DocumentationSummary `json:"documentation_summary_list"`
	}

	GetProfileResponse struct {
		UserID       string `json:"user_id"`
		Username     string `json:"username"`
		Email        string `json:"email"`
		Role         string `json:"role"`
		Organization string `json:"organization"`
		LastLogin    string `json:"last_login"`
	}
)
