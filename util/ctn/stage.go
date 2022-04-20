package ctn

import "context"

type StageContext interface {
	context.Context
	// Returns compilation-wide DI container.
	GetGlobalContainer() DI
	GetStageContanier() DI
}

// Represents single compilation stage, which is implemented by user.
type Stage interface {
	// Initializes container.
	// Called before rewriting of previous container is performed.
	Initialize(ctx StageContext) (err error)

	// Performs actual compilation - registers DI types, which
	// are result of compilation.
	Compile(ctx StageContext) (err error)

	// Rewrites all entries from current DI container to next stage's DI container.
	// So that all redundant dependencies are dropped.
	//
	// Called after compile.
	Rewrite(ctx StageContext, newContainer DI) (err error)

	// Shutdowns DI container.
	// This function must perform any kind of cleanup.
	Shutdown(ctx StageContext) (err error)
}

// Stage, which implements it's functionality using functions provided by user.
// Functions, which are not provided are noop.
type SimpleStage struct {
	InitializeImpl func(ctx StageContext) (err error)
	CompileImpl    func(ctx StageContext, c DI) (err error)
	RewriteImpl    func(ctx StageContext, current DI, newContainer DI) (err error)
	ShutdownImpl   func(ctx StageContext, c DI) (err error)
}

func (ss *SimpleStage) Initialize(ctx StageContext) (err error) {
	if ss.InitializeImpl != nil {
		return ss.InitializeImpl(ctx)
	}
	return
}

func (ss *SimpleStage) Compile(ctx StageContext, c DI) (err error) {
	if ss.CompileImpl != nil {
		return ss.CompileImpl(ctx, c)
	}
	return
}

func (ss *SimpleStage) Rewrite(ctx StageContext, current DI, newContainer DI) (err error) {
	if ss.RewriteImpl != nil {
		return ss.RewriteImpl(ctx, current, newContainer)
	}
	return
}

func (ss *SimpleStage) Shutdown(ctx StageContext, c DI) (err error) {
	if ss.ShutdownImpl != nil {
		return ss.ShutdownImpl(ctx, c)
	}
	return
}

type stagedCompilationContext struct {
	context.Context
	stageDI  DI
	globalDI DI
}

func (sc *stagedCompilationContext) GetGlobalContainer() DI {
	return sc.globalDI
}
func (sc *stagedCompilationContext) GetStageContanier() DI {
	return sc.stageDI
}

type StagedCompiler struct {
	InitializeGlobalContainer func(ctx context.Context, globalDI DI) (err error)
	Stages                    []Stage
}

type StagedCompilerResult struct {
	LastStageDI DI
	GlobalDI    DI
}

func (sm *StagedCompiler) RunCompilation(ctx context.Context) (err error) {
	globalDI := NewDI()
	if sm.InitializeGlobalContainer != nil {
		sm.InitializeGlobalContainer(ctx, globalDI)
	}

	nextStageDI := NewDI()
	stageDI := NewDI()

	for _, s := range sm.Stages {
		stageContext := &stagedCompilationContext{
			Context:  ctx,
			stageDI:  stageDI,
			globalDI: globalDI,
		}

		err = s.Initialize(stageContext)
		if err != nil {
			return
		}

		err = s.Compile(stageContext)
		if err != nil {
			return
		}

		err = s.Rewrite(stageContext, nextStageDI)
		if err != nil {
			return
		}

		stageDI = nextStageDI
		nextStageDI = NewDI()
	}

	return
}
