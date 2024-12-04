package cidbot

type ResponseData[T any] struct {
	Data []T `json:"data"`
}