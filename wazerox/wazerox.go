package wazerox

import (
	"context"
	"github.com/tetratelabs/wazero"
	"github.com/tetratelabs/wazero/api"
	"os"
)

// RunnableModule is a convenient wrapper around a pair of wazero.Runtime and wazero.CompiledModule
type RunnableModule struct {
	runtime wazero.Runtime
	module  wazero.CompiledModule
}

func (m RunnableModule) Instantiate(ctx context.Context, wconf wazero.ModuleConfig) (api.Module, error) {
	return m.runtime.InstantiateModule(ctx, m.module, wconf)
}

func ReadModuleFromPath(ctx context.Context, rt wazero.Runtime, wasmPath string) (*RunnableModule, error) {
	wasm, err := os.ReadFile(wasmPath)
	if err != nil {
		return nil, err
	}

	code, err := rt.CompileModule(ctx, wasm)
	if err != nil {
		return nil, err
	}

	module := RunnableModule{
		runtime: rt,
		module:  code,
	}
	return &module, nil
}
