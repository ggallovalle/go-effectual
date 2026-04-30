package luagen

func Classify(info *TypeInfo, cfg *GenConfig) {
	for i := range info.Methods {
		m := &info.Methods[i]
		if m.IsSkipped {
			continue
		}

		if cfg.IsSkipped(m.Name) {
			m.IsSkipped = true
			continue
		}

		if cfg.IsNilMapped(m.Name) {
			m.IsNilMap = true
		}

		if info.Module != "" {
			classifyModuleMethod(m)
		} else {
			m.IsGetter = isGetter(m, cfg)
		}
	}

	for i := range info.Fields {
		if cfg.IsFieldSkipped(info.Fields[i].Name) {
			info.Fields[i].IsSkipped = true
		}
	}
}

func isGetter(m *MethodInfo, cfg *GenConfig) bool {
	if m.IsSkipped {
		return false
	}
	if cfg.IsForceMethod(m.Name) || m.IsForceMethod {
		return false
	}
	if len(m.Params) != 0 {
		return false
	}
	if m.ReturnKind == ReturnVoid || m.ReturnKind == ReturnComplex {
		return false
	}
	return true
}

func isGoLuaFunction(m *MethodInfo) bool {
	return m.PtrType == "lua.State" && len(m.Params) == 1 && m.ReturnKind == ReturnPointer
}

func classifyModuleMethod(m *MethodInfo) {
	if m.Method {
		return
	}
	if isGoLuaFunction(m) {
		m.IsVerbatim = true
		return
	}
	if len(m.Params) == 0 && m.ReturnKind != ReturnVoid && m.ReturnKind != ReturnComplex {
		m.IsField = true
	}
}
