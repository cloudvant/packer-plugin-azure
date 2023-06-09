// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package dtl

// Code to resolve resources that are required by the API.  These resources
// can most likely be resolved without asking the user, thereby reducing the
// amount of configuration they need to provide.
//
// Resource resolver differs from config retriever because resource resolver
// requires a client to communicate with the Azure API.  A config retriever is
// used to determine values without use of a client.

import (
	"context"
	"fmt"
	"strings"

	"github.com/Azure/azure-sdk-for-go/services/compute/mgmt/2018-04-01/compute"
)

type resourceResolver struct {
	client                          *AzureClient
	findVirtualNetworkResourceGroup func(*AzureClient, string) (string, error)
	findVirtualNetworkSubnet        func(*AzureClient, string, string) (string, error)
}

func newResourceResolver(client *AzureClient) *resourceResolver {
	return &resourceResolver{
		client:                          client,
		findVirtualNetworkResourceGroup: findVirtualNetworkResourceGroup,
		findVirtualNetworkSubnet:        findVirtualNetworkSubnet,
	}
}

func (s *resourceResolver) Resolve(c *Config) error {
	// if s.shouldResolveResourceGroup(c) {
	// 	resourceGroupName, err := s.findVirtualNetworkResourceGroup(s.client, c.VirtualNetworkName)
	// 	if err != nil {
	// 		return err
	// 	}

	// 	subnetName, err := s.findVirtualNetworkSubnet(s.client, resourceGroupName, c.VirtualNetworkName)
	// 	if err != nil {
	// 		return err
	// 	}

	// 	c.VirtualNetworkResourceGroupName = resourceGroupName
	// 	c.VirtualNetworkSubnetName = subnetName
	// }

	if s.shouldResolveManagedImageName(c) {
		image, err := findManagedImageByName(s.client, c.CustomManagedImageName, c.CustomManagedImageResourceGroupName)
		if err != nil {
			return err
		}

		c.customManagedImageID = *image.ID
	}

	return nil
}

// func (s *resourceResolver) shouldResolveResourceGroup(c *Config) bool {
// 	return c.VirtualNetworkName != "" && c.VirtualNetworkResourceGroupName == ""
// }

func (s *resourceResolver) shouldResolveManagedImageName(c *Config) bool {
	return c.CustomManagedImageName != ""
}

func getResourceGroupNameFromId(id string) string {
	// "/subscriptions/3f499422-dd76-4114-8859-86d526c9deb6/resourceGroups/packer-Resource-Group-yylnwsl30j/providers/...
	xs := strings.Split(id, "/")
	return xs[4]
}

func findManagedImageByName(client *AzureClient, name, resourceGroupName string) (*compute.Image, error) {
	images, err := client.ImagesClient.ListByResourceGroupComplete(context.TODO(), resourceGroupName)
	if err != nil {
		return nil, err
	}

	for images.NotDone() {
		image := images.Value()
		if strings.EqualFold(name, *image.Name) {
			return &image, nil
		}
		if err = images.Next(); err != nil {
			return nil, err
		}
	}

	return nil, fmt.Errorf("Cannot find an image named '%s' in the resource group '%s'", name, resourceGroupName)
}

func findVirtualNetworkResourceGroup(client *AzureClient, name string) (string, error) {
	virtualNetworks, err := client.VirtualNetworksClient.ListAllComplete(context.TODO())
	if err != nil {
		return "", err
	}

	resourceGroupNames := make([]string, 0)
	for virtualNetworks.NotDone() {
		virtualNetwork := virtualNetworks.Value()
		if strings.EqualFold(name, *virtualNetwork.Name) {
			rgn := getResourceGroupNameFromId(*virtualNetwork.ID)
			resourceGroupNames = append(resourceGroupNames, rgn)
		}
		if err = virtualNetworks.Next(); err != nil {
			return "", err
		}
	}

	if len(resourceGroupNames) == 0 {
		return "", fmt.Errorf("Cannot find a resource group with a virtual network called %q", name)
	}

	if len(resourceGroupNames) > 1 {
		return "", fmt.Errorf("Found multiple resource groups with a virtual network called %q, please use virtual_network_subnet_name and virtual_network_resource_group_name to disambiguate", name)
	}

	return resourceGroupNames[0], nil
}

func findVirtualNetworkSubnet(client *AzureClient, resourceGroupName string, name string) (string, error) {
	subnets, err := client.SubnetsClient.List(context.TODO(), resourceGroupName, name)
	if err != nil {
		return "", err
	}

	subnetList := subnets.Values() // only first page of subnets, but only interested in ==0 or >1

	if len(subnetList) == 0 {
		return "", fmt.Errorf("Cannot find a subnet in the resource group %q associated with the virtual network called %q", resourceGroupName, name)
	}

	if len(subnetList) > 1 {
		return "", fmt.Errorf("Found multiple subnets in the resource group %q associated with the virtual network called %q, please use virtual_network_subnet_name and virtual_network_resource_group_name to disambiguate", resourceGroupName, name)
	}

	subnet := subnetList[0]
	return *subnet.Name, nil
}
