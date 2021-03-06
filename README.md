validator
---------

WIP validating admission webhook

By default, this will create the validating webhook configuration and run in the `default` namespace. It will target any create or update operations with namespaces that have the `validator: enabled` label. The example configuration (`deploy/configmap.yaml`) uses regex to match patterns for container images. The default only allows images from ECR. It will not allow load balancers to be created in the namespace as well. Any labels in the label sections are required for resource creation, the example config requires the `foo` label on Pod and Service objects. 

Installation Instructions
=========

1. Clone this repository
2. Run `make gen-pki`
3. Run `make deploy`
4. Label namespaces with `validator: enabled` that the webhook should be used on.

Uninstall Instructions
=========

1. Run `kubectl delete -f deploy/`

Troubleshooting
=========

TLS errors in the validator pod usually indicate an issue with the CA_BUNDLE in the validatingwebhookconfiguration or mistake in the CN of the server certificate.
