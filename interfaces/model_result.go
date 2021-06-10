package interfaces

type Result interface {
	Value() map[string]interface{}
	SetValue(key string, value interface{})
	GetValue(key string) (value interface{})
}
