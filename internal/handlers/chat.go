package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/nulldiego/oh-back/config"
	"github.com/nulldiego/oh-back/internal/database"
	"github.com/nulldiego/oh-back/internal/model"
	"github.com/nulldiego/oh-back/internal/utils"
	"gorm.io/gorm"
)

func GetUserChats(c *gin.Context) {
	var chats []model.Chat

	userId := utils.GetCurrentUser(c).ID

	pageStr := c.DefaultQuery("page", "1")
	page, _ := strconv.Atoi(pageStr)

	limitStr := c.DefaultQuery("limit", config.Conf.DefaultPageLimit)
	limit, _ := strconv.Atoi(limitStr)

	// TO-DO: Esto no est´a bien, hay que ordenarlo por fecha del ´ultimo message
	queryFunc := func(query *gorm.DB) *gorm.DB {
		query = query.Where("user_id = ?", userId).Order("id DESC")
		return query.Preload("Messages", func(db *gorm.DB) *gorm.DB {
			return db.Select("id, chat_id, author, text, created_at").Limit(limit)
		})
	}

	result, err := utils.Paginated(database.DB, page, limit, queryFunc, &chats)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

func GetMessages(c *gin.Context) {
	chatIdStr := c.Param("chatId")
	chatId, _ := strconv.Atoi(chatIdStr)
	if !chatBelongsToCurrentUser(c, chatId) {
		c.JSON(http.StatusNotFound, gin.H{})
		return
	}

	var messages []model.Message

	pageStr := c.DefaultQuery("page", "1")
	page, _ := strconv.Atoi(pageStr)

	limitStr := c.DefaultQuery("limit", config.Conf.DefaultPageLimit)
	limit, _ := strconv.Atoi(limitStr)

	queryFunc := func(query *gorm.DB) *gorm.DB {
		return query.Where("chat_id = ?", chatId).Order("id DESC")
	}

	result, err := utils.Paginated(database.DB, page, limit, queryFunc, &messages)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

func SendMessage(c *gin.Context) {
	var reqBody struct {
		ChatId int    `json:"chat_id"`
		Text   string `json:"message" binding:"required"`
	}

	if err := c.ShouldBindJSON(&reqBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if reqBody.ChatId == 0 {
		responseMessage, err := newChat(c, reqBody.Text)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, responseMessage)
		return
	}

	if !chatBelongsToCurrentUser(c, reqBody.ChatId) {
		c.JSON(http.StatusNotFound, gin.H{})
		return
	}
	responseMessage, err := newMessage(reqBody.ChatId, reqBody.Text)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, responseMessage)
}

func newMessage(chatId int, message string) (*model.Message, error) {
	uChatId := uint(chatId)

	userMessage := model.Message{
		Author: "HUMAN",
		Text:   message,
		ChatID: uChatId,
	}

	if err := database.DB.Create(&userMessage).Error; err != nil {
		return nil, err
	}

	payload, _ := json.Marshal(map[string]string{"id": string(chatId), "prompt": message})
	payloadBuffer := bytes.NewBuffer(payload)
	resp, err := http.Post("http://192.168.1.117:8000/analyze", "application/json", payloadBuffer)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var respBody struct {
		Text string `json:"text"`
	}

	if err = json.NewDecoder(resp.Body).Decode(&respBody); err != nil {
		return nil, err
	}

	response := model.Message{
		Author: "GPT",
		Text:   respBody.Text,
		ChatID: uChatId,
	}

	if err := database.DB.Create(&response).Error; err != nil {
		return nil, err
	}

	return &response, nil
}

func newChat(c *gin.Context, message string) (*model.Message, error) {
	userId := utils.GetCurrentUser(c).ID

	chat := model.Chat{
		UserID:   userId,
		Messages: []model.Message{{Author: "HUMAN", Text: message}},
	}
	if err := database.DB.Create(&chat).Error; err != nil {
		return nil, err
	}

	payload, _ := json.Marshal(map[string]string{"id": string(chat.ID), "prompt": message})
	payloadBuffer := bytes.NewBuffer(payload)
	resp, err := http.Post("http://192.168.1.117:8000/analyze", "application/json", payloadBuffer)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var respBody struct {
		Text string `json:"text"`
	}

	if err = json.NewDecoder(resp.Body).Decode(&respBody); err != nil {
		return nil, err
	}

	response := model.Message{
		Author: "GPT",
		Text:   respBody.Text,
		ChatID: chat.ID,
	}

	if err := database.DB.Create(&response).Error; err != nil {
		return nil, err
	}

	return &response, nil
}

func chatFromCurrentUser(c *gin.Context, chatId int) *model.Chat {
	var chat model.Chat
	userId := utils.GetCurrentUser(c).ID
	if err := database.DB.First(&chat, chatId).Error; err != nil || chat.UserID != userId {
		return nil
	}
	return &chat
}

func chatBelongsToCurrentUser(c *gin.Context, chatId int) bool {
	if chat := chatFromCurrentUser(c, chatId); chat != nil {
		return true
	}
	return false
}
