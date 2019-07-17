/*
 * Copyright © 2019 Hedzr Yeh.
 */

package ags

import (
	"github.com/hedzr/awesome-tool/ags/gql"
	"github.com/russross/blackfriday"
	"github.com/sirupsen/logrus"
	"io"
)

func newMarkdownRenderer() blackfriday.Renderer {
	return &mdRenderer{}
}

type mdRenderer struct {
	itNode      *blackfriday.Node
	sectionNode *blackfriday.Node
	listNode    *blackfriday.Node
	liNode      *blackfriday.Node
	listItem    *gql.ListItem
	section     *gql.Section
	sections    []*gql.Section
}

func (s *mdRenderer) RenderNode(w io.Writer, node *blackfriday.Node, entering bool) blackfriday.WalkStatus {
	switch node.Type {
	case blackfriday.Text:
		if s.listItem != nil {
			if len(s.listItem.Url) > 0 {
				s.listItem.Desc += string(node.Literal)
			} else {
				s.listItem.Title += string(node.Literal)
			}
			// logrus.Infof("    %v TEXT: li title = %v", entering, string(node.Literal))
		} else if s.section != nil {
			if len(s.section.Header) == 0 && len(node.Literal) != 0 {
				s.section.Header = string(node.Literal)
				// logrus.Infof("=== SECTION HEADER %v: %v", s.section.level, string(node.Literal))
			}
		} else {
			logrus.Tracef("    %v TEXT: %v | %v", entering, node, string(node.Literal))
		}
	case blackfriday.Link:
		// mark it but don't link it if it is not a safe link: no smartypants
		// dest := node.LinkData.Destination
		// if entering {
		// logrus.Infof("%v link: %v -> %v", entering, string(node.LinkData.Title), string(node.LinkData.Destination))
		// }
		if entering {
			s.itNode = node
			if s.section != nil && s.listItem != nil {

			}
		} else {
			if s.section != nil && s.listItem != nil {
				s.listItem.Url = string(node.LinkData.Destination)
			}
		}
	case blackfriday.Document:
		break
	case blackfriday.Paragraph:
	case blackfriday.Heading:
		// logrus.Infof("%v heading: %v", entering, string(node.Title))
		if entering {
			if s.section != nil {
				s.sections = append(s.sections, s.section)
				s.section = nil
			}

			s.itNode = node
			s.liNode = node
			s.section = new(gql.Section)
			s.section.Level = node.Level
			s.listItem = nil
		} else {

		}
	case blackfriday.List:
		// logrus.Infof("%v list: %v", entering, node.ListData)
		if entering {
			s.itNode = node
			s.listNode = node
		} else {
			if s.section != nil {
				// if len(s.section.list) > 0 {
				s.sections = append(s.sections, s.section)
				// }
				s.section = nil
			}
		}
	case blackfriday.Item:
		if entering {
			// logrus.Infof("%v li: %v", entering, node.ListData)
			s.itNode = node
			s.liNode = node
			s.listItem = new(gql.ListItem)
		} else {
			if s.section != nil && s.listItem != nil {
				// logrus.Infof("    LI GOT: %v, %v", s.listItem.title, s.listItem.url)
				s.section.List = append(s.section.List, s.listItem)
				s.listItem = nil
			}
		}
	default:
		logrus.Tracef("    %v %v: %v", entering, node.Type, node)
	}
	return blackfriday.GoToNext
}

func (*mdRenderer) RenderHeader(w io.Writer, ast *blackfriday.Node) {
	logrus.Tracef("header: %v", ast)
}

func (*mdRenderer) RenderFooter(w io.Writer, ast *blackfriday.Node) {
	logrus.Tracef("footer: %v", ast)
}
