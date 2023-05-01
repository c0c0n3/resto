package msg

import "github.com/c0c0n3/resto/yoorel"

// HTTP message control data.
// See:
// - https://www.rfc-editor.org/rfc/rfc9110.html#name-message-abstraction

// An HTTP message starts with control data to tell whether it's a
// request or a response. Request control data includes a request
// method, request target, and protocol version.
type RequestControlData struct {
	Method  Method
	Target  yoorel.HttpUrl
	Version ProtocolVersion
}

// An HTTP message starts with control data to tell whether it's a
// request or a response. Response control data includes a a status
// code, optional reason phrase, and protocol version.
type ResponseControlData struct {
	Status  StatusCode
	Reason  string
	Version ProtocolVersion
}
