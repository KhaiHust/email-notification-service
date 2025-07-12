package impl

import (
	"crypto/tls"
	"github.com/KhaiHust/email-notification-service/core/msg/constant"
	coreProperties "github.com/KhaiHust/email-notification-service/core/msg/properties"
	"github.com/Shopify/sarama"
	"github.com/golibs-starter/golib-message-bus/kafka/core"
	"github.com/golibs-starter/golib-message-bus/kafka/utils"
	"github.com/pkg/errors"
)

type CommonProperties interface {
	GetClientId() string
	GetSecurityProtocol() string
	GetTls() *coreProperties.Tls
	GetSasl() *coreProperties.Sasl
}

func CreateCommonSaramaConfig(version string, props CommonProperties) (*sarama.Config, error) {
	config := sarama.NewConfig()
	configVersion, err := sarama.ParseKafkaVersion(version)
	if err != nil {
		return nil, errors.WithMessage(err, "Error parsing Kafka version")
	}
	config.Version = configVersion

	if props.GetClientId() != "" {
		config.ClientID = props.GetClientId()
	}

	if props.GetSecurityProtocol() == core.SecurityProtocolTls {
		tlsConfig, err := createTlsConfiguration(props.GetTls())
		if err != nil {
			return nil, errors.WithMessage(err, "Error when create tls properties")
		}
		config.Net.TLS.Enable = true
		config.Net.TLS.Config = tlsConfig
	}
	if props.GetSecurityProtocol() == constant.SecurityProtocolSaslSSL {
		saslProps := props.GetSasl()

		config.Net.SASL.Enable = true
		config.Net.SASL.User = saslProps.Username
		config.Net.SASL.Password = saslProps.Password

		switch saslProps.Mechanism {
		case constant.SecurityMechanismSaslPlain:
			config.Net.SASL.Mechanism = sarama.SASLTypePlaintext
		case constant.SecurityMechanismSaslScramSha256:
			config.Net.SASL.Mechanism = sarama.SASLTypeSCRAMSHA256
		case constant.SecurityMechanismSaslScramSha512:
			config.Net.SASL.Mechanism = sarama.SASLTypeSCRAMSHA512
		default:
			return nil, errors.Errorf("Unsupported SASL mechanism: %s", saslProps.Mechanism)
		}
		config.Net.TLS.Enable = true
		tlsConfig := &tls.Config{
			InsecureSkipVerify: true,
			ClientAuth:         0,
		}
		config.Net.TLS.Config = tlsConfig
	}
	return config, nil
}

func createTlsConfiguration(tlsProps *coreProperties.Tls) (*tls.Config, error) {
	if tlsProps == nil {
		return nil, errors.New("Tls properties not found when using SecurityProtocol=TLS")
	}
	tlsConfig, err := utils.NewTLSConfig(
		tlsProps.CertFileLocation,
		tlsProps.KeyFileLocation,
		tlsProps.CaFileLocation,
	)
	if err != nil {
		return nil, errors.WithMessage(err, "Error when load TLS properties")
	}
	tlsConfig.InsecureSkipVerify = tlsProps.InsecureSkipVerify
	return tlsConfig, nil
}
