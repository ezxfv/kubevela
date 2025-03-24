annotations: {
	type: "trait"
	annotations: {}
	description: "Add annotations on K8s pod for your workload"
	attributes: {
		podDisruptive: true
		appliesToWorkloads: ["deployments.apps", "statefulsets.apps", "daemonsets.apps", "jobs.batch", "clonesets.apps.kruise.io", "statefulsets.apps.kruise.io", "daemonsets.apps.kruise.io", "broadcastjobs.apps.kruise.io", "uniteddeployments.apps.kruise.io"]
	}
}
template: {
	// +patchStrategy=jsonMergePatch
	patch: {
		let annotationsContent = {
			for k, v in parameter {
				(k): v
			}
		}

		metadata: {
			annotations: annotationsContent
		}
		if context.output.spec != _|_ if context.output.spec.template != _|_ {
			spec: template: metadata: annotations: annotationsContent
		}
		if context.output.spec != _|_ if context.output.spec.jobTemplate != _|_ {
			spec: jobTemplate: metadata: annotations: annotationsContent
		}
		if context.output.spec != _|_ if context.output.spec.jobTemplate != _|_ if context.output.spec.jobTemplate.spec != _|_ if context.output.spec.jobTemplate.spec.template != _|_ {
			spec: jobTemplate: spec: template: metadata: annotations: annotationsContent
		}
	}
	parameter: [string]: string | null
}
