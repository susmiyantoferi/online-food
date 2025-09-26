package handling

import (
	"errors"
	"net/http"
	"online-food/utils/response"

	"github.com/gin-gonic/gin"
)

var (
	ErrorIdNotFound    = errors.New("id not found")
	ErrorEmailNotFound = errors.New("email not found")
	ErrorEmailExist    = errors.New("email already exist")
	ErrNotEnoughStock  = errors.New("not enough stock")
	ErrorValidation    = errors.New("validation failed")
	ErrFailedLogin     = errors.New("email or password wrong")
	ErrInvalidToken    = errors.New("invalid token refresh")
	ErrEmptyItems      = errors.New("cart has no items")
	ErrMenuNotFound    = errors.New("menu not found")
	ErrCheckoutCart    = errors.New("cart already checkout")
)

var errorMapping = map[error]struct {
	Code    int
	Status  string
	Message string
	Data    interface{}
}{
	ErrorEmailExist:    {http.StatusConflict, "Conflict", "email already exists", nil},
	ErrorValidation:    {http.StatusBadRequest, "Bad Request", "invalid input", nil},
	ErrNotEnoughStock:  {http.StatusBadRequest, "Bad Request", "not enough stock", nil},
	ErrFailedLogin:     {http.StatusBadRequest, "Bad Request", "email or password wrong", nil},
	ErrInvalidToken:    {http.StatusBadRequest, "Bad Request", "invalid token refresh", nil},
	ErrorEmailNotFound: {http.StatusNotFound, "Not Found", "email not found", nil},
	ErrorIdNotFound:    {http.StatusNotFound, "Not Found", "id not found", nil},
	ErrMenuNotFound:    {http.StatusNotFound, "Not Found", "menu not found", nil},
	ErrEmptyItems:      {http.StatusBadRequest, "Bad Request", "cart has no items", nil},
	ErrCheckoutCart:    {http.StatusBadRequest, "Bad Request", "cart already checkout", nil},
}

func HandleError(ctx *gin.Context, err error) {
	for key, v := range errorMapping {
		if errors.Is(err, key) {
			response.ToResponseJson(ctx, v.Code, v.Status, v.Message, v.Data)
			return
		}
	}

	response.ToResponseJson(ctx, http.StatusInternalServerError, "Internal Server Error", "internal server error", nil)
}
