package ins

// 指令同一格式

type Instruction interface {
	New(originCmd string) Instruction
	Prepare(originCmd string, ctrls ...interface{}) (values []string)
	Exec(values []string) (stop bool)
}
