package gewe

import "testing"

func TestSplitGroupContent(t *testing.T) {
	cases := []struct {
		in           string
		wantMember   string
		wantRest     string
		wantParsed   bool
	}{
		// 文本群消息：wxid + 昵称 + 群内名
		{"wxid_phyyedw9xap22（朝夕。）（昵称）:\n你好", "wxid_phyyedw9xap22", "你好", true},
		// 仅 wxid 无昵称
		{"wxid_phyyedw9xap22:\n今天天气不错", "wxid_phyyedw9xap22", "今天天气不错", true},
		// 多媒体群消息：wxid + 冒号 + XML
		{"wxid_phyyedw9xap22:\n<msg><img/></msg>", "wxid_phyyedw9xap22", "<msg><img/></msg>", true},
		// 私聊内容（不应被当成群成员前缀，但函数本身只负责群聊分支，无 wxid 前缀应原样返回 rest）
		{"hello world", "", "hello world", false},
		// 群 id 前缀（含 @chatroom，不应解析）
		{"39238473509@chatroom:\n内容", "", "39238473509@chatroom:\n内容", false},
		// 纯数字 wxid（当前模式不匹配，退回不解析——保持 content 原样）
		{"2345678901:\n内容", "", "2345678901:\n内容", false},
	}
	for _, c := range cases {
		gotMember, gotRest := splitGroupContent(c.in)
		if c.wantParsed {
			if gotMember != c.wantMember || gotRest != c.wantRest {
				t.Errorf("splitGroupContent(%q) = (%q,%q), want (%q,%q)", c.in, gotMember, gotRest, c.wantMember, c.wantRest)
			}
		} else {
			if gotRest != c.wantRest {
				t.Errorf("splitGroupContent(%q) rest = %q, want %q", c.in, gotRest, c.wantRest)
			}
		}
	}
}
