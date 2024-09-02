package cmd

type FileConfig struct {
	Op          string   `yaml:"op"`
	SrcFilePath string   `yaml:"srcFilePath"`
	DstFilePath string   `yaml:"dstFilePath"`
	TargetOS    []string `yaml:"targetOS"`
	TargetArch  []string `yaml:"targetArch"`
}
