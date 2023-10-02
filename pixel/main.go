package main

import (
    "github.com/faiface/pixel"
    "github.com/faiface/pixel/pixelgl"
)

func main() {
    pixelgl.Run(run)
}

func run() {
    cfg := pixelgl.WindowConfig{
        Title:  "Zineb!",
        Bounds: pixel.R(0, 0, 640, 320), // Modifiez les dimensions de la fenêtre ici
        VSync:  true,
    }

    win, err := pixelgl.NewWindow(cfg)
    if err != nil {
        panic(err)
    }

    for !win.Closed() {
        win.Update()
    }
}
