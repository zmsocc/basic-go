package service

import (
	"gitee.com/zmsoc/gogogo/webook/internal/repository"
	"gitee.com/zmsoc/gogogo/webook/internal/service/sms/service"
	"reflect"
	"testing"
)

func TestNewCodeService(t *testing.T) {
	type args struct {
		repo   repository.CodeRepository
		smsSvc service.Service
	}
	tests := []struct {
		name string
		args args
		want CodeService
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewCodeService(tt.args.repo, tt.args.smsSvc); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewCodeService() = %v, want %v", got, tt.want)
			}
		})
	}
}
