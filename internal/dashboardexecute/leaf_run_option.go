package dashboardexecute

type LeafRunOption = func(target *LeafRun)

func withName(name string) LeafRunOption {
	return func(target *LeafRun) {
		target.Name = name
	}
}
