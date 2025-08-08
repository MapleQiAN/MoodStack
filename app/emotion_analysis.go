package app

import (
	"bytes"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"sort"
	"strings"
	"time"

	"gorm.io/gorm"
)

// EmotionKeywords defines emotion keywords for programmatic analysis
var EmotionKeywords = map[string][]string{
	"joy": {
		"开心", "快乐", "高兴", "愉悦", "兴奋", "满足", "幸福", "欣喜", "喜悦", "舒心",
		"畅快", "欢乐", "惊喜", "满意", "放松", "轻松", "安心", "温暖", "甜蜜", "美好",
		"棒", "太好了", "成功", "胜利", "完美", "优秀", "精彩", "wonderful", "amazing",
	},
	"sadness": {
		"难过", "伤心", "悲伤", "痛苦", "沮丧", "失落", "郁闷", "忧郁", "抑郁", "孤独",
		"寂寞", "空虚", "无助", "绝望", "心酸", "心痛", "眼泪", "哭", "流泪", "想哭",
		"失望", "失落", "遗憾", "可惜", "不幸", "糟糕", "terrible", "awful", "sad",
	},
	"anger": {
		"愤怒", "生气", "恼火", "烦躁", "暴躁", "愤恨", "憎恨", "讨厌", "厌恶", "反感",
		"不爽", "火大", "气死", "烦死", "讨厌死", "恨", "愤慨", "激愤", "不满", "抱怨",
		"该死", "混蛋", "白痴", "蠢", "垃圾", "操", "妈的", "angry", "hate", "furious",
	},
	"fear": {
		"害怕", "恐惧", "担心", "紧张", "焦虑", "不安", "慌张", "恐慌", "惊慌", "胆怯",
		"畏惧", "忧虑", "忐忑", "心慌", "紧张兮兮", "提心吊胆", "惶恐", "惊恐", "战栗",
		"颤抖", "吓", "怕", "可怕", "恐怖", "吓死", "scared", "afraid", "terrified",
	},
	"love": {
		"爱", "喜欢", "爱情", "恋爱", "暗恋", "表白", "约会", "亲", "吻", "拥抱",
		"想念", "思念", "牵挂", "在乎", "关心", "温柔", "甜蜜", "浪漫", "心动", "迷恋",
		"深爱", "热爱", "钟爱", "疼爱", "宠爱", "爱意", "love", "romantic", "crush",
	},
	"surprise": {
		"惊讶", "震惊", "吃惊", "惊奇", "意外", "出乎意料", "想不到", "没想到", "突然",
		"忽然", "竟然", "居然", "原来", "天哪", "我的天", "不会吧", "真的吗", "哇",
		"amazing", "surprising", "unexpected", "wow", "omg",
	},
	"disgust": {
		"恶心", "讨厌", "厌恶", "反感", "嫌弃", "恶劣", "肮脏", "龌龊", "污秽", "臭",
		"难闻", "难看", "丑", "恶心死", "受不了", "无法忍受", "gross", "disgusting", "yuck",
	},
}

// SentimentKeywords defines positive and negative sentiment keywords
var SentimentKeywords = map[string][]string{
	"positive": {
		"好", "棒", "优秀", "完美", "成功", "胜利", "满意", "开心", "快乐", "幸福",
		"美好", "温暖", "甜蜜", "感谢", "赞", "喜欢", "爱", "支持", "鼓励", "positive",
		"great", "good", "excellent", "awesome", "wonderful", "amazing", "perfect",
	},
	"negative": {
		"坏", "糟", "失败", "错误", "问题", "困难", "麻烦", "痛苦", "难过", "伤心",
		"失望", "讨厌", "恨", "愤怒", "害怕", "担心", "焦虑", "压力", "negative",
		"bad", "terrible", "awful", "horrible", "wrong", "failed", "problem",
	},
}

// OllamaRequest represents a request to Ollama API
type OllamaRequest struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
	Stream bool   `json:"stream"`
}

// OllamaResponse represents a response from Ollama API
type OllamaResponse struct {
	Response string `json:"response"`
	Done     bool   `json:"done"`
}

// EmotionAnalysisResult represents the result of emotion analysis
type EmotionAnalysisResult struct {
	Joy             float64  `json:"joy"`
	Sadness         float64  `json:"sadness"`
	Anger           float64  `json:"anger"`
	Fear            float64  `json:"fear"`
	Love            float64  `json:"love"`
	Surprise        float64  `json:"surprise"`
	Disgust         float64  `json:"disgust"`
	DominantEmotion string   `json:"dominantEmotion"`
	Confidence      float64  `json:"confidence"`
	SentimentScore  float64  `json:"sentimentScore"`
	SentimentLabel  string   `json:"sentimentLabel"`
	Keywords        []string `json:"keywords"`
	AnalysisMethod  string   `json:"analysisMethod"`
}

// AnalyzeEmotionProgrammatically performs rule-based emotion analysis
func AnalyzeEmotionProgrammatically(content string) (*EmotionAnalysisResult, error) {
	content = strings.ToLower(content)

	// Initialize emotion scores
	emotions := map[string]float64{
		"joy":      0,
		"sadness":  0,
		"anger":    0,
		"fear":     0,
		"love":     0,
		"surprise": 0,
		"disgust":  0,
	}

	// Count emotion keywords
	totalWords := len(strings.Fields(content))
	if totalWords == 0 {
		totalWords = 1 // Prevent division by zero
	}

	var allKeywords []string

	for emotion, keywords := range EmotionKeywords {
		matchCount := 0
		for _, keyword := range keywords {
			// Use word boundaries to match complete words
			pattern := `\b` + regexp.QuoteMeta(keyword) + `\b`
			matched, _ := regexp.MatchString(pattern, content)
			if matched {
				matchCount++
				allKeywords = append(allKeywords, keyword)
			}
		}
		// Normalize score based on content length
		emotions[emotion] = float64(matchCount) / float64(totalWords) * 10
		if emotions[emotion] > 1 {
			emotions[emotion] = 1 // Cap at 1.0
		}
	}

	// Calculate sentiment
	positiveCount := 0
	negativeCount := 0

	for _, keyword := range SentimentKeywords["positive"] {
		pattern := `\b` + regexp.QuoteMeta(keyword) + `\b`
		matched, _ := regexp.MatchString(pattern, content)
		if matched {
			positiveCount++
		}
	}

	for _, keyword := range SentimentKeywords["negative"] {
		pattern := `\b` + regexp.QuoteMeta(keyword) + `\b`
		matched, _ := regexp.MatchString(pattern, content)
		if matched {
			negativeCount++
		}
	}

	sentimentScore := (float64(positiveCount) - float64(negativeCount)) / float64(totalWords) * 10
	if sentimentScore > 1 {
		sentimentScore = 1
	} else if sentimentScore < -1 {
		sentimentScore = -1
	}

	var sentimentLabel string
	if sentimentScore > 0.1 {
		sentimentLabel = "positive"
	} else if sentimentScore < -0.1 {
		sentimentLabel = "negative"
	} else {
		sentimentLabel = "neutral"
	}

	// Find dominant emotion
	var dominantEmotion string
	var maxScore float64
	for emotion, score := range emotions {
		if score > maxScore {
			maxScore = score
			dominantEmotion = emotion
		}
	}

	if dominantEmotion == "" {
		dominantEmotion = "neutral"
	}

	// Calculate confidence based on how much stronger the dominant emotion is
	confidence := maxScore
	if confidence < 0.1 {
		confidence = 0.1
		dominantEmotion = "neutral"
	}

	return &EmotionAnalysisResult{
		Joy:             emotions["joy"],
		Sadness:         emotions["sadness"],
		Anger:           emotions["anger"],
		Fear:            emotions["fear"],
		Love:            emotions["love"],
		Surprise:        emotions["surprise"],
		Disgust:         emotions["disgust"],
		DominantEmotion: dominantEmotion,
		Confidence:      confidence,
		SentimentScore:  sentimentScore,
		SentimentLabel:  sentimentLabel,
		Keywords:        allKeywords,
		AnalysisMethod:  "programmatic",
	}, nil
}

// AnalyzeEmotionWithAI performs AI-based emotion analysis using Ollama
func AnalyzeEmotionWithAI(content string, ollamaURL string) (*EmotionAnalysisResult, error) {
	if ollamaURL == "" {
		ollamaURL = "http://localhost:11434"
	}

	prompt := fmt.Sprintf(`请分析以下中文文本的情感，并按照JSON格式返回结果。只返回JSON，不要其他文字。

文本：%s

请按照以下格式返回分析结果：
{
  "joy": 0.0-1.0之间的数值,
  "sadness": 0.0-1.0之间的数值,
  "anger": 0.0-1.0之间的数值,
  "fear": 0.0-1.0之间的数值,
  "love": 0.0-1.0之间的数值,
  "surprise": 0.0-1.0之间的数值,
  "disgust": 0.0-1.0之间的数值,
  "dominantEmotion": "主导情感（joy/sadness/anger/fear/love/surprise/disgust/neutral）",
  "confidence": 0.0-1.0之间的置信度,
  "sentimentScore": -1.0到1.0之间的情感分数（负数表示消极，正数表示积极）,
  "sentimentLabel": "positive/negative/neutral",
  "keywords": ["从文本中提取的情感关键词"]
}`, content)

	reqBody := OllamaRequest{
		Model:  "qwen2.5:7b", // 使用通义千问模型，你可以根据需要修改
		Prompt: prompt,
		Stream: false,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %v", err)
	}

	client := &http.Client{
		Timeout: 60 * time.Second,
	}

	resp, err := client.Post(ollamaURL+"/api/generate", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to call Ollama API: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Ollama API returned status: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %v", err)
	}

	var ollamaResp OllamaResponse
	if err := json.Unmarshal(body, &ollamaResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal Ollama response: %v", err)
	}

	// Parse the AI response
	var result EmotionAnalysisResult

	// Try to extract JSON from the response
	jsonStart := strings.Index(ollamaResp.Response, "{")
	jsonEnd := strings.LastIndex(ollamaResp.Response, "}") + 1

	if jsonStart == -1 || jsonEnd <= jsonStart {
		// If no JSON found, fall back to programmatic analysis
		return AnalyzeEmotionProgrammatically(content)
	}

	jsonStr := ollamaResp.Response[jsonStart:jsonEnd]
	if err := json.Unmarshal([]byte(jsonStr), &result); err != nil {
		// If JSON parsing fails, fall back to programmatic analysis
		return AnalyzeEmotionProgrammatically(content)
	}

	result.AnalysisMethod = "ai"

	// Validate and normalize scores
	emotions := []float64{result.Joy, result.Sadness, result.Anger, result.Fear, result.Love, result.Surprise, result.Disgust}
	for i, score := range emotions {
		if score < 0 {
			emotions[i] = 0
		} else if score > 1 {
			emotions[i] = 1
		}
	}
	result.Joy, result.Sadness, result.Anger, result.Fear, result.Love, result.Surprise, result.Disgust = emotions[0], emotions[1], emotions[2], emotions[3], emotions[4], emotions[5], emotions[6]

	if result.Confidence < 0 {
		result.Confidence = 0
	} else if result.Confidence > 1 {
		result.Confidence = 1
	}

	if result.SentimentScore < -1 {
		result.SentimentScore = -1
	} else if result.SentimentScore > 1 {
		result.SentimentScore = 1
	}

	return &result, nil
}

// GenerateID generates a random ID for emotion analysis
func generateAnalysisID() string {
	bytes := make([]byte, 8)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}

// SaveEmotionAnalysis saves emotion analysis result to database
func SaveEmotionAnalysis(diaryID string, userID uint, result *EmotionAnalysisResult) error {
	analysis := &EmotionAnalysis{
		ID:              generateAnalysisID(),
		DiaryID:         diaryID,
		UserID:          userID,
		Joy:             result.Joy,
		Sadness:         result.Sadness,
		Anger:           result.Anger,
		Fear:            result.Fear,
		Love:            result.Love,
		Surprise:        result.Surprise,
		Disgust:         result.Disgust,
		DominantEmotion: result.DominantEmotion,
		Confidence:      result.Confidence,
		SentimentScore:  result.SentimentScore,
		SentimentLabel:  result.SentimentLabel,
		AnalysisMethod:  result.AnalysisMethod,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	if err := analysis.SetKeywords(result.Keywords); err != nil {
		return fmt.Errorf("failed to set keywords: %v", err)
	}

	// Use GORM's Save method which handles both create and update
	if err := gormDB.Save(analysis).Error; err != nil {
		return fmt.Errorf("failed to save emotion analysis: %v", err)
	}

	return nil
}

// GetEmotionAnalysis gets emotion analysis for a diary
func GetEmotionAnalysis(diaryID string, userID uint) (*EmotionAnalysis, error) {
	var analysis EmotionAnalysis
	if err := gormDB.Where("diary_id = ? AND user_id = ?", diaryID, userID).First(&analysis).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil // Not found, but not an error
		}
		return nil, fmt.Errorf("failed to get emotion analysis: %v", err)
	}

	return &analysis, nil
}

// GetUserEmotionTrends gets emotion trends for a user over time
func GetUserEmotionTrends(userID uint, days int) ([]EmotionAnalysis, error) {
	var analyses []EmotionAnalysis

	query := gormDB.Where("user_id = ?", userID)
	if days > 0 {
		since := time.Now().AddDate(0, 0, -days)
		query = query.Where("created_at >= ?", since)
	}

	if err := query.Order("created_at ASC").Find(&analyses).Error; err != nil {
		return nil, fmt.Errorf("failed to get emotion trends: %v", err)
	}

	return analyses, nil
}

// GetEmotionStatistics gets aggregated emotion statistics for a user
func GetEmotionStatistics(userID uint, days int) (map[string]interface{}, error) {
	analyses, err := GetUserEmotionTrends(userID, days)
	if err != nil {
		return nil, err
	}

	if len(analyses) == 0 {
		return map[string]interface{}{
			"totalEntries": 0,
			"averageEmotions": map[string]float64{
				"joy": 0, "sadness": 0, "anger": 0, "fear": 0,
				"love": 0, "surprise": 0, "disgust": 0,
			},
			"sentimentDistribution": map[string]int{
				"positive": 0, "negative": 0, "neutral": 0,
			},
			"emotionDistribution": map[string]int{},
			"topKeywords":         []string{},
		}, nil
	}

	// Calculate averages
	emotionSums := map[string]float64{
		"joy": 0, "sadness": 0, "anger": 0, "fear": 0,
		"love": 0, "surprise": 0, "disgust": 0,
	}

	sentimentCounts := map[string]int{
		"positive": 0, "negative": 0, "neutral": 0,
	}

	emotionCounts := map[string]int{}
	allKeywords := map[string]int{}

	for _, analysis := range analyses {
		emotionSums["joy"] += analysis.Joy
		emotionSums["sadness"] += analysis.Sadness
		emotionSums["anger"] += analysis.Anger
		emotionSums["fear"] += analysis.Fear
		emotionSums["love"] += analysis.Love
		emotionSums["surprise"] += analysis.Surprise
		emotionSums["disgust"] += analysis.Disgust

		sentimentCounts[analysis.SentimentLabel]++
		emotionCounts[analysis.DominantEmotion]++

		keywords := analysis.GetKeywords()
		for _, keyword := range keywords {
			allKeywords[keyword]++
		}
	}

	count := float64(len(analyses))
	averageEmotions := map[string]float64{}
	for emotion, sum := range emotionSums {
		averageEmotions[emotion] = sum / count
	}

	// Get top keywords
	type keywordCount struct {
		Keyword string
		Count   int
	}

	var sortedKeywords []keywordCount
	for keyword, count := range allKeywords {
		sortedKeywords = append(sortedKeywords, keywordCount{keyword, count})
	}

	sort.Slice(sortedKeywords, func(i, j int) bool {
		return sortedKeywords[i].Count > sortedKeywords[j].Count
	})

	topKeywords := []string{}
	maxKeywords := 20
	if len(sortedKeywords) < maxKeywords {
		maxKeywords = len(sortedKeywords)
	}

	for i := 0; i < maxKeywords; i++ {
		topKeywords = append(topKeywords, sortedKeywords[i].Keyword)
	}

	return map[string]interface{}{
		"totalEntries":          len(analyses),
		"averageEmotions":       averageEmotions,
		"sentimentDistribution": sentimentCounts,
		"emotionDistribution":   emotionCounts,
		"topKeywords":           topKeywords,
	}, nil
}

// AnalyzeDiaryEmotion analyzes emotion for a diary and saves the result
func AnalyzeDiaryEmotion(diaryID string, userID uint, content string, useAI bool, ollamaURL string) (*EmotionAnalysisResult, error) {
	// Check if analysis already exists
	existing, err := GetEmotionAnalysis(diaryID, userID)
	if err != nil {
		return nil, err
	}

	var result *EmotionAnalysisResult

	if useAI {
		result, err = AnalyzeEmotionWithAI(content, ollamaURL)
		if err != nil {
			// Fall back to programmatic analysis if AI fails
			result, err = AnalyzeEmotionProgrammatically(content)
			if err != nil {
				return nil, err
			}
		}
	} else {
		result, err = AnalyzeEmotionProgrammatically(content)
		if err != nil {
			return nil, err
		}
	}

	// Save the analysis
	if err := SaveEmotionAnalysis(diaryID, userID, result); err != nil {
		return nil, err
	}

	// If there was an existing analysis, we've updated it
	if existing != nil {
		result.AnalysisMethod = fmt.Sprintf("%s (updated)", result.AnalysisMethod)
	}

	return result, nil
}
