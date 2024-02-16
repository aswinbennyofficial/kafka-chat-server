package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	"log"
	"time"

	"github.com/segmentio/kafka-go"
)

func main() {
    TOPIC_NAME := "TOPIC_NAME"

    keypair, err := tls.LoadX509KeyPair("./keys/service.cert", "./keys/service.key")
    if err != nil {
        log.Fatalf("Failed to load access key and/or access certificate: %s", err)
    }

    caCert, err := ioutil.ReadFile("./keys/ca.pem")
    if err != nil {
        log.Fatalf("Failed to read CA certificate file: %s", err)
    }

    caCertPool := x509.NewCertPool()
    ok := caCertPool.AppendCertsFromPEM(caCert)
    if !ok {
        log.Fatalf("Failed to parse CA certificate file: %s", err)
    }

    dialer := &kafka.Dialer{
        Timeout:   10 * time.Second,
        DualStack: true,
        TLS: &tls.Config{
            Certificates: []tls.Certificate{keypair},
            RootCAs:      caCertPool,
        },
    }

    // init consumer
    consumer := kafka.NewReader(kafka.ReaderConfig{
        Brokers: []string{"kafka-3b4095c8-kafka-chat.a.aivencloud.com:27451"},
        Topic:   TOPIC_NAME,
        Dialer:  dialer,
    })

    for {
        message, err := consumer.ReadMessage(context.Background())

        if err != nil {
            log.Printf("Could not read message: %s", err)
        } else {
            log.Printf("Got message using SSL: %s", message.Value)
        }
    }
}