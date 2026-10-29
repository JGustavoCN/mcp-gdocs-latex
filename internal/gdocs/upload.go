package gdocs

import (
	"context"
	"fmt"
)

// SyncPDFToDrive atualiza o PDF compilado no Google Drive com verificação idempotente de hash MD5 e preservação de histórico
func SyncPDFToDrive(ctx context.Context, fileID string, localPDFPath string) (string, error) {
	report, err := SynchronizePDF(ctx, localPDFPath, fileID)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("SUCESSO [%s]: %s (MD5: %s | Tamanho: %d bytes)",
		report.Status, report.Message, report.LocalMD5, report.SizeBytes), nil
}
