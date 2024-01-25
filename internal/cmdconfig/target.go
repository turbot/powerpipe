package cmdconfig

// UpdateTargetConnectionParams determines if the target resource is from a dependency mod
// and if so, checks if the dependency mod has a search path, search path prefix or database configured
// if so, it sets these values in viper
// NOTE: subsequent runs for  different resources may need the original values so save them in viper
//func UpdateTargetConnectionParams(target modconfig.ModTreeItem, workspaceMod *modconfig.Mod) error {
//	// NOTE: if the target is in a dependency mod, check whether a search path has been specificed for ut
//	depName := target.GetMod().DependencyName
//
//	if depName == "" {
//		return nil
//	}
//	// look for this mod in teh workspace mod require
//
//	modRequirement := workspaceMod.Require.GetModDependency(depName)
//	if modRequirement == nil {
//		// not expected
//		return sperr.New("could not find mod requirement for %s", depName)
//	}
//	// if the mod requirement has a search path, prefix or database, set it in viper,
//	// overriding whatever value sth, use it
//	if modRequirement.SearchPath != nil {
//		viper.Set(constants.ArgSearchPath, modRequirement.SearchPath)
//	}
//	if modRequirement.SearchPathPrefix != nil {
//		viper.Set(constants.ArgSearchPathPrefix, modRequirement.SearchPathPrefix)
//	}
//	if modRequirement.Database != nil {
//		viper.Set(constants.ArgDatabase, modRequirement.Database)
//	}
//
//	return nil
//}
