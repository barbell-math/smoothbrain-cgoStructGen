package main

import (
	sbbs "github.com/barbell-math/smoothbrain-bs"
)

func main() {
	sbbs.RegisterBsBuildTarget()
	sbbs.RegisterUpdateDepsTarget()
	sbbs.RegisterGoMarkDocTargets()
	sbbs.RegisterGoEnumTargets()
	sbbs.RegisterCommonGoCmdTargets(sbbs.NewGoTargets().
		DefaultFmtTarget().
		DefaultGenerateTarget().
		DefaultTestTarget(),
	)
	sbbs.RegisterMergegateTarget(sbbs.MergegateTargets{
		PreStages: []sbbs.StageFunc{
			sbbs.TargetAsStage("goenumInstall"),
		},
		CheckDepsUpdated:     true,
		CheckReadmeGomarkdoc: true,
		FmtTarget:            sbbs.DefaultFmtTargetName,
		TestTarget:           sbbs.DefaultFmtTargetName,
	})
	sbbs.Main("bs")
}
