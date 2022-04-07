package main

import (
	"github.com/go-gl/glfw/v3.3/glfw"
)

func main() {
	renderer, err := NewRenderer()
	if err != nil {
		panic(err)
	}

	for !renderer.glfwWindow.ShouldClose() {
		glfw.PollEvents()
	}

	defer renderer.Delete()
}
