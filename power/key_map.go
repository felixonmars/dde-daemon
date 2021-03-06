package power

func queryPowerAction(action string) int32 {
	switch action {
	case "blank":
		return 0
	case "suspend":
		return 1
	case "shutdown":
		return 2
	case "hibernate":
		return 3
	case "interactive":
		return 4
	case "nothing":
		return 5
	case "logout":
		return 6
	}
	return -1
}

func queryPowerPlan(plan string) int32 {
	switch plan {
	case "custom":
		return 0
	case "powersaver":
		return 1
	case "balanced":
		return 2
	case "hightperformance":
		return 3
	}
	return -1
}
