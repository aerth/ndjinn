package checkout

var (
	z CheckoutInfo
)

// Checkout stores the PayPal information
type CheckoutInfo struct {
	PayPalC string `json:"PayPalC"` // Server name
	PayPalK string `json:"PayPalK"` // Listen on HTTP
}

// Configure the Checkout Info
func Configure(c CheckoutInfo) {
	z = c
}

// ReadConfig returns the SMTP information
func ReadConfig() CheckoutInfo {
	return z
}
