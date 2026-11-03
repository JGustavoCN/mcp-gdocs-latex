package gdocs

import (
	"bytes"
	"context"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"strings"

	"google.golang.org/api/drive/v3"
)

// SyncReport encapsula o relatório auditável de integridade e sincronização do PDF
type SyncReport struct {
	FileID    string `json:"file_id"`
	LocalMD5  string `json:"local_md5"`
	RemoteMD5 string `json:"remote_md5"`
	SizeBytes int64  `json:"size_bytes"`
	Status    string `json:"status"` // "IN_SYNC", "UPDATED", "ERROR"
	Message   string `json:"message"`
}

// CheckPDFSync apenas diagnostica a paridade de hash entre o PDF local e o PDF no Google Drive sem fazer upload
func CheckPDFSync(ctx context.Context, localPath, fileID string) (*SyncReport, error) {
	localMD5, size, err := CalculateLocalPDFMD5(localPath)
	if err != nil {
		return nil, err
	}

	driveSvc, err := GetDriveService(ctx)
	if err != nil {
		return nil, fmt.Errorf("erro ao obter serviço do Google Drive: %w", err)
	}

	remoteFile, err := driveSvc.Files.Get(fileID).
		Fields("id, name, md5Checksum, size, modifiedTime").
		Do()
	if err != nil {
		return nil, fmt.Errorf("falha ao consultar metadados do arquivo '%s' no Drive: %w", fileID, err)
	}

	inSync := strings.EqualFold(localMD5, remoteFile.Md5Checksum)
	status := "OUT_OF_SYNC"
	msg := fmt.Sprintf("PDF local (%s) difere da versão na nuvem (%s). Sincronização necessária.", localMD5, remoteFile.Md5Checksum)
	if inSync {
		status = "IN_SYNC"
		msg = "PDF local e arquivo no Google Drive são criptograficamente idênticos (MD5 coincidentes)."
	}

	return &SyncReport{
		FileID:    fileID,
		LocalMD5:  localMD5,
		RemoteMD5: remoteFile.Md5Checksum,
		SizeBytes: size,
		Status:    status,
		Message:   msg,
	}, nil
}

// CalculateLocalPDFMD5 calcula o hash MD5 por streaming com verificação rigorosa de integridade (ISO 32000-1) e TOCTOU
func CalculateLocalPDFMD5(localPath string) (string, int64, error) {
	f, err := os.Open(localPath)
	if err != nil {
		return "", 0, fmt.Errorf("falha ao abrir arquivo PDF local: %w", err)
	}
	defer f.Close()

	stat, err := f.Stat()
	if err != nil {
		return "", 0, fmt.Errorf("falha ao obter atributos do PDF local: %w", err)
	}

	// 1. Verificação do trailer ISO 32000-1 (%%EOF nos últimos 1024 bytes)
	readSize := int64(1024)
	if stat.Size() < readSize {
		readSize = stat.Size()
	}
	if readSize <= 0 {
		return "", 0, fmt.Errorf("arquivo PDF local está completamente vazio (0 bytes)")
	}

	tailBuf := make([]byte, readSize)
	if _, err := f.ReadAt(tailBuf, stat.Size()-readSize); err != nil {
		return "", 0, fmt.Errorf("falha ao inspecionar cauda do PDF: %w", err)
	}

	if !bytes.Contains(tailBuf, []byte("%%EOF")) {
		return "", 0, fmt.Errorf("arquivo PDF incompleto ou corrompido: marcador %%%%EOF ausente nos últimos 1024 bytes (possível compilação LaTeX em andamento)")
	}

	// 2. Cálculo do hash MD5 por streaming
	hasher := md5.New()
	size, err := io.Copy(hasher, f)
	if err != nil {
		return "", 0, fmt.Errorf("falha ao computar streaming MD5 do arquivo: %w", err)
	}

	// 3. Mitigação de corrida TOCTOU
	if size != stat.Size() {
		return "", 0, fmt.Errorf("inconsistência de concorrência (TOCTOU): tamanho inspecionado (%d) difere dos bytes lidos (%d)", stat.Size(), size)
	}

	return hex.EncodeToString(hasher.Sum(nil)), size, nil
}

// SynchronizePDF sincroniza o PDF local de forma idempotente com o Google Drive
func SynchronizePDF(ctx context.Context, localPath, fileID string) (*SyncReport, error) {
	localMD5, size, err := CalculateLocalPDFMD5(localPath)
	if err != nil {
		return nil, err
	}

	driveSvc, err := GetDriveService(ctx)
	if err != nil {
		return nil, fmt.Errorf("erro ao obter serviço do Google Drive: %w", err)
	}

	remoteFile, err := driveSvc.Files.Get(fileID).
		Fields("id, name, md5Checksum, size, modifiedTime").
		Do()
	if err != nil {
		return nil, fmt.Errorf("falha ao consultar metadados do arquivo '%s' no Drive: %w", fileID, err)
	}

	// Se já estiver em sincronia, suprime o upload desnecessário
	if strings.EqualFold(localMD5, remoteFile.Md5Checksum) {
		return &SyncReport{
			FileID:    fileID,
			LocalMD5:  localMD5,
			RemoteMD5: remoteFile.Md5Checksum,
			SizeBytes: size,
			Status:    "IN_SYNC",
			Message:   "O PDF local já é criptograficamente idêntico ao arquivo na nuvem. Upload suprimido.",
		}, nil
	}

	f, err := os.Open(localPath)
	if err != nil {
		return nil, fmt.Errorf("falha ao reabrir arquivo PDF para upload: %w", err)
	}
	defer f.Close()

	mediaMeta := &drive.File{}

	updatedFile, err := driveSvc.Files.Update(fileID, mediaMeta).
		KeepRevisionForever(true).
		Media(f).
		Fields("id, name, md5Checksum, size, modifiedTime").
		Do()
	if err != nil {
		return nil, fmt.Errorf("falha ao transmitir mídia atualizada para o Drive: %w", err)
	}

	// Validação estrita pós-upload
	if !strings.EqualFold(localMD5, updatedFile.Md5Checksum) {
		return nil, fmt.Errorf("alerta crítico de integridade: hash local (%s) difere do hash remoto pós-upload (%s)", localMD5, updatedFile.Md5Checksum)
	}

	return &SyncReport{
		FileID:    fileID,
		LocalMD5:  localMD5,
		RemoteMD5: updatedFile.Md5Checksum,
		SizeBytes: size,
		Status:    "UPDATED",
		Message:   "SUCESSO: Upload concluído com sucesso e integridade MD5 verificada com retenção permanente de revisão.",
	}, nil
}
