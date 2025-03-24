package definition_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"

	"github.com/oam-dev/kubevela/pkg/cue/definition"
	"github.com/oam-dev/kubevela/pkg/cue/process"
)

func TestExternalPackageWithMyLabelsTrait(t *testing.T) {
	// Read the template files
	traitTemplateBytes, err := os.ReadFile(filepath.Join("testdata", "my_labels_trait.cue"))
	require.NoError(t, err)

	// Define a simple futuservice template for testing
	futuserviceTemplate, err := os.ReadFile(filepath.Join("testdata", "futuservice.cue"))
	require.NoError(t, err)

	// Create a test case for the my-labels trait
	t.Run("my-labels trait with external utils.Sum function", func(t *testing.T) {
		// Create a shared context for both workload and trait
		ctx := process.NewContext(process.ContextData{
			AppName:         "first-vela-app",
			CompName:        "express-server",
			Namespace:       "default",
			AppRevisionName: "first-vela-app-v1",
		})

		// First apply the workload
		workload := definition.NewWorkloadAbstractEngine("futuservice")
		err := workload.Complete(ctx, string(futuserviceTemplate), map[string]interface{}{})
		require.NoError(t, err)

		// Then apply the trait to the same context
		trait := definition.NewTraitAbstractEngine("my-labels")
		params := map[string]interface{}{
			"x.io/xxx": "v1",
		}
		err = trait.Complete(ctx, string(traitTemplateBytes), params)
		require.NoError(t, err)

		// Get the final output with the trait applied
		finalBase, outputs := ctx.Output()
		finalObj, err := finalBase.Unstructured()
		require.NoError(t, err)

		// Verify the expected output
		// The deployment should have the labels from the parameters
		// plus the x.io/debug label with the sum result (10 + 2 = 12)
		labels, found, err := unstructured.NestedStringMap(finalObj.Object, "metadata", "labels")
		require.NoError(t, err)
		require.True(t, found)
		assert.Equal(t, "v1", labels["x.io/xxx"])
		assert.Equal(t, "12", labels["x.io/debug"])

		// Also check the pod template labels
		podLabels, found, err := unstructured.NestedStringMap(finalObj.Object, "spec", "template", "metadata", "labels")
		require.NoError(t, err)
		require.True(t, found)
		assert.Equal(t, "v1", podLabels["x.io/xxx"])
		assert.Equal(t, "12", podLabels["x.io/debug"])
		assert.Equal(t, "express-server", podLabels["app"])

		// Check that the outputs also have the labels applied
		require.NotEmpty(t, outputs)
		for name, output := range outputs {
			outputObj, err := output.Ins.Unstructured()
			require.NoError(t, err, "Failed to convert output %s to unstructured", name)

			outputLabels, found, err := unstructured.NestedStringMap(outputObj.Object, "metadata", "labels")
			require.NoError(t, err, "Error getting labels for output %s", name)
			require.True(t, found, "No labels found for output %s", name)
			assert.Equal(t, "v1", outputLabels["x.io/xxx"], "Missing parameter label in output %s", name)
			assert.Equal(t, "12", outputLabels["x.io/debug"], "Missing debug label in output %s", name)
		}
	})
}
