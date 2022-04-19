package main

import (
	"bytes"
	"fmt"
	"unsafe"

	"github.com/go-gl/glfw/v3.3/glfw"
	vk "github.com/vulkan-go/vulkan"
)

type DebugUtilsMessengerCallbackDataFlags uint32

type DebugUtilsMessengerCallbackData struct {
	SType            vk.StructureType
	PNext            unsafe.Pointer
	Flags            DebugUtilsMessengerCallbackDataFlags
	PMessageIdName   string
	MessageIdNumber  int32
	PMessage         string
	QueueLabelCount  uint32
	PQueueLabels     []vk.DebugUtilsLabel
	CmdBufLabelCount uint32
	PCmdBufLabels    []vk.DebugUtilsLabel
	ObjectCount      uint32
	PObject          vk.DebugUtilsObjectNameInfo
}

type DebugUtilsMessengerCreateFlags uint32
type DebugUtilsMessengerCallback func(messageSeverity vk.DebugUtilsMessageSeverityFlagBits, messageTypes vk.DebugUtilsMessageTypeFlags, pCallbackData DebugUtilsMessengerCallbackData, pUserData unsafe.Pointer) vk.Result

type DebugUtilsMessengerCreateInfo struct {
	SType           vk.StructureType
	PNext           unsafe.Pointer
	Flags           DebugUtilsMessengerCreateFlags
	MessageSeverity vk.DebugUtilsMessageSeverityFlags
	MessageType     vk.DebugUtilsMessageTypeFlags
	PfnUserCallback DebugUtilsMessengerCallback
	PUserData       unsafe.Pointer
}

type DebugUtilsMessenger uint64

type Renderer struct {
	glfwWindow *glfw.Window
	vkInstance vk.Instance
}

func NewRenderer() (*Renderer, error) {

	if err := glfw.Init(); err != nil {
		return nil, err
	}

	glfw.WindowHint(glfw.ClientAPI, glfw.NoAPI)
	glfw.WindowHint(glfw.Resizable, glfw.False)

	window, err := glfw.CreateWindow(1280, 720, "My Vulkan Renderer in Go", nil, nil)
	if err != nil {
		return nil, err
	}

	vk.SetGetInstanceProcAddr(glfw.GetVulkanGetInstanceProcAddress())

	if err := vk.Init(); err != nil {
		window.Destroy()
		glfw.Terminate()
		return nil, err
	}

	appInfo := vk.ApplicationInfo{
		SType:              vk.StructureTypeApplicationInfo,
		PNext:              nil,
		PApplicationName:   "My Renderer",
		PEngineName:        "No Engine",
		EngineVersion:      vk.MakeVersion(1, 0, 0),
		ApplicationVersion: vk.MakeVersion(1, 0, 0),
		ApiVersion:         vk.ApiVersion11,
	}

	glfwExtensions := glfw.GetCurrentContext().GetRequiredInstanceExtensions()

	extensions := glfwExtensions
	extensions = append(extensions, vk.ExtDebugUtilsExtensionName+"\x00")

	createInfo := vk.InstanceCreateInfo{
		SType:                   vk.StructureTypeInstanceCreateInfo,
		PNext:                   nil,
		PApplicationInfo:        &appInfo,
		Flags:                   0,
		EnabledExtensionCount:   uint32(len(extensions)),
		PpEnabledExtensionNames: extensions,
		EnabledLayerCount:       0,
	}

	validationLayers := getValidationLayers()

	createInfo.EnabledLayerCount = uint32(len(validationLayers))
	createInfo.PpEnabledLayerNames = validationLayers

	var instance vk.Instance = nil
	res := vk.CreateInstance(&createInfo, nil, &instance)

	if res != vk.Success {
		window.Destroy()
		return nil, vk.Error(res)
	}

	r := Renderer{
		glfwWindow: window,
		vkInstance: instance,
	}

	return &r, nil
}

func (renderer Renderer) Delete() {
	vk.DestroyInstance(renderer.vkInstance, nil)
	renderer.glfwWindow.Destroy()
	glfw.Terminate()
}

func getValidationLayers() []string {

	validationLayers := []string{
		"VK_LAYER_KHRONOS_validation\x00",
	}

	if checkValidationLayersSupport(validationLayers) {
		return validationLayers
	} else {
		return []string{}
	}
}

func checkValidationLayersSupport(neededLayers []string) bool {

	layerMap := map[string]struct{}{}

	var layerCount uint32 = 0
	vk.EnumerateInstanceLayerProperties(&layerCount, nil)

	var layers []vk.LayerProperties = make([]vk.LayerProperties, layerCount)
	res := vk.EnumerateInstanceLayerProperties(&layerCount, layers)

	for _, layer := range layers {
		layer.Deref()
		layerName := string(bytes.Trim(layer.LayerName[:], "\x00")) + "\x00"
		layerMap[layerName] = struct{}{}
	}

	if res != vk.Success {
		return false
	}

	allPresent := true

	for _, ext := range neededLayers {
		if _, ok := layerMap[ext]; !ok {
			allPresent = false
			fmt.Printf("Coult not find layer %v\n", ext)
			break
		} else {
			fmt.Printf("Found layer %v\n", ext)
		}
	}

	return allPresent
}

func debugCallback(messageSeverity vk.DebugUtilsMessageSeverityFlagBits, messageTypes vk.DebugUtilsMessageTypeFlags, pCallbackData DebugUtilsMessengerCallbackData, pUserData unsafe.Pointer) vk.Result {
	fmt.Println("msg")
	return vk.False
}

func setupDebugMessenger(instance vk.Instance) {
	var messenger DebugUtilsMessenger

	createInfo := DebugUtilsMessengerCreateInfo{
		SType:           vk.StructureTypeDebugUtilsMessengerCreateInfo,
		PNext:           nil,
		MessageSeverity: vk.DebugUtilsMessageSeverityFlags(vk.DebugUtilsMessageSeverityVerboseBit | vk.DebugUtilsMessageSeverityWarningBit | vk.DebugUtilsMessageSeverityErrorBit),
		MessageType:     vk.DebugUtilsMessageTypeFlags(vk.DebugUtilsMessageTypeGeneralBit | vk.DebugUtilsMessageTypeValidationBit | vk.DebugUtilsMessageTypePerformanceBit),
		PfnUserCallback: debugCallback,
		PUserData:       nil,
	}

	createDebugUtilsMessengerProc := vk.GetInstanceProcAddr(instance, "vkCreateDebugUtilsMessengerEXT")
}

func CreateDebugUtilsMessengerEXT(instance vk.Instance, const vk.DebugUtilsMess) {
}