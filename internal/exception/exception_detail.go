package exception

type ErrorDetail struct {
	StatusCode int    `json:"statusCode"`
	Message    string `json:"message"`
}

func buildErrorDetail(m string, status int) ErrorDetail {
	return ErrorDetail{StatusCode: status, Message: m}
}

const (
	// commons
	UnknownError   = "unknown_error"
	DbCommandError = "database_error"

	// events
	ErrorPublishingEvent      = "event.publish_error"
	ErrorConvertMessageToJson = "event.convert_message_to_json"
	ErrorInstanceEventBus     = "event.instance_event_bus"

	// account
	AccountEmailAlreadyExists = "account.email_already_exists"
	AccountPhoneAlreadyExists = "account.phone_already_exists"
	AccountNotFound           = "account.not_found"
	InvalidUserOrPassword     = "account.invalid_user_or_password"

	// webhooks
	WebhookTransactionNotFound            = "webhook_transaction.transaction_not_found"
	WebhookTransactionNotificationTimeout = "webhook_transaction.notification_timeout"

	// auth
	ErrorParseToken   = "auth.error_parse_token"
	InvalidToken      = "auth.invalid_token"
	ErrorConvertToken = "auth.convert_token"

	// password hash
	FailureToCreateHash = "hash.create_password"
	FailureToComparHash = "hash.failure_to_compare"

	// webhook subscription
	WebhookSubscriptionNotFound = "webhook_subscription.not_found"
)

var messages = map[string]ErrorDetail{
	// commons
	UnknownError:   buildErrorDetail("An error occurred, please try again later", 500),
	DbCommandError: buildErrorDetail("An error occurred while processing the request", 500),

	// account
	AccountEmailAlreadyExists: buildErrorDetail("The email is already in use", 400),
	AccountPhoneAlreadyExists: buildErrorDetail("The phone is already in use", 400),
	AccountNotFound:           buildErrorDetail("Account not found", 404),
	InvalidUserOrPassword:     buildErrorDetail("Invalid user or password", 400),

	// webhooks
	WebhookTransactionNotFound:            buildErrorDetail("Transaction not found", 404),
	WebhookTransactionNotificationTimeout: buildErrorDetail("Webhook notification timeout", 408),

	// events
	ErrorPublishingEvent:      buildErrorDetail("Error publishing event", 500),
	ErrorConvertMessageToJson: buildErrorDetail("Error converting message to json", 500),
	ErrorInstanceEventBus:     buildErrorDetail("Error creating instance of event bus", 500),

	// auth
	ErrorParseToken:   buildErrorDetail("Error parsing token", 401),
	InvalidToken:      buildErrorDetail("Invalid token", 401),
	ErrorConvertToken: buildErrorDetail("Error converting token", 500),

	// password hash
	FailureToCreateHash: buildErrorDetail("Failure to create hash", 500),
	FailureToComparHash: buildErrorDetail("Failure to compare hash", 500),

	// webhook subscription
	WebhookSubscriptionNotFound: buildErrorDetail("Webhook subscription not found", 404),
}

func GetMessageByCode(code string) (string, int) {
	m, ok := messages[code]
	if !ok {
		m = messages[UnknownError]
		return m.Message, m.StatusCode
	}
	return m.Message, m.StatusCode
}
