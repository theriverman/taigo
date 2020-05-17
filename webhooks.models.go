package gotaiga

import "time"

// Webhook => https://taigaio.github.io/taiga-doc/dist/api.html#object-webhook-detail
type Webhook struct {
	ID          int    `json:"id,omitempty"`
	Key         string `json:"key,omitempty"`
	LogsCounter int    `json:"logs_counter,omitempty"`
	Name        string `json:"name,omitempty"`
	Project     int    `json:"project,omitempty"`
	URL         string `json:"url,omitempty"`
}

// WebhookLog => https://taigaio.github.io/taiga-doc/dist/api.html#object-webhook-log-detail
type WebhookLog struct {
	ID          int    `json:"id"`
	Webhook     int    `json:"webhook"`
	URL         string `json:"url"`
	Status      int    `json:"status"`
	RequestData struct {
		By struct {
			ID         int         `json:"id"`
			Photo      interface{} `json:"photo"`
			Username   string      `json:"username"`
			FullName   string      `json:"full_name"`
			Permalink  string      `json:"permalink"`
			GravatarID string      `json:"gravatar_id"`
		} `json:"by"`
		Data struct {
			Test string `json:"test"`
		} `json:"data"`
		Date   time.Time `json:"date"`
		Type   string    `json:"type"`
		Action string    `json:"action"`
	} `json:"request_data"`
	RequestHeaders struct {
		ContentType            string `json:"Content-Type"`
		ContentLength          string `json:"Content-Length"`
		XHubSignature          string `json:"X-Hub-Signature"`
		XTAIGAWEBHOOKSIGNATURE string `json:"X-TAIGA-WEBHOOK-SIGNATURE"`
	} `json:"request_headers"`
	ResponseData    string `json:"response_data"`
	ResponseHeaders struct {
		Date                       string `json:"Date"`
		Vary                       string `json:"Vary"`
		Pragma                     string `json:"Pragma"`
		Server                     string `json:"Server"`
		Expires                    string `json:"Expires"`
		Connection                 string `json:"Connection"`
		SetCookie                  string `json:"Set-Cookie"`
		ContentType                string `json:"Content-Type"`
		CacheControl               string `json:"Cache-Control"`
		ReferrerPolicy             string `json:"Referrer-Policy"`
		TransferEncoding           string `json:"Transfer-Encoding"`
		AccessControlAllowOrigin   string `json:"Access-Control-Allow-Origin"`
		AccessControlExposeHeaders string `json:"Access-Control-Expose-Headers"`
	} `json:"response_headers"`
	Duration     float64   `json:"duration"`
	Created      time.Time `json:"created"`
	ErrorMessage string    `json:"_error_message,omitempty"`
}

// WebhookQueryParameters represents URL query parameters to filter responses
type WebhookQueryParameters struct {
	ProjectID int `url:"project,omitempty"`
	WebhookID int `url:"webhook,omitempty"`
}
