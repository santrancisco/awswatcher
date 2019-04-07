package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
)

var (
	ErrParsingJson          = errors.New("Cannot unmarshal Message")
	ErrInvalidStatusCode    = errors.New("invalid status code")
	ErrSlackWebhookNotFound = errors.New("slack webhook not found in env variables")
)

type SNSRequest struct {
	Records []struct {
		SNS struct {
			Type       string `json:"Type"`
			Timestamp  string `json:"Timestamp"`
			SNSMessage string `json:"Message"`
			Subject    string `json:"Subject"`
		} `json:"Sns"`
	} `json:"Records"`
}

type Alarm struct {
	AlarmName        string `json:"AlarmName"`
	AWSAccountId     string `json:"AWSAccountId"`
	NewStateReason   string `json:"NewStateReason"`
	Region           string `json:"Region"`
	NewStateValue    string `json:"NewStateValue"`
	StateChangeTime  string `json:"StateChangeTime"`
	AlarmDescription string `json:"AlarmDescription"`
}

// The Slack payload we prepare to be marshaled into JSON later and send to Slack server
type SlackPayload struct {
	Text      string `json:"text"` // To create a link in your text, enclose the URL in <> angle brackets
	Username  string `json:"username,omitempty"`
	IconURL   string `json:"icon_url,omitempty"`
	IconEmoji string `json:"icon_emoji,omitempty"`
	Channel   string `json:"channel,omitempty"`
}

func DirectHandler(alarm Alarm) error {
	// log.Printf("processing message from SNS: %v\n", signinmsg)
	slackURL, found := os.LookupEnv("SLACK_WEBHOOK")
	if !found {
		fmt.Println("No Slack webhook found")
		return ErrSlackWebhookNotFound
	}
	payload := SlackPayload{
		Text: fmt.Sprintf("Alarm for AWS Account %s was triggered: \n    Alarm: %s\n    Time:%s\n    Reason:%s\n    Region:%s\n    Description:%s", alarm.AWSAccountId, alarm.AlarmName, alarm.StateChangeTime, alarm.NewStateReason, alarm.Region, alarm.AlarmDescription),
	}
	log.Printf("Sending to Slack: %s\n", payload.Text)
	payloadJSON, _ := json.Marshal(payload)
	resp, err := http.Post(slackURL, "application/json", bytes.NewBuffer(payloadJSON))
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return ErrInvalidStatusCode
	}
	log.Printf("Slack notified")
	return nil
}

// SNSHandler - Handle SNS event coming from Cloudwatch Log
func SNSHandler(sns SNSRequest) error {
	for i, _ := range sns.Records {
		var alarm Alarm
		err := json.Unmarshal([]byte(sns.Records[i].SNS.SNSMessage), &alarm)
		if err != nil {
			fmt.Println("Error unmarshalling alarm")
			return err
		}
		DirectHandler(alarm)
	}
	return nil
}

func main() {
	lambda.Start(SNSHandler)
}
