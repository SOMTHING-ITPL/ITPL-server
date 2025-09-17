package scheduler

import (
	"encoding/json"
	"fmt"

	"github.com/SOMTHING-ITPL/ITPL-server/internal/api"
)

type GPTResponse struct {
	Genre   int    `json:"genre"`
	Keyword string `json:"keyword"`
	Cast    string `json:"cast"`
}

func getSystemPrompt() string {
	prompt := fmt.Sprintf(`You are a data analyst for a performance recommendation system.
	Classify the genre of a performance and extract 20 core keywords in Korean.
	
	For the "keyword" and "cast" fields, do not rely only on the text.
	You must infer and generate values by analyzing the provided title, the given cast field, and the title together.
	
	- The "cast" field must always be filled: 
	  * If the given cast field is empty, infer from the title or text context.
	  * If no information at all can be inferred, return an empty string "".
	
	Use this genre list (return integer only):
	1: KPOP, 2: Rock/Metal, 3: Ballad, 4: Rap/Hip-hop, 5: Folk/Trot,
	6: Fan Meeting, 7: Indie, 8: Jazz/Soul,
	10: R&B, 11: EDM, 12: Dinner Show, 13: Others
	
	- The "keyword" field must contain exactly 20 core keywords joined with the literal string "|"
	- Example: "keyword": "키워드1|키워드2|키워드3|...|키워드20"
	
	Output strictly in JSON:
	{
	"genre": int,
	"keyword": "keyword1|keyword2| ... |keyword20",
	"cast" : "cast1, cast2, cast3" ...
	}
	`)
	return prompt
}

func makeUserPrompt(name string, cast string) string {
	prompt := fmt.Sprintf(
		"here is performance Information\nPerformance title: %s\nCast: %s",
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
	var recentErr error
	token := 3 // 최대 3번의 토큰 가지고 3번의 시도 -> 없으면 그냥 title만 담는걸로
	for token > 0 {
		if err := json.Unmarshal([]byte(content), &gptResp); err != nil {
			recentErr = fmt.Errorf("invalid JSON response: %w", err)
			token--
			continue
		} else {
			return &gptResp, nil
		}
	}
	return &gptResp, recentErr
}
