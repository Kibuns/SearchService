package main

import (
	"encoding/json"
	"log"
	"regexp"
	"strings"

	"github.com/Kibuns/SearchService/DAL"
	"github.com/Kibuns/SearchService/Models"
	amqp "github.com/rabbitmq/amqp091-go"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}

func receiveTwoots() {
    conn, err := amqp.Dial("amqp://guest:guest@rabbitmq:5672/") //locally change rabbitmq to localhost
    failOnError(err, "Failed to connect to RabbitMQ")
    defer conn.Close()

    ch, err := conn.Channel()
    failOnError(err, "Failed to open a channel")
    defer ch.Close()

    twootQueue, err := ch.QueueDeclare(
        "twoot", // name
        false,   // durable
        false,   // delete when unused
        false,   // exclusive
        false,   // no-wait
        nil,     // arguments
    )
    failOnError(err, "Failed to declare a queue")

    twootMsgs, err := ch.Consume(
        twootQueue.Name, // queue
        "",     // consumer
        true,   // auto-ack
        false,  // exclusive
        false,  // no-local
        false,  // no-wait
        nil,    // args
    )
    failOnError(err, "Failed to register a consumer")

    var forever chan struct{}

    go func() {
        for d := range twootMsgs {
            var twoot Models.Twoot
            err := json.Unmarshal(d.Body, &twoot)
            failOnError(err, "Error deserializing message body")

            twoot.Hashtags = extractHashtags(twoot.Content);

            log.Printf("Received a message: %+v", twoot)
            DAL.InsertTwoot(twoot)
        }
    }()

    log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
    <-forever
}

func receiveDeleted(){
    conn, err := amqp.Dial("amqp://guest:guest@rabbitmq:5672/") //locally change rabbitmq to localhost
    failOnError(err, "Failed to connect to RabbitMQ")
    defer conn.Close()

    ch, err := conn.Channel()
    failOnError(err, "Failed to open a channel")
    defer ch.Close()

    deleteQueue, err := ch.QueueDeclare(
        "delete_hashtagtwoots", // name
        false,   // durable
        false,   // delete when unused
        false,   // exclusive
        false,   // no-wait
        nil,     // arguments
    )
    failOnError(err, "Failed to declare a queue")

    deleteMsgs, err := ch.Consume(
        deleteQueue.Name, // queue
        "",     // consumer
        true,   // auto-ack
        false,  // exclusive
        false,  // no-local
        false,  // no-wait
        nil,    // args
    )
    failOnError(err, "Failed to register a consumer")

    var forever chan struct{}

    go func() {
        for t := range deleteMsgs {
            var username string
            err := json.Unmarshal(t.Body, &username)
            failOnError(err, "Error deserializing message body")
            log.Printf("Received a message to delete everything regarding user: %+v", username)
            DAL.DeleteAllOfUser(username)
        }
    }()

    log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
    <-forever
}

func extractHashtags(tweetContent string) []string {
    // regular expression to match hashtags
    re := regexp.MustCompile(`#[^\s]+`)
    
    // find all hashtags in the tweetContent
    hashtags := re.FindAllString(tweetContent, -1)
    
    // remove the "#" character from each hashtag
    for i := range hashtags {
        hashtags[i] = strings.TrimPrefix(hashtags[i], "#")
    }
    
    // remove duplicates from the hashtags array
    uniqueHashtags := make(map[string]bool)
    for _, hashtag := range hashtags {
        uniqueHashtags[hashtag] = true
    }
    var result []string
    for hashtag := range uniqueHashtags {
        result = append(result, hashtag)
    }
    
    return result
}