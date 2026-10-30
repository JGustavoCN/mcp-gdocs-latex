package gdocs

import (
	"testing"
)

func TestRegexHeadingsDetection(t *testing.T) {
	primaryCases := []string{
		"1 INTRODUÇÃO",
		"2 REFERENCIAL TEÓRICO",
		"3 METODOLOGIA",
		"4 RESULTADOS E DISCUSSÃO",
		"5 CONCLUSÃO",
		"REFERÊNCIAS",
		"RESUMO",
		"ABSTRACT",
	}

	for _, c := range primaryCases {
		if !regexPrimaryHeading.MatchString(c) {
			t.Errorf("Esperava que %q fosse reconhecido como título primário", c)
		}
	}

	secondaryCases := []string{
		"1.1 Contextualização",
		"1.2 Problema de Pesquisa",
		"2.1 Model Context Protocol",
		"3.1 Delineamento Experimental",
		"1.1 CONTEXTUALIZAÇÃO",
	}

	for _, c := range secondaryCases {
		if !regexSecondaryHeading.MatchString(c) {
			t.Errorf("Esperava que %q fosse reconhecido como título secundário", c)
		}
	}

	tertiaryCases := []string{
		"1.1.1 Escopo do Projeto",
		"2.1.1 Arquitetura Híbrida",
	}

	for _, c := range tertiaryCases {
		if !regexTertiaryHeading.MatchString(c) {
			t.Errorf("Esperava que %q fosse reconhecido como título terciário", c)
		}
	}

	negativeCases := []string{
		"Este é um parágrafo comum que fala sobre a Introdução e outras coisas.",
		"A taxa de conformidade foi de 95% segundo os testes realizados.",
	}

	for _, c := range negativeCases {
		if regexPrimaryHeading.MatchString(c) || regexSecondaryHeading.MatchString(c) || regexTertiaryHeading.MatchString(c) {
			t.Errorf("Texto comum %q não deveria ser reconhecido como título", c)
		}
	}
}

func TestRegexLatexLeakDetection(t *testing.T) {
	leaks := []string{
		`\chapter{Introdução}`,
		`\chapter{Introdução)`,
		`\section{Contextualização}`,
		`\label{sec:contexto}`,
		`\begin{table}`,
		`\textit{Retrieval-Augmented}`,
	}

	for _, l := range leaks {
		if !regexLatexLeakCmd.MatchString(l) {
			t.Errorf("Esperava detectar vazamento de comando LaTeX em %q", l)
		}
		if !reLatexLeakInline.MatchString(l) {
			t.Errorf("Esperava detectar inline leak em %q", l)
		}
	}

	pureCmds := []string{
		`\chapter{Introdução}`,
		`\chapter{Introdução)`,
		`  \label{cap:teste}  `,
	}
	for _, pc := range pureCmds {
		if !rePureLatexCmd.MatchString(pc) {
			t.Errorf("Esperava que %q fosse reconhecido como comando puro espúrio", pc)
		}
	}

	nonLeaks := []string{
		"Texto normal sem nenhum comando do TeX.",
		"Avaliação realizada no Instituto Federal de Sergipe.",
	}

	for _, nl := range nonLeaks {
		if regexLatexLeakCmd.MatchString(nl) {
			t.Errorf("Não deveria detectar vazamento em %q", nl)
		}
	}
}

func TestPreTextualClassification(t *testing.T) {
	note := "Trabalho de Conclusão de Curso apresentado ao curso de Bacharelado em Sistemas de Informação..."
	if !isPreTextualNote(note) {
		t.Errorf("Esperava que nota institucional fosse reconhecida")
	}

	advisor := "Orientador: Prof. Me. Diego Armando de Oliveira Meneses"
	if !isPreTextualNote(advisor) {
		t.Errorf("Esperava que linha de orientador fosse reconhecida")
	}

	body := "Neste capítulo apresentamos o referencial teórico."
	if isPreTextualNote(body) {
		t.Errorf("Texto de corpo não deveria ser reconhecido como nota institucional")
	}
}
