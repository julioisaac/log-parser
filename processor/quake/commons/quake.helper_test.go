package commons

import "testing"

func Test_getTypeOfProcess(t *testing.T) {
	type args struct {
		line string
	}
	tests := []struct {
		name string
		args args
		want TypeOfProcess
	}{
		{
			name: "should return InitGame process type",
			args: args{line: ` 16:53 InitGame: \capturelimit\8\g_maxGameClients\0\timelimit\15\fraglimit\20\dmflags\0\bot_minplayers\0\sv_allowDownload\0\sv_maxclients\16\sv_privateClients\2\g_gametype\4\sv_hostname\Code Miner Server\sv_minRate\0\sv_maxRate\10000\sv_minPing\0\sv_maxPing\0\sv_floodProtect\1\version\ioq3 1.36 linux-x86_64 Apr 12 2009\protocol\68\mapname\Q3TOURNEY6_CTF\gamename\baseq3\g_needpass\0`},
			want: InitGame,
		},
		{
			name: "should return MatchKill process type",
			args: args{line: ` 17:16 MatchKill: 1022 4 22: <world> killed Zeh by MOD_TRIGGER_HURT`},
			want: Kill,
		},
		{
			name: "should return ShutdownGame process type",
			args: args{line: `981:27 ShutdownGame:`},
			want: ShutdownGame,
		},
		{
			name: "should return Player process type",
			args: args{line: `  0:04 ClientUserinfoChanged: 6 n\Zeh\t\0\model\sarge/default\hmodel\sarge/default\g_redteam\\g_blueteam\\c1\1\c2\5\hc\100\w\0\l\0\tt\0\tl\0`},
			want: Player,
		},
		{
			name: "should return Ignore process type",
			args: args{line: `ClientConnect: 2`},
			want: Ignore,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetTypeOfProcess(tt.args.line); got != tt.want {
				t.Errorf("getTypeOfProcess() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestContainsString(t *testing.T) {
	type args struct {
		arr   []string
		input string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "should return false if input is not in the array",
			args: args{
				arr: []string{
					"Fasano Again",
					"Isgalamido",
				},
				input: "Oootsimo",
			},
			want: false,
		},
		{
			name: "should return true if input is in the array",
			args: args{
				arr: []string{
					"Zeh",
				},
				input: "Zeh",
			},
			want: true,
		},
		{
			name: "should return false if array is null",
			args: args{
				arr:   nil,
				input: "Zeh",
			},
			want: false,
		},
		{
			name: "should return false if array is empty",
			args: args{
				arr:   []string{},
				input: "Zeh",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ContainsString(tt.args.arr, tt.args.input); got != tt.want {
				t.Errorf("ContainsString() = %v, want %v", got, tt.want)
			}
		})
	}
}
