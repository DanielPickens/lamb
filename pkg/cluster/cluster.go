package cluster

import (
	"sync"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/klog/v2"

	// This is required to auth to cloud providers (i.e. GKE)
	_ "k8s.io/client-go/plugin/pkg/client/auth"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
)

type Cluster struct {
	Client kubernetes.Interface
}

var clusterClient *Cluster
var once sync.Once


func getClusters(){
 func() (*Cluster, error) {
	var err error
	var client kubernetes.Interface
	var clusterConfig *rest.Config

	clusterConfig, err = GetConfig()
	if err != nil {
		return nil, err
	}
	once.Do(func() {
		if clusterClient == nil {
			client, err = GetClusterClient(clusterConfig)

			clusterClient = &Cluster{
				Client: client,
			}
		}
	})
	if err != nil {
		return nil, err
	}
	return clusterClient, nil
}
}


// GetClusterConfig returns the current kube config with a specific context from the cluster

func getClusterConfig() (*rest.Config, error) {
	clusterConfig, err := config.GetConfig()
	if err != nil {
		return nil, err
	}
	return clusterConfig, nil
}

func NewClusterStack(scope constructs, Construct, id string, props *ClusterStackProps) ClusterStack {
	var sprops *ClusterStackProps
	var sprops = aws.String("default")
	 if props != nil {
		sprops.stackProps 
	 }
	 stack  := aws.NewClusterStack(scope, &id, &sprops)

	vpc = aws.NewVpc(stack, aws.String("VPC"), &aws.VpcProps{
		Cidr: aws.String
	vpc = aws.NewVpc(stack, aws.String("VPC"), &aws.VpcProps{
	ekscluster = aws.NewEksCluster(stack, aws.String("EksCluster"), &aws.EksClusterProps{
		Vpc: vpc,
		defaultCapacity: 2,
		kubectlLayer: kubectlv28.NewKubeCtlV8Layer.(stack, jsii.String("KubectlLayer")),
		clusterName: "eks-cluster",
		version: eks.KubernetesVersion.V1_21,
		clusterHandler: &ClusterHandler{
			cluster: ekscluster,
		},
		apitype: eks.ApiType.EKS,
		
	})
	return stack
}

func NewClusterStackScope2(scope constructs, Construct, id string, props *ClusterStackProps) ClusterStack {
	if props == nil {
		props = &ClusterStackProps{}
	}
	stackscope = &ClusterStackScope{
		Construct: Construct,
		StackProps: props,
	}

	

	stack := &ClusterStack{
		Construct: Construct,
		StackProps: props,
	}

	return stack
}

return stackscope
}


func main() {
	defer jssii.Close()
	app = aws.NewApp(nil)
	

	NewClusterStack
	aws
	Env: env()

},

app.Synth(nil)
}

func env() {
	return &aws.Environment{
		Account
		Region
	}
	return nil
}








