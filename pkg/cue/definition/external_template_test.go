package definition

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"

	"github.com/oam-dev/kubevela/pkg/cue/process"
)

func TestExternalPackageWithMyLabelsTrait(t *testing.T) {
	// Read the template files
	traitTemplateBytes, err := os.ReadFile(filepath.Join("testdata", "my_labels_trait.cue"))
	require.NoError(t, err)

	// Define a simple webservice template for testing
	webserviceTemplate, err := os.ReadFile(filepath.Join("testdata", "webservice.cue"))
	require.NoError(t, err)

	// Create a test case for the my-labels trait
	t.Run("my-labels trait with external utils.Sum function", func(t *testing.T) {
		// Create a shared context for both workload and trait
		ctx := process.NewContext(process.ContextData{
			AppName:         "first-vela-app-3",
			CompName:        "express-server-3",
			Namespace:       "default",
			AppRevisionName: "first-vela-app-3-v1",
		})

		// First apply the workload
		workload := NewWorkloadAbstractEngine("webservice")
		err := workload.Complete(ctx, string(webserviceTemplate), map[string]interface{}{})
		require.NoError(t, err)

		// Then apply the trait to the same context
		trait := NewTraitAbstractEngine("my-labels")
		params := map[string]interface{}{
			"x.io/xxx": "v1",
		}
		err = trait.Complete(ctx, string(traitTemplateBytes), params)
		require.NoError(t, err)

		// Get the final output with the trait applied
		finalBase, _ := ctx.Output()
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
		assert.Equal(t, "express-server-3", podLabels["app"])
	})
}
