package channel

import (
	"testing"

	"github.com/voip-app/internal/models"
)

func TestParseTypeText(t *testing.T) {
	result := ParseType("text")
	if result != models.ChannelTypeText {
		t.Errorf("expected ChannelTypeText, got %s", result)
	}
}

func TestParseTypeTextCaseInsensitive(t *testing.T) {
	tests := []string{"Text", "TEXT", "text", "tExT"}
	for _, input := range tests {
		result := ParseType(input)
		if result != models.ChannelTypeText {
			t.Errorf("ParseType(%s): expected ChannelTypeText, got %s", input, result)
		}
	}
}

func TestParseTypeVoice(t *testing.T) {
	result := ParseType("voice")
	if result != models.ChannelTypeVoice {
		t.Errorf("expected ChannelTypeVoice, got %s", result)
	}
}

func TestParseTypeVoiceCaseInsensitive(t *testing.T) {
	tests := []string{"Voice", "VOICE", "voice", "VoIcE"}
	for _, input := range tests {
		result := ParseType(input)
		if result != models.ChannelTypeVoice {
			t.Errorf("ParseType(%s): expected ChannelTypeVoice, got %s", input, result)
		}
	}
}

func TestParseTypeInvalid(t *testing.T) {
	tests := []string{"", "invalid", "video", "audio", "123"}
	for _, input := range tests {
		result := ParseType(input)
		if result != "" {
			t.Errorf("ParseType(%s): expected empty string, got %s", input, result)
		}
	}
}

func TestErrorSentinels(t *testing.T) {
	if ErrNotFound == nil {
		t.Error("ErrNotFound should not be nil")
	}
	if ErrDefaultDelete == nil {
		t.Error("ErrDefaultDelete should not be nil")
	}
	if ErrDuplicateName == nil {
		t.Error("ErrDuplicateName should not be nil")
	}
	if ErrEmptyName == nil {
		t.Error("ErrEmptyName should not be nil")
	}
	if ErrInvalidType == nil {
		t.Error("ErrInvalidType should not be nil")
	}

	if ErrNotFound.Error() != "channel not found" {
		t.Errorf("unexpected ErrNotFound message: %s", ErrNotFound.Error())
	}
	if ErrDefaultDelete.Error() != "cannot delete default channel" {
		t.Errorf("unexpected ErrDefaultDelete message: %s", ErrDefaultDelete.Error())
	}
	if ErrDuplicateName.Error() != "channel name already exists" {
		t.Errorf("unexpected ErrDuplicateName message: %s", ErrDuplicateName.Error())
	}
	if ErrEmptyName.Error() != "channel name cannot be empty" {
		t.Errorf("unexpected ErrEmptyName message: %s", ErrEmptyName.Error())
	}
	if ErrInvalidType.Error() != "invalid channel type" {
		t.Errorf("unexpected ErrInvalidType message: %s", ErrInvalidType.Error())
	}
}

func TestChannelInfo(t *testing.T) {
	info := ChannelInfo{
		ID:        "ch1",
		Name:      "#general",
		Type:      "text",
		IsDefault: true,
	}

	if info.ID != "ch1" {
		t.Errorf("expected ID ch1, got %s", info.ID)
	}
	if info.Name != "#general" {
		t.Errorf("expected name #general, got %s", info.Name)
	}
	if info.Type != "text" {
		t.Errorf("expected type text, got %s", info.Type)
	}
	if !info.IsDefault {
		t.Error("expected is_default true")
	}
}
