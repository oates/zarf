package k8s

import (
	"bytes"
	"fmt"
	"github.com/go-logr/logr"
	"io"
	"io/ioutil"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/klog/v2"
	"os"
	"regexp"

	"github.com/defenseunicorns/zarf/cli/internal/message"
	"github.com/defenseunicorns/zarf/cli/internal/template"
	"github.com/defenseunicorns/zarf/cli/internal/utils"
	"github.com/go-logr/logr/funcr"
	kubeyaml "k8s.io/apimachinery/pkg/util/yaml"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"sigs.k8s.io/yaml"
)

// ImageSwap Pre-compute all the replacements for the embedded registry
type ImageSwap struct {
	find    string
	replace string
}

func init() {
	klog.SetLogger(generateLogShim())
}

func getRestConfig() *rest.Config {
	homePath, err := os.UserHomeDir()
	if err != nil {
		message.Fatal(nil, "Unable to load the current user's home directory")
	}

	// use the current context in kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", homePath+"/.kube/config")
	if err != nil {
		message.Fatalf(err, "Unable to connect to the K8s cluster")
	}
	return config
}

func getClientset() *kubernetes.Clientset {
	config := getRestConfig()
	// create the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		message.Fatal(err, "Unable to connect to the K8s cluster")
	}

	return clientset
}

// readFile just reads a file into a byte array.
func readFile(file string) ([]byte, error) {
	b, err := ioutil.ReadFile(file)
	if err != nil {
		message.Debug(err)
		return []byte{}, fmt.Errorf("cannot read file %v, %v", file, err)
	}
	return b, nil
}

func GetContext() (string, error) {
	kubeconfig := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
		clientcmd.NewDefaultClientConfigLoadingRules(),
		&clientcmd.ConfigOverrides{},
	)
	kubeconfig.ConfigAccess().GetLoadingPrecedence()
	kubeConf, err := kubeconfig.ConfigAccess().GetStartingConfig()
	if err != nil {
		return "", fmt.Errorf("unable to load the default kube config")
	}

	return kubeConf.CurrentContext, nil
}

// ProcessYamlFilesInPath iterates over all yaml files in a given path and performs Zarf templating + image swapping
func ProcessYamlFilesInPath(path string, componentImages []string) []string {
	message.Debugf("k8s.ProcessYamlFilesInPath(%s, %v)", path, componentImages)

	// Only pull in yml and yaml files
	pattern := regexp.MustCompile(`(?mi)\.ya?ml$`)
	manifests := utils.RecursiveFileList(path, pattern)
	valueTemplate := template.Generate()

	// Match images in the given list and replace if found in the given files
	var imageSwap []ImageSwap
	for _, image := range componentImages {
		imageSwap = append(imageSwap, ImageSwap{
			find:    image,
			replace: utils.SwapHost(image, valueTemplate.GetRegistry()),
		})
	}

	for _, manifest := range manifests {
		message.Debugf("Processing k8s manifest files %s", manifest)
		// Iterate over each image swap to see if it exists in the manifest
		for _, swap := range imageSwap {
			utils.ReplaceText(manifest, swap.find, swap.replace)
		}
		valueTemplate.Apply(manifest)
	}

	return manifests
}

func generateLogShim() logr.Logger {
	message.Debug("k8s.generateLogShim()")
	return funcr.New(func(prefix, args string) {
		message.Debug(args)
	}, funcr.Options{})
}

// SplitYAML splits a YAML file into unstructured objects. Returns list of all unstructured objects
// found in the yaml. If an error occurs, returns objects that have been parsed so far too.
// Source: https://github.com/argoproj/gitops-engine/blob/v0.5.2/pkg/utils/kube/kube.go#L286
func SplitYAML(yamlData []byte) ([]*unstructured.Unstructured, error) {
	var objs []*unstructured.Unstructured
	ymls, err := splitYAMLToString(yamlData)
	if err != nil {
		return nil, err
	}
	for _, yml := range ymls {
		u := &unstructured.Unstructured{}
		if err := yaml.Unmarshal([]byte(yml), u); err != nil {
			return objs, fmt.Errorf("failed to unmarshal manifest: %v", err)
		}
		objs = append(objs, u)
	}
	return objs, nil
}

// splitYAMLToString splits a YAML file into strings. Returns list of yamls
// found in the yaml. If an error occurs, returns objects that have been parsed so far too.
// Source: https://github.com/argoproj/gitops-engine/blob/v0.5.2/pkg/utils/kube/kube.go#L304
func splitYAMLToString(yamlData []byte) ([]string, error) {
	// Similar way to what kubectl does
	// https://github.com/kubernetes/cli-runtime/blob/master/pkg/resource/visitor.go#L573-L600
	// Ideally k8s.io/cli-runtime/pkg/resource.Builder should be used instead of this method.
	// E.g. Builder does list unpacking and flattening and this code does not.
	d := kubeyaml.NewYAMLOrJSONDecoder(bytes.NewReader(yamlData), 4096)
	var objs []string
	for {
		ext := runtime.RawExtension{}
		if err := d.Decode(&ext); err != nil {
			if err == io.EOF {
				break
			}
			return objs, fmt.Errorf("failed to unmarshal manifest: %v", err)
		}
		ext.Raw = bytes.TrimSpace(ext.Raw)
		if len(ext.Raw) == 0 || bytes.Equal(ext.Raw, []byte("null")) {
			continue
		}
		objs = append(objs, string(ext.Raw))
	}
	return objs, nil
}
