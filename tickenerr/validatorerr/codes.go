package validatorerr

const (
	ValidatorNotFoundErrorCode = iota + 100
)

func GetErrMessage(code uint32) string {
	switch code {
	case ValidatorNotFoundErrorCode:
		return "validator not found"
	default:
		return "an error has occurred"
	}
}
