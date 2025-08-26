package scheduler

import (
	"encoding/json"
	"fmt"

	"github.com/SOMTHING-ITPL/ITPL-server/internal/api"
)

type GPTResponse struct {
	Genre   int    `json:"genre"`
	Keyword string `json:"keyword"`
}

func getSystemPrompt() string {
	prompt := fmt.Sprintf(`You are a data analyst for a performance recommendation system.
Classify the genre of a performance and extract 20 core keywords in Korean.
Use this genre list (return integer only):
0: KPOP, 1: Rock/Metal, 2: Ballad, 3: Rap/Hip-hop, 4: Folk/Trot,
5: Fan Meeting, 6: Indie, 7: Jazz/Soul, 8: International Artist (Visit Korea),
9: R&B, 10: EDM, 11: Dinner Show, 12: Others
Output strictly in JSON:
{
"genre": int,
"keyword": "keyword1 keyword2 ... keyword20"
}
`)
	return prompt
}

func makeUserPrompt(name string, cast string) string {
	prompt := fmt.Sprintf(
		"here is performance Information\n",
		"Performance title: %s\nCast: %s",
		name,
		cast,
	)
	return prompt
}

func PreProcessPerformance(name string, cast string) (*GPTResponse, error) {
	content, err := api.SendPromptToModel(makeUserPrompt(name, cast), getSystemPrompt())
	if err != nil {
		return &GPTResponse{}, fmt.Errorf("Fail to send prompt to gpt : %s", err)
	}

	var gptResp GPTResponse
	if err := json.Unmarshal([]byte(content), &gptResp); err != nil {
		return &GPTResponse{}, fmt.Errorf("invalid JSON response: %w", err)
	}
	return &gptResp, nil
}
