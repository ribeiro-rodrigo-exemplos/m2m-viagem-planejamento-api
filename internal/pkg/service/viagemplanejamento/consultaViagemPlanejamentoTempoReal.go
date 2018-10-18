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

	f1 := filtro.Clone()
	f2 := filtro.Clone()

	if filtro.TempoRealInicio != "" {
		f1.AtualizarParaTempoReal(time.Now(), vps.cacheCliente.Cache[filtro.IDCliente].Location)
		f2.AtualizarParaTempoReal(time.Now(), vps.cacheCliente.Cache[filtro.IDCliente].Location)
	}

	go func(f *dto.FilterDTO) {
		f.Complemento.Instancia = "Padrão"
		c, err := vps.ConsultarPeriodo(*f)
		chConsulta <- c
		chErr <- err
	}(&f1)

	var c *dto.ConsultaViagemPlanejamentoDTO
	var crt *dto.ConsultaViagemPlanejamentoDTO

	if f2.TempoRealInicio != "" {
		tempoReal = true

		filtroRT := &f2
		filtroRT.Complemento.Cliente = vps.cacheCliente.Cache[f2.IDCliente]

		dtIni := *filtroRT.GetDataInicio()
		dtFim := *f2.GetDataInicio()

		dtIni = dtIni.Add(-1 * (72 * time.Hour))
		filtroRT.DataInicio = util.FormatarAMDHMS(&dtIni)

		dtFim = dtFim.Add(-1 * (1 * time.Second))
		filtroRT.DataFim = util.FormatarAMDHMS(&dtFim)

		filtroRT.Complemento.ApenasViagemExecutada = true
		filtroRT.Complemento.Instancia = "RealTime"

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
