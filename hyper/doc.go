/*

Package hyper and its sub-packages let you exchange HTTP messages.

Package hyper.wire defines the interfaces to read and write HTTP
messages as well as the interfaces to exchange them. Plus, it comes
with a default implementation backed by the standard lib's "net/http"
package. Notice this package is low-level stuff you probably don't
want to use, but it lets us easily build EDSLs on top of it.

In fact, that's exactly what the hyper.client package does. It uses
hyper.wire's prim ops to implement a high-level EDSL to make HTTP
requests which is independent of the underlying client lib you can
use. It works with "net/http" but can also happily work with any
other library that has common, run-of-the-mill HTTP functionality.

*/
package hyper
