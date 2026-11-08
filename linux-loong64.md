# linux/loong64 GTK backend

Fork of [egoist/mygo](https://github.com/egoist/mygo) **v0.7.4**. The module path stays `github.com/egoist/mygo`.

```
replace github.com/egoist/mygo => github.com/yibaiba/mygo v0.7.4-loong64.1
```

Upstream still treats `linux/loong64` as unsupported. This tree enables the linux GTK 3 path for `loong64` (`linux && (amd64 || arm64 || loong64)`).

## What changed vs v0.7.4

- Linux GTK 3 + WebKitGTK + OpenGL backend build tags include `loong64`.
- Backend is pure Go + [purego](https://github.com/ebitengine/purego) (`Dlopen` GTK 3 / GDK / GLib / Pango / Cairo / libepoxy). No `#cgo`.
- Vendored `third_party/purego` is v0.11.1 with `sys_loong64.s` `STACK_SIZE` **216** (upstream 208). Go’s loong64 assembler prepends an 8-byte LR save on `TEXT $0` that contains `JAL`; 208+8 is 8-byte aligned and GTK SIGSEGVs on LoongArch ABI 1.0. 216+8 = 224 keeps the C stack 16-byte aligned.
- `MYGO_GPU=0` uses grayscale font AA (GtkDrawingArea cannot dual-source LCD; subpixel coverage shows as mint/cyan).
- `gpuGL()` does not create a Gdk/GLX context unless Mesa has a DRI `.so` for the DRM driver. A render node alone (`gsgpu` without `gsgpu_dri.so`) can take down the X server (`libGL: driver (null)`).

darwin / windows / linux-amd64 / linux-arm64 tags are unchanged.

## Toolchain (LoongArch ABI 1.0)

ABI 1.0 hosts use `ld.so.1`, ELF Flags=`0x3`, `_NSIG=65`. Official go.dev `linux-loong64` is ABI 2.0 and will not run there. Use an ABI 1.0 Go 1.27.3:

https://github.com/yibaiba/go-loong64-abi1

```bash
export GOROOT=/usr/local/go PATH=$GOROOT/bin:$PATH GOTOOLCHAIN=local CGO_ENABLED=0
go build
```

Do not lower this module’s `go 1.27.3` line.

## Runtime libs

GTK 3, GLib, GDK, Pango, Cairo, fontconfig, FreeType matching ABI 1.0. libepoxy optional; without a usable GPU the CPU rasterizer is used (`MYGO_GPU=0`).
