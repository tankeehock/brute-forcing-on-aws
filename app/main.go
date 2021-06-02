package main

import (
	"context"
	"crypto/tls"
	"flag"
	"fmt"
	"net/http"
	"math"
	"os"
	"runtime"
	"runtime/pprof"
	"strconv"
	"strings"
	"sync"
	"time"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sns"
)

var (
	region  string
	bufSize int
	workers int
	verbose bool
)

func guessKey(ctx context.Context, cancel context.CancelFunc, c <-chan int, length int, format, secret string) (result string) {
	var wg sync.WaitGroup
	keyParts := strings.Split(format, "%s")

	wg.Add(workers)
	for i := 0; i < workers; i++ {
		go func() {
			defer wg.Done()
			buf := make([]byte, length, length)
			guesser := &Guesser{}
			for j := range c {
				select {
				case <-ctx.Done():
					return
				default:
					offset2combination(buf, j, length)
					key := strings.Join(keyParts, string(buf))
					if verbose {
						fmt.Println("Offset:", j, ", Key:", key)
					}
					err := guesser.verifyKey(key, secret)
					if err == nil {
						cancel()
						result = key
						return
					}
				}
			}
		}()
	}
	wg.Wait()
	return
}

func assertWithFail(condition bool, err string) {
	if condition {
		fmt.Println("ERROR: ", err)
		flag.Usage()
		os.Exit(1)
	}
}

func main() {
	start := time.Now()
	n := flag.Int("n", 4, "No. of missing characters in key")
	format := flag.String("format", "", "The format string containing partial access key (eg. AKIA%sXXXXXXXXXXXX)")
	secret := flag.String("secret", "", "The secret key")
	phoneNumber := flag.String("phone-number", "", "Phone number to recieve the SMS notification upon completion (+65XXXXXXXX)")
	random := flag.Bool("random", false, "Generate combinations randomly instead of incrementally")
	nodeIndex := flag.Int("node-index", 0, "Index value of the node")
	numberOfNodes := flag.Int("number-of-nodes", 1, "Total number of node(s)")
	lcg := flag.Bool("lcg", false, "Experimental LCG random mode")
	fairDistribution := flag.Bool("fair-distribution", false, "Ensures work are distribution fairly - Excess combination space will be discarded")
	cpuProfile := flag.String("cpu-profile", "", "Generates a pprof CPU profile")
	flag.StringVar(&region, "region", "us-east-1", "AWS Region")
	flag.IntVar(&bufSize, "bufsize", 100, "Buffer size when generating combinations")
	flag.IntVar(&workers, "workers", runtime.NumCPU(), "Number of Workers")
	flag.BoolVar(&verbose, "verbose", false, "Verbose mode")
	flag.Parse()

	assertWithFail(len(strings.ReplaceAll(*format, "%s", ""))+*n != 20, "AWS access key length must be exactly 20 characters long.")
	assertWithFail(len(*secret) != 40, "AWS secret key length must be exactly 40 characters long.")

	http.DefaultTransport = &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}
	if *cpuProfile != "" {
		f, err := os.Create(*cpuProfile)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}
	maxNumberOfCombination := int(math.Pow(float64(32), float64(*n)))
	combinationSpacePerNode := maxNumberOfCombination / *numberOfNodes
	offset :=  combinationSpacePerNode * (*nodeIndex)
	limit := offset + combinationSpacePerNode
	if (*nodeIndex + 1) == *numberOfNodes && *fairDistribution == false {
		limit = maxNumberOfCombination
	}
	message := "Node Index: " + strconv.Itoa(*nodeIndex) + ", Offset: " + strconv.Itoa(offset) + ", limit: " + strconv.Itoa(limit) + " - "
	ctx, cancel := context.WithCancel(context.Background())
	var c <-chan int
	if limit > 0 && offset < limit {
		c = generateSequenceInRange(ctx, offset, limit, *n, *random, *lcg)
	} else {
		c = generateSequence(ctx, *n, *random, *lcg)
	}
	if key := guessKey(ctx, cancel, c, *n, *format, *secret); key != "" {
		elapsed := time.Since(start)
		message += "Found access key: " + key + " - Time Taken(s): " + elapsed.String()
	} else {
		elapsed := time.Since(start)
		message += "Unable to find access key - Time Taken(s): " + elapsed.String()
	}
	fmt.Println(message)
	if *phoneNumber != "" {
		sess, sessionErr := session.NewSession(&aws.Config{
			Region: aws.String("ap-southeast-1")},
		)
		svc := sns.New(sess)
		if sessionErr == nil {
			params := &sns.PublishInput{
				Message: &message,
				PhoneNumber: phoneNumber,
			}
			_, publishErr := svc.Publish(params)
	
			if publishErr != nil {
				fmt.Println(publishErr.Error())
			} else if verbose{
				fmt.Println("Notified end user via SMS")
			}
		} else {
			fmt.Println(sessionErr.Error())
		}
	}
}
