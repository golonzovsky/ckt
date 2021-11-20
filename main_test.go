package main

import (
	"io/ioutil"
	"k8s.io/client-go/tools/clientcmd"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
	"os"
	"testing"
)

func TestUnsetTokenString(t *testing.T) {
	conf := &clientcmdapi.Config{
		Kind:       "Config",
		APIVersion: "v1",
		Clusters: map[string]*clientcmdapi.Cluster{
			"minikube":   {Server: "https://192.168.99.100:8443"},
			"my-cluster": {Server: "https://192.168.0.1:3434"},
		},
		Contexts: map[string]*clientcmdapi.Context{
			"minikube":   {AuthInfo: "minikube", Cluster: "minikube"},
			"gcp":        {AuthInfo: "gke", Cluster: "my-cluster"},
		},
		CurrentContext: "gcp",
		AuthInfos: map[string]*clientcmdapi.AuthInfo{
			"gke": {AuthProvider: &clientcmdapi.AuthProviderConfig{Name: "gcp", Config: map[string]string{"access-token": "tokenValue"}}},
		},
	}

	doTestRemoveToken(t, conf)

}

func doTestRemoveToken(t *testing.T, testConfig *clientcmdapi.Config){
	fakeKubeFile, err := ioutil.TempFile(os.TempDir(), "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer os.Remove(fakeKubeFile.Name())
	err = clientcmd.WriteToFile(*testConfig, fakeKubeFile.Name())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	os.Args = append(os.Args, "--kubeconfig=" + fakeKubeFile.Name())
	main()

	config, err := clientcmd.LoadFromFile(fakeKubeFile.Name())
	if err != nil {
		t.Fatalf("unexpected error loading kubeconfig file: %v", err)
	}
	for _, user := range config.AuthInfos {
		if user.AuthProvider == nil || user.AuthProvider.Config == nil {
			continue
		}
		authConfig := user.AuthProvider.Config
		_, actualExist := authConfig["access-token"]
		if actualExist {
			t.Errorf("Token present despite expected to be removed")
		}
	}
}
