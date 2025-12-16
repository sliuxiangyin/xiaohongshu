package utils

import (
	"testing"
)

func TestGetPath(t *testing.T) {
	// 测试 dev 环境
	devPath, err := GetPath("dev")
	if err != nil {
		t.Errorf("GetPath(\"dev\") returned error: %v", err)
	}
	if devPath == "" {
		t.Error("GetPath(\"dev\") returned empty string")
	}

	// 测试 prod 环境
	prodPath, err := GetPath("prod")
	if err != nil {
		t.Errorf("GetPath(\"prod\") returned error: %v", err)
	}
	if prodPath == "" {
		t.Error("GetPath(\"prod\") returned empty string")
	}

	// 测试无效环境参数
	_, err = GetPath("invalid")
	if err == nil {
		t.Error("GetPath(\"invalid\") should return error but didn't")
	}
}
