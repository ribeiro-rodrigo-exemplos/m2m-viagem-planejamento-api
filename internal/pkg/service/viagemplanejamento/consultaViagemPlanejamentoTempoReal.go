package viagemplanejamento

import (
	"time"

	"git.m2mfacil.com.br/golang/m2m-viagem-planejamento-api/internal/pkg/dto"
	"git.m2mfacil.com.br/golang/m2m-viagem-planejamento-api/internal/pkg/util"
)

//Consultar -
func (vps *Service) Consultar(filtro dto.FilterDTO) (consulta *dto.ConsultaViagemPlanejamentoDTO, err error) {

	var tempoReal bool
	chConsulta := make(chan *dto.ConsultaViagemPlanejamentoDTO)
	chErr := make(chan error)

	if filtro.TempoRealInicio != "" {
		filtro.AtualizarParaTempoReal(time.Now(), vps.cacheCliente.Cache[filtro.IDCliente].Location)
	}

	go func(f dto.FilterDTO) {
		c, err := vps.ConsultarPeriodo(filtro)
		chConsulta <- c
		chErr <- err
	}(filtro)

	var c *dto.ConsultaViagemPlanejamentoDTO
	var crt *dto.ConsultaViagemPlanejamentoDTO

	if filtro.TempoRealInicio != "" {
		tempoReal = true

		filtroRT := &filtro
		filtroRT.Complemento.Cliente = vps.cacheCliente.Cache[filtro.IDCliente]

		dtIni := *filtroRT.GetDataInicio()
		dtFim := *filtro.GetDataInicio()

		dtIni = dtIni.Add(-1 * (72 * time.Hour))
		filtroRT.DataInicio = util.FormatarAMDHMS(&dtIni)

		dtFim = dtFim.Add(-1 * (1 * time.Second))
		filtroRT.DataFim = util.FormatarAMDHMS(&dtFim)

		filtroRT.Complemento.ApenasViagemExecutada = true

		crt, err = vps.serviceRealTime.ConsultarPeriodo(*filtroRT)
	}

	select {
	case c = <-chConsulta:
	case err = <-chErr:
	}

	if tempoReal {
		var viagens = []*dto.ViagemDTO{}

		for _, vg := range crt.Viagens {
			if *vg.Status == dto.StatusViagem.EmAndamento {
				viagens = append(viagens, vg)
			}
		}

		crt.Viagens = viagens
		crt.Viagens = append(crt.Viagens, c.Viagens...)
		crt.Totalizadores = c.Totalizadores
		crt.Informacoes = c.Informacoes
		consulta = crt
	} else {
		consulta = c
	}

	return
}
