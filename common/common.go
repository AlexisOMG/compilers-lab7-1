package common

const (
	eps   = "eps"
	term  = "term"
	nterm = "nterm"
)

type Expr struct {
	Value string
	Kind  string
}

var (
	Epsilon = Expr{
		Value: eps,
		Kind:  eps,
	}
)

type Rules map[Expr][][]Expr

func unionWithEps(a, b map[Expr]struct{}) map[Expr]struct{} {
	res := make(map[Expr]struct{}, len(a)-1+len(b))
	for e := range a {
		if e != Epsilon {
			res[e] = struct{}{}
		}
	}

	for e := range b {
		res[e] = struct{}{}
	}

	return res
}

func F(seq []Expr, first map[Expr]map[Expr]struct{}) map[Expr]struct{} {
	if len(seq) == 0 {
		return map[Expr]struct{}{
			Epsilon: {},
		}
	}

	if seq[0].Kind == term {
		return map[Expr]struct{}{
			seq[0]: {},
		}
	}

	if _, ok := first[seq[0]][Epsilon]; !ok {
		// copy?
		return first[seq[0]]
	}

	return unionWithEps(first[seq[0]], F(seq[1:], first))
}

func First(seq []Expr, rls Rules) map[Expr]map[Expr]struct{} {
	res := make(map[Expr]map[Expr]struct{}, len(rls))

	for l := range rls {
		res[l] = make(map[Expr]struct{})
	}

	changed := true

	for changed {
		changed = false

		for l := range res {
			for _, exprs := range rls[l] {
				f := F(exprs, res)
				for e := range f {
					if _, ok := res[l][e]; !ok {
						res[l][e] = struct{}{}
						changed = true
					}
				}
			}
		}
	}

	return res
}
