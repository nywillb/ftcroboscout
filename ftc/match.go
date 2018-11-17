package ftc

import "time"

// Match holds the information regarding a match.
type Match struct {
	Key             string             `json:"match_key"`
	EventKey        string             `json:"event_key"`
	TournamentLevel int                `json:"tournament_level"`
	ScheduledTime   time.Time          `json:"scheduled_time"`
	MatchName       string             `json:"match_name"`
	PlayNumber      int                `json:"play_number"`
	FieldNumber     int                `json:"field_number"`
	PrestartTime    time.Time          `json:"prestart_time"`
	RedScore        int                `json:"red_score"`
	BlueScore       int                `json:"blue_score"`
	RedPenalty      int                `json:"red_penalty"`
	BluePenalty     int                `json:"blue_penalty"`
	RedAutoScore    int                `json:"red_auto_score"`
	BlueAutoScore   int                `json:"blue_auto_score"`
	RedTeleOpScore  int                `json:"red_tele_score"`
	BlueTeleOpScore int                `json:"blue_tele_score"`
	RedEndScore     int                `json:"red_end_score"`
	BlueEndScore    int                `json:"blue_end_score"`
	VideoURL        string             `json:"video_url"`
	Participants    []MatchParticipant `json:"participants"`
}

// MatchParticipant holds the information regarding a participant in a match.
type MatchParticipant struct {
	Key  int `json:"match_participant_key"`
	Team int `json:"team_key"`
}
