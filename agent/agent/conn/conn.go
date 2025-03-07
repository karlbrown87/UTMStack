package conn

import (
	"crypto/tls"
	"fmt"
	"time"

	"github.com/threatwinds/logger"
	"github.com/utmstack/UTMStack/agent/agent/configuration"
	"github.com/utmstack/UTMStack/agent/agent/utils"
	grpc "google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

const (
	maxMessageSize        = 1024 * 1024 * 1024
	maxConnectionAttempts = 3
	initialReconnectDelay = 10 * time.Second
	maxReconnectDelay     = 60 * time.Second
)

func ConnectToServer(cnf *configuration.Config, h *logger.Logger, addrs, port string) (*grpc.ClientConn, error) {
	connectionAttemps := 0
	reconnectDelay := initialReconnectDelay

	// Connect to the gRPC server
	serverAddress := addrs + ":" + port
	var conn *grpc.ClientConn
	var err error

	for {
		if connectionAttemps >= maxConnectionAttempts {
			return nil, fmt.Errorf("failed to connect to Server")
		}

		h.Info("trying to connect to Server...")
		var opts grpc.DialOption
		if !cnf.SkipCertValidation {
			creds, err := credentials.NewClientTLSFromFile(configuration.GetCaPath(), "")
			if err != nil {
				return nil, fmt.Errorf("failed to load CA trust certificate: %v", err)
			}
			opts = grpc.WithTransportCredentials(creds)
		} else {
			tlsConfig := &tls.Config{InsecureSkipVerify: true}
			creds := credentials.NewTLS(tlsConfig)
			opts = grpc.WithTransportCredentials(creds)
		}

		conn, err = grpc.NewClient(serverAddress, opts, grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(maxMessageSize)))
		if err != nil {
			connectionAttemps++
			h.Info("error connecting to Server, trying again in %.0f seconds", reconnectDelay.Seconds())
			time.Sleep(reconnectDelay)
			reconnectDelay = utils.IncrementReconnectDelay(reconnectDelay, maxReconnectDelay)
			continue
		}

		break
	}

	return conn, nil
}
