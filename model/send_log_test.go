package model

import "testing"

func TestSendLog_TableName(t *testing.T) {
	l := SendLog{}
	if l.TableName() != "email_send_logs" {
		t.Errorf("table name = %s", l.TableName())
	}
}

func TestSendLog_MarkSent(t *testing.T) {
	l := &SendLog{Status: SendStatusPending}
	l.MarkSent()
	if l.Status != SendStatusSent {
		t.Errorf("status = %s, want sent", l.Status)
	}
	if l.SentAt == nil {
		t.Error("sent_at should be set")
	}
}

func TestSendLog_MarkFailed(t *testing.T) {
	l := &SendLog{Status: SendStatusPending}
	l.MarkFailed("timeout")
	if l.Status != SendStatusFailed {
		t.Errorf("status = %s, want failed", l.Status)
	}
	if l.ErrorMessage != "timeout" {
		t.Errorf("error_message = %s", l.ErrorMessage)
	}
}
