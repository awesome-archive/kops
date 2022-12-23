/*
Copyright 2017 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package ops

import (
	"fmt"

	"k8s.io/kops/pkg/apis/kops"
	"k8s.io/kops/pkg/resources"
	"k8s.io/kops/pkg/resources/aws"
	"k8s.io/kops/pkg/resources/azure"
	"k8s.io/kops/pkg/resources/digitalocean"
	"k8s.io/kops/pkg/resources/gce"
	"k8s.io/kops/pkg/resources/hetzner"
	"k8s.io/kops/pkg/resources/openstack"
	"k8s.io/kops/pkg/resources/scaleway"
	"k8s.io/kops/upup/pkg/fi"
	"k8s.io/kops/upup/pkg/fi/cloudup/awsup"
	cloudazure "k8s.io/kops/upup/pkg/fi/cloudup/azure"
	clouddo "k8s.io/kops/upup/pkg/fi/cloudup/do"
	cloudgce "k8s.io/kops/upup/pkg/fi/cloudup/gce"
	cloudhetzner "k8s.io/kops/upup/pkg/fi/cloudup/hetzner"
	cloudopenstack "k8s.io/kops/upup/pkg/fi/cloudup/openstack"
	cloudscaleway "k8s.io/kops/upup/pkg/fi/cloudup/scaleway"
)

// ListResources collects the resources from the specified cloud
func ListResources(cloud fi.Cloud, cluster *kops.Cluster) (map[string]*resources.Resource, error) {
	clusterInfo := resources.ClusterInfo{
		Name:        cluster.Name,
		UsesNoneDNS: cluster.UsesNoneDNS(),
	}

	switch cloud.ProviderID() {
	case kops.CloudProviderAWS:
		return aws.ListResourcesAWS(cloud.(awsup.AWSCloud), clusterInfo)
	case kops.CloudProviderDO:
		return digitalocean.ListResources(cloud.(clouddo.DOCloud), clusterInfo)
	case kops.CloudProviderGCE:
		return gce.ListResourcesGCE(cloud.(cloudgce.GCECloud), clusterInfo)
	case kops.CloudProviderHetzner:
		return hetzner.ListResources(cloud.(cloudhetzner.HetznerCloud), clusterInfo)
	case kops.CloudProviderOpenstack:
		return openstack.ListResources(cloud.(cloudopenstack.OpenstackCloud), clusterInfo)
	case kops.CloudProviderAzure:
		clusterInfo.AzureResourceGroupName = cluster.AzureResourceGroupName()
		clusterInfo.AzureResourceGroupShared = cluster.IsSharedAzureResourceGroup()
		clusterInfo.AzureNetworkShared = cluster.SharedVPC()
		clusterInfo.AzureRouteTableShared = cluster.IsSharedAzureRouteTable()
		return azure.ListResourcesAzure(cloud.(cloudazure.AzureCloud), clusterInfo)
	case kops.CloudProviderScaleway:
		return scaleway.ListResources(cloud.(cloudscaleway.ScwCloud), clusterInfo)
	default:
		return nil, fmt.Errorf("delete on clusters on %q not (yet) supported", cloud.ProviderID())
	}
}
