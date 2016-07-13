package fbmessenger

/*------------------------------------------------------
Send API
------------------------------------------------------*/

// Text message is a fluent helper method for creating a SendRequest containing a text message.
func TextMessage(text string) *SendRequest {
	return &SendRequest{
		Message: Message{
			Text: text,
		},
	}
}

// ImageMessage is a fluent helper method for creating a SendRequest containing a message with
// an image attachment that has a MediaPayload.
func ImageMessage(url string) *SendRequest {
	return &SendRequest{
		Message: Message{
			Attachment: &Attachment{
				Type: "image",
				Payload: MediaPayload{
					Url: url,
				},
			},
		},
	}
}

// ButtonTemplateMessage is a fluent helper method for creating a SendRequest containing text
// and buttons to request input from the user.
func ButtonTemplateMessage(text string, buttons ...*Button) *SendRequest {
	return &SendRequest{
		Message: Message{
			Attachment: &Attachment{
				Type: "template",
				Payload: ButtonPayload{
					TemplateType: "button",
					Text:         text,
					Buttons:      buttons,
				},
			},
		},
	}
}

// To is a fluent helper method for setting Recipient. It is a mutator
// and returns the same SendRequest on which it is called to support method chaining.
func (sr *SendRequest) To(userId string) *SendRequest {
	sr.Recipient = Recipient{Id: userId}

	return sr
}

// ToPhoneNumber is a fluent helper method for setting Recipient. It
// is a mutator and returns the same SendRequest on which it is called to support method chaining.
func (sr *SendRequest) ToPhoneNumber(phoneNumber string) *SendRequest {
	sr.Recipient = Recipient{PhoneNumber: phoneNumber}
	return sr
}

// Regular is a fluent helper method for setting NotificationType. It is a mutator and
// returns the same SendRequest on which it is called to support method chaining.
func (sr *SendRequest) Regular() *SendRequest {
	sr.NotificationType = "REGULAR"

	return sr
}

// SilentPush is a fluent helper method for setting NotificationType. It is a mutator and
// returns the same SendRequest on which it is called to support method chaining.
func (sr *SendRequest) SilentPush() *SendRequest {
	sr.NotificationType = "SILENT_PUSH"

	return sr
}

// NoPush is a fluent helper method for setting NotificationType. It is a mutator and
// returns the same SendRequest on which it is called to support method chaining.
func (sr *SendRequest) NoPush() *SendRequest {
	sr.NotificationType = "NO_PUSH"

	return sr
}

/*
SendRequest is the top level structure for representing any type of message to send.

See https://developers.facebook.com/docs/messenger-platform/send-api-reference#request
*/
type SendRequest struct {
	Recipient        Recipient `json:"recipient" binding:"required"`
	Message          Message   `json:"message" binding:"required"`
	NotificationType string    `json:"notification_type,omitempty"`
}

// Recipient identifies the user to send to. Either Id or PhoneNumber must be set, but not both.
type Recipient struct {
	Id          string `json:"id,omitempty"`
	PhoneNumber string `json:"phone_number,omitempty"`
}

// Message can represent either a text message, or a message with an attachment. Either
// Text or Attachment must be set, but not both.
type Message struct {
	Text       string      `json:"text,omitempty"`
	Attachment *Attachment `json:"attachment,omitempty"`
}

// Attachment is used to build a message with attached media, or a structured message.
type Attachment struct {
	Type    string      `json:"type" binding:"required"`
	Payload interface{} `json:"payload" binding:"required"`
}

/*
MediaPayload is used to hold the URL of media attached to a message.

See https://developers.facebook.com/docs/messenger-platform/send-api-reference/image-attachment
*/
type MediaPayload struct {
	Url string `json:"url" binding:"required"`
}

/*
ButtonPayload is used to build a structured message using the button template.

See https://developers.facebook.com/docs/messenger-platform/send-api-reference/button-template
*/
type ButtonPayload struct {
	TemplateType string    `json:"template_type" binding:"required"`
	Text         string    `json:"text" binding:"required"`
	Buttons      []*Button `json:"buttons" binding:"required"`
}

// Button represents a single button in a structured message using the button template.
type Button struct {
	Type    string `json:"type" binding:"required"`
	Title   string `json:"title" binding:"required"`
	Url     string `json:"url,omitempty"`
	Payload string `json:"payload,omitempty"`
}

/*
SendResponse is returned when sending a SendRequest.

See https://developers.facebook.com/docs/messenger-platform/send-api-reference#response
*/
type SendResponse struct {
	RecipientId string     `json:"recipient_id" binding:"required"`
	MessageId   string     `json:"message_id" binding:"required"`
	Error       *SendError `json:"error"`
}

/*
SendError indicates an error returned from Facebook.

See https://developers.facebook.com/docs/messenger-platform/send-api-reference#errors
*/
type SendError struct {
	Message   string `json:"message" binding:"required"`
	Type      string `json:"type" binding:"required"`
	Code      int    `json:"code" binding:"required"`
	ErrorData string `json:"error_data" binding:"required"`
	FBTraceId string `json:"fbtrace_id" binding:"required"`
}

/*------------------------------------------------------
Webhook
------------------------------------------------------*/

/*
Callback is the top level structure that represents a callback received by your
webhook endpoint.

See https://developers.facebook.com/docs/messenger-platform/webhook-reference#format
*/
type Callback struct {
	Object  string   `json:"object" binding:"required"`
	Entries []*Entry `json:"entry" binding:"required"`
}

// Entry is part of the common format of callbacks.
type Entry struct {
	PageId    string            `json:"id" binding:"required"`
	Time      int               `json:"time" binding:"required"`
	Messaging []*MessagingEntry `json:"messaging"`
}

/*
MessagingEntry is an individual interaction a user has with a page.
The Sender and Recipient fields are common to all types of callbacks and the
other fields only apply to specific types of callbacks.
*/
type MessagingEntry struct {
	Sender    Principal        `json:"sender" binding:"required"`
	Recipient Principal        `json:"recipient" binding:"required"`
	Timestamp int              `json:"timestamp"`
	Message   *CallbackMessage `json:"message"`
	Delivery  *Delivery        `json:"delivery"`
	Postback  *Postback        `json:"postback"`
	OptIn     *OptIn           `json:"optin"`
}

// Principal holds the Id of a sender or recipient.
type Principal struct {
	Id string `json:"id" binding:"required"`
}

/*
CallbackMessage represents a message a user has sent to your page.
Either the Text or Attachments field will be set, but not both.

See https://developers.facebook.com/docs/messenger-platform/webhook-reference/message-received
*/
type CallbackMessage struct {
	MessageId   string                `json:"mid" binding:"required"`
	Sequence    int                   `json:"seq" binding:"required"`
	Text        string                `json:"text"`
	Attachments []*CallbackAttachment `json:"attachments"`
}

// CallbackAttachment holds the type and payload of an attachment sent by a user.
type CallbackAttachment struct {
	Type    string                    `json:"type" binding:"required"`
	Payload CallbackAttachmentPayload `json:"payload" binding:"required"`
}

// CallbackAttachmentPayload holds the URL of an attachment sent by the user.
type CallbackAttachmentPayload struct {
	Url string `json:"url" binding:"required"`
}

/*
Delivery holds information about which of the messages that you've sent have been delivered.

See https://developers.facebook.com/docs/messenger-platform/webhook-reference/message-delivered
*/
type Delivery struct {
	MessageIds []string `json:"mids"`
	Watermark  int      `json:"watermark" binding:"required"`
	Sequence   int      `json:"seq" bindging:"required"`
}

/*
Postback holds the data defined for buttons the user taps.

See https://developers.facebook.com/docs/messenger-platform/webhook-reference/postback-received
*/
type Postback struct {
	Payload string `json:"payload" binding:"required"`
}

/*
OptIn holds the data defined for the Send-to-Messenger plugin.

See https://developers.facebook.com/docs/messenger-platform/webhook-reference/authentication
*/
type OptIn struct {
	Ref string `json:"ref" binding:"required"`
}

/*------------------------------------------------------
User Profile
------------------------------------------------------*/

/*
UserProfile represents additional information about the user.

See https://developers.facebook.com/docs/messenger-platform/user-profile
*/
type UserProfile struct {
	FirstName       string `json:"first_name"`
	LastName        string `json:"last_name"`
	ProfilePhotoUrl string `json:"profile_pic"`
	Locale          string `json:"locale"`
	Timezone        int    `json:"timezone"`
	Gender          string `json:"gender"`
}
