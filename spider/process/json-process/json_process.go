package process

import (
	"YiSpider/spider/model"
	"YiSpider/spider/process/rule"
)

type TemplateProcess struct {

}

func NewTemplateProcess() *TemplateProcess{
	return &TemplateProcess{}
}

func (h *TemplateProcess)Process(bytes []byte,task *model.Task) (*model.Page,error){
	return rule.TemplateProcess(task,bytes)
}