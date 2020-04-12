package opts

type Opts struct {
	Quality     string `short:"q" long:"quality" choice:"high" choice:"medium" choice:"low" default:"high"`
	Manual      bool   `short:"m" long:"manual"`
	AudioFormat string `short:"a" long:"audio-format" default:"flac"`
	Codec       string `short:"c" long:"codec" default:"flac"`
}
