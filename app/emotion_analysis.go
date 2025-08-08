package app

import (
	"bytes"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/http"
	"os"
	"regexp"
	"sort"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
	"gorm.io/gorm"
)

// EmotionKeywords will be loaded from JSON file
var EmotionKeywords map[string][]string

// SentimentKeywords will be loaded from JSON file
var SentimentKeywords map[string][]string

// Dictionary structure for JSON loading
type DictionaryData struct {
	Emotions         map[string]EmotionCategory   `json:"emotions"`
	Sentiment        map[string]SentimentCategory `json:"sentiment"`
	ContextModifiers ContextModifiers             `json:"context_modifiers"`
}

type EmotionCategory struct {
	Keywords  []string        `json:"keywords"`
	Intensity IntensityLevels `json:"intensity"`
}

type SentimentCategory struct {
	Keywords []string `json:"keywords"`
}

var (
	dictionaryData        *DictionaryData
	dictionaryInitialized bool
)

// initializeDictionary initializes the emotion and sentiment dictionaries from JSON file
func initializeDictionary() error {
	if dictionaryInitialized {
		return nil
	}

	// Try to load from different possible paths
	possiblePaths := []string{
		"emotion_dictionaries/emotion_keywords.json",
		"./emotion_keywords.json",
		"../emotion_dictionaries/emotion_keywords.json",
	}

	var data []byte
	var err error
	var loadedPath string

	for _, path := range possiblePaths {
		data, err = os.ReadFile(path)
		if err == nil {
			loadedPath = path
			break
		}
	}

	if err != nil {
		return fmt.Errorf("无法加载词典文件: %v", err)
	}

	dictionaryData = &DictionaryData{}
	if err := json.Unmarshal(data, dictionaryData); err != nil {
		return fmt.Errorf("解析词典文件失败 (%s): %v", loadedPath, err)
	}

	// Initialize EmotionKeywords
	EmotionKeywords = make(map[string][]string)
	for emotion, category := range dictionaryData.Emotions {
		EmotionKeywords[emotion] = category.Keywords
	}

	// Initialize SentimentKeywords
	SentimentKeywords = make(map[string][]string)
	for sentiment, category := range dictionaryData.Sentiment {
		SentimentKeywords[sentiment] = category.Keywords
	}

	dictionaryInitialized = true
	return nil
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
	// Initialize dictionary if not already done
	if err := initializeDictionary(); err != nil {
		return nil, fmt.Errorf("初始化词典失败: %v", err)
	}

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

// AnalyzeDiaryEmotion 分析日记情绪（主要接口，推荐使用）
// 自动使用最佳的分析方法，提供准确的情绪分析结果
func AnalyzeDiaryEmotion(diaryID string, userID uint, content string, useAI bool, ollamaURL string) (*EmotionAnalysisResult, error) {
	// 直接使用增强版分析，它包含了所有优化和改进
	return AnalyzeDiaryEmotionEnhanced(diaryID, userID, content, useAI, ollamaURL)
}

// AnalyzeDiaryEmotionBasic 基础版情绪分析（仅在需要简单分析时使用）
func AnalyzeDiaryEmotionBasic(diaryID string, userID uint, content string, useAI bool, ollamaURL string) (*EmotionAnalysisResult, error) {
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

// AnalyzeDiaryEmotionEnhanced 使用增强版分析引擎（推荐使用）
// 这是主要的情绪分析接口，提供最准确的分析结果
func AnalyzeDiaryEmotionEnhanced(diaryID string, userID uint, content string, useAI bool, ollamaURL string) (*EmotionAnalysisResult, error) {
	// 首先尝试使用增强版分析
	if enhancedResult, err := performEnhancedAnalysis(content, useAI, ollamaURL); err == nil {
		// 保存分析结果
		if saveErr := SaveEmotionAnalysis(diaryID, userID, enhancedResult); saveErr != nil {
			return nil, saveErr
		}
		return enhancedResult, nil
	}

	// 如果增强版失败，回退到基础版本
	return AnalyzeDiaryEmotion(diaryID, userID, content, useAI, ollamaURL)
}

// performEnhancedAnalysis 执行增强版情绪分析
func performEnhancedAnalysis(content string, useAI bool, ollamaURL string) (*EmotionAnalysisResult, error) {
	// 初始化增强版词典
	if err := initializeEnhancedDictionary(); err != nil {
		return nil, err
	}

	// 执行增强版分析
	result, err := analyzeWithEnhancedEngine(content)
	if err != nil {
		return nil, err
	}

	// 如果启用AI，进行混合分析
	if useAI {
		if aiResult, aiErr := AnalyzeEmotionWithAI(content, ollamaURL); aiErr == nil {
			result = blendResults(result, aiResult)
		}
	}

	return result, nil
}

// === 增强版情绪分析核心功能 ===

// EnhancedDictionary 增强版情绪词典结构
type EnhancedDictionary struct {
	Emotions         map[string]EnhancedCategory `json:"emotions"`
	ContextModifiers ContextModifiers            `json:"context_modifiers"`
}

// EnhancedCategory 增强版情绪类别
type EnhancedCategory struct {
	Keywords  []string        `json:"keywords"`
	Intensity IntensityLevels `json:"intensity"`
}

// IntensityLevels 强度级别
type IntensityLevels struct {
	High   []string `json:"high"`
	Medium []string `json:"medium"`
	Low    []string `json:"low"`
}

// ContextModifiers 上下文修饰符
type ContextModifiers struct {
	Negation     []string `json:"negation"`
	Intensifiers []string `json:"intensifiers"`
	Diminishers  []string `json:"diminishers"`
}

// KeywordMatch 关键词匹配结果
type KeywordMatch struct {
	Keyword   string  `json:"keyword"`
	Emotion   string  `json:"emotion"`
	Intensity string  `json:"intensity"`
	Weight    float64 `json:"weight"`
	Modified  bool    `json:"modified"`
	Modifier  string  `json:"modifier"`
}

var (
	enhancedDict *EnhancedDictionary
	regexCache   = make(map[string]*regexp.Regexp)

	// 权重配置
	intensityWeights = map[string]float64{
		"high":   1.0,
		"medium": 0.6,
		"low":    0.3,
	}

	modifierWeights = map[string]float64{
		"negation":    -1.0,
		"intensifier": 1.5,
		"diminisher":  0.5,
	}
)

// initializeEnhancedDictionary 初始化增强版词典
func initializeEnhancedDictionary() error {
	if enhancedDict != nil {
		return nil // 已经初始化
	}

	// 尝试加载外部词典文件
	possiblePaths := []string{
		"emotion_dictionaries/emotion_keywords.json",
		"./emotion_keywords.json",
	}

	for _, path := range possiblePaths {
		if dict, err := loadEnhancedDictionary(path); err == nil {
			enhancedDict = dict
			return nil
		}
	}

	// 使用内置词典
	enhancedDict = buildEnhancedDictionary()
	return nil
}

// loadEnhancedDictionary 加载增强版词典文件
func loadEnhancedDictionary(filepath string) (*EnhancedDictionary, error) {
	data, err := os.ReadFile(filepath)
	if err != nil {
		return nil, err
	}

	dict := &EnhancedDictionary{}

	// 根据文件扩展名选择解析方式
	if strings.HasSuffix(strings.ToLower(filepath), ".yaml") || strings.HasSuffix(strings.ToLower(filepath), ".yml") {
		err = yaml.Unmarshal(data, dict)
	} else {
		err = json.Unmarshal(data, dict)
	}

	if err != nil {
		return nil, err
	}

	return dict, nil
}

// buildEnhancedDictionary 构建增强版词典（从JSON数据）
func buildEnhancedDictionary() *EnhancedDictionary {
	// 确保基础词典已初始化
	if err := initializeDictionary(); err != nil {
		// 如果初始化失败，返回基本结构
		return &EnhancedDictionary{
			Emotions: make(map[string]EnhancedCategory),
			ContextModifiers: ContextModifiers{
				Negation:     []string{"不", "没", "无", "非", "未", "not", "no", "never", "none"},
				Intensifiers: []string{"非常", "特别", "极其", "超级", "相当", "很", "十分", "extremely", "very", "really", "quite", "so", "too"},
				Diminishers:  []string{"有点", "稍微", "略", "轻微", "一点", "slightly", "somewhat", "a bit", "a little", "kind of"},
			},
		}
	}

	dict := &EnhancedDictionary{
		Emotions:         make(map[string]EnhancedCategory),
		ContextModifiers: dictionaryData.ContextModifiers,
	}

	// 直接使用JSON文件中的数据
	for emotion, category := range dictionaryData.Emotions {
		dict.Emotions[emotion] = EnhancedCategory{
			Keywords:  category.Keywords,
			Intensity: category.Intensity,
		}
	}

	return dict
}

// analyzeWithEnhancedEngine 使用增强版引擎分析
func analyzeWithEnhancedEngine(content string) (*EmotionAnalysisResult, error) {
	if enhancedDict == nil {
		return nil, fmt.Errorf("enhanced dictionary not initialized")
	}

	content = strings.ToLower(content)
	words := strings.Fields(content)
	totalWords := len(words)

	if totalWords == 0 {
		totalWords = 1
	}

	// 查找关键词匹配
	matches := findEnhancedMatches(content, words, enhancedDict)

	// 计算情绪分数
	emotions := map[string]float64{
		"joy": 0, "sadness": 0, "anger": 0, "fear": 0,
		"love": 0, "surprise": 0, "disgust": 0,
	}

	var allKeywords []string

	for _, match := range matches {
		emotions[match.Emotion] += math.Abs(match.Weight) / float64(totalWords) * 10
		allKeywords = append(allKeywords, match.Keyword)
	}

	// 标准化分数
	for emotion := range emotions {
		if emotions[emotion] > 1 {
			emotions[emotion] = 1
		}
	}

	// 计算情感分数
	sentimentScore := calculateSentiment(matches, totalWords)
	var sentimentLabel string
	if sentimentScore > 0.1 {
		sentimentLabel = "positive"
	} else if sentimentScore < -0.1 {
		sentimentLabel = "negative"
	} else {
		sentimentLabel = "neutral"
	}

	// 找到主导情绪
	var dominantEmotion string
	var maxScore float64
	for emotion, score := range emotions {
		if score > maxScore {
			maxScore = score
			dominantEmotion = emotion
		}
	}

	if dominantEmotion == "" || maxScore < 0.1 {
		dominantEmotion = "neutral"
		maxScore = 0.1
	}

	// 去重关键词
	allKeywords = removeDuplicatesStr(allKeywords)

	return &EmotionAnalysisResult{
		Joy:             emotions["joy"],
		Sadness:         emotions["sadness"],
		Anger:           emotions["anger"],
		Fear:            emotions["fear"],
		Love:            emotions["love"],
		Surprise:        emotions["surprise"],
		Disgust:         emotions["disgust"],
		DominantEmotion: dominantEmotion,
		Confidence:      maxScore,
		SentimentScore:  sentimentScore,
		SentimentLabel:  sentimentLabel,
		Keywords:        allKeywords,
		AnalysisMethod:  "enhanced",
	}, nil
}

// findEnhancedMatches 查找增强版关键词匹配
func findEnhancedMatches(content string, words []string, dict *EnhancedDictionary) []KeywordMatch {
	var matches []KeywordMatch

	// 创建词汇位置索引
	wordPositions := make(map[string][]int)
	for i, word := range words {
		wordPositions[word] = append(wordPositions[word], i)
	}

	// 查找情绪关键词
	for emotion, category := range dict.Emotions {
		// 检查所有关键词
		for _, keyword := range category.Keywords {
			keyword = strings.ToLower(keyword)
			if positions, found := wordPositions[keyword]; found {
				for _, pos := range positions {
					// 确定强度级别
					intensity := getKeywordIntensity(keyword, category)

					match := KeywordMatch{
						Keyword:   keyword,
						Emotion:   emotion,
						Intensity: intensity,
						Weight:    intensityWeights[intensity],
						Modified:  false,
					}

					// 应用上下文修饰符
					match = applyContextModifiersToMatch(match, words, pos, dict.ContextModifiers)
					matches = append(matches, match)
				}
			}
		}
	}

	return matches
}

// getKeywordIntensity 获取关键词强度级别
func getKeywordIntensity(keyword string, category EnhancedCategory) string {
	for _, highKeyword := range category.Intensity.High {
		if strings.ToLower(highKeyword) == keyword {
			return "high"
		}
	}
	for _, lowKeyword := range category.Intensity.Low {
		if strings.ToLower(lowKeyword) == keyword {
			return "low"
		}
	}
	return "medium"
}

// applyContextModifiersToMatch 应用上下文修饰符
func applyContextModifiersToMatch(match KeywordMatch, words []string, position int, modifiers ContextModifiers) KeywordMatch {
	// 检查否定词
	if hasModifierNear(words, position, modifiers.Negation, 3) {
		match.Weight *= modifierWeights["negation"]
		match.Modified = true
		match.Modifier = "negation"
		return match
	}

	// 检查强化词
	if hasModifierNear(words, position, modifiers.Intensifiers, 2) {
		match.Weight *= modifierWeights["intensifier"]
		match.Modified = true
		match.Modifier = "intensifier"
		return match
	}

	// 检查弱化词
	if hasModifierNear(words, position, modifiers.Diminishers, 2) {
		match.Weight *= modifierWeights["diminisher"]
		match.Modified = true
		match.Modifier = "diminisher"
		return match
	}

	return match
}

// hasModifierNear 检查附近是否有修饰符
func hasModifierNear(words []string, position int, modifiers []string, distance int) bool {
	start := maxInt(0, position-distance)
	end := minInt(len(words), position+distance+1)

	for i := start; i < end; i++ {
		if i == position {
			continue
		}
		for _, modifier := range modifiers {
			if strings.EqualFold(words[i], modifier) {
				return true
			}
		}
	}
	return false
}

// calculateSentiment 计算情感分数
func calculateSentiment(matches []KeywordMatch, totalWords int) float64 {
	positiveScore := 0.0
	negativeScore := 0.0

	for _, match := range matches {
		switch match.Emotion {
		case "joy", "love", "surprise":
			positiveScore += math.Abs(match.Weight)
		case "sadness", "anger", "fear", "disgust":
			negativeScore += math.Abs(match.Weight)
		}
	}

	sentimentScore := (positiveScore - negativeScore) / float64(totalWords) * 10

	if sentimentScore > 1 {
		sentimentScore = 1
	} else if sentimentScore < -1 {
		sentimentScore = -1
	}

	return sentimentScore
}

// blendResults 混合AI和增强版分析结果
func blendResults(enhancedResult, aiResult *EmotionAnalysisResult) *EmotionAnalysisResult {
	// 使用加权平均
	aiWeight := 0.6
	enhancedWeight := 0.4

	result := &EmotionAnalysisResult{
		Joy:            enhancedResult.Joy*enhancedWeight + aiResult.Joy*aiWeight,
		Sadness:        enhancedResult.Sadness*enhancedWeight + aiResult.Sadness*aiWeight,
		Anger:          enhancedResult.Anger*enhancedWeight + aiResult.Anger*aiWeight,
		Fear:           enhancedResult.Fear*enhancedWeight + aiResult.Fear*aiWeight,
		Love:           enhancedResult.Love*enhancedWeight + aiResult.Love*aiWeight,
		Surprise:       enhancedResult.Surprise*enhancedWeight + aiResult.Surprise*aiWeight,
		Disgust:        enhancedResult.Disgust*enhancedWeight + aiResult.Disgust*aiWeight,
		Confidence:     enhancedResult.Confidence*enhancedWeight + aiResult.Confidence*aiWeight,
		SentimentScore: enhancedResult.SentimentScore*enhancedWeight + aiResult.SentimentScore*aiWeight,
		AnalysisMethod: "enhanced+ai",
	}

	// 重新确定主导情绪
	emotions := map[string]float64{
		"joy": result.Joy, "sadness": result.Sadness, "anger": result.Anger,
		"fear": result.Fear, "love": result.Love, "surprise": result.Surprise, "disgust": result.Disgust,
	}

	var dominantEmotion string
	var maxScore float64
	for emotion, score := range emotions {
		if score > maxScore {
			maxScore = score
			dominantEmotion = emotion
		}
	}

	if dominantEmotion == "" || maxScore < 0.1 {
		dominantEmotion = "neutral"
	}

	result.DominantEmotion = dominantEmotion

	// 确定情感标签
	if result.SentimentScore > 0.1 {
		result.SentimentLabel = "positive"
	} else if result.SentimentScore < -0.1 {
		result.SentimentLabel = "negative"
	} else {
		result.SentimentLabel = "neutral"
	}

	// 合并关键词
	keywordSet := make(map[string]bool)
	for _, keyword := range enhancedResult.Keywords {
		keywordSet[keyword] = true
	}
	for _, keyword := range aiResult.Keywords {
		keywordSet[keyword] = true
	}

	var combinedKeywords []string
	for keyword := range keywordSet {
		combinedKeywords = append(combinedKeywords, keyword)
	}
	result.Keywords = combinedKeywords

	return result
}

// 辅助函数
func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func removeDuplicatesStr(slice []string) []string {
	keys := make(map[string]bool)
	var result []string

	for _, item := range slice {
		if !keys[item] {
			keys[item] = true
			result = append(result, item)
		}
	}

	return result
}
