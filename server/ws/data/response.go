package data

type ResponseMethod string

const (
	ResponseNewMessage      = ResponseMethod("new_message")
	ResponseNewFriend       = ResponseMethod("new_friend")
	ResponseNewConversation = ResponseMethod("new_conversation")
	ResponseError           = ResponseMethod("error")
)

type Response struct {
	Method ResponseMethod `json:"method"`
	Data   interface{}    `json:"data"`
}
