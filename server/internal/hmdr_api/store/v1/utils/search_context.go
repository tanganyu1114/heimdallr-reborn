package utils

import (
	"cmp"
	metav1 "github.com/tanganyu1114/heimdallr-reborn/server/internal/pkg/meta/v1"
	"slices"

	"github.com/ClessLi/bifrost/api/bifrost/v1"
	"github.com/ClessLi/bifrost/pkg/resolv/V3/nginx/configuration"
	nginx_context "github.com/ClessLi/bifrost/pkg/resolv/V3/nginx/configuration/context"
	"github.com/ClessLi/bifrost/pkg/resolv/V3/nginx/configuration/context/local"
	"github.com/ClessLi/bifrost/pkg/resolv/V3/nginx/configuration/context_type"
	"github.com/marmotedu/errors"
)

func SearchContextPoses(set nginx_context.PosSet, isonlycurrent bool, keywords string, isregexp bool) ([]metav1.ConfigContextPos, error) {
	result := make([]metav1.ConfigContextPos, 0)
	kw := nginx_context.NewKeyWords(func(targetCtx nginx_context.Context) bool {
		return true
	}).SetSkipQueryFilter(func(targetCtx nginx_context.Context) bool {
		return isonlycurrent && targetCtx.Type() == context_type.TypeConfig
	}).AppendMatchingFilter(func(targetCtx nginx_context.Context) bool {
		return targetCtx.Type() != context_type.TypeConfig
	})
	if isregexp {
		kw.SetRegexpMatchingValue(keywords)
	} else {
		kw.SetStringMatchingValue(keywords)
	}
	parsedPosMap := make(map[nginx_context.Context]bool)
	err := set.QueryAll(kw).
		Filter(
			func(pos nginx_context.Pos) bool {
				if parsedPosMap[pos.Target()] {
					return false
				}
				parsedPosMap[pos.Target()] = true
				return true
			},
		).
		Map(
			func(pos nginx_context.Pos) (nginx_context.Pos, error) {
				ctxPos, err := local.PosBasedOnConfig(pos.Target())
				if err != nil {
					return pos, err
				}
				result = append(result, ConvertToConfigContextPos(ctxPos))
				return pos, nil
			},
		).
		Error()
	if err != nil {
		return nil, err
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

func ConvertToBifrostContextPos(in metav1.ConfigContextPos) (out *v1.ContextPos) {
	out = new(v1.ContextPos)
	out.ConfigPath = in.Config
	for _, idx := range in.ContextPosPath {
		out.PosIndex = append(out.PosIndex, int32(idx))
	}
	return
}

func ConvertToConfigContextPos(in *v1.ContextPos) (out metav1.ConfigContextPos) {
	if in == nil {
		return
	}
	out.Config = in.ConfigPath
	for _, idx := range in.PosIndex {
		out.ContextPosPath = append(out.ContextPosPath, int(idx))
	}
	return
}

func ParseContext(nginxconfig configuration.NginxConfig, configPath string, ctxPosPath []int) (nginx_context.Context, error) {
	posConfigPath, err := nginx_context.NewRelConfigPath(nginxconfig.Main().MainConfig().BaseDir(), configPath)
	if err != nil {
		return nginx_context.NullContext(), errors.Errorf("failed to parse the nginx config path(%s), cased by: %s", configPath, err)
	}
	target := nginx_context.NullContext()
	target, err = nginxconfig.Main().GetConfig(posConfigPath.FullPath())
	if err != nil {
		return nginx_context.NullContext(), err
	}
	for _, idx := range ctxPosPath {
		target = target.Child(idx)
	}
	return target, target.Error()
}

func ParseContextPos(nginxconfig configuration.NginxConfig, pos metav1.ConfigContextPos) (nginx_context.Pos, error) {
	if len(pos.ContextPosPath) == 0 {
		return nginx_context.NullPos(), errors.New("the nginx config context pos path is null")
	}

	ctx, err := ParseContext(nginxconfig, pos.Config, pos.ContextPosPath)
	if err != nil {
		return nginx_context.NullPos(), err
	}
	p := nginx_context.GetPos(ctx)
	return p, p.Target().Error()
}

func ParseContextPosModifyTO(nginxconfig configuration.NginxConfig, pos metav1.ConfigContextPos) (nginx_context.Pos, error) {
	if len(pos.ContextPosPath) == 0 {
		return nginx_context.NullPos(), errors.New("the nginx config context pos path is null")
	}
	var fatherPosPath []int
	if len(pos.ContextPosPath) > 1 {
		fatherPosPath = pos.ContextPosPath[:len(pos.ContextPosPath)-1]
	}
	fatherCtx, err := ParseContext(nginxconfig, pos.Config, fatherPosPath)
	if err != nil {
		return nginx_context.NullPos(), err
	}

	var targetIdx = pos.ContextPosPath[len(pos.ContextPosPath)-1]
	if targetIdx > fatherCtx.Len() {
		targetIdx = fatherCtx.Len()
	}

	return nginx_context.SetPos(fatherCtx, targetIdx), nil
}
