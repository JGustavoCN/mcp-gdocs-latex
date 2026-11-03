package gdocs

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// HeadingSequenceAnomaly descreve uma descontinuidade ou salto na numeração progressiva ABNT
type HeadingSequenceAnomaly struct {
	Title        string `json:"title"`
	ParentPrefix string `json:"parent_prefix"`
	Expected     string `json:"expected"`
	Actual       string `json:"actual"`
	Details      string `json:"details"`
}

// HeadingRecord auxilia a ordenação e auditoria hierárquica dos títulos
type HeadingRecord struct {
	Text         string
	Number       string
	Parts        []int
	StartIndex   int64
	EndIndex     int64
	ElementIndex int
	StyleName    string
}

var reSectionNumberExtractor = regexp.MustCompile(`^([1-9][0-9]*(?:\.[0-9]+)*)\b`)

// ParseSectionNumber decompõe "2.3.1" em []int{2, 3, 1}
func ParseSectionNumber(numStr string) []int {
	parts := strings.Split(numStr, ".")
	var res []int
	for _, p := range parts {
		if val, err := strconv.Atoi(strings.TrimSpace(p)); err == nil {
			res = append(res, val)
		} else {
			return nil
		}
	}
	return res
}

// AuditHeadingNumberingContinuity analisa a sequência de números de seções identificadas no documento
// e reporta qualquer quebra na numeração progressiva ABNT NBR 6024 (ex: salto de 2.2 para 2.5).
func AuditHeadingNumberingContinuity(records []HeadingRecord) []HeadingSequenceAnomaly {
	var anomalies []HeadingSequenceAnomaly
	if len(records) < 2 {
		return anomalies
	}

	// Agrupa seções consecutivas por nível hierárquico pai
	// Por exemplo, para "2.1", "2.2", "2.5", o pai é "2" e a sequência de filhos é [1, 2, 5]
	type siblingTracker struct {
		lastSubNumber int
		lastFullStr   string
	}
	parentTrackers := make(map[string]*siblingTracker)

	for _, rec := range records {
		if len(rec.Parts) == 0 {
			continue
		}

		// Nível 1 (Capítulo: 1, 2, 3...)
		if len(rec.Parts) == 1 {
			currentChap := rec.Parts[0]
			tracker, exists := parentTrackers["root"]
			if !exists {
				parentTrackers["root"] = &siblingTracker{lastSubNumber: currentChap, lastFullStr: rec.Number}
			} else {
				if currentChap > tracker.lastSubNumber+1 {
					var missing []string
					for m := tracker.lastSubNumber + 1; m < currentChap; m++ {
						missing = append(missing, strconv.Itoa(m))
					}
					anomalies = append(anomalies, HeadingSequenceAnomaly{
						Title:        rec.Text,
						ParentPrefix: "Capítulos",
						Expected:     strings.Join(missing, ", "),
						Actual:       rec.Number,
						Details: fmt.Sprintf("Quebra de sequência nos Capítulos: detectado '%s' após '%s', ausentes capítulo(s) intermediário(s): %s (NBR 6024).",
							rec.Number, tracker.lastFullStr, strings.Join(missing, ", ")),
					})
				}
				tracker.lastSubNumber = currentChap
				tracker.lastFullStr = rec.Number
			}
			continue
		}

		// Nível 2 ou superior (ex: 2.1, 2.2, 2.3.1...)
		parentPrefix := strings.Join(strings.Split(rec.Number, ".")[:len(rec.Parts)-1], ".")
		currentSubNumber := rec.Parts[len(rec.Parts)-1]

		tracker, exists := parentTrackers[parentPrefix]
		if !exists {
			// Se o primeiro filho de uma seção começar maior que 1 (ex: primeiro filho de 2 for 2.2 em vez de 2.1)
			if currentSubNumber > 1 {
				var missing []string
				for m := 1; m < currentSubNumber; m++ {
					missing = append(missing, fmt.Sprintf("%s.%d", parentPrefix, m))
				}
				anomalies = append(anomalies, HeadingSequenceAnomaly{
					Title:        rec.Text,
					ParentPrefix: parentPrefix,
					Expected:     strings.Join(missing, ", "),
					Actual:       rec.Number,
					Details: fmt.Sprintf("Quebra de sequência na Seção %s: a numeração iniciou em '%s' em vez de '%s.1', ausente(s): %s (NBR 6024).",
						parentPrefix, rec.Number, parentPrefix, strings.Join(missing, ", ")),
				})
			}
			parentTrackers[parentPrefix] = &siblingTracker{lastSubNumber: currentSubNumber, lastFullStr: rec.Number}
		} else {
			// Salto entre irmãos (ex: de 2.2 para 2.5)
			if currentSubNumber > tracker.lastSubNumber+1 {
				var missing []string
				for m := tracker.lastSubNumber + 1; m < currentSubNumber; m++ {
					missing = append(missing, fmt.Sprintf("%s.%d", parentPrefix, m))
				}
				anomalies = append(anomalies, HeadingSequenceAnomaly{
					Title:        rec.Text,
					ParentPrefix: parentPrefix,
					Expected:     strings.Join(missing, ", "),
					Actual:       rec.Number,
					Details: fmt.Sprintf("Quebra de sequência nos títulos da Seção %s: detectado '%s' após '%s', ausente(s) seção(ões) intermediária(s): %s (NBR 6024).",
						parentPrefix, rec.Number, tracker.lastFullStr, strings.Join(missing, ", ")),
				})
			}
			tracker.lastSubNumber = currentSubNumber
			tracker.lastFullStr = rec.Number
		}
	}

	return anomalies
}

