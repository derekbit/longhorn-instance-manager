package cmd

import (
	"crypto/tls"
	"net"
	"os"
	"os/signal"
	"path/filepath"
	"strconv"
	"syscall"
	"time"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
	"google.golang.org/grpc"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/reflection"

	"github.com/longhorn/longhorn-instance-manager/pkg/disk"
	"github.com/longhorn/longhorn-instance-manager/pkg/health"
	rpc "github.com/longhorn/longhorn-instance-manager/pkg/imrpc"
	"github.com/longhorn/longhorn-instance-manager/pkg/instance"
	"github.com/longhorn/longhorn-instance-manager/pkg/process"
	"github.com/longhorn/longhorn-instance-manager/pkg/proxy"
	"github.com/longhorn/longhorn-instance-manager/pkg/types"
	"github.com/longhorn/longhorn-instance-manager/pkg/util"
)

func StartCmd() cli.Command {
	return cli.Command{
		Name: "daemon",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "listen",
				Value: "tcp://localhost:8500",
				Usage: "specifies the server endpoint to listen on supported protocols are 'tcp' and 'unix'. The proxy server will be listening on the next port.",
			},
			cli.StringFlag{
				Name:  "logs-dir",
				Value: "/var/log/instances",
			},
			cli.StringFlag{
				Name:  "port-range",
				Value: "10000-30000",
			},
			cli.BoolFlag{
				Name:  "spdk-enabled",
				Usage: "enable SPDK support",
			},
		},
		Action: func(c *cli.Context) {
			if err := start(c); err != nil {
				logrus.WithError(err).Fatal("Failed to run start command")
			}
		},
	}
}

func cleanup(pm *process.Manager) {
	logrus.Info("Trying to gracefully shut down Instance Manager")
	pmResp, err := pm.ProcessList(nil, &rpc.ProcessListRequest{})
	if err != nil {
		logrus.WithError(err).Error("Failed to list processes before shutdown")
		return
	}
	for _, p := range pmResp.Processes {
		pm.ProcessDelete(nil, &rpc.ProcessDeleteRequest{
			Name: p.Spec.Name,
		})
	}

	for i := 0; i < types.WaitCount; i++ {
		pmResp, err := pm.ProcessList(nil, &rpc.ProcessListRequest{})
		if err != nil {
			logrus.WithError(err).Error("Failed to list instance processes when shutting down")
			break
		}
		if len(pmResp.Processes) == 0 {
			logrus.Info("Shutdown all instance processes successfully")
			break
		}
		time.Sleep(types.WaitInterval)
	}

	logrus.Error("Failed to clean up all processes for Instance Manager graceful shutdown")
}

func start(c *cli.Context) (err error) {
	listen := c.String("listen")
	logsDir := c.String("logs-dir")
	portRange := c.String("port-range")
	spdkEnabled := c.Bool("spdk-enabled")

	if err := util.SetUpLogger(logsDir); err != nil {
		return err
	}

	// setup tls config
	var tlsConfig *tls.Config
	tlsDir := c.GlobalString("tls-dir")
	if tlsDir != "" {
		tlsConfig, err = util.LoadServerTLS(
			filepath.Join(tlsDir, "ca.crt"),
			filepath.Join(tlsDir, "tls.crt"),
			filepath.Join(tlsDir, "tls.key"),
			"longhorn-backend.longhorn-system")
		if err != nil {
			logrus.WithError(err).Warnf("Failed to add TLS key pair from %v", tlsDir)
		}
	}

	if tlsConfig != nil {
		logrus.Info("Creating gRPC server with mtls auth")
	} else {
		logrus.Info("Creating gRPC server with no auth")
	}

	shutdownCh := make(chan error)

	processManagerServiceAddress, proxyServiceAddress, diskServiceAddress, instanceServiceAddress, err := getServiceAddresses(listen)
	if err != nil {
		return err
	}

	// Start instance server
	instanceRpcServer, instanceRpcListener, err := setupInstanceGrpcServer(logsDir,
		instanceServiceAddress, processManagerServiceAddress, diskServiceAddress, tlsConfig, spdkEnabled, shutdownCh)
	if err != nil {
		return err
	}
	go func() {
		if err := instanceRpcServer.Serve(instanceRpcListener); err != nil {
			logrus.WithError(err).Error("Stopping instance gRPC server")
		}
		// graceful shutdown before exit
		close(shutdownCh)
	}()
	logrus.Infof("Instance Manager instance gRPC server listening to %v", instanceServiceAddress)

	// Start proxy server
	proxyRpcServer, proxyRpcListener, err := setupProxyGrpcServer(logsDir, proxyServiceAddress, tlsConfig, shutdownCh)
	if err != nil {
		return err
	}
	go func() {
		if err := proxyRpcServer.Serve(proxyRpcListener); err != nil {
			logrus.WithError(err).Error("Stopping proxy gRPC server")
		}
		// graceful shutdown before exit
		close(shutdownCh)
	}()
	logrus.Infof("Instance Manager proxy gRPC server listening to %v", proxyServiceAddress)

	// Start process manager server
	pm, pmRpcServer, pmRpcListener, err := setupProcessManagerGrpcServer(portRange, logsDir, processManagerServiceAddress, tlsConfig, shutdownCh)
	if err != nil {
		return err
	}
	go func() {
		if err := pmRpcServer.Serve(pmRpcListener); err != nil {
			logrus.WithError(err).Error("Stopping process manager gRPC server")
		}
		// graceful shutdown before exit
		cleanup(pm)
		close(shutdownCh)
	}()
	logrus.Infof("Instance Manager process manager gRPC server listening to %v", listen)

	// Start disk server
	diskRpcServer, diskRpcListener, err := setupDiskGrpcServer(diskServiceAddress, tlsConfig, shutdownCh)
	if err != nil {
		return err
	}
	go func() {
		if err := diskRpcServer.Serve(diskRpcListener); err != nil {
			logrus.WithError(err).Error("Stopping disk gRPC server")
		}
		// graceful shutdown before exit
		close(shutdownCh)
	}()
	logrus.Infof("Instance Manager disk gRPC server listening to %v", diskServiceAddress)

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		sig := <-sigs
		logrus.Infof("Instance Manager received %v to exit", sig)
		pmRpcServer.Stop()
	}()

	return <-shutdownCh
}

func getServiceAddresses(listen string) (processManagerServiceAddress, proxyServiceAddress, diskServiceAddress, instanceServiceAddress string, err error) {
	host, port, err := net.SplitHostPort(listen)
	if err != nil {
		return "", "", "", "", err
	}

	intPort, err := strconv.Atoi(port)
	if err != nil {
		return "", "", "", "", err
	}

	return net.JoinHostPort(host, strconv.Itoa(intPort)),
		net.JoinHostPort(host, strconv.Itoa(intPort+1)),
		net.JoinHostPort(host, strconv.Itoa(intPort+2)),
		net.JoinHostPort(host, strconv.Itoa(intPort+3)),
		nil
}

func setupDiskGrpcServer(listen string, tlsConfig *tls.Config, shutdownCh chan error) (*grpc.Server, net.Listener, error) {
	ds, err := disk.NewServer(shutdownCh)
	if err != nil {
		return nil, nil, err
	}
	hc := health.NewDiskHealthCheckServer(ds)

	rpcServer, rpcListener, err := util.NewServer(listen, tlsConfig,
		grpc.KeepaliveEnforcementPolicy(keepalive.EnforcementPolicy{
			MinTime:             10 * time.Second,
			PermitWithoutStream: true,
		}),
	)
	if err != nil {
		return nil, nil, errors.Wrap(err, "failed to setup disk gRPC server")
	}

	rpc.RegisterDiskServiceServer(rpcServer, ds)
	healthpb.RegisterHealthServer(rpcServer, hc)
	reflection.Register(rpcServer)

	return rpcServer, rpcListener, nil
}

func setupProxyGrpcServer(logsDir, listen string, tlsConfig *tls.Config, shutdownCh chan error) (*grpc.Server, net.Listener, error) {
	// TODO: skip proxy for replica instance manager pod
	proxy, err := proxy.NewProxy(logsDir, shutdownCh)
	if err != nil {
		return nil, nil, err
	}
	hc := health.NewProxyHealthCheckServer(proxy)

	rpcProxyServer, rpcProxyListener, err := util.NewServer(listen, tlsConfig,
		grpc.KeepaliveEnforcementPolicy(keepalive.EnforcementPolicy{
			MinTime:             10 * time.Second,
			PermitWithoutStream: true,
		}),
	)
	if err != nil {
		return nil, nil, errors.Wrap(err, "failed to setup proxy gRPC server")
	}

	rpc.RegisterProxyEngineServiceServer(rpcProxyServer, proxy)
	healthpb.RegisterHealthServer(rpcProxyServer, hc)
	reflection.Register(rpcProxyServer)

	return rpcProxyServer, rpcProxyListener, nil
}

func setupProcessManagerGrpcServer(portRange, logsDir, listen string, tlsConfig *tls.Config, shutdownCh chan error) (*process.Manager, *grpc.Server, net.Listener, error) {
	pm, err := process.NewManager(portRange, logsDir, shutdownCh)
	if err != nil {
		return nil, nil, nil, err
	}
	hc := health.NewHealthCheckServer(pm)

	rpcServer, rpcListener, err := util.NewServer(listen, tlsConfig,
		grpc.KeepaliveEnforcementPolicy(keepalive.EnforcementPolicy{
			MinTime:             10 * time.Second,
			PermitWithoutStream: true,
		}),
	)
	if err != nil {
		return nil, nil, nil, errors.Wrap(err, "failed to setup process manager gRPC server")
	}

	rpc.RegisterProcessManagerServiceServer(rpcServer, pm)
	healthpb.RegisterHealthServer(rpcServer, hc)
	reflection.Register(rpcServer)

	return pm, rpcServer, rpcListener, nil
}

func setupInstanceGrpcServer(logsDir, listen, processManagerServiceAddress, diskServiceAddress string, tlsConfig *tls.Config, spdkEnabled bool, shutdownCh chan error) (*grpc.Server, net.Listener, error) {
	srv, err := instance.NewServer(logsDir, processManagerServiceAddress, diskServiceAddress, spdkEnabled, shutdownCh)
	if err != nil {
		return nil, nil, err
	}
	hc := health.NewInstanceHealthCheckServer(srv)

	rpcServer, rpcListener, err := util.NewServer(listen, tlsConfig,
		grpc.KeepaliveEnforcementPolicy(keepalive.EnforcementPolicy{
			MinTime:             10 * time.Second,
			PermitWithoutStream: true,
		}),
	)
	if err != nil {
		return nil, nil, errors.Wrap(err, "failed to setup instance gRPC server")
	}

	rpc.RegisterInstanceServiceServer(rpcServer, srv)
	healthpb.RegisterHealthServer(rpcServer, hc)
	reflection.Register(rpcServer)

	return rpcServer, rpcListener, nil
}
