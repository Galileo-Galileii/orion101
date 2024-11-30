//go:generate go run github.com/orion101-ai/nah/cmd/deepcopy ./pkg/storage/apis/orion101.orion101.ai/v1/
//go:generate go run github.com/orion101-ai/nah/cmd/deepcopy ./apiclient/types/
//go:generate go run k8s.io/kube-openapi/cmd/openapi-gen --go-header-file tools/header.txt --output-file openapi_generated.go --output-dir ./pkg/storage/openapi/generated/ --output-pkg github.com/orion101-ai/orion101/pkg/storage/openapi/generated github.com/orion101-ai/orion101/pkg/storage/apis/orion101.orion101.ai/v1 k8s.io/apimachinery/pkg/apis/meta/v1 k8s.io/apimachinery/pkg/runtime k8s.io/apimachinery/pkg/version k8s.io/apimachinery/pkg/api/resource k8s.io/apimachinery/pkg/util/intstr k8s.io/api/coordination/v1 github.com/orion101-ai/orion101/apiclient/types

package main
