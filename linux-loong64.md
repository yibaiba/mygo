# Linux GTK (Kylin / UOS / Loongnix)

Fork of [egoist/mygo](https://github.com/egoist/mygo) **v0.7.4**. The module path stays `github.com/egoist/mygo`.

```
replace github.com/egoist/mygo => github.com/yibaiba/mygo v0.7.4-loong64.1
```

Works on **linux/amd64, linux/arm64, and linux/loong64** — Kylin, UOS, Loongnix, and other GTK 3 desktops. It is not Loongson-only. Upstream already had GTK on amd64 and arm64; this tree adds **loong64** and Linux GPU/font fixes that apply on every arch.

## What changed vs v0.7.4

- Linux GTK 3 backend build tags include `loong64` (`linux && (amd64 || arm64 || loong64)`).
- Backend is pure Go + [purego](https://github.com/ebitengine/purego) (`Dlopen` GTK 3 / GDK / GLib / Pango / Cairo / libepoxy). No `#cgo`.
- `gpuGL()` does not create a Gdk/GLX context unless Mesa has a DRI `.so` for the DRM driver. A render node alone (for example `gsgpu` without `gsgpu_dri.so`, or similar on other SoCs) can take down the X server (`libGL: driver (null)`).
- `MYGO_GPU=0` uses grayscale font AA (GtkDrawingArea cannot dual-source LCD; subpixel coverage shows as mint/cyan). This is the CPU path on any Linux, not only loong64.
- Vendored `third_party/purego` is v0.11.1 with `sys_loong64.s` `STACK_SIZE` **216** (upstream 208). Needed on LoongArch ABI 1.0 so GTK does not SIGSEGV; other GOARCH ignore that file.

darwin / windows tags are unchanged.

## Toolchain

- **amd64 / arm64:** official Go 1.27.3 is fine (`GOTOOLCHAIN=go1.27.3`).
- **loong64 ABI 1.0** (`ld.so.1`, ELF Flags=`0x3`, `_NSIG=65`): official go.dev `linux-loong64` is ABI 2.0 and will not run. Use https://github.com/yibaiba/go-loong64-abi1

```bash
export GOROOT=/usr/local/go PATH=$GOROOT/bin:$PATH GOTOOLCHAIN=local CGO_ENABLED=0
go build
```

Do not lower this module’s `go 1.27.3` line.

## Runtime libs

GTK 3, GLib, GDK, Pango, Cairo, fontconfig, FreeType. Match the host ABI (including LoongArch ABI 1.0). libepoxy optional; without a usable GPU the CPU rasterizer is used (`MYGO_GPU=0`).
