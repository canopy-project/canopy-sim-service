// Copyright 2014-2015 Canopy Services, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package main

import (
    "encoding/json"
    "fmt"
    "net/http"
    "github.com/gorilla/mux"
)

type CanopySimTest struct {
    droneCnt int64
    testname string
}

type CanopySimService struct {
    tests map[string]*CanopySimTest
}

var gService = CanopySimService{
    tests: map[string]*CanopySimTest{},
}

func main() {
    r := mux.NewRouter().StrictSlash(false)

    r.HandleFunc("/drones_started", DronesStartedHandler)
    r.HandleFunc("/batch_report", BatchReportHandler)

    fmt.Println("Starting server on :8383")
    http.ListenAndServe(":8383", r)
}

/*
 * /batch_report
 * Takes payload: 
 *  {
 *      "simHostname" : string,
 *      "testname" : string,
 *      "avgReportPeriod" : numeric,
 *      "avgReportPeriodCount" : numeric,
 *      "responseAvgLatency" : numeric,
 *      "responseAvgLatencyCount" : numeric,
 *      "responseMinLatency" : numeric,
 *      "responseMaxLatency" : numeric,
 *  }
 */
func BatchReportHandler(w http.ResponseWriter, r *http.Request) {
    // Decode payload
    inPayload, err := ReadAndDecodeRequestBody(r)
    if err != nil {
        out := fmt.Sprintf("{\"error\" : \"%s\"}\n", err.Error())
        fmt.Printf("%s", out)
        fmt.Fprintf(w, "%s", out)
        return
    }

    fmt.Println(inPayload)

    // Read fields
    testname, ok := inPayload["testname"].(string)
    if !ok {
        out := fmt.Sprintf("{\"error\" : \"Expected string field \\\"testname\\\"\"}\n")
        fmt.Printf("%s", out)
        fmt.Fprintf(w, "%s", out)
        return
    }

    simHostname, ok := inPayload["simHostname"].(string)
    if !ok {
        out := fmt.Sprintf("{\"error\" : \"Expected string field \\\"simHostname\\\"\"}\n")
        fmt.Printf("%s", out)
        fmt.Fprintf(w, "%s", out)
        return
    }

    avgReportPeriod, ok := inPayload["avgReportPeriod"].(float64)
    if !ok {
        out := fmt.Sprintf("{\"error\" : \"Expected numeric field \\\"avgReportPeriod\\\"\"}\n")
        fmt.Printf("%s", out)
        fmt.Fprintf(w, "%s", out)
        return
    }

    avgReportPeriodCount, ok := inPayload["avgReportPeriodCount"].(float64)
    if !ok {
        out := fmt.Sprintf("{\"error\" : \"Expected numeric field \\\"avgReportPeriodCount\\\"\"}\n")
        fmt.Printf("%s", out)
        fmt.Fprintf(w, "%s", out)
        return
    }

    responseAvgLatency, ok := inPayload["responseAvgLatency"].(float64)
    if !ok {
        out := fmt.Sprintf("{\"error\" : \"Expected numeric field \\\"responseAvgLatency\\\"\"}\n")
        fmt.Printf("%s", out)
        fmt.Fprintf(w, "%s", out)
        return
    }

    responseAvgLatencyCount, ok := inPayload["responseAvgLatencyCount"].(float64)
    if !ok {
        out := fmt.Sprintf("{\"error\" : \"Expected numeric field \\\"responseAvgLatencyCount\\\"\"}\n")
        fmt.Printf("%s", out)
        fmt.Fprintf(w, "%s", out)
        return
    }

    responseMinLatency, ok := inPayload["responseMinLatency"].(float64)
    if !ok {
        out := fmt.Sprintf("{\"error\" : \"Expected numeric field \\\"responseMinLatency\\\"\"}\n")
        fmt.Printf("%s", out)
        fmt.Fprintf(w, "%s", out)
        return
    }

    responseMaxLatency, ok := inPayload["responseMaxLatency"].(float64)
    if !ok {
        out := fmt.Sprintf("{\"error\" : \"Expected numeric field \\\"responseMaxLatency\\\"\"}\n")
        fmt.Printf("%s", out)
        fmt.Fprintf(w, "%s", out)
        return
    }

    // Create test if necessary
    test, ok := gService.tests[testname]
    if !ok {
        fmt.Println("creating test", testname)
        test = &CanopySimTest{
            testname: testname,
        }
        gService.tests[testname] = test
    }
    fmt.Println(testname, simHostname, responseAvgLatency, responseMinLatency, responseMaxLatency, responseAvgLatencyCount, avgReportPeriod, avgReportPeriodCount)

    // Send response
    out := fmt.Sprintf("{\"result\" : \"Thanks for your data!\"}\n")
    fmt.Printf("%s", out)
    fmt.Fprintf(w, "%s", out)
}

/*
 * Takes payload: 
 *  {
 *      "cnt" : 1,
 *      "testname" : "dev02.canopy.link:myTest"
 *  }
 */
func DronesStartedHandler(w http.ResponseWriter, r *http.Request) {
    // Decode payload
    inPayload, err := ReadAndDecodeRequestBody(r)
    if err != nil {
        out := fmt.Sprintf("{\"error\" : \"%s\"}\n", err.Error())
        fmt.Printf("%s", out)
        fmt.Fprintf(w, "%s", out)
        return
    }
   
    // Read fields
    cnt_f64, ok := inPayload["cnt"].(float64)
    if !ok {
        out := fmt.Sprintf("{\"error\" : \"Expected integer field \\\"cnt\\\"\"}\n")
        fmt.Printf("%s", out)
        fmt.Fprintf(w, "%s", out)
        return
    }
    cnt := int64(cnt_f64)

    testname, ok := inPayload["testname"].(string)
    if !ok {
        out := fmt.Sprintf("{\"error\" : \"Expected string field \\\"testname\\\"\"}\n")
        fmt.Printf("%s", out)
        fmt.Fprintf(w, "%s", out)
        return
    }

    // Create test if necessary
    test, ok := gService.tests[testname]
    if !ok {
        fmt.Println("creating test", testname)
        test = &CanopySimTest{
            testname: testname,
        }
        gService.tests[testname] = test
    }
    test.droneCnt += cnt
    out := fmt.Sprintf("{\"testname\" : \"%s\", \"drone_cnt\" : %d}\n", test.testname, test.droneCnt)
    fmt.Printf("%s", out)
    fmt.Fprintf(w, "%s", out)
}

func ReadAndDecodeRequestBody(r *http.Request) (map[string]interface{}, error) {
    var out map[string]interface{}
    decoder := json.NewDecoder(r.Body)
    err := decoder.Decode(&out)
    if err != nil {
        return nil, fmt.Errorf("Error decoding body"+ err.Error()) 
    }
    return out, nil
}

