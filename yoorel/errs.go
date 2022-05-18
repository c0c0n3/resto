package yoorel

// A nil pointer error.
type NilPtr string

// An error for an invalid HTTP hostname.
type InvalidHostname string

// An error for an HTTP port which is out of range.
type InvalidPort string

// An error output by URL builders or parsers if the resulting URL isn't
// valid.
type InvalidUrl string
