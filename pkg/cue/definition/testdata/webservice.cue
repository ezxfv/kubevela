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