package harness

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type Cluster int

const (
	ClusterKind Cluster = iota
	ClusterMinikube
	ClusterRancherDesktop
	ClusterOther
)

var clusterPrefixes = map[string]Cluster{
	"kind-":            ClusterKind,
	"minikube":         ClusterMinikube,
	"rancher-desktop":  ClusterRancherDesktop,
}

func DetectCluster() (Cluster, bool) {
	ctx, err := RunQuiet("kubectl", "config", "current-context")
	if err != nil {
		return 0, false
	}

	for prefix, cluster := range clusterPrefixes {
		if strings.HasPrefix(ctx, prefix) {
			return cluster, true
		}
	}

	return ClusterOther, true
}

func LoadImages(images []string) error {
	cluster, ok := DetectCluster()
	if !ok {
		return nil
	}

	cli, err := ContainerCLI()
	if err != nil {
		return err
	}

	switch cluster {
	case ClusterMinikube:
		for _, img := range images {
			if cli == "podman" {
				if err := loadPodmanImageToMinikube(img); err != nil {
					return err
				}
			} else {
				if err := Run("minikube", "image", "load", img); err != nil {
					return err
				}
			}
		}

	case ClusterKind:
		for _, img := range images {
			if err := Run("kind", "load", "docker-image", img); err != nil {
				return err
			}
		}

	case ClusterRancherDesktop, ClusterOther:
	}

	return nil
}

func loadPodmanImageToMinikube(img string) error {
	tmpDir, err := os.MkdirTemp("", "harness-img-*")
	if err != nil {
		return fmt.Errorf("create temp dir: %w", err)
	}
	defer os.RemoveAll(tmpDir)

	tarPath := filepath.Join(tmpDir, "image.tar")
	if err := Run("podman", "save", img, "-o", tarPath); err != nil {
		return err
	}
	return Run("minikube", "image", "load", tarPath)
}