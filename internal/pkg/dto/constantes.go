package dto

//StatusViagem -
var StatusViagem = StatusViagemConst{
	Todos:                 0,
	RealizadaPlanejada:    1,
	NaoRealizada:          2,
	RealizadaNaoPlanejada: 3,
	NaoIniciada:           4,
	EmAndamento:           5,
	Cancelada:             6,
	Atrasada:              7,
	Extra:                 8,
}

//StatusViagemConst -
type StatusViagemConst struct {
	Todos                 int
	RealizadaPlanejada    int
	NaoRealizada          int
	RealizadaNaoPlanejada int
	NaoIniciada           int
	EmAndamento           int
	Cancelada             int
	Atrasada              int
	Extra                 int
}

//OrigemMensagem -
var OrigemMensagem = OrigemMensagemConst{
	Planejada: 0,
	Executada: 1,
}

//OrigemMensagemConst -
type OrigemMensagemConst struct {
	Planejada int
	Executada int
}
