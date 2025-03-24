output: {
	// CloneSet workload from OpenKruise
	apiVersion: "apps.kruise.io/v1alpha1"
	kind:       "CloneSet"
	metadata: {
		name: context.name
	}
	spec: {
		replicas: 2
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
					env: [{
						name: "CONFIG_VALUE"
						valueFrom: {
							configMapKeyRef: {
								name: context.name + "-config"
								key:  "config.properties"
							}
						}
					}, {
						name: "SECRET_VALUE"
						valueFrom: {
							secretKeyRef: {
								name: context.name + "-secret"
								key:  "credentials"
							}
						}
					}]
				}]
			}
		}
		updateStrategy: {
			type:           "InPlaceIfPossible"
			maxUnavailable: 1
			maxSurge:       0
		}
	}
}

outputs: {
	// Service for the CloneSet
	service: {
		apiVersion: "v1"
		kind:       "Service"
		metadata: {
			name: context.name
		}
		spec: {
			selector: {
				"app": context.name
			}
			ports: [{
				port:       80
				targetPort: 8000
			}]
		}
	}

	// ConfigMap for application configuration
	configmap: {
		apiVersion: "v1"
		kind:       "ConfigMap"
		metadata: {
			name: context.name + "-config"
		}
		data: {
			"config.properties": """
			app.name=FutuService
			app.version=1.0.0
			app.environment=development
			"""
		}
	}

	// Secret for sensitive information
	secret: {
		apiVersion: "v1"
		kind:       "Secret"
		metadata: {
			name: context.name + "-secret"
		}
		type: "Opaque"
		data: {
			// Base64 encoded "admin:password"
			"credentials": "YWRtaW46cGFzc3dvcmQ="
		}
	}
}
