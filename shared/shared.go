package shared

type Diff struct {
	Key        string
	Change     string
	LeftValue  any
	RightValue any
	Child      map[string]Diff
}

func NewDiff(key, change string, leftValue, rightValue any, child map[string]Diff) Diff {
	return Diff{
		key, change, leftValue, rightValue, child,
	}
}
