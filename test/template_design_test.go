package test

import (
	"testing"
	"fmt"
)

func TestTemplateDesign(t *testing.T)  {
	netDoc := NewNetDoc()
	lcDoc := NewLocalDoc()

	netDoc.DoOperate()
	lcDoc.DoOperate()
}


type DocSuper struct {
	GetContent func() string
}
func (d DocSuper) DoOperate() {
	fmt.Println("对这个文档做了一些处理,文档是:", d.GetContent())
}




type LocalDoc struct {
	DocSuper
}
func NewLocalDoc() *LocalDoc {
	c := new(LocalDoc)
	c.DocSuper.GetContent = c.GetContent
	return c
}

func (e *LocalDoc) GetContent() string {
	return "this is a LocalDoc."
}



type NetDoc struct {
	DocSuper
}

func NewNetDoc() *NetDoc {
	c := new(NetDoc)
	c.DocSuper.GetContent = c.GetContent
	return c
}

func (c *NetDoc) GetContent() string {
	return "this is net doc."
}
