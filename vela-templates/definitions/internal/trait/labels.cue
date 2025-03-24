labels: {
	type: "trait"
	annotations: {}
	description: "Add labels on K8s pod for your workload"
	attributes: {
		podDisruptive: true
		appliesToWorkloads: ["deployments.apps", "statefulsets.apps", "daemonsets.apps", "jobs.batch", "clonesets.apps.kruise.io", "statefulsets.apps.kruise.io", "daemonsets.apps.kruise.io", "broadcastjobs.apps.kruise.io", "uniteddeployments.apps.kruise.io"]
	}
}
template: {
	// +patchStrategy=jsonMergePatch
	patch: {
		metadata: {
			labels: {
				for k, v in parameter {
					(k): v
				}
			}
		}
		if context.output.spec != _|_ && context.output.spec.template != _|_ {
			spec: template: metadata: labels: {
				for k, v in parameter {
					(k): v
				}
			}
		}
	}
	parameter: [string]: string | null
}
