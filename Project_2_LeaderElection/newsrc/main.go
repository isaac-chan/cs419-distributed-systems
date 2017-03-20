package main

import (
    "os"
    "fmt"
    "time"
    "math/rand"
    "net"
    "net/rpc"
)

type Heartbeat struct {
    LeaderID string
    Term int
}

type HeartbeatResponse struct {
    Success bool
    Term int
}

type VoteRequest struct {
    CandidateID string
    Term int
}

type VoteResponse struct {
    VoteGranted bool
    Term int
}

var leaderMsg chan Heartbeat
var candidateMsg chan VoteRequest
var voterMsg chan VoteResponse

type Message int

//
func (t *Message) AppendEntries(heartbeat Heartbeat, heartbeatResponse *HeartbeatResponse) error {
    leaderMsg <- heartbeat
    return nil
}

//
func (s *Message) RequestVote(voteRequest VoteRequest, voteResponse *VoteResponse) error {
    candidateMsg <- voteRequest
    *voteResponse = <-voterMsg
    return nil
}

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

    //
    leaderMsg = make(chan Heartbeat)
    candidateMsg = make(chan VoteRequest)
    voterMsg = make(chan VoteResponse)

    //
    rpc.Register(new(Message))

    //
    messages, error := net.Listen("tcp", thisAddress)
    if error != nil {
        fmt.Println(pid, "UNABLE TO LISTEN ON", thisAddress)
        os.Exit(1)
    }
    go rpc.Accept(messages)

    // event loop
    for {

        switch state {

        case "follower":

            select {

            // receive leader message before timeout
            case <-leaderMsg:
                fmt.Println(pid, "LEADER MESSAGE RECEIVED")
//                followerMsg <- HeartbeatResponse{Success: true, Term: term}

            // receive vote request
            case <-candidateMsg:
                fmt.Println(pid, "CANDIDATE MESSAGE RECEIVED")
                voterMsg <- VoteResponse{VoteGranted: true, Term: term}

            // otherwise begin election
            case <-time.After(electionTimeout):
                state = "candidate"
                fmt.Println(pid, "ELECTION TIMEOUT")
                fmt.Println(pid, "STATE", state)
            }

        case "candidate":

            // increment term
            term++
            fmt.Println(pid, "TERM", term)

            // vote for self
            votes = 1

            // request votes
            for _,address := range thatAddress {
                go func(address string){
                    client, error := rpc.Dial("tcp", address)
                    if error != nil {
                        fmt.Println(pid, "UNABLE TO DIAL", address)
                    } else {
                        fmt.Println(pid, "REQUEST VOTE FROM", address)
                    }
                    vreq := new(VoteRequest)
                    vreq.CandidateID = thisAddress
                    vreq.Term = term
                    vresp := new(VoteResponse)
                    client.Call("Message.RequestVote", vreq, &vresp)
                    voterMsg <- *vresp
                }(address)
            }

            election: for {
                select {

                // receive votes
                case <-voterMsg:
                    fmt.Println(pid, "VOTE RECEIVED")
                    votes++

                    // if majority of votes, go to leader state
                    if votes > clusterSize/2 {
                        state = "leader"
                        fmt.Println(pid, "STATE", state)
                        break election
                    }

                // receive leader challenge
                case <-leaderMsg:
                    fmt.Println(pid, "LEADER CHALLENGE RECEIVED")

                    // if that term >= this term, return to follower state
                    // TODO
                    if true {
                        state = "follower"
                        fmt.Println(pid, "STATE", state)
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
            for _,address := range thatAddress {
                go func(address string){
                    client, error := rpc.Dial("tcp", address)
                    if error != nil {
                        fmt.Println(pid, "UNABLE TO DIAL", address)
                    } else {
                        fmt.Println(pid, "SEND HEARTBEAT TO", address)
                    }
                    hb := new(Heartbeat)
                    hb.LeaderID = thisAddress
                    hb.Term = term
                    hbr := new(HeartbeatResponse)
                    client.Call("Message.AppendEntries", hb, &hbr)
                    fmt.Println(pid, "RECEIVE HEARTBEAT RESPONSE FROM ", address)
                }(address)
            }

            // wait
            time.Sleep(heartbeatTimeout)

        }
    }
}
