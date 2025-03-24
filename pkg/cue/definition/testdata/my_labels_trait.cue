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
