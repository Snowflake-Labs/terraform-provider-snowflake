package architest

import (
	"fmt"
	"regexp"
)

func FileNameFilterProvider(text string) FileFilter {
	regex := regexp.MustCompile(fmt.Sprintf(`^.*%s.*$`, text))
	return func(f *File) bool {
		return regex.Match([]byte(f.fileName))
	}
}

func FileNameRegexFilterProvider(regex *regexp.Regexp) FileFilter {
	return func(f *File) bool {
		return regex.Match([]byte(f.fileName))
	}
}

func FileNameFilterWithExclusionsProvider(regex *regexp.Regexp, exclusionRegex ...*regexp.Regexp) FileFilter {
	return func(f *File) bool {
		matches := regex.MatchString(f.fileName)
		for _, e := range exclusionRegex {
			matches = matches && !e.MatchString(f.fileName)
		}
		return matches
	}
}

func PackageFilterProvider(packageName string) FileFilter {
	return func(f *File) bool {
		return f.packageName == packageName
	}
}
