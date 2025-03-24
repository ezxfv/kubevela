hostalias: {
	type: "trait"
	annotations: {}
	description: "Add host aliases for your workload"
	attributes: {
		podDisruptive: true
		appliesToWorkloads: ["deployments.apps", "statefulsets.apps", "daemonsets.apps", "jobs.batch", "clonesets.apps.kruise.io", "statefulsets.apps.kruise.io", "daemonsets.apps.kruise.io", "broadcastjobs.apps.kruise.io"]
	}
}
template: {
	patch: {
		// +patchKey=ip
		spec: template: spec: hostAliases: parameter.hostAliases
	}
	parameter: {
		// +usage=Specify the hostAliases to add
		hostAliases: [...{
			ip: string
			hostnames: [...string]
		}]
	}
}
