package admin

var (
	c AdminInfo
)

// Checkout stores the PayPal information
type AdminInfo struct {
	Email    string `json:"Email"`    // Server name
	Mkpasswd string `json:"Mkpasswd"` // Listen on HTTP
}

// Configure the Checkout Info
func Configure(a AdminInfo) {
	a = c
}

// ReadConfig returns the SMTP information
func ReadConfig() AdminInfo {
	return c
}
