package hyper2

import (
	"github.com/c0c0n3/resto/util/err"
)

// The HTTP protocol version we got isn't in the format we expect.
type UnparsableProtocolVersion string

func UnparsableProtocolVersionErr(input string) err.Err[UnparsableProtocolVersion] {
	return err.Mk[UnparsableProtocolVersion]("%s", input)
}
