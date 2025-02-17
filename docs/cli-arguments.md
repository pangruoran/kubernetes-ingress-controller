<!-- This document is generated by KIC's 'generate.docs' make target, DO NOT EDIT -->

## Flags

| Flag | Type | Description | Default |
| ---- | ---- | ----------- | ------- |
| `--admission-webhook-cert` | `string` | Admission server PEM certificate value. |  |
| `--admission-webhook-cert-file` | `string` | Admission server PEM certificate file path; if both this and the cert value is unset, defaults to /admission-webhook/tls.crt. |  |
| `--admission-webhook-key` | `string` | Admission server PEM private key value. |  |
| `--admission-webhook-key-file` | `string` | Admission server PEM private key file path; if both this and the key value is unset, defaults to /admission-webhook/tls.key. |  |
| `--admission-webhook-listen` | `string` | The address to start admission controller on (ip:port).  Setting it to 'off' disables the admission controller. | `off` |
| `--anonymous-reports` | `bool` | Send anonymized usage data to help improve Kong. | `true` |
| `--apiserver-burst` | `int` | The Kubernetes API RateLimiter maximum burst queries per second. | `300` |
| `--apiserver-host` | `string` | The Kubernetes API server URL. If not set, the controller will use cluster config discovery. |  |
| `--apiserver-qps` | `int` | The Kubernetes API RateLimiter maximum queries per second. | `100` |
| `--cache-sync-timeout` | `duration` | The time limit set to wait for syncing controllers' caches. Leave this empty to use default from controller-runtime. | `0s` |
| `--dump-config` | `bool` | Enable config dumps via web interface host:10256/debug/config. | `false` |
| `--dump-sensitive-config` | `bool` | Include credentials and TLS secrets in configs exposed with --dump-config. | `false` |
| `--election-id` | `string` | Election id to use for status update. | `5b374a9e.konghq.com` |
| `--election-namespace` | `string` | Leader election namespace to use when running outside a cluster. |  |
| `--enable-controller-ingress-class-networkingv1` | `bool` | Enable the networking.k8s.io/v1 IngressClass controller. | `true` |
| `--enable-controller-ingress-class-parameters` | `bool` | Enable the IngressClassParameters controller. | `true` |
| `--enable-controller-ingress-networkingv1` | `bool` | Enable the networking.k8s.io/v1 Ingress controller. | `true` |
| `--enable-controller-kongclusterplugin` | `bool` | Enable the KongClusterPlugin controller. | `true` |
| `--enable-controller-kongconsumer` | `bool` | Enable the KongConsumer controller. . | `true` |
| `--enable-controller-kongingress` | `bool` | Enable the KongIngress controller. | `true` |
| `--enable-controller-kongplugin` | `bool` | Enable the KongPlugin controller. | `true` |
| `--enable-controller-service` | `bool` | Enable the Service controller. | `true` |
| `--enable-controller-tcpingress` | `bool` | Enable the TCPIngress controller. | `true` |
| `--enable-controller-udpingress` | `bool` | Enable the UDPIngress controller. | `true` |
| `--enable-reverse-sync` | `bool` | Send configuration to Kong even if the configuration checksum has not changed since previous update. | `false` |
| `--feature-gates` | `mapStringBool` | A set of key=value pairs that describe feature gates for alpha/beta/experimental features. See the Feature Gates documentation for information and available options: https://github.com/Kong/kubernetes-ingress-controller/blob/main/FEATURE_GATES.md. |  |
| `--gateway-api-controller-name` | `string` | The controller name to match on Gateway API resources. | `konghq.com/kic-gateway-controller` |
| `--gateway-discovery-dns-strategy` | `dns-strategy` | DNS strategy to use when creating Gateway's Admin API addresses. One of: ip, service, pod. | `"ip"` |
| `--health-probe-bind-address` | `string` | The address the probe endpoint binds to. | `:10254` |
| `--ingress-class` | `string` | Name of the ingress class to route through this controller. | `kong` |
| `--kong-admin-ca-cert` | `string` | PEM-encoded CA certificate to verify Kong's Admin SSL certificate. |  |
| `--kong-admin-ca-cert-file` | `string` | Path to PEM-encoded CA certificate file to verify Kong's Admin SSL certificate. |  |
| `--kong-admin-concurrency` | `int` | Max number of concurrent requests sent to Kong's Admin API. | `10` |
| `--kong-admin-filter-tag` | `stringSlice` | The tag used to manage and filter entities in Kong. This flag can be specified multiple times to specify multiple tags. This setting will be silently ignored if the Kong instance has no tags support. | `[managed-by-ingress-controller]` |
| `--kong-admin-header` | `stringSlice` | Add a header (key:value) to every Admin API call, this flag can be used multiple times to specify multiple headers. | `[]` |
| `--kong-admin-init-retries` | `uint` | Number of attempts that will be made initially on controller startup to connect to the Kong Admin API. | `60` |
| `--kong-admin-init-retry-delay` | `duration` | The time delay between every attempt (on controller startup) to connect to the Kong Admin API. | `1s` |
| `--kong-admin-svc` | `namespacedName` | Kong Admin API Service namespaced name in "namespace/name" format, to use for Kong Gateway service discovery. |  |
| `--kong-admin-svc-port-names` | `stringSlice` | Names of ports on Kong Admin API service to take into account when doing gateway discovery. | `[admin,admin-tls,kong-admin,kong-admin-tls]` |
| `--kong-admin-tls-client-cert` | `string` | MTLS client certificate for authentication. |  |
| `--kong-admin-tls-client-cert-file` | `string` | MTLS client certificate file for authentication. |  |
| `--kong-admin-tls-client-key` | `string` | MTLS client key for authentication. |  |
| `--kong-admin-tls-client-key-file` | `string` | MTLS client key file for authentication. |  |
| `--kong-admin-tls-server-name` | `string` | SNI name to use to verify the certificate presented by Kong in TLS. |  |
| `--kong-admin-tls-skip-verify` | `bool` | Disable verification of TLS certificate of Kong's Admin endpoint. | `false` |
| `--kong-admin-token` | `string` | The Kong Enterprise RBAC token used by the controller. |  |
| `--kong-admin-token-file` | `string` | Path to the Kong Enterprise RBAC token file used by the controller. |  |
| `--kong-admin-url` | `stringSlice` | Kong Admin URL(s) to connect to in the format "protocol://address:port". More than 1 URL can be provided, in such case the flag should be used multiple times or a corresponding env variable should use comma delimited addresses. | `[http://localhost:8001]` |
| `--kong-workspace` | `string` | Kong Enterprise workspace to configure. Leave this empty if not using Kong workspaces. |  |
| `--konnect-address` | `string` | Base address of Konnect API. | `https://us.kic.api.konghq.com` |
| `--konnect-control-plane-id` | `string` | An ID of a control plane that is to be synchronized with data plane configuration. |  |
| `--konnect-initial-license-polling-period` | `duration` | Polling period to be used before the first license is retrieved. | `1m0s` |
| `--konnect-license-polling-period` | `duration` | Polling period to be used after the first license is retrieved. | `12h0m0s` |
| `--konnect-licensing-enabled` | `bool` | Retrieve licenses from Konnect if available. Overrides licenses provided via the environment. | `false` |
| `--konnect-refresh-node-period` | `duration` | Period of uploading status of KIC and controlled kong gateway instances. | `1m0s` |
| `--konnect-sync-enabled` | `bool` | Enable synchronization of data plane configuration with a Konnect control plane. | `false` |
| `--konnect-tls-client-cert` | `string` | Konnect TLS client certificate. |  |
| `--konnect-tls-client-cert-file` | `string` | Konnect TLS client certificate file path. |  |
| `--konnect-tls-client-key` | `string` | Konnect TLS client key. |  |
| `--konnect-tls-client-key-file` | `string` | Konnect TLS client key file path. |  |
| `--kubeconfig` | `string` | Path to the kubeconfig file. |  |
| `--log-format` | `string` | Format of logs of the controller. Allowed values are text and json. | `text` |
| `--log-level` | `string` | Level of logging for the controller. Allowed values are trace, debug, info, and error. | `info` |
| `--metrics-bind-address` | `string` | The address the metric endpoint binds to. | `:10255` |
| `--profiling` | `bool` | Enable profiling via web interface host:10256/debug/pprof/. | `false` |
| `--proxy-sync-seconds` | `float32` | Define the rate (in seconds) in which configuration updates will be applied to the Kong Admin API. | `3` |
| `--proxy-timeout-seconds` | `float32` | Sets the timeout (in seconds) for all requests to Kong's Admin API. | `30` |
| `--publish-service` | `namespacedName` | Service fronting Ingress resources in "namespace/name" format. The controller will update Ingress status information with this Service's endpoints. |  |
| `--publish-service-udp` | `namespacedName` | Service fronting UDP routing resources in "namespace/name" format. The controller will update UDP route status information with this Service's endpoints. If omitted, the same Service will be used for both TCP and UDP routes. |  |
| `--publish-status-address` | `stringSlice` | User-provided addresses in comma-separated string format, for use in lieu of "publish-service" when that Service lacks useful address information (for example, in bare-metal environments). | `[]` |
| `--publish-status-address-udp` | `stringSlice` | User-provided address CSV, for use in lieu of "publish-service-udp" when that Service lacks useful address information. | `[]` |
| `--skip-ca-certificates` | `bool` | Disable syncing CA certificate syncing (for use with multi-workspace environments). | `false` |
| `--sync-period` | `duration` | Relist and confirm cloud resources this often. | `48h0m0s` |
| `--term-delay` | `duration` | The time delay to sleep before SIGTERM or SIGINT will shut down the Ingress Controller. | `0s` |
| `--update-status` | `bool` | Indicates if the ingress controller should update the status of resources (e.g. IP/Hostname for v1.Ingress, e.t.c.). | `true` |
| `--update-status-queue-buffer-size` | `int` | Buffer size of the underlying channels used to update the status of resources. | `8192` |
| `--watch-namespace` | `stringSlice` | Namespace(s) to watch for Kubernetes resources. Defaults to all namespaces. To watch multiple namespaces, use a comma-separated list of namespaces. | `[]` |

