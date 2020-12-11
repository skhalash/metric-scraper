package main

import (
	"context"
	"os"

	"github.com/go-kit/kit/log"
	"github.com/prometheus/common/model"
	"github.com/prometheus/prometheus/config"
	"github.com/prometheus/prometheus/discovery"
	"github.com/prometheus/prometheus/discovery/kubernetes"
	"github.com/prometheus/prometheus/pkg/relabel"
	"github.com/prometheus/prometheus/scrape"
)

func main() {
	discoveryManager := discovery.NewManager(context.Background(), log.NewLogfmtLogger(os.Stdout))
	discoveryManager.ApplyConfig(createDiscoveryConfigs())
	go discoveryManager.Run()

	scrapeManager := scrape.NewManager(log.NewLogfmtLogger(os.Stdout), &noOpStore{})
	scrapeManager.ApplyConfig(createScrapeConfig())
	scrapeManager.Run(discoveryManager.SyncCh())
}

func createDiscoveryConfigs() map[string]discovery.Configs {
	discoveryConfigs := make(map[string]discovery.Configs)
	sdConfig := discovery.Configs{}

	// for _, v := range []string{"node", "endpoints", "service", "pod", "ingress"} {
	for _, v := range []string{"endpoints"} {
		k8sCfg := &kubernetes.SDConfig{
			Role: kubernetes.Role(v),
		}
		sdConfig = append(sdConfig, k8sCfg)
	}
	discoveryConfigs["k8s"] = sdConfig

	return discoveryConfigs
}

func createScrapeConfig() *config.Config {
	// These scrape endpoints are hard-coded, but should be discovered as well
	config := &config.Config{
		ScrapeConfigs: []*config.ScrapeConfig{{
			JobName:     "k8s",
			HonorLabels: true,
			Params: map[string][]string{
				"module": {"http_2xx"},
				"target": {"dex-service.kyma-system:5556/healthz"},
			},
			MetricsPath: "/probe",
			Scheme:      "http",
			RelabelConfigs: []*relabel.Config{
				{
					SourceLabels: []model.LabelName{"__meta_kubernetes_service_label_app_kubernetes_io_instance"},
					Regex:        relabel.MustNewRegexp("blackbox"),
					Replacement:  "$1",
					Action:       relabel.Action("keep"),
				},
				{
					SourceLabels: []model.LabelName{"__meta_kubernetes_service_label_app_kubernetes_io_name"},
					Regex:        relabel.MustNewRegexp("blackbox-exporter"),
					Replacement:  "$1",
					Action:       relabel.Action("keep"),
				},
				{
					SourceLabels: []model.LabelName{"__meta_kubernetes_endpoint_port_name"},
					Regex:        relabel.MustNewRegexp("http"),
					Replacement:  "$1",
					Action:       relabel.Action("keep"),
				},
			},
		}},
	}

	// Set some default scrape intervals and timeouts
	config.ScrapeConfigs[0].ScrapeInterval.Set("10s")
	config.ScrapeConfigs[0].ScrapeTimeout.Set("10s")

	return config
}
