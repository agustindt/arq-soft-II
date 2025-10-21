package errorspkg

import "errors"

var (
	ErrNotFound       = errors.New("No Encontrado")
	ErrOwnerNotValid  = errors.New("Usuario no Valido")
	ErrConcurrentCalc = errors.New("Fallo en concurrencia")
)
