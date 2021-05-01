package color

import "github.com/crawlab-team/crawlab-core/store"

func InitColor() (err error) {
	store.ColorService = NewColorService()
	return nil
}
