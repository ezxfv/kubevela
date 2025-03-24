lifecycle: {
	type: "trait"
	annotations: {}
	description: "Add lifecycle hooks for the container"
	attributes: {
		podDisruptive: true
		appliesToWorkloads: ["deployments.apps", "statefulsets.apps", "daemonsets.apps", "jobs.batch", "clonesets.apps.kruise.io", "statefulsets.apps.kruise.io", "daemonsets.apps.kruise.io", "broadcastjobs.apps.kruise.io"]
	}
}
template: {
	patch: spec: template: spec: containers: [...{
		lifecycle: {
			if parameter.postStart != _|_ {
				postStart: parameter.postStart
			}
			if parameter.preStop != _|_ {
				preStop: parameter.preStop
			}
		}
	}]
	parameter: {
		postStart?: #LifeCycleHandler
		preStop?:   #LifeCycleHandler
	}
	#Port: int & >=1 & <=65535
	#LifeCycleHandler: {
		exec?: command: [...string]
		httpGet?: {
			path?:  string
			port:   #Port
			host?:  string
			scheme: *"HTTP" | "HTTPS"
			httpHeaders?: [...{
				name:  string
				value: string
			}]
		}
		tcpSocket?: {
			port:  #Port
			host?: string
		}
	}
}
