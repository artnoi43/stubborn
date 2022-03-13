package cacher

func AppendRootServer(d string) string {
	if d[len(d)-1] != '.' {
		return d + "."
	}
	return d
}

func AppendRootServerKey(k Key) Key {
	if k.Dom[len(k.Dom)-1] != '.' {
		k.Dom += "."
	}
	return k
}
