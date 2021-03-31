package controllers

func NewPostActionControllerDelegate(id ControllerId, actions []PostAction) (d *PostActionControllerDelegate) {
	return &PostActionControllerDelegate{
		id:      id,
		actions: actions,
	}
}

type PostActionControllerDelegate struct {
	id      ControllerId
	actions []PostAction
}

func (ctr *PostActionControllerDelegate) Actions() (actions []PostAction) {
	return ctr.actions
}
