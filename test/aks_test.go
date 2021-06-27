package test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/gruntwork-io/terratest/modules/azure"
	"github.com/gruntwork-io/terratest/modules/random"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTerraformAzureAKSExample(t *testing.T) {
	t.Parallel()
	// MC_+ResourceGroupName_ClusterName_AzureRegion must be no greater than 80 chars.
	// https://docs.microsoft.com/en-us/azure/aks/troubleshooting#what-naming-restrictions-are-enforced-for-aks-resources-and-parameters
	expectedClusterName := fmt.Sprintf("cluster-%s", random.UniqueId())
	expectedResourceGroupName := fmt.Sprintf("rg-aks-%s-eastus", strings.ToLower(expectedClusterName))
	expectedAagentCount := 3

	terraformOptions := &terraform.Options{
		TerraformDir: "../terraform",
		Vars: map[string]interface{}{
			"cluster_name": expectedClusterName,
			"location":     "East US",
			"node_count":   expectedAagentCount,
		},
	}

	// apparently this runs terraform destroy at the end of the test
	defer terraform.Destroy(t, terraformOptions)

	terraform.InitAndApply(t, terraformOptions)

	cluster, err := azure.GetManagedClusterE(t, expectedResourceGroupName, expectedClusterName, "")
	require.NoError(t, err)
	actualCount := *(*cluster.ManagedClusterProperties.AgentPoolProfiles)[0].Count

	// Test that the Node count matches the Terraform specification
	assert.Equal(t, int32(expectedAagentCount), actualCount)

}
