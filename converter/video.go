// Copyright 2013 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package converter

import (
	"fmt"
	"github.com/andreaskoch/allmark/path"
	"github.com/andreaskoch/allmark/repository"
	"github.com/andreaskoch/allmark/util"
	"regexp"
	"strings"
)

var (
	// video: [*description text*](*a link to a youtube video or to a video file*)
	videoPattern = regexp.MustCompile(`video: \[([^\]]+)\]\(([^)]+)\)`)

	// youtube video link pattern
	youTubeVideoPattern = regexp.MustCompile(`http[s]?://www\.youtube\.com/watch\?v=([^&]+)`)
)

func newVideoRenderer(markdown string, fileIndex *repository.FileIndex, pathProvider *path.Provider) func(text string) string {
	return func(text string) string {
		return renderVideo(markdown, fileIndex, pathProvider)
	}
}

func renderVideo(markdown string, fileIndex *repository.FileIndex, pathProvider *path.Provider) string {

	for {

		found, matches := util.IsMatch(markdown, videoPattern)
		if !found || (found && len(matches) != 3) {
			break
		}

		// parameters
		originalText := strings.TrimSpace(matches[0])
		title := strings.TrimSpace(matches[1])
		path := strings.TrimSpace(matches[2])

		// get a renderer
		renderer := getVideoRenderer(title, path, fileIndex, pathProvider)

		// execute the renderer
		renderedCode := renderer()

		// replace markdown with link list
		markdown = strings.Replace(markdown, originalText, renderedCode, 1)

	}

	return markdown
}

func getVideoRenderer(title, path string, fileIndex *repository.FileIndex, pathProvider *path.Provider) func() string {

	// youtube
	if isYouTube, youTubeVideoId := isYouTubeLink(path); isYouTube {
		return func() string {
			return renderYouTubeVideo(title, youTubeVideoId)
		}
	}

	// html5 video file
	if isVideoFile := isVideoFileLink(path); isVideoFile {
		return func() string {
			return renderVideoFileLink(title, path)
		}
	}

	// return the fallback handler
	return func() string {
		return fmt.Sprintf(`<a href="%s" target="_blank" title="%s">%s</a>`, path, title, title)
	}
}

func isYouTubeLink(link string) (isYouTubeLink bool, videoId string) {
	if found, matches := util.IsMatch(link, youTubeVideoPattern); found && len(matches) == 2 {
		return true, matches[1]
	}

	return false, ""
}

func renderYouTubeVideo(title, videoId string) string {
	return fmt.Sprintf(`<section class="video video-youtube">
		<h1>YouTube Video: %s</h1>
		<p>
			<a href="http://www.youtube.com/watch?%s" target="_blank" title="%s">http://www.youtube.com/watch?%s</a>
		</p>
		<iframe width="560" height="315" src="http://www.youtube.com/embed/%s" frameborder="0" allowfullscreen></iframe>
	</section>`, title, videoId, title, videoId, videoId)
}

func renderVideoFileLink(title, link string) string {
	return fmt.Sprintf(`<section class="video video-file">
		<h1>Video: %s</h1>
		<p>
			<a href="%s" target="_blank" title="%s">%s</a>
		</p>
		<video width="560" height="315" controls>
			<source src="%s" type="video/mp4">
		</video>
	</section>`, title, link, title, link, link)
}

func isVideoFileLink(link string) bool {
	normalizedLink := strings.ToLower(link)
	fileExtension := normalizedLink[strings.LastIndex(normalizedLink, "."):]

	switch fileExtension {
	case ".mp4", ".ogg", "webm":
		return true
	default:
		return false
	}

	panic("Unreachable")
}
