// SiGG-Satellite-Network-SII  //

package kafka

import (
	"fmt"
	"os"
	"time"

	"crypto/tls"
	"crypto/x509"

	"github.com/Shopify/sarama"

	"github.com/apache/skywalking-satellite/internal/pkg/log"
)

// loadConfig use the client params to build the kafka producer config.
func (c *Client) loadConfig() (*sarama.Config, error) {
	cfg := sarama.NewConfig()
	cfg.Producer.Return.Successes = true
	cfg.Producer.Return.Errors = true
	cfg.Producer.Idempotent = c.IdempotentWrites
	cfg.Producer.RequiredAcks = sarama.RequiredAcks(c.RequiredAcks)
	cfg.Producer.Compression = sarama.CompressionCodec(c.CompressionCodec)
	if c.ProducerMaxRetry > 0 {
		cfg.Producer.Retry.Max = c.ProducerMaxRetry
	}
	if c.MetaMaxRetry > 0 {
		cfg.Metadata.Retry.Max = c.MetaMaxRetry
	}
	if c.RetryBackoff > 0 {
		cfg.Producer.Retry.Backoff = time.Millisecond * time.Duration(c.RetryBackoff)
	}
	if c.RefreshPeriod > 0 {
		cfg.Metadata.RefreshFrequency = time.Duration(c.RefreshPeriod) * time.Minute
	}
	if c.MaxMessageBytes > 0 {
		cfg.Producer.MaxMessageBytes = c.MaxMessageBytes
	}
	if c.ClientID != "" {
		cfg.ClientID = c.ClientID
	}
	if c.Version != "" {
		if version, err := sarama.ParseKafkaVersion(c.Version); err != nil {
			log.Logger.Errorf("error in parsing the kafka version, the kafka version would be set as default value: %v", err)
		} else {
			cfg.Version = version
		}
	}
	cfg.Net.TLS.Enable = c.EnableTLS
	if c.EnableTLS {
		configTLS, err := c.configTLS()
		if err != nil {
			return nil, err
		}
		cfg.Net.TLS.Config = configTLS
	}
	return cfg, nil
}

// configTLS loads and parse the TLS configs.
func (c *Client) configTLS() (tc *tls.Config, tlsErr error) {
	if err := checkTLSFile(c.CaPemPath); err != nil {
		return nil, err
	}
	if err := checkTLSFile(c.ClientKeyPath); err != nil {
		return nil, err
	}
	if err := checkTLSFile(c.ClientPemPath); err != nil {
		return nil, err
	}
	tlsConfig := new(tls.Config)
	tlsConfig.Renegotiation = tls.RenegotiateNever
	tlsConfig.InsecureSkipVerify = c.InsecureSkipVerify
	caPem, err := os.ReadFile(c.CaPemPath)
	if err != nil {
		return nil, err
	}
	certPool := x509.NewCertPool()
	certPool.AppendCertsFromPEM(caPem)
	tlsConfig.RootCAs = certPool

	clientPem, err := tls.LoadX509KeyPair(c.ClientPemPath, c.ClientKeyPath)
	if err != nil {
		return nil, err
	}
	tlsConfig.Certificates = []tls.Certificate{clientPem}
	return tlsConfig, nil
}

// checkTLSFile checks the TLS files.
func checkTLSFile(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	stat, err := file.Stat()
	if err != nil {
		return err
	}
	if stat.Size() == 0 {
		return fmt.Errorf("the TLS file is illegal: %s", path)
	}
	return nil
}
