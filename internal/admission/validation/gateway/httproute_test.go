package gateway

import (
	"context"
	"fmt"
	"testing"

	"github.com/kong/go-kong/kong"
	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/kong/kubernetes-ingress-controller/v2/internal/dataplane/parser"
	"github.com/kong/kubernetes-ingress-controller/v2/internal/gatewayapi"
)

func TestValidateHTTPRoute(t *testing.T) {
	var (
		nonexistentListener = gatewayapi.SectionName("listener-that-doesnt-exist")
		group               = gatewayapi.Group("gateway.networking.k8s.io")
		defaultGWNamespace  = gatewayapi.Namespace(corev1.NamespaceDefault)
		exampleGroup        = gatewayapi.Group("example")
		podKind             = gatewayapi.Kind("Pod")
	)

	for _, tt := range []struct {
		msg           string
		route         *gatewayapi.HTTPRoute
		gateways      []*gatewayapi.Gateway
		valid         bool
		validationMsg string
		err           error
	}{
		{
			msg: "if you provide errant gateways for validation, it fails validation",
			route: &gatewayapi.HTTPRoute{
				ObjectMeta: metav1.ObjectMeta{
					Namespace: corev1.NamespaceDefault,
					Name:      "testing-httproute",
				},
			}, // no parentRefs
			gateways: []*gatewayapi.Gateway{{
				ObjectMeta: metav1.ObjectMeta{
					Namespace: corev1.NamespaceDefault,
					Name:      "testing-gateway",
				},
				Spec: gatewayapi.GatewaySpec{
					Listeners: []gatewayapi.Listener{{
						Name:     "http",
						Port:     80,
						Protocol: (gatewayapi.HTTPProtocolType),
						AllowedRoutes: &gatewayapi.AllowedRoutes{
							Kinds: []gatewayapi.RouteGroupKind{{
								Group: &group,
								Kind:  "HTTPRoute",
							}},
						},
					}},
				},
			}},
			valid:         false,
			validationMsg: "couldn't determine parentRefs for httproute",
			err:           fmt.Errorf("no parentRef matched gateway default/testing-gateway"),
		},
		{
			msg: "if you use sectionname to attach to a non-existent gateway listener, it fails validation",
			route: &gatewayapi.HTTPRoute{
				ObjectMeta: metav1.ObjectMeta{
					Namespace: corev1.NamespaceDefault,
					Name:      "testing-httproute",
				},
				Spec: gatewayapi.HTTPRouteSpec{
					CommonRouteSpec: gatewayapi.CommonRouteSpec{
						ParentRefs: []gatewayapi.ParentReference{{
							Name:        "testing-gateway",
							SectionName: &nonexistentListener,
						}},
					},
				},
			},
			gateways: []*gatewayapi.Gateway{{
				ObjectMeta: metav1.ObjectMeta{
					Namespace: corev1.NamespaceDefault,
					Name:      "testing-gateway",
				},
				Spec: gatewayapi.GatewaySpec{
					Listeners: []gatewayapi.Listener{{
						Name:     "not-the-right-listener",
						Port:     80,
						Protocol: (gatewayapi.HTTPProtocolType),
						AllowedRoutes: &gatewayapi.AllowedRoutes{
							Kinds: []gatewayapi.RouteGroupKind{{
								Group: &group,
								Kind:  "HTTPRoute",
							}},
						},
					}},
				},
			}},
			valid:         false,
			validationMsg: "couldn't find gateway listeners for httproute",
			err:           fmt.Errorf("sectionname referenced listener listener-that-doesnt-exist was not found on gateway default/testing-gateway"),
		},
		{
			msg: "if the provided gateway has NO listeners, the HTTPRoute fails validation",
			route: &gatewayapi.HTTPRoute{
				ObjectMeta: metav1.ObjectMeta{
					Namespace: corev1.NamespaceDefault,
					Name:      "testing-httproute",
				},
				Spec: gatewayapi.HTTPRouteSpec{
					CommonRouteSpec: gatewayapi.CommonRouteSpec{
						ParentRefs: []gatewayapi.ParentReference{{
							Name: "testing-gateway",
						}},
					},
				},
			},
			gateways: []*gatewayapi.Gateway{{
				ObjectMeta: metav1.ObjectMeta{
					Namespace: corev1.NamespaceDefault,
					Name:      "testing-gateway",
				},
				Spec: gatewayapi.GatewaySpec{
					Listeners: []gatewayapi.Listener{},
				},
			}},
			valid:         false,
			validationMsg: "couldn't find gateway listeners for httproute",
			err:           fmt.Errorf("no listeners could be found for gateway default/testing-gateway"),
		},
		{
			msg: "parentRefs which omit the namespace pass validation in the same namespace",
			route: &gatewayapi.HTTPRoute{
				ObjectMeta: metav1.ObjectMeta{
					Namespace: corev1.NamespaceDefault,
					Name:      "testing-httproute",
				},
				Spec: gatewayapi.HTTPRouteSpec{
					CommonRouteSpec: gatewayapi.CommonRouteSpec{
						ParentRefs: []gatewayapi.ParentReference{{
							Name: "testing-gateway",
						}},
					},
				},
			},
			gateways: []*gatewayapi.Gateway{{
				ObjectMeta: metav1.ObjectMeta{
					Namespace: corev1.NamespaceDefault,
					Name:      "testing-gateway",
				},
				Spec: gatewayapi.GatewaySpec{
					Listeners: []gatewayapi.Listener{{
						Name:     "http",
						Port:     80,
						Protocol: (gatewayapi.HTTPProtocolType),
						AllowedRoutes: &gatewayapi.AllowedRoutes{
							Kinds: []gatewayapi.RouteGroupKind{{
								Group: &group,
								Kind:  "HTTPRoute",
							}},
						},
					}},
				},
			}},
			valid: true,
		},
		{
			msg: "if the gateway listener doesn't support HTTPRoute, validation fails",
			route: &gatewayapi.HTTPRoute{
				ObjectMeta: metav1.ObjectMeta{
					Namespace: corev1.NamespaceDefault,
					Name:      "testing-httproute",
				},
				Spec: gatewayapi.HTTPRouteSpec{
					CommonRouteSpec: gatewayapi.CommonRouteSpec{
						ParentRefs: []gatewayapi.ParentReference{{
							Name: "testing-gateway",
						}},
					},
				},
			},
			gateways: []*gatewayapi.Gateway{{
				ObjectMeta: metav1.ObjectMeta{
					Namespace: corev1.NamespaceDefault,
					Name:      "testing-gateway",
				},
				Spec: gatewayapi.GatewaySpec{
					Listeners: []gatewayapi.Listener{{
						Name:     "http-alternate",
						Port:     8000,
						Protocol: (gatewayapi.HTTPProtocolType),
						AllowedRoutes: &gatewayapi.AllowedRoutes{
							Kinds: []gatewayapi.RouteGroupKind{{
								Group: &group,
								Kind:  "TCPRoute",
							}},
						},
					}},
				},
			}},
			valid:         false,
			validationMsg: "httproute linked gateway listeners did not pass validation",
			err:           fmt.Errorf("HTTPRoute not supported by listener http-alternate"),
		},
		{
			msg: "if an HTTPRoute is using queryparams matching it fails validation due to only supporting expression router",
			route: &gatewayapi.HTTPRoute{
				ObjectMeta: metav1.ObjectMeta{
					Namespace: corev1.NamespaceDefault,
					Name:      "testing-httproute",
				},
				Spec: gatewayapi.HTTPRouteSpec{
					CommonRouteSpec: gatewayapi.CommonRouteSpec{
						ParentRefs: []gatewayapi.ParentReference{{
							Name: "testing-gateway",
						}},
					},
					Rules: []gatewayapi.HTTPRouteRule{{
						Matches: []gatewayapi.HTTPRouteMatch{{
							QueryParams: []gatewayapi.HTTPQueryParamMatch{{
								Name:  "user-agent",
								Value: "netscape navigator",
							}},
						}},
						BackendRefs: []gatewayapi.HTTPBackendRef{{
							BackendRef: gatewayapi.BackendRef{
								BackendObjectReference: gatewayapi.BackendObjectReference{
									Namespace: &defaultGWNamespace,
								},
							},
						}},
					}},
				},
			},
			gateways: []*gatewayapi.Gateway{{
				ObjectMeta: metav1.ObjectMeta{
					Namespace: corev1.NamespaceDefault,
					Name:      "testing-gateway",
				},
				Spec: gatewayapi.GatewaySpec{
					Listeners: []gatewayapi.Listener{{
						Name:     "http",
						Port:     80,
						Protocol: (gatewayapi.HTTPProtocolType),
						AllowedRoutes: &gatewayapi.AllowedRoutes{
							Kinds: []gatewayapi.RouteGroupKind{{
								Group: &group,
								Kind:  "HTTPRoute",
							}},
						},
					}},
				},
			}},
			valid:         false,
			validationMsg: "httproute spec did not pass validation",
			err:           fmt.Errorf("queryparam matching is supported with expression router only"),
		},
		{
			msg: "we don't support any group except core kubernetes for backendRefs",
			route: &gatewayapi.HTTPRoute{
				ObjectMeta: metav1.ObjectMeta{
					Namespace: corev1.NamespaceDefault,
					Name:      "testing-httproute",
				},
				Spec: gatewayapi.HTTPRouteSpec{
					CommonRouteSpec: gatewayapi.CommonRouteSpec{
						ParentRefs: []gatewayapi.ParentReference{{
							Name: "testing-gateway",
						}},
					},
					Rules: []gatewayapi.HTTPRouteRule{{
						Matches: []gatewayapi.HTTPRouteMatch{{
							Headers: []gatewayapi.HTTPHeaderMatch{{
								Name:  "Content-Type",
								Value: "audio/vorbis",
							}},
						}},
						BackendRefs: []gatewayapi.HTTPBackendRef{
							{
								BackendRef: gatewayapi.BackendRef{
									BackendObjectReference: gatewayapi.BackendObjectReference{
										Group:     &exampleGroup,
										Kind:      &podKind,
										Namespace: &defaultGWNamespace,
										Name:      "service1",
									},
								},
							},
						},
					}},
				},
			},
			gateways: []*gatewayapi.Gateway{{
				ObjectMeta: metav1.ObjectMeta{
					Namespace: corev1.NamespaceDefault,
					Name:      "testing-gateway",
				},
				Spec: gatewayapi.GatewaySpec{
					Listeners: []gatewayapi.Listener{{
						Name:     "http",
						Port:     80,
						Protocol: (gatewayapi.HTTPProtocolType),
						AllowedRoutes: &gatewayapi.AllowedRoutes{
							Kinds: []gatewayapi.RouteGroupKind{{
								Group: &group,
								Kind:  "HTTPRoute",
							}},
						},
					}},
				},
			}},
			valid:         false,
			validationMsg: "httproute spec did not pass validation",
			err:           fmt.Errorf("example is not a supported group for httproute backendRefs, only core is supported"),
		},
		{
			msg: "we don't support any core kind except Service for backendRefs",
			route: &gatewayapi.HTTPRoute{
				ObjectMeta: metav1.ObjectMeta{
					Namespace: corev1.NamespaceDefault,
					Name:      "testing-httproute",
				},
				Spec: gatewayapi.HTTPRouteSpec{
					CommonRouteSpec: gatewayapi.CommonRouteSpec{
						ParentRefs: []gatewayapi.ParentReference{{
							Name: "testing-gateway",
						}},
					},
					Rules: []gatewayapi.HTTPRouteRule{{
						Matches: []gatewayapi.HTTPRouteMatch{{
							Headers: []gatewayapi.HTTPHeaderMatch{{
								Name:  "Content-Type",
								Value: "audio/vorbis",
							}},
						}},
						BackendRefs: []gatewayapi.HTTPBackendRef{
							{
								BackendRef: gatewayapi.BackendRef{
									BackendObjectReference: gatewayapi.BackendObjectReference{
										Kind:      &podKind,
										Namespace: &defaultGWNamespace,
										Name:      "service1",
									},
								},
							},
						},
					}},
				},
			},
			gateways: []*gatewayapi.Gateway{{
				ObjectMeta: metav1.ObjectMeta{
					Namespace: corev1.NamespaceDefault,
					Name:      "testing-gateway",
				},
				Spec: gatewayapi.GatewaySpec{
					Listeners: []gatewayapi.Listener{{
						Name:     "http",
						Port:     80,
						Protocol: (gatewayapi.HTTPProtocolType),
						AllowedRoutes: &gatewayapi.AllowedRoutes{
							Kinds: []gatewayapi.RouteGroupKind{{
								Group: &group,
								Kind:  "HTTPRoute",
							}},
						},
					}},
				},
			}},
			valid:         false,
			validationMsg: "httproute spec did not pass validation",
			err:           fmt.Errorf("Pod is not a supported kind for httproute backendRefs, only Service is supported"),
		},
	} {
		// Passed routesValidator is irrelevant for the above test cases.
		valid, validMsg, err := ValidateHTTPRoute(
			context.Background(), mockRoutesValidator{}, parser.FeatureFlags{}, tt.route, tt.gateways...,
		)
		assert.Equal(t, tt.valid, valid, tt.msg)
		assert.Equal(t, tt.validationMsg, validMsg, tt.msg)
		assert.Equal(t, tt.err, err, tt.msg)
	}
}

type mockRoutesValidator struct{}

func (mockRoutesValidator) Validate(_ context.Context, _ *kong.Route) (bool, string, error) {
	return true, "", nil
}
