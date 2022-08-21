package wire

// A function to write the core parts of an HTTP message.
// The implementation uses the given MessageWriter to do that and returns
// an error if something goes wrong.
type MessageBuilder func(msg MessageWriter) error

// A function to write an HTTP request.
// The implementation uses the given RequestWriter to do that and returns
// an error if something goes wrong.
type RequestBuilder func(req RequestWriter) error

// A function to handle an HTTP response.
// The implementation reads the server's response from the given ResponseReader
// and returns an error if something goes wrong.
type ResponseHandler func(res ResponseReader) error

// A function to handle an HTTP client request.
// The implementation reads the request from the given RequestReader and
// replies to the client by using the given ResponseWriter. The implementation
// returns an error if something goes wrong.
type RequestHandler func(req RequestReader, res ResponseWriter) error

// A function to send an HTTP request to a server.
// The implementation calls the given RequestBuilder to write the request
// and returns a ResponseReader to read the server's response or an error
// if something goes wrong.
type Sender func(RequestBuilder) (ResponseReader, error)
