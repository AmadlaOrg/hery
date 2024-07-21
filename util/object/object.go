package object

// MergeMultilevel
func MergeMultilevel(dTo, dFrom map[string]interface{}, inplace bool) map[string]interface{} {
	if !inplace {
		dTo = DeepCopy(dTo)
	}
	for k, vFrom := range dFrom {
		if vTo, ok := dTo[k]; ok {
			if vToMap, ok := vTo.(map[string]interface{}); ok {
				if vFromMap, ok := vFrom.(map[string]interface{}); ok {
					dTo[k] = MergeMultilevel(vToMap, vFromMap, true)
					continue
				}
			}
		}
		dTo[k] = vFrom
	}
	return dTo
}

/*
func MergeMultilevel(dst, src map[string]interface{}, overwrite bool) map[string]interface{} {
	for k, v := range src {
		if vMap, ok := v.(map[string]interface{}); ok {
			if dstV, ok := dst[k].(map[string]interface{}); ok {
				dst[k] = MergeMultilevel(dstV, vMap, overwrite)
				continue
			}
		}
		if overwrite {
			dst[k] = v
		} else {
			if _, ok := dst[k]; !ok {
				dst[k] = v
			}
		}
	}
	return dst
}
*/

// DeepCopy
func DeepCopy(d map[string]interface{}) map[string]interface{} {
	copied := make(map[string]interface{})
	for k, v := range d {
		if vMap, ok := v.(map[string]interface{}); ok {
			copied[k] = DeepCopy(vMap)
		} else {
			copied[k] = v
		}
	}
	return copied
}
