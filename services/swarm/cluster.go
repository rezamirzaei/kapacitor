package swarm

import (
	"log"
	"sync/atomic"

	"github.com/influxdata/kapacitor/services/swarm/client"
	"github.com/pkg/errors"
)

type Cluster struct {
	configValue atomic.Value // Config
	client      client.Client
	logger      *log.Logger
}

func NewCluster(c Config, l *log.Logger) (*Cluster, error) {
	clientConfig, err := c.ClientConfig()
	if err != nil {
		return nil, errors.Wrap(err, "failed to create k8s client config")
	}
	cli, err := client.New(clientConfig)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create k8s client")
	}

	s := &Cluster{
		client: cli,
		logger: l,
	}
	s.configValue.Store(c)
	return s, nil
}

func (s *Cluster) Update(c Config) error {
	s.configValue.Store(c)
	clientConfig, err := c.ClientConfig()
	if err != nil {
		return errors.Wrap(err, "failed to create swarm client config")
	}
	return s.client.Update(clientConfig)
}

func (s *Cluster) Test() error {
	cli, err := s.Client()
	if err != nil {
		return errors.Wrap(err, "failed to get client")
	}
	version, err := cli.Version()
	if err != nil {
		return errors.Wrap(err, "failed to query server version")
	}
	if version == "" {
		return errors.New("got empty version from server")
	}
	return nil
}

func (s *Cluster) Client() (client.Client, error) {
	config := s.configValue.Load().(Config)
	if !config.Enabled {
		return nil, errors.New("service is not enabled")
	}
	return s.client, nil
}
