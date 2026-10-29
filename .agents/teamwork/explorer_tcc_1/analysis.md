# Relatório de Análise Técnica e Mapeamento dos Artefatos do TCC

**Autor:** explorer_tcc_1 (`teamwork_explorer`)  
**Data:** 29 de Outubro de 2026  
**Ambiente:** `c:\Projetos\golang\mcp-gdocs-latex`  
**Escopo:** Investigação detalhada dos artefatos em `resources/tcc/`, `resources/` e subsidiação arquitetural para o servidor Go `mcp-gdocs-latex`.

---

## 1. Visão Geral e Estrutura de Diretórios

O diretório `resources/tcc/` armazena a totalidade dos artefatos-fonte e compilados do Trabalho de Conclusão de Curso (TCC I) intitulado:
> *"Avaliação de uma Arquitetura Híbrida RAG-MCP na Validação Normativa de Requerimentos Acadêmicos"*  
> **Autor:** José Gustavo Correia Nascimento  
> **Orientador:** Prof. Diego Armando de Oliveira Meneses  
> **Instituição:** Instituto Federal de Sergipe (IFS) — *Campus* Lagarto — Bacharelado em Sistemas de Informação

### Árvore de Arquivos Mapeada

```text
resources/
├── Avaliação de uma Arquitetura Híbrida RAG-MCP na Validação Normativa... (2).md [52.239 B]
├── Avaliação de uma Arquitetura Híbrida RAG-MCP na Validação Normativa... .rtf   [166.425 B]
├── mcp_sync_history_and_improvements.md                                          [19.310 B]
├── outras infromações.md                                                         [19.310 B]
└── tcc/
    ├── config.tex           [6.868 B]   Configurações institucionais, pacotes e macros ABNT
    ├── metadados.tex        [1.470 B]   Dados do autor, título, instituição e orientador
    ├── pre-textuais.tex     [5.934 B]   Resumo, Abstract e seções pré-textuais
    ├── referencias.bib      [19.131 B]  Base de dados bibliográfica (BibLaTeX)
    ├── logo_ifs.png         [13.583 B]  Brasão institucional em alta resolução
    ├── main.tex             [2.019 B]   Documento mestre orquestrador
    ├── main.pdf             [1.098.114 B] PDF final compilado do TCC
    ├── capitulos/
    │   ├── 01-introducao.tex      [18.971 B]  Capítulo 1: Introdução e delimitação
    │   ├── 02-referencial.tex     [44.178 B]  Capítulo 2: Fundamentação e trabalhos correlatos
    │   ├── 03-metodologia.tex     [46.983 B]  Capítulo 3: Metodologia e desenho fatorial 2x2
    │   ├── 04-desenvolvimento.tex [73.341 B]  Capítulo 4: Engenharia do protótipo e ADRs
    │   ├── 05-resultados.tex      [829 B]     Capítulo 5: Resultados (planejamento TCC II)
    │   └── 06-conclusao.tex       [2.553 B]   Capítulo 6: Considerações finais
    └── figuras/
        ├── arquitetura.mmd                [1.816 B]   Diagrama Mermaid da arquitetura
        ├── arquitetura.pdf                [194.551 B] Vetor PDF compilado da arquitetura
        ├── casos_uso.mmd                  [1.380 B]   Diagrama Mermaid de casos de uso
        ├── casos_uso.pdf                  [122.920 B] Vetor PDF compilado de casos de uso
        ├── procedimento_experimental.mmd  [1.558 B]   Diagrama Mermaid do procedimento
        ├── procedimento_experimental.pdf  [102.954 B] Vetor PDF compilado do procedimento
        ├── sequencia.mmd                  [1.847 B]   Diagrama Mermaid de sequência temporal
        └── sequencia.pdf                  [109.619 B] Vetor PDF compilado de sequência
```

---

## 2. Estrutura e Formatação LaTeX (ABNT NBR 14724 e NBR 6023)

### 2.1 Preâmbulo e Pacotes Carregados (`config.tex` e `main.tex`)

O documento utiliza a classe **`abntex2`** com as opções:
- `12pt`, `oneside`, `a4paper`, `english`, `brazilian`, `sumario=tradicional`, `chapter=TITLE`.
- **Tipografia:** Fonte Palatino via pacote `mathpazo` com microtipografia ativada via `microtype` (`activate={true,nocompatibility}, stretch=10, shrink=10`).
- **Tabelas:** Pacotes `booktabs`, `tabularx`, `multirow`, `array`, `longtable`, `siunitx`, `[section]{placeins}`.
- **Códigos-fonte:** `listings`, `xcolor`, `courier` (fonte monoespaçada).
- **Listas:** `enumitem` com `topsep=0pt, itemsep=0pt, parsep=4pt, leftmargin=\parindent`.
- **Bibliografia:** `biblatex` com estilo oficial `style=abnt, justify, maxbibnames=99, sorting=nty` e espaçamento entre entradas `\setlength{\bibitemsep}{\baselineskip}` (conforme NBR 6023).

### 2.2 Regras de Formatação ABNT e Equivalência com Google Docs API

| Parâmetro ABNT | Valor no LaTeX (`config.tex`) | Especificação Técnica Google Docs API (`batchUpdate`) |
|---|---|---|
| **Margem Esquerda** | `3,0 cm` (`\setlrmarginsandblock{3cm}{2cm}{*}`) | `UpdateDocumentStyleRequest`: `marginLeft: { magnitude: 85.04, unit: "PT" }` |
| **Margem Superior** | `3,0 cm` (`\setulmarginsandblock{3cm}{2cm}{*}`) | `UpdateDocumentStyleRequest`: `marginTop: { magnitude: 85.04, unit: "PT" }` |
| **Margem Direita** | `2,0 cm` | `UpdateDocumentStyleRequest`: `marginRight: { magnitude: 56.69, unit: "PT" }` |
| **Margem Inferior** | `2,0 cm` | `UpdateDocumentStyleRequest`: `marginBottom: { magnitude: 56.69, unit: "PT" }` |
| **Espaçamento entre Linhas** | `1,5 linha` (`\OnehalfSpacing`) | `UpdateParagraphStyleRequest`: `lineSpacing: 150` |
| **Recuo de Parágrafo** | `1,25 cm` (`\setlength{\parindent}{1.25cm}`) | `UpdateParagraphStyleRequest`: `indentFirstLine: { magnitude: 35.43, unit: "PT" }`, `indentStart: { magnitude: 0, unit: "PT" }` |
| **Alinhamento do Texto** | Justificado (`\frenchspacing`, padrão abntex2) | `UpdateParagraphStyleRequest`: `alignment: "JUSTIFIED"` |
| **Tamanho da Fonte (Corpo)** | `12 pt` (`\normalsize`) | `UpdateTextStyleRequest`: `fontSize: { magnitude: 12, unit: "PT" }` |
| **Citação Longa (> 3 linhas)** | Recuo `4,0 cm`, fonte menor, entrelinha simples | `indentStart: { magnitude: 113.39, unit: "PT" }`, `fontSize: { magnitude: 10, unit: "PT" }`, `lineSpacing: 100` |
| **Títulos de Capítulos** | Caixa Alta, Negrito, 12pt (`chapter=TITLE`) | `UpdateParagraphStyleRequest`: `namedStyleType: "HEADING_1"`, texto em `UPPERCASE`, `bold: true` |
| **Títulos de Seções** | Caixa Baixa, Negrito, 12pt (`\ABNTEXuppersection{false}`) | `UpdateParagraphStyleRequest`: `namedStyleType: "HEADING_2"`, `bold: true` |
| **Legenda de Tabela/Figura** | `\small\bfseries`, delimitador ` -- ` | `alignment: "CENTER"`, `fontSize: { magnitude: 10, unit: "PT" }`, label em negrito |
| **Fonte de Tabela/Figura** | `\small`, alinhamento sob o objeto | `alignment: "CENTER"`, `fontSize: { magnitude: 10, unit: "PT" }` |

---

## 3. Mapeamento Exaustivo de TODAS as 22 Tabelas do TCC

Foram identificadas exatamente **22 tabelas** no corpo do trabalho, distribuídas em 4 capítulos:
- **Capítulo 1:** 1 tabela
- **Capítulo 2:** 1 tabela
- **Capítulo 3:** 5 tabelas
- **Capítulo 4:** 15 tabelas

### 3.1 Padrão Construtivo e Sintaxe no LaTeX
- **21 tabelas** utilizam o ambiente `tabularx` com largura total de página `\textwidth`, colunas dinâmicas `X` e colunas prefixadas `>{\raggedright\arraybackslash}p{...}` para quebra automática de texto sem hifenação defeituosa.
- **1 tabela** (Tabela 14 - `tab:resumo_situacao_mvp`) utiliza o ambiente `tabular` tradicional com colunas `@{}l c c@{}`.
- Todas as 22 tabelas aplicam rigorosamente as réguas do pacote `booktabs`:
  - `\toprule`: linha superior reforçada (0,08 em)
  - `\midrule`: linha divisória de cabeçalho (0,05 em)
  - `\bottomrule`: linha de fechamento inferior (0,08 em)
  - `\addlinespace`: microespaçamento semântico (0,5 em) para agrupamentos de linhas.
- Em conformidade com a ABNT NBR 14724, todas as tabelas contêm:
  - Topo: `\caption{...}` e `\label{...}`
  - Rodapé: `\fonte{...}` logo antes de `\end{table}`.

---

### 3.2 Tabela Analítica Consolidada (As 22 Tabelas do TCC)

| # | Arquivo e Linhas | Rótulo (`\label`) | Legenda (`\caption`) | Ambiente | Dimensões | Cabeçalhos | Fonte |
|---|---|---|---|---|---|---|---|
| **1** | `capitulos/01-introducao.tex`<br>L. 256–271 | `tab:escopo_recorte` | Delimitação do escopo experimental em relação ao protótipo funcional | `tabularx`<br>`\small` | 3 colunas<br>3 linhas de dados | Dimensão Funcional \| Descrição do Escopo no Protótipo \| Papel na Avaliação Experimental | Elaborado pelo autor (2026). |
| **2** | `capitulos/02-referencial.tex`<br>L. 113–133 | `tab:trabalhos_relacionados` | Quadro comparativo entre os trabalhos relacionados e a presente pesquisa | `tabularx`<br>`\footnotesize` | 5 colunas<br>7 linhas de dados | Trabalho \| Paradigma \| Base de Dados \| Acesso ao Sistema \| Procedimento e Avaliação | Elaborado pelo autor (2026). |
| **3** | `capitulos/03-metodologia.tex`<br>L. 21–36 | `quadro:escopo_metodologico` | Distinção entre a solução tecnológica ampla e o objeto de investigação | `tabularx` | 3 colunas<br>4 linhas de dados | Dimensão \| Protótipo Desenvolvido \| Objeto de Investigação Experimental | Elaborado pelo autor (2026). |
| **4** | `capitulos/03-metodologia.tex`<br>L. 214–231 | `tab:distribuicao-requisitos` | Distribuição quantitativa dos requisitos de software consolidados | `tabularx`<br>`\small` | 3 colunas<br>4 linhas (3 + total) | Categoria \| Quantidade \| Descrição do Escopo | Elaborado pelo autor (2026). |
| **5** | `capitulos/03-metodologia.tex`<br>L. 270–288 | `tab:moscow` | Distribuição dos requisitos pelo método MoSCoW | `tabularx`<br>`\small` | 3 colunas<br>5 linhas (4 + total) | Prioridade \| Quantidade \| Papel no Desenvolvimento do Protótipo | Elaborado pelo autor (2026). |
| **6** | `capitulos/03-metodologia.tex`<br>L. 326–341 | `tab:fatorial_2x2` | Delineamento experimental fatorial $2 \times 2$ e condições de teste | `tabularx` | 4 colunas<br>4 linhas de dados | Condição \| Fator N \| Fator D \| Contexto e Ferramentas Disponibilizadas ao Modelo | Elaborado pelo autor (2026). |
| **7** | `capitulos/03-metodologia.tex`<br>L. 551–567 | `tab:rubrica_operacional` | Definição operacional dos quatro critérios da rubrica TCD | `tabularx`<br>`\small` | 4 colunas<br>4 linhas de dados | Critério \| Peso \| Condição para Nota 1 (Conforme) \| Condição para Nota 0 (Não Conforme) | Elaborado pelo autor (2026). |
| **8** | `capitulos/04-desenvolvimento.tex`<br>L. 177–195 | `tab:rf_dados` | Requisitos funcionais de acesso e agregação de dados discentes | `tabularx`<br>`\small` | 5 colunas<br>6 linhas (RF-01 a 06) | ID \| Descrição Operacional \| MoSCoW \| Status \| Componente Oficial | Elaborado pelo autor (2026). |
| **9** | `capitulos/04-desenvolvimento.tex`<br>L. 203–218 | `tab:rf_normas` | Requisitos funcionais de recuperação de conhecimento normativo e calendário | `tabularx`<br>`\small` | 5 colunas<br>3 linhas (RF-07 a 09) | ID \| Descrição Operacional \| MoSCoW \| Status \| Componente Oficial | Elaborado pelo autor (2026). |
| **10** | `capitulos/04-desenvolvimento.tex`<br>L. 226–243 | `tab:rf_validacao` | Requisitos funcionais de validação normativa e emissão documental | `tabularx`<br>`\small` | 5 colunas<br>5 linhas (RF-10 a 14) | ID \| Descrição Operacional \| MoSCoW \| Status \| Componente Oficial | Elaborado pelo autor (2026). |
| **11** | `capitulos/04-desenvolvimento.tex`<br>L. 251–269 | `tab:rf_ui` | Requisitos funcionais da interface visual e supervisão discente (HITL) | `tabularx`<br>`\small` | 5 colunas<br>6 linhas (RF-15 a 20) | ID \| Descrição Operacional \| MoSCoW \| Status \| Componente Oficial | Elaborado pelo autor (2026). |
| **12** | `capitulos/04-desenvolvimento.tex`<br>L. 277–293 | `tab:rf_dialogo` | Requisitos funcionais de mediação dialógica proativa | `tabularx`<br>`\small` | 5 colunas<br>4 linhas (RF-21 a 24) | ID \| Descrição Operacional \| MoSCoW \| Status \| Componente Oficial | Elaborado pelo autor (2026). |
| **13** | `capitulos/04-desenvolvimento.tex`<br>L. 302–322 | `tab:rnf_qualidade` | Catálogo de requisitos não funcionais e atributos de qualidade (ISO/IEC 25010) | `tabularx`<br>`\small` | 5 colunas<br>8 linhas (RNF-01 a 08) | ID \| Atributo ISO 25010 \| Descrição do Requisito de Qualidade \| MoSCoW \| Status | Elaborado pelo autor (2026). |
| **14** | `capitulos/04-desenvolvimento.tex`<br>L. 339–356 | `tab:resumo_situacao_mvp` | Distribuição dos requisitos de software por situação no protótipo | `tabular`<br>`\small` | 3 colunas<br>4 linhas (3 + total) | Situação no Protótipo \| Quantidade \| Percentual | Elaborado pelo autor (2026). |
| **15** | `capitulos/04-desenvolvimento.tex`<br>L. 386–411 | `quadro:trajetoria_decisoes` | Dilemas de engenharia, alternativas avaliadas e decisões consolidadas | `tabularx`<br>`\small` | 4 colunas<br>7 linhas de decisões | Momento / Eixo \| Alternativa Considerada / Testada \| Comportamento / Risco Identificado \| Decisão Consolidada | Elaborado pelo autor (2026). |
| **16** | `capitulos/04-desenvolvimento.tex`<br>L. 504–525 | `tab:sintese_adrs` | Síntese do catálogo de Decisões Arquiteturais (ADRs) do projeto | `tabularx`<br>`\small` | 4 colunas<br>9 linhas de ADRs | ID \| Domínio Técnico \| Decisão Adotada no Sistema \| Status | Elaborado pelo autor (2026). |
| **17** | `capitulos/04-desenvolvimento.tex`<br>L. 652–669 | `tab:componentes` | Módulos de software implementados no servidor oficial | `tabularx`<br>`\small` | 3 colunas<br>5 componentes | Componente Oficial \| Implementação Técnica \| Função Operacional no Sistema | Elaborado pelo autor (2026). |
| **18** | `capitulos/04-desenvolvimento.tex`<br>L. 712–731 | `tab:rn_trancamento` | Regras de negócio normativas referentes a Trancamento de Matrícula (ROD) | `tabularx`<br>`\small` | 5 colunas<br>7 linhas (RN-01 a 07) | ID \| Dispositivo \| Regra de Negócio Derivada \| MoSCoW \| Integração | Elaborado pelo autor a partir do ROD (IFS, 2016). |
| **19** | `capitulos/04-desenvolvimento.tex`<br>L. 739–758 | `tab:rn_proficiencia` | Regras de negócio normativas referentes a Exame de Proficiência (ROD) | `tabularx`<br>`\small` | 5 colunas<br>7 linhas (RN-08 a 14) | ID \| Dispositivo \| Regra de Negócio Derivada \| MoSCoW \| Integração | Elaborado pelo autor a partir do ROD e resoluções complementares. |
| **20** | `capitulos/04-desenvolvimento.tex`<br>L. 768–783 | `tab:rn_calendario` | Regras de negócio temporais baseadas no Calendário Acadêmico Oficial | `tabularx`<br>`\small` | 5 colunas<br>3 linhas (RN-15 a 17) | ID \| Marco Oficial \| Regra Regimental de Prazo \| MoSCoW \| Integração | Elaborado pelo autor a partir do Calendário e SUAP (IFS, 2026). |
| **21** | `capitulos/04-desenvolvimento.tex`<br>L. 856–883 | `tab:tradeoffs_rod_json` | Matriz analítica de trade-offs da modelagem das regras do ROD em JSON estruturado | `tabularx`<br>`\small` | 3 colunas<br>4 dimensões | Dimensão Analítica \| Vantagens Determinísticas (Prós) \| Desvantagens e Complexidades (Contras) | Elaborado pelo autor (2026). |
| **22** | `capitulos/04-desenvolvimento.tex`<br>L. 1101–1123 | `tab:rastreabilidade_solucao` | Matriz de rastreabilidade entre necessidades, requisitos e validação | `tabularx`<br>`\small` | 5 colunas<br>10 linhas de regras | ID \| Requisito / Regra \| Prioridade \| Componente Oficial \| Mecanismo de Validação | Elaborado pelo autor (2026). |

---

### 3.3 Mecanismo de Conversão LaTeX -> Google Docs Native Tables

Na Google Docs API v1, a inserção de tabelas nativas opera através do comando `InsertTableRequest`:
1. **Passo 1 (Criação da Grade):** Enviar `InsertTableRequest` no índice de inserção desejado (`endOfSegmentLocation` ou `index`), especificando `rows` e `columns`.
2. **Passo 2 (Mapeamento de Índices):** A criação de uma tabela $R \times C$ insere imediatamente células vazias no documento, cada uma contendo um parágrafo delimitador.
3. **Passo 3 (Injeção de Células):** Para preencher cada célula sem corrupção posicional, deve-se ler o documento atualizado ou calcular os deslocamentos estritos em ordem reversa (ou incremental controlada).
4. **Passo 4 (Estilo das Células e Bordas):**
   - Cabeçalhos: Aplicar `UpdateTableCellStyleRequest` com preenchimento em tom cinza claro (`#F2F2F2`) e `UpdateTextStyleRequest` com `bold: true`.
   - Bordas ABNT: A ABNT recomenda tabelas abertas nas laterais (sem bordas verticais externas/internas) e apenas réguas horizontais no topo, cabeçalho e rodapé. No Google Docs, configura-se via `tableCellStyle.borderTop`, `borderBottom` com espessura 1pt / cor preta, e `borderLeft = nil`, `borderRight = nil`.

---

## 4. Figuras e Diagramas Mermaid

### 4.1 Mapeamento e Análise dos Diagramas

Todos os 4 diagramas visuais do TCC foram desenvolvidos em sintaxe Mermaid (`.mmd`) e estão situados em `resources/tcc/figuras/`:

| Arquivo `.mmd` | Vetor Gerado | Tamanho `.mmd` | Tipo Mermaid | Figura no LaTeX | Legenda no TCC |
|---|---|---|---|---|---|
| `arquitetura.mmd` | `arquitetura.pdf` (194 KB) | 1.816 B | `flowchart TB` | Figura 2 (`04-desenvolvimento.tex:89-95`) | Fluxo de processamento e recuperação de contexto via ferramentas MCP para validação normativa. |
| `casos_uso.mmd` | `casos_uso.pdf` (122 KB) | 1.380 B | `flowchart LR` | Figura 3 (`04-desenvolvimento.tex:943-949`) | Diagrama de Casos de Uso do assistente acadêmico baseado em MCP e recuperação de contexto. |
| `procedimento_experimental.mmd` | `procedimento_experimental.pdf` (102 KB) | 1.558 B | `flowchart LR` | Figura 1 (`03-metodologia.tex:624-630`) | Fluxograma do procedimento experimental da pesquisa |
| `sequencia.mmd` | `sequencia.pdf` (109 KB) | 1.847 B | `sequenceDiagram` | Figura 4 (`04-desenvolvimento.tex:968-974`) | Diagrama de Sequência da validação normativa e emissão de requerimento. |

### 4.2 Análise de Conteúdo e Complexidade de Cada Diagrama

1. **`arquitetura.mmd` (Flowchart Top-to-Bottom):**
   - **Estrutura:** 4 subgrafos aninhados (`Host`, `Server`, `Tools`, `Providers`, `GroundTruth`).
   - **Fluxos:** Conexão STDIO JSON-RPC do LLM com 4 Tools (`obter-dados-discente`, `gerar-requerimento`, `consultar-calendario-academico`, `consultar-regras-rod`), mapeadas a 4 Providers de serviço e 3 bases de referência canônica (`fixtures-tcd-v0.2.json`, `calendarioLagarto.json`, `rodRegras.json`).
   - **Complexidade:** Alta densidade de nós e orientações ortogonais mistas (`TB` e `LR`).

2. **`casos_uso.mmd` (Flowchart Left-to-Right com Semântica UML):**
   - **Estrutura:** Ator `Discente`, subsistema `Sistema` (5 Casos de Uso: UC01 a UC05), subgrafo `Validacoes` e banco `SUAP_API`.
   - **Relações:** Aplicação de estereótipos UML (`include` para UC_ROD, UC_CAL e UC_SUAP em UC02/UC03; `extends` para UC04 sobre UC02/UC03).
   - **Complexidade:** Média, exige renderização nítida de rótulos em conexões pontilhadas.

3. **`procedimento_experimental.mmd` (Flowchart Left-to-Right em 2 Colunas e 4 Fases):**
   - **Estrutura:** Divisão em `Col1` (Etapa Inicial: Planejamento e Calibração) e `Col2` (Etapa Final: Execução e Análise).
   - **Nós:** 4 Fases sequenciais (F1 Preparatória, F2 Calibração com nó de decisão condicional em losango `F{"Protocolo e isolamento OK?"}`, F3 Coleta Sistemática com 204 execuções, F4 Avaliação com Rubrica TCD).
   - **Complexidade:** Média/Alta, diagramação balanceada horizontalmente.

4. **`sequencia.mmd` (Sequence Diagram):**
   - **Estrutura:** `autonumber`, 1 Ator (`Discente`), 5 Participantes (`Interface Web`, `Orquestrador LLM`, `Servidor MCP`, `IAlunoDataProvider`, `Motor de Regras`).
   - **Fases com zonas coloridas:** Utiliza 3 blocos `rect` destacados:
     * `rect rgb(240, 248, 255)` — Fase 1: Recuperação Transacional de Prontuário
     * `rect rgb(255, 250, 240)` — Fase 2: Auditoria Normativa e Tempestividade
     * `rect rgb(240, 255, 240)` — Fase 3: Geração da Minuta e Supervisão HITL
   - **Complexidade:** Alta, exige renderizador que suporte caixas de agrupamento de fases coloridas e notas contextuais.

### 4.3 Requisitos e Pipeline para Inserção no Google Docs

#### O Gargalo Técnico:
- O repositório LaTeX atualmente compila os arquivos Mermaid diretamente para arquivos **`.pdf`** vetoriais.
- **A API do Google Docs (`InsertInlineImageRequest`) NÃO aceita arquivos PDF vetoriais nem código Mermaid puro.**
- A API aceita estritamente imagens rasterizadas públicas ou URIs com protocolo de mídia acessível pelo Google Docs (PNG, JPEG, GIF) ou arquivos referenciados no Drive.

#### Pipeline Recomendado para Conversão e Inserção:
```text
[Arquivo .mmd]
      │
      ▼  (Mermaid CLI: mmdc -i input.mmd -o output.png -s 3 -b transparent)
[Arquivo .png em Alta Resolução (300 DPI)]
      │
      ▼  (Google Drive API v3: Files.Create / Media Upload)
[Imagem hospedada no Google Drive com permissão de leitura]
      │
      ▼  (Google Docs API v1: batchUpdate -> InsertInlineImageRequest)
[Imagem nativa inserida inline no Google Doc com Caption e Fonte ABNT]
```

- **Comando de Compilação Mermaid CLI:**
  ```bash
  npx @mermaid-js/mermaid-cli -i resources/tcc/figuras/arquitetura.mmd -o resources/tcc/figuras/arquitetura.png -s 3 -b transparent
  ```
- **Alternativa via Conversão de PDF existente:**
  Caso o ambiente local não disponha do Node.js/Mermaid CLI, os PDFs vetoriais já gerados em `resources/tcc/figuras/*.pdf` podem ser convertidos para PNG em alta resolução (300 DPI) usando ferramentas de renderização gráfica (como `pdftoppm -png -r 300`).

---

## 5. PDF Compilado Principal e Validação de Integridade via Hash MD5

### 5.1 Identificação do Arquivo Local
- **Nome do Arquivo:** `main.pdf`
- **Caminho Absoluto:** `c:\Projetos\golang\mcp-gdocs-latex\resources\tcc\main.pdf`
- **Tamanho Exato em Disco:** **1.098.114 bytes** (~1,05 MB)

### 5.2 Mapeamento e Referência no Google Drive
Nos arquivos `resources/Avaliação de uma Arquitetura Híbrida RAG-MCP na Validação Normativa de Requerimentos Acadêmicos (2).md` e `.rtf`, o PDF final está explicitamente documentado:
- **URL no Google Drive:** `https://drive.google.com/file/d/1mbBbWNz-N1P4n7SYnj28i9yHCFz9PqZH/view`
- **ID do Arquivo no Google Drive (`fileID`):** `1mbBbWNz-N1P4n7SYnj28i9yHCFz9PqZH`

### 5.3 Análise da Implementação Atual em Go (`internal/gdocs/upload.go`)
No servidor Go atual, a função `SyncPDFToDrive` executa o upload cego sem qualquer verificação prévia de integridade:
```go
// Implementação atual em internal/gdocs/upload.go:
func SyncPDFToDrive(ctx context.Context, fileID string, localPDFPath string) (string, error) {
    ...
    // Faz APENAS o update da mídia no ID fornecido
    _, err = svc.Files.Update(fileID, &drive.File{}).Media(file).Do()
    ...
}
```
**Deficiências observadas:**
1. Não calcula o hash do arquivo local antes do envio.
2. Não consulta o hash remoto do Google Drive (`md5Checksum`), resultando em upload redundante de 1,1 MB mesmo se o arquivo for exatamente igual.
3. Não valida após o upload se os bytes foram gravados com integridade perfeita na nuvem.

### 5.4 Especificação Algorítmica da Sincronização Inteligente por MD5

A Google Drive API v3 disponibiliza nativamente para qualquer arquivo binário (não convertível em Doc) o metadado canônico `md5Checksum`.

#### Fluxo Algorítmico Recomendado para o Servidor Go:
1. **Cálculo Local:**
   ```go
   func CalculateMD5(filePath string) (string, error) {
       f, err := os.Open(filePath)
       if err != nil {
           return "", err
       }
       defer f.Close()
       h := md5.New()
       if _, err := io.Copy(h, f); err != nil {
           return "", err
       }
       return hex.EncodeToString(h.Sum(nil)), nil
   }
   ```
2. **Consulta de Metadados no Drive v3:**
   ```go
   remoteFile, err := svc.Files.Get(fileID).Fields("id, name, md5Checksum, size").Do()
   ```
3. **Decisão e Validação Atômica:**
   - Se `localMD5 == remoteFile.Md5Checksum`:
     Retornar status `"SYNC_SKIPPED: O PDF local já é estritamente idêntico ao arquivo na nuvem (MD5: <hash>)"`.
   - Se `localMD5 != remoteFile.Md5Checksum`:
     Executar `svc.Files.Update(fileID, nil).Media(file).Fields("id, md5Checksum").Do()`.
     Verificar se `updatedFile.Md5Checksum == localMD5`. Se idênticos, reportar `"SYNC_SUCCESS"`, senão emitir erro de integridade de rede.

---

## 6. Conclusões da Investigação para a Arquitetura do MCP

1. **LaTeX como Fonte Única da Verdade:** A estrutura em `resources/tcc` é altamente modular (`config.tex`, `metadados.tex`, capítulos isolados). Qualquer alteração de texto promovida pela IA no Google Docs deve ser retropropagada no `.tex` correspondente preservando a sintaxe de macros e pacotes.
2. **Fidelidade Tabular:** As 22 tabelas possuem complexidade estrutural significativa (cabeçalhos técnicos, notas de rodapé, réguas booktabs). A criação de tabelas nativas via API do Google Docs é viável e recomendada em substituição ao texto puro ou ASCII Art corrompido.
3. **Conversão de Diagramas:** Os 4 diagramas Mermaid compilados para PDF precisam passar por conversão raster (PNG 300 DPI) para serem inseridos de forma fluida no Google Docs via `InsertInlineImageRequest`.
4. **Governança de Sincronia de PDF:** A integração com o ID `1mbBbWNz-N1P4n7SYnj28i9yHCFz9PqZH` deve adotar comparação prévia de `md5Checksum` local vs remoto, garantindo economia de banda e integridade verificável de ponta a ponta.
