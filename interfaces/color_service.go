package interfaces

type ColorService interface {
	GetByName(name string) (res Color, err error)
	GetRandom() (res Color, err error)
}
