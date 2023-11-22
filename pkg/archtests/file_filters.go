package archtests

import "regexp"

type FileFilter = func(*File) bool

func filterFiles(files []File, filter FileFilter) []File {
	filteredFiles := make([]File, 0)
	for _, f := range files {
		if filter(&f) {
			filteredFiles = append(filteredFiles, f)
		}
	}
	return filteredFiles
}

func fileNameFilterProvider(regex *regexp.Regexp) FileFilter {
	return func(f *File) bool {
		return regex.Match([]byte(f.fileName))
	}
}

func fileNameFilterWithExclusionsProvider(regex *regexp.Regexp, exclusionRegex ...*regexp.Regexp) FileFilter {
	return func(f *File) bool {
		matches := regex.MatchString(f.fileName)
		for _, e := range exclusionRegex {
			matches = matches && !e.MatchString(f.fileName)
		}
		return matches
	}
}

func packageFilterProvider(packageName string) FileFilter {
	return func(f *File) bool {
		return f.packageName == packageName
	}
}
