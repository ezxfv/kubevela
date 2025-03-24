package definition

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"sigs.k8s.io/yaml"

	"github.com/oam-dev/kubevela/pkg/cue/process"
)

// YAMLToUnstructured converts a YAML string to an unstructured.Unstructured object
// It will panic if the YAML cannot be parsed
func YAMLToUnstructured(yamlStr string) *unstructured.Unstructured {
	var jsonObj any
	err := yaml.Unmarshal([]byte(yamlStr), &jsonObj)
	if err != nil {
		panic(err)
	}

	// Convert to JSON and back to ensure numbers are parsed as float64
	jsonBytes, err := json.Marshal(jsonObj)
	if err != nil {
		panic(err)
	}

	obj := &unstructured.Unstructured{}
	err = json.Unmarshal(jsonBytes, &obj.Object)
	if err != nil {
		panic(err)
	}

	return obj
}

func TestExternalPackageWithMyLabelsTrait(t *testing.T) {
	// Define the trait template that uses the external package
	myLabelsTraitTemplate := `
import "ext/utils"

parameter: [string]: string | null

// Call the Sum function and get the result
sumResult: utils.#Sum & {
	$params: {
		x: 10
		y: 2
	}
}

// +patchStrategy=jsonMergePatch
patch: {
	metadata: {
		labels: {
			for k, v in parameter {
				(k): v
			}
			// Use the Sum function's return value
			"x.io/debug": "\(sumResult.$returns.result)"
		}
	}
	if context.output.spec != _|_ && context.output.spec.template != _|_ {
		spec: template: metadata: labels: {
			for k, v in parameter {
				(k): v
			}
			// Also add the result to the pod template
			"x.io/debug": "\(sumResult.$returns.result)"
		}
	}
}
`

	// Define a simple workload to apply the trait to
	webserviceTemplate := `
output: {
	apiVersion: "apps/v1"
	kind:       "Deployment"
	metadata: {
		name: context.name
	}
	spec: {
		selector: {
			matchLabels: {
				"app": context.name
			}
		}
		template: {
			metadata: {
				labels: {
					"app": context.name
				}
			}
			spec: {
				containers: [{
					name:  context.name
					image: "oamdev/hello-world"
					ports: [{
						containerPort: 8000
					}]
				}]
			}
		}
	}
}
`

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
		err := workload.Complete(ctx, webserviceTemplate, map[string]interface{}{})
		require.NoError(t, err)

		// Then apply the trait to the same context
		trait := NewTraitAbstractEngine("my-labels")
		params := map[string]interface{}{
			"x.io/xxx": "v1",
		}
		err = trait.Complete(ctx, myLabelsTraitTemplate, params)
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

	// Test with YAML definitions
	t.Run("my-labels trait from YAML definitions", func(t *testing.T) {
		// Define the trait definition YAML
		traitDefYAML := `
apiVersion: core.oam.dev/v1beta1
kind: TraitDefinition
metadata:
  name: my-labels
spec:
  schematic:
    cue:
      template: |
        import "ext/utils"

        parameter: [string]: string | null

        // Call the Sum function and get the result
        sumResult: utils.#Sum & {
          $params: {
            x: 10
            y: 2
          }
        }

        // +patchStrategy=jsonMergePatch
        patch: {
          metadata: {
            labels: {
              for k, v in parameter {
                (k): v
              }
              // Use the Sum function's return value
              "x.io/debug": "\(sumResult.$returns.result)"
            }
          }
          if context.output.spec != _|_ && context.output.spec.template != _|_ {
            spec: template: metadata: labels: {
              for k, v in parameter {
                (k): v
              }
              // Also add the result to the pod template
              "x.io/debug": "\(sumResult.$returns.result)"
            }
          }
        }
`

		// Define the application YAML
		appYAML := `
apiVersion: core.oam.dev/v1beta1
kind: Application
metadata:
  name: first-vela-app-3
spec:
  components:
    - name: express-server-3
      type: webservice
      properties:
        image: oamdev/hello-world
        ports:
         - port: 8000
           expose: true
      traits:
        - type: scaler
          properties:
            replicas: 0
        - type: my-labels
          properties:
            x.io/xxx: v1
`

		// Parse the trait definition
		traitDefObj := YAMLToUnstructured(traitDefYAML)
		require.Equal(t, "TraitDefinition", traitDefObj.GetKind())
		require.Equal(t, "my-labels", traitDefObj.GetName())

		// Extract the template from the trait definition
		template, found, err := unstructured.NestedString(traitDefObj.Object, "spec", "schematic", "cue", "template")
		require.NoError(t, err)
		require.True(t, found)
		require.Contains(t, template, "import \"ext/utils\"")
		require.Contains(t, template, "sumResult: utils.#Sum")

		// Parse the application
		appObj := YAMLToUnstructured(appYAML)
		require.Equal(t, "Application", appObj.GetKind())
		require.Equal(t, "first-vela-app-3", appObj.GetName())

		// Extract components and traits
		components, found, err := unstructured.NestedSlice(appObj.Object, "spec", "components")
		require.NoError(t, err)
		require.True(t, found)
		require.Len(t, components, 1)

		component := components[0].(map[string]interface{})
		require.Equal(t, "express-server-3", component["name"])
		require.Equal(t, "webservice", component["type"])

		traits, found, err := unstructured.NestedSlice(component, "traits")
		require.NoError(t, err)
		require.True(t, found)
		require.Len(t, traits, 2)

		// Find the my-labels trait
		var myLabelsTrait map[string]interface{}
		for _, t := range traits {
			trait := t.(map[string]interface{})
			if trait["type"] == "my-labels" {
				myLabelsTrait = trait
				break
			}
		}
		require.NotNil(t, myLabelsTrait, "my-labels trait not found in application")

		properties, found, err := unstructured.NestedMap(myLabelsTrait, "properties")
		require.NoError(t, err)
		require.True(t, found)
		require.Equal(t, "v1", properties["x.io/xxx"])

		// Test the template with the extracted properties
		// Create a shared context for both workload and trait
		ctx := process.NewContext(process.ContextData{
			AppName:         "first-vela-app-3",
			CompName:        "express-server-3",
			Namespace:       "default",
			AppRevisionName: "first-vela-app-3-v1",
		})

		// First apply the workload
		workload := NewWorkloadAbstractEngine("webservice")
		err = workload.Complete(ctx, webserviceTemplate, map[string]interface{}{})
		require.NoError(t, err)

		// Then apply the trait to the same context
		trait := NewTraitAbstractEngine("my-labels")
		params := map[string]interface{}{
			"x.io/xxx": "v1",
		}
		err = trait.Complete(ctx, template, params)
		require.NoError(t, err)

		// Get the final output with the trait applied
		finalBase, _ := ctx.Output()
		finalObj, err := finalBase.Unstructured()
		require.NoError(t, err)

		// Verify the expected output
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
