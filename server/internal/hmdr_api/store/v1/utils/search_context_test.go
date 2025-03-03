package utils

import (
	metav1 "gin-vue-admin/internal/pkg/meta/v1"
	nginx_context "github.com/ClessLi/bifrost/pkg/resolv/V3/nginx/configuration/context"
	"github.com/ClessLi/bifrost/pkg/resolv/V3/nginx/configuration/context/local"
	"github.com/ClessLi/bifrost/pkg/resolv/V3/nginx/configuration/context_type"
	"reflect"
	"testing"
)

func TestSearchContextPoses(t *testing.T) {
	testMain, err := local.NewMain("C:\\test\\nginx.conf")
	if err != nil {
		t.Fatal(err)
	}

	testMain.Insert(
		local.NewContext(context_type.TypeHttp, "").
			Insert(
				local.NewContext(context_type.TypeServer, "").
					Insert(
						local.NewContext(context_type.TypeInlineComment, "enabled server with enabled children configs"),
						0,
					).
					Insert(
						local.NewContext(context_type.TypeInclude, "conf.d/enabled.conf"),
						1,
					),
				0,
			).
			Insert(
				local.NewContext(context_type.TypeServer, "").
					Insert(
						local.NewContext(context_type.TypeInlineComment, "enabled server with disabled children configs"),
						0,
					).
					Insert(
						local.NewContext(context_type.TypeInclude, "conf.d/disabled.conf"),
						1,
					),
				1,
			).
			Insert(
				local.NewContext(context_type.TypeServer, "").
					Insert(
						local.NewContext(context_type.TypeInlineComment, "enabled server with disabled include context"),
						0,
					).
					Insert(
						local.NewContext(context_type.TypeInclude, "conf.d/enabled.conf").Disable(),
						1,
					),
				2,
			).
			Insert(
				local.NewContext(context_type.TypeServer, "").Disable().
					Insert(
						local.NewContext(context_type.TypeInlineComment, "disabled server with enabled children configs"),
						0,
					).
					Insert(
						local.NewContext(context_type.TypeInclude, "conf.d/enabled.conf"),
						1,
					),
				3,
			).
			Insert(
				local.NewContext(context_type.TypeServer, "").Disable().
					Insert(
						local.NewContext(context_type.TypeInlineComment, "disabled server with disabled children configs"),
						0,
					).
					Insert(
						local.NewContext(context_type.TypeInclude, "conf.d/disabled.conf"),
						1,
					),
				4,
			).
			Insert(
				local.NewContext(context_type.TypeServer, "").Disable().
					Insert(
						local.NewContext(context_type.TypeInlineComment, "disabled server with disabled include context"),
						0,
					).
					Insert(
						local.NewContext(context_type.TypeInclude, "conf.d/enabled.conf").Disable(),
						1,
					),
				5,
			),
		0,
	)
	err = testMain.AddConfig(
		local.NewContext(context_type.TypeConfig, "conf.d/enabled.conf").
			Insert(
				local.NewContext(context_type.TypeLocation, "~ /test").
					Insert(local.NewContext(context_type.TypeDirective, "return 200 'test'"), 0),
				0,
			).(*local.Config),
	)
	if err != nil {
		t.Fatal(err)
	}
	err = testMain.AddConfig(
		local.NewContext(context_type.TypeConfig, "conf.d/disabled.conf").Disable().
			Insert(
				local.NewContext(context_type.TypeComment, "disabled config"),
				0,
			).
			Insert(
				local.NewContext(context_type.TypeLocation, "~ /test").
					Insert(local.NewContext(context_type.TypeDirective, "return 404"), 0),
				1,
			).(*local.Config),
	)
	if err != nil {
		t.Fatal(err)
	}
	type args struct {
		set           nginx_context.PosSet
		isonlycurrent bool
		keywords      string
		isregexp      bool
	}
	tests := []struct {
		name    string
		args    args
		want    []metav1.ConfigContextPos
		wantErr bool
	}{
		{
			name: "string match rule searching, only in the current config",
			args: args{
				set:           testMain.ChildrenPosSet(),
				isonlycurrent: true,
				keywords:      "enabled.conf",
				isregexp:      false,
			},
			want: []metav1.ConfigContextPos{
				{"C:\\test\\nginx.conf", []int{0, 0, 1}},
				{"C:\\test\\nginx.conf", []int{0, 2, 1}},
				{"C:\\test\\nginx.conf", []int{0, 3, 1}},
				{"C:\\test\\nginx.conf", []int{0, 5, 1}},
			},
		},
		{
			name: "regexp match rule searching, only in the current config",
			args: args{
				set:           testMain.ChildrenPosSet(),
				isonlycurrent: true,
				keywords:      "disabled.*",
				isregexp:      true,
			},
			want: []metav1.ConfigContextPos{
				{"C:\\test\\nginx.conf", []int{0, 1, 0}},
				{"C:\\test\\nginx.conf", []int{0, 1, 1}},
				{"C:\\test\\nginx.conf", []int{0, 2, 0}},
				{"C:\\test\\nginx.conf", []int{0, 3, 0}},
				{"C:\\test\\nginx.conf", []int{0, 4, 0}},
				{"C:\\test\\nginx.conf", []int{0, 4, 1}},
				{"C:\\test\\nginx.conf", []int{0, 5, 0}},
			},
		},
		{
			name: "string match rule searching, not only in the current config",
			args: args{
				set:           testMain.ChildrenPosSet(),
				isonlycurrent: false,
				keywords:      "~ /test",
				isregexp:      false,
			},
			want: []metav1.ConfigContextPos{
				{"conf.d/disabled.conf", []int{1}},
				{"conf.d/enabled.conf", []int{0}},
			},
		},
		{
			name: "regexp match rule searching, not only in the current config",
			args: args{
				set:           testMain.ChildrenPosSet(),
				isonlycurrent: false,
				keywords:      ".*?test",
				isregexp:      true,
			},
			want: []metav1.ConfigContextPos{
				{"conf.d/disabled.conf", []int{1}},
				{"conf.d/enabled.conf", []int{0}},
				{"conf.d/enabled.conf", []int{0, 0}},
			},
		},
		{
			name: "context not found",
			args: args{
				set:           testMain.ChildrenPosSet(),
				isonlycurrent: false,
				keywords:      ".*not found.*",
				isregexp:      true,
			},
			want: []metav1.ConfigContextPos{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := SearchContextPoses(tt.args.set, tt.args.isonlycurrent, tt.args.keywords, tt.args.isregexp)
			if (err != nil) != tt.wantErr {
				t.Errorf("SearchContextPoses() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SearchContextPoses() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_parseContextPos(t *testing.T) {
	testMain, err := local.NewMain("C:\\test\\nginx.conf")
	if err != nil {
		t.Fatal(err)
	}

	testMain.Insert(
		local.NewContext(context_type.TypeHttp, "").
			Insert(
				local.NewContext(context_type.TypeServer, "").
					Insert(
						local.NewContext(context_type.TypeInlineComment, "enabled server with enabled children configs"),
						0,
					).
					Insert(
						local.NewContext(context_type.TypeInclude, "conf.d/enabled.conf"),
						1,
					),
				0,
			).
			Insert(
				local.NewContext(context_type.TypeServer, "").
					Insert(
						local.NewContext(context_type.TypeInlineComment, "enabled server with disabled children configs"),
						0,
					).
					Insert(
						local.NewContext(context_type.TypeInclude, "conf.d/disabled.conf"),
						1,
					),
				1,
			).
			Insert(
				local.NewContext(context_type.TypeServer, "").
					Insert(
						local.NewContext(context_type.TypeInlineComment, "enabled server with disabled include context"),
						0,
					).
					Insert(
						local.NewContext(context_type.TypeInclude, "conf.d/enabled.conf").Disable(),
						1,
					),
				2,
			).
			Insert(
				local.NewContext(context_type.TypeServer, "").Disable().
					Insert(
						local.NewContext(context_type.TypeInlineComment, "disabled server with enabled children configs"),
						0,
					).
					Insert(
						local.NewContext(context_type.TypeInclude, "conf.d/enabled.conf"),
						1,
					),
				3,
			).
			Insert(
				local.NewContext(context_type.TypeServer, "").Disable().
					Insert(
						local.NewContext(context_type.TypeInlineComment, "disabled server with disabled children configs"),
						0,
					).
					Insert(
						local.NewContext(context_type.TypeInclude, "conf.d/disabled.conf"),
						1,
					),
				4,
			).
			Insert(
				local.NewContext(context_type.TypeServer, "").Disable().
					Insert(
						local.NewContext(context_type.TypeInlineComment, "disabled server with disabled include context"),
						0,
					).
					Insert(
						local.NewContext(context_type.TypeInclude, "conf.d/enabled.conf").Disable(),
						1,
					),
				5,
			),
		0,
	)
	err = testMain.AddConfig(
		local.NewContext(context_type.TypeConfig, "conf.d/enabled.conf").
			Insert(
				local.NewContext(context_type.TypeLocation, "~ /test").
					Insert(local.NewContext(context_type.TypeDirective, "return 200 'test'"), 0),
				0,
			).(*local.Config),
	)
	if err != nil {
		t.Fatal(err)
	}
	err = testMain.AddConfig(
		local.NewContext(context_type.TypeConfig, "conf.d/disabled.conf").Disable().
			Insert(
				local.NewContext(context_type.TypeComment, "disabled config"),
				0,
			).
			Insert(
				local.NewContext(context_type.TypeLocation, "~ /test").
					Insert(local.NewContext(context_type.TypeDirective, "return 404"), 0),
				1,
			).(*local.Config),
	)
	if err != nil {
		t.Fatal(err)
	}
	type args struct {
		pos nginx_context.Pos
	}
	tests := []struct {
		name    string
		args    args
		want    metav1.ConfigContextPos
		wantErr bool
	}{
		{
			name: "normal test",
			args: args{
				pos: testMain.ChildrenPosSet().QueryOne(nginx_context.NewKeyWords(context_type.TypeInlineComment).SetStringMatchingValue("disabled server with disabled children configs")),
			},
			want: metav1.ConfigContextPos{
				Config:         "C:\\test\\nginx.conf",
				ContextPosPath: []int{0, 4, 0},
			},
		},
		{
			name:    "error test",
			args:    args{pos: nginx_context.NullPos()},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseContextPos(tt.args.pos)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseContextPos() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseContextPos() got = %v, want %v", got, tt.want)
			}
		})
	}
}
