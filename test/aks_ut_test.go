package test

import (
	"fmt"
	"testing"

	"github.com/gruntwork-io/terratest/modules/random"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
	"github.com/tidwall/gjson"
)

func TestTerraformAzureAKSUnitExample(t *testing.T) {
	t.Parallel()
	// MC_+ResourceGroupName_ClusterName_AzureRegion must be no greater than 80 chars.
	// https://docs.microsoft.com/en-us/azure/aks/troubleshooting#what-naming-restrictions-are-enforced-for-aks-resources-and-parameters
	expectedClusterName := fmt.Sprintf("cluster-%s", random.UniqueId())
	//expectedResourceGroupName := fmt.Sprintf("rg-aks-%s-eastus", strings.ToLower(expectedClusterName))
	expectedAagentCount := 3
	tfPlanOutput := "terraform.tfplan"

	terraformOptions := &terraform.Options{
		TerraformDir: "../terraform",
		Vars: map[string]interface{}{
			"cluster_name": expectedClusterName,
			"location":     "East US",
			"node_count":   expectedAagentCount,
		},
	}

	// apparently this runs terraform destroy at the end of the test
	//defer terraform.Destroy(t, terraformOptions)

	terraform.Init(t, terraformOptions)

	terraform.RunTerraformCommand(t, terraformOptions, terraform.FormatArgs(terraformOptions, "plan", "-out="+tfPlanOutput)...)
	terraformOptions.Vars = nil
	planjsonstr := terraform.RunTerraformCommand(t, terraformOptions, terraform.FormatArgs(terraformOptions, "show", "-json", tfPlanOutput)...)

	assert.Equal(t, float64(3), gjson.Get(planjsonstr, "planned_values.root_module.child_modules.0.resources.0.values.default_node_pool.0.node_count").Float())
}
