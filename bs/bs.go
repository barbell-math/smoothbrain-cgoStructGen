package main

import (
	sbbs "github.com/barbell-math/smoothbrain-bs"
)

func main() {
	sbbs.RegisterBsBuildTarget()
	sbbs.RegisterUpdateDepsTarget()
	sbbs.RegisterGoMarkDocTargets()
	sbbs.RegisterCommonGoCmdTargets(sbbs.GoTargets{
		GenericTestTarget:     true,
		GenericBenchTarget:    true,
		GenericFmtTarget:      true,
		GenericGenerateTarget: true,
	})
	sbbs.RegisterMergegateTarget(sbbs.MergegateTargets{
		PreStages: []sbbs.StageFunc{
			sbbs.TargetAsStage("goenumInstall"),
		},
		CheckDepsUpdated:     true,
		CheckReadmeGomarkdoc: true,
		CheckFmt:             true,
		CheckUnitTests:       true,
	})
	sbbs.Main("bs")
}
