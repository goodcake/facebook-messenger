package messenger

// ButtonType for buttons, it can be ButtonTypeWebURL or ButtonTypePostback
type ButtonType string

//ContentType for quick replies, it can be text, location, user_phone_number or user_email
type ContentType string

// AttachmentType describes attachment type in GenericMessage
type AttachmentType string

// TemplateType of template in GenericMessage
type TemplateType string

// NotificationType for sent messages
type NotificationType string

// Message interface that represents all type of messages that we can send to Facebook Messenger
type Message interface {
	foo()
}

func (m TextMessage) foo()       {} // Message interface
func (m GenericMessage) foo()    {} // Message interface
func (m QuickReplyMessage) foo() {} // Message interface

const (
	// ButtonTypeWebURL is type for web links
	ButtonTypeWebURL = ButtonType("web_url")

	//ButtonTypePostback is type for postback buttons that sends data back to webhook
	ButtonTypePostback = ButtonType("postback")

	// AttachmentTypeTemplate for template attachments
	AttachmentTypeTemplate = AttachmentType("template")

	// TemplateTypeGeneric for generic message templates
	TemplateTypeGeneric = TemplateType("generic")

	// NotificationTypeRegular for regular notification type
	NotificationTypeRegular = NotificationType("REGULAR")

	// NotificationTypeSilentPush for silent push
	NotificationTypeSilentPush = NotificationType("SILENT_PUSH")

	// NotificationTypeNoPush for no push
	NotificationTypeNoPush = NotificationType("NO_PUSH")

	TextQuickReply = ContentType("text")

	LocationQuickReply = ContentType("location")

	PhoneNumberQuickReply = ContentType("user_phone_number")

	EmailQuickReply = ContentType("user_email")
)

// TextMessage struct used for sending text messages to messenger
type TextMessage struct {
	Message          textMessageContent `json:"message"`
	Recipient        recipient          `json:"recipient"`
	NotificationType NotificationType   `json:"notification_type,omitempty"`
}

type QuickReplyMessage struct {
	Message          quickReplyContent `json:"message"`
	Recipient        recipient         `json:"recipient"`
	NotificationType NotificationType  `json:"notification_type,omitempty"`
}

// GenericMessage struct used for sending structural messages to messenger (messages with images, links, and buttons)
type GenericMessage struct {
	Message          genericMessageContent `json:"message"`
	Recipient        recipient             `json:"recipient"`
	NotificationType NotificationType      `json:"notification_type,omitempty"`
}

type recipient struct {
	ID int64 `json:"id,string"`
}

type textMessageContent struct {
	Text string `json:"text,omitempty"`
}

type quickReplyContent struct {
	Text         string       `json:"text,omitempty"`
	QuickReplies []QuickReply `json:"quick_replies,omitempty"`
}

type genericMessageContent struct {
	Attachment *attachment `json:"attachment,omitempty"`
}

type attachment struct {
	Type    string  `json:"type,omitempty"`
	Payload payload `json:"payload,omitempty"`
}

type payload struct {
	TemplateType string    `json:"template_type,omitempty"`
	Elements     []Element `json:"elements,omitempty"`
}

// Element in Generic Message template attachment
type Element struct {
	Title    string   `json:"title"`
	Subtitle string   `json:"subtitle,omitempty"`
	ItemURL  string   `json:"item_url,omitempty"`
	ImageURL string   `json:"image_url,omitempty"`
	Buttons  []Button `json:"buttons,omitempty"`
}

type QuickReply struct {
	ContentType ContentType `json:"content_type"`
	Title       string      `json:"title,omitempty"`
	Payload     string      `json:"payload,omitempty"`
	ImageURL    string      `json:"image_url,omitempty"`
}

// Button on Generic Message template element
type Button struct {
	Type    ButtonType `json:"type"`
	URL     string     `json:"url,omitempty"`
	Title   string     `json:"title"`
	Payload string     `json:"payload,omitempty"`
}

// NewTextMessage creates new text message for userID
// This function is here for convenient reason, you will
// probably use shorthand version SentTextMessage which sends message immediatly
func (msng Messenger) NewTextMessage(userID int64, text string) TextMessage {
	return TextMessage{
		Recipient: recipient{ID: userID},
		Message:   textMessageContent{Text: text},
	}
}

func (msng Messenger) NewQuickReplyMessage(userID int64, text string) QuickReplyMessage {
	return QuickReplyMessage{
		Recipient: recipient{ID: userID},
		Message:   quickReplyContent{Text: text},
	}
}

// NewGenericMessage creates new Generic Template message for userID
// Generic template messages are used for structured messages with images, links, buttons and postbacks
func (msng Messenger) NewGenericMessage(userID int64) GenericMessage {
	return GenericMessage{
		Recipient: recipient{ID: userID},
		Message: genericMessageContent{
			Attachment: &attachment{
				Type:    "template",
				Payload: payload{TemplateType: "generic"},
			},
		},
	}
}

// AddNewElement adds element to Generic template message with defined title, subtitle, link url and image url
// Title param is mandatory. If not used set "" for other params and nil for buttons param
// Generic messages can have up to 10 elements which are scolled horizontaly in Facebook messenger
func (m *GenericMessage) AddNewElement(title, subtitle, itemURL, imageURL string, buttons []Button) {
	m.AddElement(newElement(title, subtitle, itemURL, imageURL, buttons))
}

// AddElement adds element e to Generic Message
// Generic messages can have up to 10 elements which are scolled horizontaly in Facebook messenger
func (m *GenericMessage) AddElement(e Element) {
	m.Message.Attachment.Payload.Elements = append(m.Message.Attachment.Payload.Elements, e)
}

// NewElement creates new element with defined title, subtitle, link url and image url
// Title param is mandatory. If not used set "" for other params and nil for buttons param
// Instead of calling this function you can also initialize Element struct, depends what you prefere
func (msng Messenger) NewElement(title, subtitle, itemURL, imageURL string, buttons []Button) Element {
	return newElement(title, subtitle, itemURL, imageURL, buttons)
}

func newElement(title, subtitle, itemURL, imageURL string, buttons []Button) Element {
	return Element{
		Title:    title,
		Subtitle: subtitle,
		ItemURL:  itemURL,
		ImageURL: imageURL,
		Buttons:  buttons,
	}
}

//AddNewQuickReply add quick reply to Quick Reply Message template
func (m *QuickReplyMessage) AddNewQuickReply(contentType ContentType, title, payload, imgURL string) {
	m.AddQuickReply(newQuickReply(contentType, title, payload, imgURL))
}

func (m *QuickReplyMessage) AddQuickReply(qr QuickReply) {
	m.Message.QuickReplies = append(m.Message.QuickReplies, qr)
}

func (msng Messenger) NewQuickReply(contentType ContentType, title, payload, imgURL string) QuickReply {
	return newQuickReply(contentType, title, payload, imgURL)
}

func newQuickReply(contentType ContentType, title, payload, imgURL string) QuickReply {
	return QuickReply{
		ContentType: contentType,
		Title:       title,
		Payload:     payload,
		ImageURL:    imgURL,
	}
}

// NewWebURLButton creates new web url button
func (msng Messenger) NewWebURLButton(title, URL string) Button {
	return Button{
		Type:  ButtonTypeWebURL,
		Title: title,
		URL:   URL,
	}
}

// NewPostbackButton creates new postback button that sends payload string back to webhook when pressed
func (msng Messenger) NewPostbackButton(title, payload string) Button {
	return Button{
		Type:    ButtonTypePostback,
		Title:   title,
		Payload: payload,
	}
}

// AddWebURLButton creates and adds web link URL button to the element
func (e *Element) AddWebURLButton(title, URL string) {
	b := Button{
		Type:  ButtonTypeWebURL,
		Title: title,
		URL:   URL,
	}
	e.Buttons = append(e.Buttons, b)
}

// AddPostbackButton creates and adds button that sends payload string back to webhook when pressed
func (e *Element) AddPostbackButton(title, payload string) {
	b := Button{
		Type:    ButtonTypePostback,
		Title:   title,
		Payload: payload,
	}
	e.Buttons = append(e.Buttons, b)
}
