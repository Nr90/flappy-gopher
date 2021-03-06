// Copyright 2017 Google Inc. All rights reserved.
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to writing, software distributed
// under the License is distributed on a "AS IS" BASIS, WITHOUT WARRANTIES OR
// CONDITIONS OF ANY KIND, either express or implied.
//
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"fmt"
	"os"
	"time"

	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

func main() {
	sdl.Main(func() {
		if err := run(); err != nil {
			fmt.Fprintf(os.Stderr, "%v", err)
			os.Exit(2)
		}
	})
}

func run() error {
	var err error
	sdl.Do(func() {
		err = sdl.Init(sdl.INIT_EVERYTHING)
	})
	if err != nil {
		return fmt.Errorf("could not initialize SDL: %v", err)
	}
	defer func() {
		sdl.Do(func() {
			sdl.Quit()
		})
	}()

	sdl.Do(func() {
		err = ttf.Init()
	})
	if err != nil {
		return fmt.Errorf("could not initialize TTF: %v", err)
	}
	defer func() {
		sdl.Do(func() {
			ttf.Quit()
		})
	}()

	var w *sdl.Window
	var r *sdl.Renderer
	sdl.Do(func() {
		w, r, err = sdl.CreateWindowAndRenderer(800, 600, sdl.WINDOW_SHOWN)
	})
	if err != nil {
		return fmt.Errorf("could not create window: %v", err)
	}
	defer func() {
		sdl.Do(func() {
			w.Destroy()
		})
	}()

	sdl.Do(func() {
		err = drawTitle(r, "Flappy Gopher")
	})
	if err != nil {
		return fmt.Errorf("could not draw title: %v", err)
	}

	time.Sleep(1 * time.Second)

	var s *scene
	sdl.Do(func() {
		s, err = newScene(r)
	})
	if err != nil {
		return fmt.Errorf("could not create scene: %v", err)
	}
	defer func() {
		sdl.Do(func() {
			s.destroy()
		})
	}()

	err = s.run(r)
	if err != nil {
		return fmt.Errorf("error while running: %v", err)
	}
	return nil
}

func drawTitle(r *sdl.Renderer, text string) error {
	r.Clear()

	f, err := ttf.OpenFont("res/fonts/Flappy.ttf", 20)
	if err != nil {
		return fmt.Errorf("could not load font: %v", err)
	}
	defer f.Close()

	c := sdl.Color{R: 255, G: 100, B: 0, A: 255}
	s, err := f.RenderUTF8_Solid(text, c)
	if err != nil {
		return fmt.Errorf("could not render title: %v", err)
	}
	defer s.Free()

	t, err := r.CreateTextureFromSurface(s)
	if err != nil {
		return fmt.Errorf("could not create texture: %v", err)
	}
	defer t.Destroy()

	if err := r.Copy(t, nil, nil); err != nil {
		return fmt.Errorf("could not copy texture: %v", err)
	}

	r.Present()

	return nil
}
