package gdocs

import (
	"context"
	"fmt"
	"os"

	"google.golang.org/api/drive/v3"
)

func SyncPDFToDrive(ctx context.Context, fileID string, localPDFPath string) (string, error) {
	svc, err := GetDriveService(ctx)
	if err != nil {
		return "", fmt.Errorf("erro ao obter serviço do Drive: %w", err)
	}

	file, err := os.Open(localPDFPath)
	if err != nil {
		return "", fmt.Errorf("erro ao abrir PDF local: %w", err)
	}
	defer file.Close()

	// Faz APENAS o update da mídia no ID fornecido
	_, err = svc.Files.Update(fileID, &drive.File{}).Media(file).Do()
	if err != nil {
		return "", fmt.Errorf("erro ao atualizar o PDF no Drive: %w", err)
	}

	return fmt.Sprintf("SUCESSO: PDF atualizado silenciosamente no ID %s", fileID), nil
}
