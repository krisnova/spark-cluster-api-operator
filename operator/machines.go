package operator

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/kubicorn/kubicorn/pkg/namer"

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

	cfg := &ServiceConfiguration{
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

	totalMachines := len(machines.Items)
	if totalMachines != n {
		logger.Always("Total Machines [%d] Expected Machines [%d]")
		for totalMachines != n {
			if totalMachines < n {
				err := addMachine(cm)
				if err != nil {
					return err
				}
			} else if totalMachines > n {
				removeMachine(cm)
				if err != nil {
					return err
				}
			} else {
				break
			}
		}
	}

	return nil
}

func addMachine(cm *crdClientMeta) error {
	listOptions := metav1.ListOptions{}
	machines, err := cm.client.Machines().List(listOptions)
	if err != nil {
		return fmt.Errorf("Unable to list machines: %v", err)

	}
	if len(machines.Items) < 3 {
		return fmt.Errorf("Unable to find base machine")
	}
	base := machines.Items[2] // Grab the third machine

	newMachine := base
	newMachine.Name = namer.RandomName()
	_, err = cm.client.Machines().Create(&newMachine)
	return err
}

func removeMachine(cm *crdClientMeta) error {
	listOptions := metav1.ListOptions{}
	machines, err := cm.client.Machines().List(listOptions)
	if err != nil {
		return fmt.Errorf("Unable to list machines: %v", err)

	}
	if len(machines.Items) == 3 {
		// Always leave one hanging around
		return nil
	} else if len(machines.Items) > 3 {
		machineToDelete := machines.Items[3]
		_, err = cm.client.Machines().Create(&machineToDelete)
		return err
	}
	return fmt.Errorf("Invalid length of machines")
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

func getClientMeta(cfg *ServiceConfiguration) (*crdClientMeta, error) {
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
