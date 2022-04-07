package main

import (
	"github.com/go-gl/glfw/v3.3/glfw"
	vk "github.com/vulkan-go/vulkan"
)

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

	createInfo := vk.InstanceCreateInfo{
		SType:                   vk.StructureTypeInstanceCreateInfo,
		PNext:                   nil,
		PApplicationInfo:        &appInfo,
		Flags:                   0,
		EnabledExtensionCount:   uint32(len(glfwExtensions)),
		PpEnabledExtensionNames: glfwExtensions,
		EnabledLayerCount:       0,
	}

	var instance vk.Instance = nil
	res := vk.CreateInstance(&createInfo, nil, &instance)

	if res != vk.Success {
		window.Destroy()
		return nil, vk.Error(res)
	}

	if err = checkExtensions(glfwExtensions); err != nil {
		vk.DestroyInstance(instance, nil)
		window.Destroy()
		glfw.Terminate()
		return nil, err
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

func checkExtensions(availableExts []string) error {

	extMap := map[string]struct{}{}
	for _, ext := range availableExts {
		extMap[ext] = struct{}{}
	}

	var extensionCount uint32 = 0
	vk.EnumerateInstanceExtensionProperties("", &extensionCount, nil)

	var extensions []vk.ExtensionProperties = make([]vk.ExtensionProperties, extensionCount)
	res := vk.EnumerateInstanceExtensionProperties("", &extensionCount, extensions)

	if res != vk.Success {
		return vk.Error(res)
	}

	allPresent := true
	extName := ""

	for _, ext := range extensions {
		ext.Deref()
		if _, ok := extMap[string(C.GoString(ext.ExtensionName[:]))]; !ok {
			allPresent = false
			extName = string(ext.ExtensionName[:])
			break
		}
	}

	if !allPresent {
		panic(extName)
	}

	return nil
}
