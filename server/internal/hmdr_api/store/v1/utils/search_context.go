package utils

import (
	"cmp"
	metav1 "gin-vue-admin/internal/pkg/meta/v1"
	nginx_context "github.com/ClessLi/bifrost/pkg/resolv/V3/nginx/configuration/context"
	"github.com/ClessLi/bifrost/pkg/resolv/V3/nginx/configuration/context_type"
	"slices"
)

func parseContextPos(pos nginx_context.Pos) (metav1.ConfigContextPos, error) {
	father, childIdx := pos.Position()
	if father.Error() != nil {
		return metav1.ConfigContextPos{}, father.Error()
	}
	if father.Type() == context_type.TypeMain || father.Type() == context_type.TypeConfig {
		return metav1.ConfigContextPos{
			Config:         father.Value(),
			ContextPosPath: []int{childIdx},
		}, nil
	}
	posmeta, err := parseContextPos(father.Father().ChildrenPosSet().QueryOne(nginx_context.NewKeyWords(father.Type()).
		SetCascaded(false).
		SetStringMatchingValue(father.Value()).
		SetSkipQueryFilter(func(targetCtx nginx_context.Context) bool { return targetCtx != father })))
	if err != nil {
		return metav1.ConfigContextPos{}, err
	}
	posmeta.ContextPosPath = append(posmeta.ContextPosPath, childIdx)
	return posmeta, nil
}

func SearchContextPoses(set nginx_context.PosSet, isonlycurrent bool, keywords string, isregexp bool) ([]metav1.ConfigContextPos, error) {
	result := make([]metav1.ConfigContextPos, 0)
	for _, contextType := range []context_type.ContextType{
		context_type.TypeEvents,
		context_type.TypeGeo,
		context_type.TypeHttp,
		context_type.TypeIf,
		context_type.TypeInclude,
		context_type.TypeDirective,
		context_type.TypeLimitExcept,
		context_type.TypeLocation,
		context_type.TypeMap,
		context_type.TypeServer,
		context_type.TypeStream,
		context_type.TypeTypes,
		context_type.TypeUpstream,
		context_type.TypeComment,
		context_type.TypeInlineComment,
	} {
		kw := nginx_context.NewKeyWords(contextType)
		if isregexp {
			kw.SetRegexpMatchingValue(keywords)
		} else {
			kw.SetStringMatchingValue(keywords)
		}

		if isonlycurrent {
			kw.SetSkipQueryFilter(
				func(targetCtx nginx_context.Context) bool { return targetCtx.Type() == context_type.TypeConfig },
			)
		}

		err := set.QueryAll(kw).
			Map(
				func(pos nginx_context.Pos) (nginx_context.Pos, error) {
					ctxPos, err := parseContextPos(pos)
					if err != nil {
						return pos, err
					}
					result = append(result, ctxPos)
					return pos, nil
				},
			).
			Error()
		if err != nil {
			return result, err
		}
	}
	slices.SortFunc(result, func(a, b metav1.ConfigContextPos) int {
		if n := cmp.Compare(a.Config, b.Config); n != 0 {
			return n
		}
		aplen := len(a.ContextPosPath)
		bplen := len(b.ContextPosPath)
		minlen := aplen
		if bplen < aplen {
			minlen = bplen
		}
		for i := 0; i < minlen; i++ {
			if n := cmp.Compare(a.ContextPosPath[i], b.ContextPosPath[i]); n != 0 {
				return n
			}
		}
		return cmp.Compare(aplen, bplen)
	})
	return result, nil
}
