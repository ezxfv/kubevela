virtualgroup: {
	type: "trait"
	annotations: {}
	labels: {}
	description: "Add virtual group labels"
	attributes: {
		appliesToWorkloads: ["deployments.apps", "statefulsets.apps", "daemonsets.apps", "clonesets.apps.kruise.io", "statefulsets.apps.kruise.io", "daemonsets.apps.kruise.io", "broadcastjobs.apps.kruise.io", "uniteddeployments.apps.kruise.io"]
		podDisruptive: false
	}
}
template: {
	patch: spec: template: metadata: labels: {
		if parameter.type == "namespace" {
			"app.namespace.virtual.group": parameter.group
		}
		if parameter.type == "cluster" {
			"app.cluster.virtual.group": parameter.group
		}
	}
	parameter: {
		group: *"default" | string
		type:  *"namespace" | string
	}
}
