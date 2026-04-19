package campus

import (
	"strings"
	"testing"

	campusModel "github.com/flipped-aurora/gin-vue-admin/server/model/campus"
	campusReq "github.com/flipped-aurora/gin-vue-admin/server/model/campus/request"
)

func TestBuildCampusAuthTaskPayload(t *testing.T) {
	payload := buildCampusAuthTaskPayload(campusReq.SubmitCampusAuthReq{
		UserID:    456,
		StudentID: "20230001",
		RealName:  "张三",
		College:   "计算机学院",
	})

	if payload.UserID != "456" {
		t.Fatalf("expected string user id, got %q", payload.UserID)
	}
	if payload.StudentID != "20230001" || payload.RealName != "张三" || payload.College != "计算机学院" {
		t.Fatalf("unexpected payload content: %+v", payload)
	}
}

func TestNewCampusAuthTaskIdentifiers(t *testing.T) {
	taskID, traceID := newCampusAuthTaskIdentifiers()

	if !strings.HasPrefix(taskID, "task_") {
		t.Fatalf("unexpected task id: %s", taskID)
	}
	if !strings.HasPrefix(traceID, "trace_") {
		t.Fatalf("unexpected trace id: %s", traceID)
	}
}

func TestNormalizeSubmitCampusAuthReq(t *testing.T) {
	req, err := normalizeSubmitCampusAuthReq(campusReq.SubmitCampusAuthReq{
		UserID:    1,
		StudentID: " 20230001 ",
		RealName:  " 张三 ",
		College:   " 计算机学院 ",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if req.StudentID != "20230001" || req.RealName != "张三" || req.College != "计算机学院" {
		t.Fatalf("normalize did not trim fields: %+v", req)
	}
}

func TestBuildAuthStatusText(t *testing.T) {
	cases := map[int]string{
		campusModel.CampusUserAuthStatusUnverified: "未认证",
		campusModel.CampusUserAuthStatusRejected:   "已拒绝",
		campusModel.CampusUserAuthStatusProcessing: "审核中",
		campusModel.CampusUserAuthStatusApproved:   "已认证",
	}

	for authStatus, expected := range cases {
		got := buildAuthStatusText(campusModel.CampusUser{AuthStatus: authStatus})
		if got != expected {
			t.Fatalf("auth status %d expected %q, got %q", authStatus, expected, got)
		}
	}
}
