package models

type variableService struct {
	*Service
}

var VariableService = variableService{NewService(ModelIdVariable)}
