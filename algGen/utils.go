package algGen

import "hash/fnv"

func hash(s []string) string {
	var token string
	for _, t := range s {
		token += t
	}
	h := fnv.New32a()
	h.Write([]byte(token))
	return string(h.Sum32())
}
