package request

import "testing"

func TestAgentReviewCallbackReqValidate(t *testing.T) {
	tests := []struct {
		name    string
		req     AgentReviewCallbackReq
		wantErr bool
	}{
		{
			name: "processing without final action",
			req: AgentReviewCallbackReq{
				TaskStatus: "processing",
			},
			wantErr: false,
		},
		{
			name: "processing with final action should fail",
			req: AgentReviewCallbackReq{
				TaskStatus:  "processing",
				FinalAction: "approve",
			},
			wantErr: true,
		},
		{
			name: "completed approve",
			req: AgentReviewCallbackReq{
				TaskStatus:  "completed",
				FinalAction: "approve",
			},
			wantErr: false,
		},
		{
			name: "completed missing final action should fail",
			req: AgentReviewCallbackReq{
				TaskStatus: "completed",
			},
			wantErr: true,
		},
		{
			name: "failed without final action",
			req: AgentReviewCallbackReq{
				TaskStatus: "failed",
			},
			wantErr: false,
		},
		{
			name: "failed with final action should fail",
			req: AgentReviewCallbackReq{
				TaskStatus:  "failed",
				FinalAction: "reject",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.req.Validate()
			if tt.wantErr && err == nil {
				t.Fatalf("expected validation error")
			}
			if !tt.wantErr && err != nil {
				t.Fatalf("unexpected validation error: %v", err)
			}
		})
	}
}

func TestAgentReviewCallbackReqNormalize(t *testing.T) {
	req := AgentReviewCallbackReq{
		TaskID:      " task-1 ",
		BizType:     " campus_auth ",
		BizID:       " 123 ",
		TraceID:     " trace-1 ",
		TaskStatus:  " COMPLETED ",
		FinalAction: " APPROVE ",
	}

	req.Normalize()

	if req.TaskID != "task-1" || req.BizType != "campus_auth" || req.BizID != "123" || req.TraceID != "trace-1" {
		t.Fatalf("normalize did not trim fields correctly: %+v", req)
	}
	if req.TaskStatus != "completed" || req.FinalAction != "approve" {
		t.Fatalf("normalize did not lowercase status fields: %+v", req)
	}
}
