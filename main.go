package main

import (
	"fmt"

	"github.com/go-gl/glfw/v3.3/glfw"
	vk "github.com/vulkan-go/vulkan"
)

func main() {
	renderer, err := NewRenderer()
	if err != nil {
		panic(err)
	}

	var extensionCount uint32 = 0
	vk.EnumerateInstanceExtensionProperties("", &extensionCount, nil)

	var extensions []vk.ExtensionProperties = make([]vk.ExtensionProperties, extensionCount)
	res := vk.EnumerateInstanceExtensionProperties("", &extensionCount, extensions)

	if res != vk.Success {
		panic(err)
	}

	for _, ext := range extensions {
		ext.Deref()
		fmt.Println(vk.ToString(ext.ExtensionName[:]))
	}

	for !renderer.glfwWindow.ShouldClose() {
		glfw.PollEvents()
	}

	defer renderer.Delete()
}
