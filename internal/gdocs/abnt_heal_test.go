package gdocs

import (
	"testing"
)

func TestShouldMergeParagraphs(t *testing.T) {
	tests := []struct {
		name        string
		curr        string
		next        string
		isSpecial1  bool
		isSpecial2  bool
		expected    bool
	}{
		{
			name:       "Hard wrap quebrado no meio da frase (caso clássico LaTeX)",
			curr:       "que retornam dados estruturados, permitindo ao modelo de linguagem recuperar os",
			next:       "elementos pertinentes à solicitação e utilizá-los na elaboração do documento.",
			isSpecial1: false,
			isSpecial2: false,
			expected:   true,
		},
		{
			name:       "Frase com vírgula no final que continua em minúscula",
			curr:       "Além da aplicação específica à elaboração de requerimentos, a disponibilização",
			next:       "estruturada de informações e funcionalidades institucionais pode favorecer o",
			isSpecial1: false,
			isSpecial2: false,
			expected:   true,
		},
		{
			name:       "Parágrafo com ponto final seguido de nova frase iniciada em maiúscula",
			curr:       "Essa abordagem possibilita investigar como a ampliação do contexto informacional disponível ao modelo influencia a conformidade das minutas geradas.",
			next:       "Além da aplicação específica à elaboração de requerimentos, a disponibilização estruturada de informações pode favorecer o desenvolvimento.",
			isSpecial1: false,
			isSpecial2: false,
			expected:   false,
		},
		{
			name:       "Parágrafo terminando em dois pontos (introdução de lista/alínea)",
			curr:       "Os seguintes requisitos foram levantados:",
			next:       "a) Autenticação via OAuth2 com PKCE;",
			isSpecial1: false,
			isSpecial2: false,
			expected:   false,
		},
		{
			name:       "Título primário seguido de texto de abertura",
			curr:       "1 INTRODUÇÃO",
			next:       "Este trabalho investiga a integração de IA na redação acadêmica.",
			isSpecial1: true,
			isSpecial2: false,
			expected:   false,
		},
		{
			name:       "Título secundário",
			curr:       "2.1 CONTEXTUALIZAÇÃO DO PROBLEMA",
			next:       "O contexto atual apresenta desafios de conformidade.",
			isSpecial1: true,
			isSpecial2: false,
			expected:   false,
		},
		{
			name:       "Legenda seguida de fonte",
			curr:       "Figura 1 – Diagrama de arquitetura da solução",
			next:       "Fonte: Autor (2026)",
			isSpecial1: true,
			isSpecial2: true,
			expected:   false,
		},
		{
			name:       "Hifenização de quebra de linha (sis- / tema)",
			curr:       "A arquitetura proposta para o sis-",
			next:       "tema de geração de documentos utiliza o protocolo MCP.",
			isSpecial1: false,
			isSpecial2: false,
			expected:   true,
		},
		{
			name:       "Item com marcador seguido de outro item",
			curr:       "• Primeiro requisito funcional da plataforma.",
			next:       "• Segundo requisito funcional da plataforma.",
			isSpecial1: false,
			isSpecial2: false,
			expected:   false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := shouldMergeParagraphs(tc.curr, tc.next, tc.isSpecial1, tc.isSpecial2)
			if got != tc.expected {
				t.Errorf("shouldMergeParagraphs(%q, %q, %v, %v) = %v; esperado %v",
					tc.curr, tc.next, tc.isSpecial1, tc.isSpecial2, got, tc.expected)
			}
		})
	}
}
