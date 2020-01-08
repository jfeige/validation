package validator

import (
	"sync"
)

var wg sync.WaitGroup

type ValidateInterface interface {
	Check()
	IsFail() bool
	AllErrors() []error
	FirstError() error
	LastError() error
	SafeData() map[string]interface{}
}

type Validate struct {
	data     map[string]interface{}
	safeData map[string]interface{}
	rules    []map[string]interface{}
	errCnt   int
	errors   []error
}

func NewValidation(data map[string]interface{}, rules []map[string]interface{}) ValidateInterface {
	validate := &Validate{
		data:  data,
		rules: rules,
	}
	return validate
}

func (this *Validate) Check() {
	//errChan := make(chan string,100)		//最大支持100条规则
	this.errors = make([]error, 0)
	for _, rule := range this.rules {
		wg.Add(1)
		this.handleRule(rule)
	}
	wg.Wait()
}

func (this *Validate) IsFail() bool {
	if len(this.errors) > 0 {
		return true
	}
	return false
}

func (this *Validate) AllErrors() []error {
	return this.errors
}

func (this *Validate) FirstError() error {
	cnt := len(this.errors)
	if cnt > 0 {
		return this.errors[0]
	}
	return nil
}

func (this *Validate) LastError() error {
	cnt := len(this.errors)
	if cnt > 0 {
		return this.errors[cnt-1]
	}
	return nil
}

func (this *Validate) SafeData() map[string]interface{} {
	return this.safeData
}
