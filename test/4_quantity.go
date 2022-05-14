package main

//
//import (
//	"fmt"
//	apiv1 "k8s.io/api/core/v1"
//	"k8s.io/apimachinery/pkg/api/resource"
//)
//
//func main() {
//	strin := "Cpu:1"
//	var res apiv1.ResourceRequirements
//	resourceLimit := make(map[apiv1.ResourceName]resource.Quantity)
//	totalClaimedQuant := resource.MustParse(strin)
//	res.Limits = resourceLimit
//	fmt.Println(res)
//}
//resource
//var res apiv1.ResourceRequirements
//resourceLimit := make(map[apiv1.ResourceName]resource.Quantity)
//var tmpcpu apiv1.ResourceName = "cpu"
//var tmpmem apiv1.ResourceName = "memory"
//var tmpgpu apiv1.ResourceName = "nvidia.com/gpu"
//var relcpu string = data.Spec.Template.Spec.Containers[i].Resources.Limits.Cpu
//var relmem string = data.Spec.Template.Spec.Containers[i].Resources.Limits.Memory
//var relgpu string = data.Spec.Template.Spec.Containers[i].Resources.Limits.NvidiaGpu
//resourceLimit[tmpcpu] = resource.MustParse(relcpu)
//resourceLimit[tmpmem] = resource.MustParse(relmem)
//resourceLimit[tmpgpu] = resource.MustParse(relgpu)
//res.Limits = resourceLimit
//
////request
//resourceRequest := make(map[apiv1.ResourceName]resource.Quantity)
//resourceRequest[tmpcpu] = resource.MustParse(relcpu)
//resourceRequest[tmpmem] = resource.MustParse(relmem)
//resourceRequest[tmpgpu] = resource.MustParse(relgpu)
//res.Requests = resourceRequest
