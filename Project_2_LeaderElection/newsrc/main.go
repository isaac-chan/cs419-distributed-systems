package main

import (
    "os"
    "fmt"
    "time"
    "math/rand"
//    "net/rpc"
)

func main() {

    // validate arguments or print usage
    if len(os.Args) < 2 {
        fmt.Println("usage:", os.Args[0], "thisAddress [thatAddress]...")
        os.Exit(1)
    }

    // process id
    pid := os.Getpid()

    // state
    state := "follower"
    fmt.Println(pid, "INITIAL STATE", state)

    // term number
    term := 0

    // address of this server
    thisAddress := os.Args[1]
    fmt.Println(pid, "LISTEN", thisAddress)

    // addresses of other servers
    thatAddress := os.Args[2:]
    for _,address := range thatAddress {
        fmt.Println(pid, "PEER", address)
    }

    // address of leader
//    leadAddress := ""

    // cluster size
    clusterSize := len(os.Args[1:])
    fmt.Println(pid, "CLUSTER SIZE", clusterSize)

    // votes
    votes := 0

    // election timeout between 1500 and 3000ms
    rand.Seed(int64(pid))
    number :=  1500 + rand.Intn(1500)
    electionTimeout := time.Millisecond * time.Duration(number)
    fmt.Println(pid, "RANDOM TIMEOUT", electionTimeout)

    // heartbeat timeout
    heartbeatTimeout := time.Millisecond * time.Duration(1000)

    // vote timeout
    voteTimeout := time.Millisecond * time.Duration(1000)

    // receive messages from leader on channel
    // TODO
    leaderMsg := make(chan error)
//    leaderMsg := client.Go("", args, foo, nil)

    // event loop
    for {

        switch state {

        case "follower":

            select {

            // receive leader message before timeout
            case <-leaderMsg:
                fmt.Println(pid, "LEADER MESSAGE RECEIVED")

            // otherwise begin election
            case <-time.After(electionTimeout):
                state = "candidate"
                fmt.Println(pid, "ELECTION TIMEOUT")

            }

        case "candidate":

            // increment term
            term++
            fmt.Println(pid, "TERM", term)

            // vote for self
            votes = 1

            // request votes
            // TODO

            // receive messages from voters on channel
            // TODO
            voterMsg := make(chan error)
//          voterMsg := client.Go("", args, foo, nil)

            election: for {
                select {

                // receive votes
                case <-voterMsg:
                    fmt.Println(pid, "VOTE RECEIVED")
                    votes++

                    // if majority of votes, go to leader state
                    if votes > clusterSize/2 {
                        state = "leader"
                        break election
                    }

                // receive leader challenge
                case <-leaderMsg:
                    fmt.Println(pid, "LEADER CHALLENGE RECEIVED")

                    // if that term >= this term, return to follower state
                    // TODO
                    if true {
                        state = "follower"
                        break election
                    }

                // time out and start new election
                case <-time.After(voteTimeout):
                    fmt.Println(pid, "VOTE TIMEOUT")
                    break election
                }
            }

        case "leader":

            // send heartbeat
            // TODO

            // wait
            time.Sleep(heartbeatTimeout)

        }
    }
}
