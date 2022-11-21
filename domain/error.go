package domain

import "fmt"

type Error struct {
	Category string
	Status   string
	Title    interface{}
}

func NewError(category string, status string, msg interface{}) Error {
	return Error{
		Category: category,
		Status:   status,
		Title:    msg,
	}
}

func (e Error) Error() string {
	return fmt.Sprintf("%s: %v", e.Status, e.Title)
}

func (e Error) SetMessage(msg interface{}) Error {
	e.Title = msg
	return e
}

const (
	BAD_REQUEST           string = "bad-request"
	UNAUTHORIZED          string = "unauthorized"
	FORBIDDEN             string = "forbidden"
	NOT_FOUND             string = "not-found"
	CONFLICT              string = "conflict"
	INTERNAL_SERVER_ERROR string = "internal-server-error"
	UNKNOWN               string = "unknown"
)

var (
	ErrorUnknown             Error = NewError(UNKNOWN, "unknown", "")
	ErrorInternalServerError Error = NewError(INTERNAL_SERVER_ERROR, "internal-server-error", "")
	ErrorForbidden           Error = NewError(FORBIDDEN, "forbidden", "")
	ErrorBadRequest          Error = NewError(BAD_REQUEST, "bad-request", "")
	ErrorValidationFailed    Error = NewError(BAD_REQUEST, "validation-failed", "")
	ErrorBindStructFailed    Error = NewError(BAD_REQUEST, "bind-struct-failed", "")
	ErrorInvalidUUID         Error = NewError(BAD_REQUEST, "invalid-uuid", "")
	ErrorExpired             Error = NewError(UNAUTHORIZED, "expired", "")
	ErrorSignInTokenFailed   Error = NewError(UNAUTHORIZED, "sign-in-token-failed", "")
)

// scb_enet
var (
	ErrorTransactionNotFound Error = NewError(NOT_FOUND, "transaction/not-found", "")
)

type ErrorResponse struct {
	Title       interface{} `json:"title"`
	Status      int         `json:"status"`
	Description string      `json:"description"`
	Result      interface{} `json:"result"`
}

const (
	Cache             string = "รายการแคช"
	SignInTokenFailed string = "ชื่อผู้ใช้งาน หรือ รหัสผ่านไม่ถูกต้อง"
)
const (
	SignIn            string = "เข้าสู่ระบบ"
	GetAccountBalance string = "ยอดเงินในบัญชี"
	GetTransaction    string = "รายการเดินบัญชี"
)

const (
	Success  string = "สำเร็จ"
	Failed   string = "ไม่สำเร็จ"
	Close    string = "ไม่สามารถให้บริการได้ในขณะนี้"
	Suspend  string = "ถูกระงับชั่วคราว"
	NotFound string = "ไม่พบข้อมูล"
)
