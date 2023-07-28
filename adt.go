
//DSL
stageCopyOptionsOnError := b.QueryAdt("StageCopyOptionsOnError")
	.With("Continue", d.String("CONTINUE")))
	.With("SkipFile", d.String("SKIP_FILE")))
	.With("SkipFileNum", d.String("SKIP_FILE").Var[int]("a"))
	.With("SkipFileNumPercent", d.String("'SKIP_FILE_").Var[int]("a").String("%%'"))
	.With("AbortStatement", d.String("ABORT_STATEMENT")))
// WYNIK

type StageCopyOptionsOnError interface {
	stageCopyOptionsOnError()
	String() string
}

type StageCopyOptionsOnErrorContinue struct{}

func (StageCopyOptionsOnErrorContinue) stageCopyOptionsOnError() {}
func (StageCopyOptionsOnErrorContinue) String() string {
	return "CONTINUE"
}

type StageCopyOptionsOnErrorSkipFile struct{}

func (StageCopyOptionsOnErrorSkipFile) stageCopyOptionsOnError() {}

func (StageCopyOptionsOnErrorSkipFile) String() string {
	return "SKIP_FILE"
}

type StageCopyOptionsOnErrorSkipFileNum struct {
	a int 
}

func (StageCopyOptionsOnErrorSkipFileNum) stageCopyOptionsOnError() {}

func (stageCopyOptionsOnErrorSkipFileNum StageCopyOptionsOnErrorSkipFileNum) String() string {
	return fmt.Sprintf("SKIP_FILE_%v", stageCopyOptionsOnErrorSkipFileNum.a)
}

type StageCopyOptionsOnErrorSkipFileNumPercent struct {
	a int
}

func (StageCopyOptionsOnErrorSkipFileNumPercentage) stageCopyOptionsOnError() {}

func (stageCopyOptionsOnErrorSkipFileNumPercentage StageCopyOptionsOnErrorSkipFileNumPercentage) String() string {
	return fmt.Sprintf("'SKIP_FILE_%d%%'", stageCopyOptionsOnErrorSkipFileNumPercentage.a)
}

type StageCopyOptionsOnErrorAbortStatement struct{}

func (StageCopyOptionsOnErrorAbortStatement) stageCopyOptionsOnError() {}

func (StageCopyOptionsOnErrorAbortStatement) String() string {
	return "ABORT_STATEMENT"
}

