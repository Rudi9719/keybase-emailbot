package main

// Email for sending
type Email struct {
	Recipients []string
	Subject    string
	Cc         []string
	Bcc        []string
	Body       string
}

// Config struct
type Config struct {
	PrivateKey string `json:"private_key"`
	KeyPass    string `json:"key_pass,omitempty"`
	MyEmail    string `json:"email"`
	EmailPass  string `json:"email_pass,omitempty"`
	SmtpServer string `json:"smtp_server"`
	AuthServer string `json:"auth_server"`
}
