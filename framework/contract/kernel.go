package contract

import "net/http"

const KernelKey = "eehe:kernel"

// Kernel provides the core struct of framework.
type Kernel interface {
	// HttpEngine provides gin engine.
	HttpEngine() http.Handler
}
