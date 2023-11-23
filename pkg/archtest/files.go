package archtest

type FileFilter = func(*File) bool
type FileReceiver = func(*File)
type Files []File

func (files Files) Filter(filter FileFilter) Files {
	filteredFiles := make([]File, 0)
	for _, f := range files {
		if filter(&f) {
			filteredFiles = append(filteredFiles, f)
		}
	}
	return filteredFiles
}

func (files Files) All(receiver FileReceiver) {
	for _, file := range files {
		receiver(&file)
	}
}
