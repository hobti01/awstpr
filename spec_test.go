package awstpr

import (
	"io/ioutil"
	"net"
	"testing"

	yaml "gopkg.in/yaml.v2"

	"github.com/giantswarm/awstpr/aws"
	"github.com/giantswarm/awstpr/aws/hostedzones"
	"github.com/giantswarm/awstpr/aws/vpc"
	"github.com/giantswarm/clustertpr"
	"github.com/giantswarm/clustertpr/calico"
	"github.com/giantswarm/clustertpr/cluster"
	"github.com/giantswarm/clustertpr/customer"
	"github.com/giantswarm/clustertpr/docker"
	"github.com/giantswarm/clustertpr/docker/daemon"
	"github.com/giantswarm/clustertpr/docker/registry"
	"github.com/giantswarm/clustertpr/etcd"
	"github.com/giantswarm/clustertpr/kubernetes"
	"github.com/giantswarm/clustertpr/kubernetes/api"
	"github.com/giantswarm/clustertpr/kubernetes/dns"
	"github.com/giantswarm/clustertpr/kubernetes/hyperkube"
	hyperkubedocker "github.com/giantswarm/clustertpr/kubernetes/hyperkube/docker"
	"github.com/giantswarm/clustertpr/kubernetes/ingress"
	"github.com/giantswarm/clustertpr/kubernetes/kubectl"
	kubectldocker "github.com/giantswarm/clustertpr/kubernetes/kubectl/docker"
	"github.com/giantswarm/clustertpr/kubernetes/kubelet"
	"github.com/giantswarm/clustertpr/kubernetes/networksetup"
	networksetupdocker "github.com/giantswarm/clustertpr/kubernetes/networksetup/docker"
	"github.com/giantswarm/clustertpr/kubernetes/ssh"
	"github.com/giantswarm/clustertpr/node"
	"github.com/giantswarm/clustertpr/vault"
	"github.com/kylelemons/godebug/pretty"
	"github.com/stretchr/testify/require"
)

func TestSpecYamlEncoding(t *testing.T) {
	spec := Spec{
		Cluster: clustertpr.Cluster{
			Calico: calico.Calico{
				CIDR:   16,
				Domain: "giantswarm.io",
				MTU:    1500,
				Subnet: "10.1.2.3",
			},
			Cluster: cluster.Cluster{
				ID: "abc12",
			},
			Customer: customer.Customer{
				ID: "BooYa",
			},
			Docker: docker.Docker{
				Daemon: daemon.Daemon{
					ExtraArgs: "--log-opt max-file=1",
				},
				ImageNamespace: "giantswarm",
				Registry: registry.Registry{
					Endpoint: "http://giantswarm.io",
				},
			},
			Etcd: etcd.Etcd{
				AltNames: "",
				Domain:   "etcd.giantswarm.io",
				Port:     2379,
				Prefix:   "giantswarm.io",
			},
			Kubernetes: kubernetes.Kubernetes{
				API: api.API{
					AltNames:       "kubernetes,kubernetes.default",
					ClusterIPRange: "172.31.0.0/24",
					Domain:         "api.giantswarm.io",
					IP:             net.ParseIP("172.31.0.1"),
					InsecurePort:   8080,
					SecurePort:     443,
				},
				CloudProvider: "aws",
				DNS: dns.DNS{
					IP: net.ParseIP("172.31.0.10"),
				},
				Domain: "cluster.giantswarm.io",
				Hyperkube: hyperkube.Hyperkube{
					Docker: hyperkubedocker.Docker{
						Image: "quay.io/giantswarm/hyperkube",
					},
				},
				IngressController: ingress.IngressController{
					Domain:         "ingress.giantswarm.io",
					WildcardDomain: "*.giantswarm.io",
					InsecurePort:   30010,
					SecurePort:     30011,
				},
				Kubectl: kubectl.Kubectl{
					Docker: kubectldocker.Docker{
						Image: "quay.io/giantswarm/docker-kubectl",
					},
				},
				Kubelet: kubelet.Kubelet{
					AltNames: "kubernetes,kubernetes.default,kubernetes.default.svc",
					Domain:   "worker.giantswarm.io",
					Labels:   "etcd.giantswarm.io",
					Port:     10250,
				},
				NetworkSetup: networksetup.NetworkSetup{
					networksetupdocker.Docker{
						Image: "quay.io/giantswarm/k8s-setup-network-environment",
					},
				},
				SSH: ssh.SSH{
					PublicKeys: []string{
						"ssh-rsa AAAAB3NzaC1yc",
					},
				},
			},
			Masters: []node.Node{
				{
					ID: "fyz88",
				},
			},
			Vault: vault.Vault{
				Address: "vault.giantswarm.io",
				Token:   "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx",
			},
			Workers: []node.Node{
				{
					ID: "axx99",
				},
				{
					ID: "cdd88",
				},
			},
		},
		AWS: aws.AWS{
			Region: "eu-central-1",
			AZ:     "eu-central-1a",
			VPC: vpc.VPC{
				CIDR:              "10.0.0.0/16",
				PrivateSubnetCIDR: "10.0.0.0/19",
				PublicSubnetCIDR:  "10.0.128.0/20",
			},
			HostedZones: hostedzones.HostedZones{
				API:     "xxxxxxxxxxxxxx",
				Etcd:    "yyyyyyyyyyyyyy",
				Ingress: "zzzzzzzzzzzzzz",
			},
			Masters: []aws.Node{
				{
					ImageID:      "ami-d60ad6b9",
					InstanceType: "m3.large",
				},
			},
			Workers: []aws.Node{
				{
					ImageID:      "ami-d60ad6b9",
					InstanceType: "m3.large",
				},
				{
					ImageID:      "ami-d60ad6b9",
					InstanceType: "m3.large",
				},
			},
		},
	}

	var got map[string]interface{}
	{
		bytes, err := yaml.Marshal(&spec)
		require.NoError(t, err, "marshaling spec")
		err = yaml.Unmarshal(bytes, &got)
		require.NoError(t, err, "unmarshaling spec to map")
	}

	var want map[string]interface{}
	{
		bytes, err := ioutil.ReadFile("testdata/spec.yaml")
		require.NoError(t, err)
		err = yaml.Unmarshal(bytes, &want)
		require.NoError(t, err, "unmarshaling fixture to map")
	}

	diff := pretty.Compare(want, got)
	require.Equal(t, "", diff, "diff: (-want +got)\n%s", diff)
}