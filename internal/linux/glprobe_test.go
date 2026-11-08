//go:build linux && (amd64 || arm64 || loong64)

package linux

import (
	"os"
	"path/filepath"
	"testing"
)

func TestGPUDevice(t *testing.T) {
	system := func(t *testing.T, devices []string, drivers map[string]string) string {
		root := t.TempDir()
		for _, d := range devices {
			p := filepath.Join(root, "dev", d)
			os.MkdirAll(filepath.Dir(p), 0o755)
			os.WriteFile(p, nil, 0o644)
		}
		for node, driver := range drivers {
			dev := filepath.Join(root, "sys", "class", "drm", node, "device")
			os.MkdirAll(dev, 0o755)
			if err := os.Symlink(filepath.Join("..", "..", "bus", "pci", "drivers", driver), filepath.Join(dev, "driver")); err != nil {
				t.Fatal(err)
			}
		}
		return root
	}
	for _, c := range []struct {
		name    string
		devices []string
		drivers map[string]string
		gallium string
		want    bool
	}{
		{"nothing", nil, nil, "", false},
		{"Intel", []string{"dri/renderD128"}, map[string]string{"renderD128": "i915"}, "", true},
		{"AMD and simpledrm", nil, map[string]string{"renderD128": "simpledrm", "renderD129": "amdgpu"}, "", true},
		{"a virtual machine", []string{"dri/card0"}, map[string]string{"renderD128": "bochs-drm"}, "", false},
		{"NVIDIA", []string{"nvidia0"}, nil, "", true},
		{"WSL", []string{"dxg"}, nil, "", false},
		{"WSL with d3d12", []string{"dxg"}, nil, "d3d12", true},
	} {
		t.Run(c.name, func(t *testing.T) {
			if got := gpuDeviceIn(system(t, c.devices, c.drivers), c.gallium); got != c.want {
				t.Errorf("gpuDevice = %v, want %v", got, c.want)
			}
		})
	}
}

func TestCanCreateGLContext(t *testing.T) {
	root := func(t *testing.T, driver, dri string, nvidia bool) string {
		t.Helper()
		r := t.TempDir()
		if nvidia {
			p := filepath.Join(r, "dev", "nvidia0")
			os.MkdirAll(filepath.Dir(p), 0o755)
			os.WriteFile(p, nil, 0o644)
			return r
		}
		dev := filepath.Join(r, "sys", "class", "drm", "renderD128", "device")
		os.MkdirAll(dev, 0o755)
		if err := os.Symlink(filepath.Join("..", "..", "bus", "pci", "drivers", driver), filepath.Join(dev, "driver")); err != nil {
			t.Fatal(err)
		}
		if dri != "" {
			dir := filepath.Join(r, "usr", "lib", "x86_64-linux-gnu", "dri")
			os.MkdirAll(dir, 0o755)
			os.WriteFile(filepath.Join(dir, dri+"_dri.so"), nil, 0o644)
		}
		return r
	}
	if canCreateGLContext(root(t, "gsgpu", "", false)) {
		t.Fatal("gsgpu without DRI must not create a GLX context")
	}
	if !canCreateGLContext(root(t, "gsgpu", "gsgpu", false)) {
		t.Fatal("gsgpu_dri.so allows probe")
	}
	if !canCreateGLContext(root(t, "i915", "iris", false)) {
		t.Fatal("i915 + iris_dri.so allows probe")
	}
	if canCreateGLContext(root(t, "i915", "", false)) {
		t.Fatal("i915 without DRI must not create a GLX context")
	}
	if !canCreateGLContext(root(t, "", "", true)) {
		t.Fatal("NVIDIA device allows probe")
	}
	if canCreateGLContext(root(t, "bochs-drm", "bochs-drm", false)) {
		t.Fatal("display-only DRM must not create a GLX context")
	}
}

func TestSoftwareGLMesaX11(t *testing.T) {
	if !softwareGL("Mesa X11") {
		t.Fatal("Mesa X11 is CPU GL")
	}
	if softwareGL("Mesa Intel(R) UHD Graphics") {
		t.Fatal("Intel Mesa is not CPU GL")
	}
}
