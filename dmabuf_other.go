//go:build !(linux && (amd64 || arm64 || loong64))

package mygo

func linkDMABUFProbe() {}
