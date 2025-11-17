package mTLS

import (
	"crypto/tls"
	"crypto/x509"
	"net/http"
)

// Manager 管理 mTLS 證書
type Manager struct {
	caCertPool *x509.CertPool
}

// NewManager 建立新的 mTLS manager（簡化版本）
func NewManager() (*Manager, error) {
	// 簡化：這裡應該載入 CA 證書
	// 實際環境中應該從檔案或 secret 載入
	caCertPool := x509.NewCertPool()
	
	// TODO: 載入 CA 證書
	// caCert, err := os.ReadFile("ca.crt")
	// if err != nil {
	//     return nil, err
	// }
	// caCertPool.AppendCertsFromPEM(caCert)

	return &Manager{
		caCertPool: caCertPool,
	}, nil
}

// VerifyClient 驗證 client certificate
func (m *Manager) VerifyClient(req *http.Request) bool {
	// 簡化版本：檢查是否有 client certificate
	if req.TLS == nil {
		return false
	}

	if len(req.TLS.PeerCertificates) == 0 {
		return false
	}

	// 驗證證書（簡化：實際應該檢查 CA 簽名）
	// opts := x509.VerifyOptions{
	//     Roots: m.caCertPool,
	// }
	// _, err := req.TLS.PeerCertificates[0].Verify(opts)
	// return err == nil

	// 簡化：只要有證書就通過
	return true
}

// GetTLSConfig 取得 TLS 配置
func (m *Manager) GetTLSConfig() *tls.Config {
	return &tls.Config{
		ClientAuth: tls.RequireAndVerifyClientCert,
		ClientCAs:  m.caCertPool,
		MinVersion: tls.VersionTLS12,
	}
}

