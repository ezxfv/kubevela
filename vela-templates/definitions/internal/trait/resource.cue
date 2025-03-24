resource: {
	type: "trait"
	annotations: {}
	description: "Specify the resource requirements for the container"
	attributes: {
		podDisruptive: true
		appliesToWorkloads: ["deployments.apps", "statefulsets.apps", "daemonsets.apps", "jobs.batch", "cronjobs.batch", "clonesets.apps.kruise.io", "statefulsets.apps.kruise.io", "daemonsets.apps.kruise.io", "broadcastjobs.apps.kruise.io"]
	}
}
template: {

	let resourceContent = {
		resources: {
			if parameter.cpu != _|_ if parameter.memory != _|_ if parameter.requests == _|_ if parameter.limits == _|_ {
				// +patchStrategy=retainKeys
				requests: {
					cpu:    parameter.cpu
					memory: parameter.memory
				}
				// +patchStrategy=retainKeys
				limits: {
					cpu:    parameter.cpu
					memory: parameter.memory
				}
			}

			if parameter.requests != _|_ {
				// +patchStrategy=retainKeys
				requests: {
					cpu:    parameter.requests.cpu
					memory: parameter.requests.memory
				}
			}
			if parameter.limits != _|_ {
				// +patchStrategy=retainKeys
				limits: {
					cpu:    parameter.limits.cpu
					memory: parameter.limits.memory
				}
			}
		}
	}

	if context.output.spec != _|_ if context.output.spec.template != _|_ {
		patch: spec: template: spec: {
			// +patchKey=name
			containers: [resourceContent]
		}
	}
	if context.output.spec != _|_ if context.output.spec.jobTemplate != _|_ {
		patch: spec: jobTemplate: spec: template: spec: {
			// +patchKey=name
			containers: [resourceContent]
		}
	}

	parameter: {
		// +usage=Specify the amount of cpu for requests and limits
		cpu?: *1 | number | string
		// +usage=Specify the amount of memory for requests and limits
		memory?: *"2048Mi" | =~"^([1-9][0-9]{0,63})(E|P|T|G|M|K|Ei|Pi|Ti|Gi|Mi|Ki)$"
		// +usage=Specify the resources in requests
		requests?: {
			// +usage=Specify the amount of cpu for requests
			cpu: *1 | number | string
			// +usage=Specify the amount of memory for requests
			memory: *"2048Mi" | =~"^([1-9][0-9]{0,63})(E|P|T|G|M|K|Ei|Pi|Ti|Gi|Mi|Ki)$"
		}
		// +usage=Specify the resources in limits
		limits?: {
			// +usage=Specify the amount of cpu for limits
			cpu: *1 | number | string
			// +usage=Specify the amount of memory for limits
			memory: *"2048Mi" | =~"^([1-9][0-9]{0,63})(E|P|T|G|M|K|Ei|Pi|Ti|Gi|Mi|Ki)$"
		}
	}
}
