package operator

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/kris-nova/logger"
	"github.com/kubicorn/kubicorn/apis/cluster"

	apiextensionsclient "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/kube-deploy/cluster-api/client"
	"k8s.io/kube-deploy/cluster-api/util"
)

type ServiceConfiguration struct {
	KubeConfigContent string
	//cloudProviderName string
	//CloudProvider     CloudProvider
}

func UpdateCRDNumberInstances(n int) error {

	// Hacky way to ensure our config is set
	kubeConfigContent := os.Getenv("KUBECONFIG_CONTENT")
	if kubeConfigContent == "" {
		logger.Critical("Missing environmental variable [KUBECONFIG_CONTENT]")
		return fmt.Errorf("Missing environmental variable [KUBECONFIG_CONTENT]")
	}

	cfg := &ServerConfiguration{
		KubeConfigContent: kubeConfigContent,
	}

	cm, err := getClientMeta(cfg)
	if err != nil {
		return err
	}
	listOptions := metav1.ListOptions{}
	machines, err := cm.client.Machines().List(listOptions)
	if err != nil {
		return fmt.Errorf("Unable to list machines: %v", err)

	}

	totalMachines := len(machines)

	return nil
}

func (s *ServiceConfiguration) GetFilePath() (string, error) {
	file, err := ioutil.TempFile(os.TempDir(), "kubicorn")
	if err != nil {
		return "", err
	}
	err = ioutil.WriteFile(file.Name(), []byte(s.KubeConfigContent), 0755)
	if err != nil {
		return "", err
	}
	return file.Name(), nil
}

type crdClientMeta struct {
	client    *client.ClusterAPIV1Alpha1Client
	clientset *apiextensionsclient.Clientset
}

func getClientMeta(cfg *ServerConfiguration) (*crdClientMeta, error) {
	kubeConfigPath, err := cfg.GetFilePath()
	if err != nil {
		return nil, err
	}
	client, err := util.NewApiClient(kubeConfigPath)
	if err != nil {
		return nil, err
	}
	config, err := clientcmd.BuildConfigFromFlags("", kubeConfigPath)
	if err != nil {
		return nil, err
	}
	cs, err := apiextensionsclient.NewForConfig(config)
	if err != nil {
		return nil, err
	}
	clientMeta := &crdClientMeta{
		client:    client,
		clientset: cs,
	}
	return clientMeta, nil
}

func getProviderConfig(providerConfig string) *cluster.MachineProviderConfig {
	//logger.Info(providerConfig)
	mp := cluster.MachineProviderConfig{
		ServerPool: &cluster.ServerPool{},
	}
	json.Unmarshal([]byte(providerConfig), &mp)
	return &mp
}
