package main

import (
    "os"
    "fmt"
    "time"
    "math/rand"
    "net/rpc"
    "net"
    "net/http"
//    "errors"
)

//methods we want to export:
//	leader heartbeats - periodically sends to followers for proof of life	
//	candidate vote requests - sends to other servers for vote response
//	follower votes 

//incoming args
type Args struct {
    msg, vote string
}

//type to export
type someMsg string

//leader heartbeat comes as a string, just send it back to reset timeout
func (t *someMsg) leader_heartbeats(args *Args, reply *string) error {
    *reply = args.msg
    return nil
}

//vote requests come as string, just send it back to increment internal vote count
func (t *someMsg) candidate_vote_requests(args *Args, reply *string) error{
    *reply = args.msg
    return nil
}

//not sure when this is called - reply might be int?
func (t *someMsg) voter_votes (args *Args, reply *string) error{
    *reply = args.vote
    return nil
}

func main() {

    // validate arguments or print usage
    if len(os.Args) < 2 {
        fmt.Println("usage:", os.Args[0], "thisAddress [thatAddress]...")
        os.Exit(1)
    }

    // server calls for HTTP service
    newSomeMsg := new(someMsg)
    rpc.Register(newSomeMsg)
    rpc.HandleHTTP()
    l, e := net.Listen("tcp", ":1234")
    if e != nil {
	fmt.Println("listen error:", e)
	os.Exit(1)
    }
    go http.Serve(l, nil)

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



    /********receive messages from leader on channel*******/
    //invoke client - client dials server. this makes the client channel
    client, err := rpc.DialHTTP("tcp", thisAddress)
    if err != nil {
	fmt.Println("dialing:", err)
	os.Exit(1)
}

    //this should asynchronously receive the message from the leader 
    //should this be in event loop?
    leaderMsg := make(chan error, 1)
    leaderMsg = client.Call("someMsg.leader_heartbeats", "msg", "msg")
    select {
	case err := <-leaderMsg:
	    fmt.Println("leader heartbeat response error:", err)
	case <-time.After(heartbeatTimeout):
	    //TODO
	    //become candidate
    }
    fmt.Println("leader heartbeat received")
    /*******************************************************/



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
