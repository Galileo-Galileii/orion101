package selectors

func RemoveEmpty(selector map[string]string) map[string]string {
	for k, v := range selector {
		if v == "" {
			delete(selector, k)
		}
	}
	return selector
}
