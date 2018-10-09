package viagemplanejamento

import (
	"fmt"
	"time"

	"git.m2mfacil.com.br/golang/m2m-viagem-planejamento-api/internal/pkg/dto"
)

//ConsultarDashboard -
func (vps *Service) ConsultarDashboard(filtro dto.FilterDTO) (*dto.ConsultaViagemPlanejamentoDTO, error) {
	var consulta *dto.ConsultaViagemPlanejamentoDTO
	var err error

	consulta, err = vps.ConsultarPeriodo(filtro)

	start := time.Now()
	t := &start
	f := &filtro

	for _, vg := range consulta.Viagens {
		applyAlert(vg, f, t)
		filtraExecucaoViagem(vg)
		filtraExecucaoViagemToolTip(vg)
		filtraStatusViagem(vg)
	}

	duracao := time.Since(start)

	informacoes := vps.consultaViagemPlanejamento.Informacoes
	if informacoes != nil {
		informacoes["duracaoConsultarDashboard"] = fmt.Sprintf("%v", duracao)
	}

	return consulta, err
}

func applyAlert(vg *dto.ViagemDTO, filtro *dto.FilterDTO, t *time.Time) {
	var classeAlerta string

	if vg.PartidaReal != nil {
		return
	}
	var tolerancia int32
	tolerancia = vg.Tolerancia.AtrasoPartida

	duration := t.Sub(*vg.Data)
	diff := int64(duration / time.Minute)

	if diff == 0 {
		classeAlerta = "m2m-green-alert"
	} else if diff > 0 && diff <= (int64(tolerancia)) {
		classeAlerta = "red-alert"
	}

	if classeAlerta != "" {
		vg.Apresentacao.ClasseAlerta = &classeAlerta
	}

}

func filtraExecucaoViagem(vg *dto.ViagemDTO) {
	var classeExecucaoViagem string

	if *vg.Status == 2 || *vg.Status == 4 {
		classeExecucaoViagem = "black-according-to-theme"

	} else if *vg.Status == 6 {
		classeExecucaoViagem = "red"

	} else if vg.ChegadaReal != nil {
		classeExecucaoViagem = "green"

	} else if vg.EmExecucao {
		classeExecucaoViagem = "blue"
	}
	vg.Apresentacao.ClasseExecucaoViagem = classeExecucaoViagem
}

func filtraExecucaoViagemToolTip(vg *dto.ViagemDTO) {
	var classeExecucaoViagemToolTip string

	if *vg.Status == 2 {
		classeExecucaoViagemToolTip = "Não Realizada"
	} else if *vg.Status == 6 {
		classeExecucaoViagemToolTip = "Viagem Cancelada"
	} else if vg.ChegadaReal != nil {
		classeExecucaoViagemToolTip = "Viagem fechada"
	} else if vg.EmExecucao {
		classeExecucaoViagemToolTip = "Em andamento"
	} else if *vg.Status == 4 {
		classeExecucaoViagemToolTip = "Não Iniciada"
	}
	vg.Apresentacao.ClasseExecucaoViagemToolTip = classeExecucaoViagemToolTip
}

func filtraStatusViagem(vg *dto.ViagemDTO) {
	var classeStatusViagem string

	if *vg.Status == 1 || vg.Apresentacao.AlertaProximo || (vg.EmExecucao && vg.PartidaPlan != nil) {
		classeStatusViagem = "m2m-green-bg min-width"
	} else if *vg.Status == 2 {
		classeStatusViagem = "white red-bg min-width"
	} else if *vg.Status == 4 && !vg.Apresentacao.AlertaProximo {
		classeStatusViagem = "white black-bg min-width"
	} else if *vg.Status == 6 {
		classeStatusViagem = "white purple-bg min-width"
	} else if *vg.Status == 7 {
		classeStatusViagem = "white yellow-bg min-width"
	} else if *vg.Status == 8 || (vg.EmExecucao && vg.PartidaPlan == nil) {
		classeStatusViagem = "white blue-bg min-width"
	}
	vg.Apresentacao.ClasseStatusViagem = classeStatusViagem
}
