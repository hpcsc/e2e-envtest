//go:build e2e
// +build e2e

package e2e

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/require"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/clientcmd/api"
	"os"
	"os/exec"
	"sigs.k8s.io/controller-runtime/pkg/envtest"
	"testing"
)

var (
	testenv        *envtest.Environment
	client         dynamic.Interface
	kubeconfigPath = os.Getenv("E2E_KUBECONFIG")
)

func setup(t *testing.T) {
	require.NoError(t, os.Setenv("KUBECONFIG", kubeconfigPath))
	testenv = &envtest.Environment{}

	cfg, err := testenv.Start()
	require.NoError(t, err)

	// create a dynamic client, to be used in this e2e test to setup/verification
	client, err = dynamic.NewForConfig(cfg)
	require.NoError(t, err)

	// write kube rest config to filesystem so that main program is able to access it
	require.NoError(t, writeRestCfgAsKubeConfigFile(cfg, os.Getenv("KUBECONFIG")))
}

func tearDown(t *testing.T) {
	require.NoError(t, testenv.Stop())
	require.NoError(t, os.Unsetenv("KUBECONFIG"))
	require.NoError(t, os.Remove(kubeconfigPath))
}

func TestE2E(t *testing.T) {
	setup(t)
	defer tearDown(t)

	t.Run("filter namespaces by created-by label", func(t *testing.T) {
		createTestNamespace(t, "ns-1")
		createTestNamespace(t, "ns-2")

		cmd := exec.Command("../bin/e2e-envtest", "-created-by", "e2e")
		out, err := cmd.CombinedOutput()

		fmt.Printf("=============== Program output ==============\n%s\n====================================\n", string(out))
		require.NoError(t, err)

		require.Contains(t, string(out), "found 2 namespaces: ns-1, ns-2")
	})
}

func createTestNamespace(t *testing.T, name string) {
	ns := unstructured.Unstructured{}
	ns.SetGroupVersionKind(schema.FromAPIVersionAndKind("v1", "Namespace"))
	ns.SetName(name)
	ns.SetLabels(map[string]string{
		"created-by": "e2e",
	})

	_, err := client.Resource(schema.GroupVersionResource{
		Group:    "",
		Version:  "v1",
		Resource: "namespaces",
	}).Create(context.TODO(), &ns, v1.CreateOptions{})

	require.NoError(t, err)
}

func writeRestCfgAsKubeConfigFile(cfg *rest.Config, kubeconfigPath string) error {
	ns := "default"
	contextName := "default-context"
	clusterName := "default-cluster"
	userName := "default-user"
	clusters := map[string]*api.Cluster{
		clusterName: {
			Server:                   cfg.Host,
			CertificateAuthorityData: cfg.CAData,
		},
	}

	contexts := map[string]*api.Context{
		contextName: {
			Cluster:   clusterName,
			Namespace: ns,
			AuthInfo:  userName,
		},
	}

	authinfos := map[string]*api.AuthInfo{
		userName: {
			ClientCertificateData: cfg.CertData,
			ClientKeyData:         cfg.KeyData,
		},
	}

	clientConfig := api.Config{
		Kind:           "Config",
		APIVersion:     "v1",
		Clusters:       clusters,
		Contexts:       contexts,
		CurrentContext: contextName,
		AuthInfos:      authinfos,
	}
	return clientcmd.WriteToFile(clientConfig, kubeconfigPath)
}
