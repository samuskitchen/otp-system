package enums

const (
	BasePath   string = "/otp-system"
	HealthPath string = "/health"

	GenerateOTPPOST  string = "/request-otp"
	ValidateOTPPOST  string = "/verify-otp"
	AddBlackListPOST string = "/black-list-otp"
)
