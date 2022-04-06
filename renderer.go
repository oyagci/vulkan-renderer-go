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
