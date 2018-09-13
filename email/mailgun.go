package email

import (
	"fmt"
	"net/url"

	mailgun "gopkg.in/mailgun/mailgun-go.v1"
)

const (
	resetBaseURL   = "http://localhost:3000/api/update_password" // Change this for production
	welcomeSubject = "Welcome!"
	resetSubject   = "Instructions for resetting your password."
)

const welcomeText = `
	Hi there!

	Welcome! We really hope you enjoy using our application!

	Best,
	Yuichi
`

const welcomeHTML = `
	Hi there!<br/>
	<br/>
	Welcome to <a href="https://www.example.com">Example</a>! We really hope you enjoy using our application!<br/>
	<br/>
	Best,<br/>
	Yuichi
`

const resetTextTmpl = `
	Hi there!

	It appears that you have requested a password reset. If this was you, please follow the link below to update your password:

	%s

	If you are asked for a token, please use the following value:

	%s

	If you didn't request a password reset you can safely ignore this email and your account will not be changed.

	Best,
	Support
`

const resetHTMLTmpl = `
	Hi there!<br/>
	<br/>
	It appears that you have requested a password reset. If this was you, please follow the link below to update your password:<br/>
	<br/>
	<a href="%s">%s</a><br/>
	<br/>
	If you are asked for a token, please use the following value:<br/>
	<br/>
	%s<br/>
	<br/>
	If you didn't request a password reset you can safely ignore this email and your account will not be changed.<br/>
	<br/>
	Best,<br/>
	Support<br/>
`

func WithMailgun(domain, apiKey, publicKey string) ClientConfig {
	return func(client *Client) {
		mg := mailgun.NewMailgun(domain, apiKey, publicKey)
		client.mg = mg
	}
}

func WithSender(username, email string) ClientConfig {
	return func(client *Client) {
		client.from = buildEmail(username, email)
	}
}

type ClientConfig func(*Client)

func NewClient(opts ...ClientConfig) *Client {
	client := Client{
		// Set a default from email address...
		from: "support@example.com",
	}
	for _, opt := range opts {
		opt(&client)
	}
	return &client
}

type Client struct {
	from string
	mg   mailgun.Mailgun
}

func (client *Client) Welcome(toUsername, toEmail string) error {
	message := mailgun.NewMessage(client.from, welcomeSubject, welcomeText, buildEmail(toUsername, toEmail))
	message.SetHtml(welcomeHTML)
	_, _, err := client.mg.Send(message)
	if err != nil {
		return err
	}
	return err
}

func buildEmail(username, email string) string {
	if username == "" {
		return email
	}
	return fmt.Sprintf("%s <%s>", username, email)
}

func (client *Client) ResetPassword(toEmail, token string) error {
	v := url.Values{}
	v.Set("token", token)
	resetUrl := resetBaseURL + "?" + v.Encode()
	resetText := fmt.Sprintf(resetTextTmpl, resetUrl, token)
	message := mailgun.NewMessage(client.from, resetSubject, resetText, toEmail)
	resetHTML := fmt.Sprintf(resetHTMLTmpl, resetUrl, resetUrl, token)
	message.SetHtml(resetHTML)
	_, _, err := client.mg.Send(message)
	return err
}
