package main

import (
	"flag"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	"log"
	"path/filepath"
)

func main() {
	var kubeconfig *string
	if home := homedir.HomeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	flag.Parse()

	config, err := clientcmd.LoadFromFile(*kubeconfig)
	if err != nil {
		log.Fatalf("unexpected error loading kubeconfig file: %v", err)
	}

	for _, user := range config.AuthInfos {
		if user.AuthProvider == nil || user.AuthProvider.Config == nil{
			continue
		}
		authConfig := user.AuthProvider.Config
		delete(authConfig, "access-token")
	}

	if err := clientcmd.WriteToFile(*config, *kubeconfig); err != nil {
		log.Fatalf("unexpected error saving kubeconfig file: %v", err)
	}
}
