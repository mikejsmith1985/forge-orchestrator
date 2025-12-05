package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

type Flow struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	Data      string `json:"data"`
	Status    string `json:"status"`
	CreatedAt string `json:"created_at"`
}

func main() {
	log.Println("=== Contract #035 E2E Test: Flow Editor API ===\n")

	// Wait for server to be ready
	time.Sleep(500 * time.Millisecond)

	allPassed := true
	var createdFlowID int

	// Test 1: Create a new flow
	log.Println("Test 1: Create new flow → Save → Appears in list")
	createdFlowID, passed := testCreateFlow()
	if !passed {
		allPassed = false
	}

	// Test 2: Load existing flow
	log.Println("\nTest 2: Load existing flow → Verify data loads correctly")
	if !testLoadFlow(createdFlowID) {
		allPassed = false
	}

	// Test 3: Modify and save flow
	log.Println("\nTest 3: Modify flow → Save → Changes persist")
	if !testUpdateFlow(createdFlowID) {
		allPassed = false
	}

	// Test 4: Verify flow appears in list
	log.Println("\nTest 4: Verify flow appears in list")
	if !testFlowInList(createdFlowID) {
		allPassed = false
	}

	// Test 5: Delete flow
	log.Println("\nTest 5: Delete flow → Removed from list")
	if !testDeleteFlow(createdFlowID) {
		allPassed = false
	}

	// Test 6: Verify flow is removed from list
	log.Println("\nTest 6: Verify flow is removed from list")
	if !testFlowNotInList(createdFlowID) {
		allPassed = false
	}

	fmt.Println("\n=================================")
	if allPassed {
		fmt.Println("🎉 All Flow Editor E2E Tests PASSED!")
		fmt.Println("✅ Create new flow and save")
		fmt.Println("✅ Load existing flow with data")
		fmt.Println("✅ Modify and save changes")
		fmt.Println("✅ Flow appears in list")
		fmt.Println("✅ Delete removes flow from database")
		fmt.Println("✅ Deleted flow removed from list")
	} else {
		fmt.Println("❌ Some tests FAILED")
		os.Exit(1)
	}
	fmt.Println("=================================")
}

func testCreateFlow() (int, bool) {
	// Create flow data matching ReactFlow format
	flowData := map[string]interface{}{
		"name": "Test Flow E2E",
		"data": `{"nodes":[{"id":"1","type":"input","data":{"label":"Start Node"},"position":{"x":250,"y":5}},{"id":"2","type":"default","data":{"label":"Agent Node"},"position":{"x":250,"y":100}}],"edges":[{"id":"e1-2","source":"1","target":"2"}]}`,
		"status": "active",
	}

	jsonData, _ := json.Marshal(flowData)
	resp, err := http.Post("http://localhost:8080/api/flows", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Printf("❌ Failed to create flow: %v", err)
		return 0, false
	}
	defer resp.Body.Close()

	if resp.StatusCode != 201 && resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		log.Printf("❌ Expected status 201/200, got %d: %s", resp.StatusCode, string(body))
		return 0, false
	}

	var createdFlow Flow
	json.NewDecoder(resp.Body).Decode(&createdFlow)

	if createdFlow.ID == 0 {
		log.Printf("❌ Flow ID not returned")
		return 0, false
	}

	log.Printf("✅ Flow created with ID: %d", createdFlow.ID)
	return createdFlow.ID, true
}

func testLoadFlow(flowID int) bool {
	resp, err := http.Get(fmt.Sprintf("http://localhost:8080/api/flows/%d", flowID))
	if err != nil {
		log.Printf("❌ Failed to load flow: %v", err)
		return false
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		log.Printf("❌ Expected status 200, got %d", resp.StatusCode)
		return false
	}

	var flow Flow
	json.NewDecoder(resp.Body).Decode(&flow)

	if flow.Name != "Test Flow E2E" {
		log.Printf("❌ Expected name 'Test Flow E2E', got '%s'", flow.Name)
		return false
	}

	// Verify data can be parsed
	var graphData map[string]interface{}
	if err := json.Unmarshal([]byte(flow.Data), &graphData); err != nil {
		log.Printf("❌ Failed to parse flow data: %v", err)
		return false
	}

	nodes, ok := graphData["nodes"].([]interface{})
	if !ok || len(nodes) != 2 {
		log.Printf("❌ Expected 2 nodes, got %v", len(nodes))
		return false
	}

	log.Printf("✅ Flow loaded with %d nodes and data correctly parsed", len(nodes))
	return true
}

func testUpdateFlow(flowID int) bool {
	// Update with new name and additional node
	updatedData := map[string]interface{}{
		"name": "Updated Test Flow",
		"data": `{"nodes":[{"id":"1","type":"input","data":{"label":"Start Node"},"position":{"x":250,"y":5}},{"id":"2","type":"default","data":{"label":"Agent Node"},"position":{"x":250,"y":100}},{"id":"3","type":"output","data":{"label":"Output Node"},"position":{"x":250,"y":200}}],"edges":[{"id":"e1-2","source":"1","target":"2"},{"id":"e2-3","source":"2","target":"3"}]}`,
		"status": "active",
	}

	jsonData, _ := json.Marshal(updatedData)
	req, _ := http.NewRequest("PUT", fmt.Sprintf("http://localhost:8080/api/flows/%d", flowID), bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("❌ Failed to update flow: %v", err)
		return false
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		log.Printf("❌ Expected status 200, got %d: %s", resp.StatusCode, string(body))
		return false
	}

	// Verify the update persisted
	getResp, _ := http.Get(fmt.Sprintf("http://localhost:8080/api/flows/%d", flowID))
	defer getResp.Body.Close()

	var flow Flow
	json.NewDecoder(getResp.Body).Decode(&flow)

	if flow.Name != "Updated Test Flow" {
		log.Printf("❌ Update didn't persist. Name is '%s'", flow.Name)
		return false
	}

	var graphData map[string]interface{}
	json.Unmarshal([]byte(flow.Data), &graphData)
	nodes := graphData["nodes"].([]interface{})

	if len(nodes) != 3 {
		log.Printf("❌ Expected 3 nodes after update, got %d", len(nodes))
		return false
	}

	log.Printf("✅ Flow updated successfully with %d nodes", len(nodes))
	return true
}

func testFlowInList(flowID int) bool {
	resp, err := http.Get("http://localhost:8080/api/flows")
	if err != nil {
		log.Printf("❌ Failed to get flow list: %v", err)
		return false
	}
	defer resp.Body.Close()

	var flows []Flow
	json.NewDecoder(resp.Body).Decode(&flows)

	for _, f := range flows {
		if f.ID == flowID {
			log.Printf("✅ Flow %d found in list with name '%s'", flowID, f.Name)
			return true
		}
	}

	log.Printf("❌ Flow %d not found in list", flowID)
	return false
}

func testDeleteFlow(flowID int) bool {
	req, _ := http.NewRequest("DELETE", fmt.Sprintf("http://localhost:8080/api/flows/%d", flowID), nil)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("❌ Failed to delete flow: %v", err)
		return false
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 && resp.StatusCode != 204 {
		body, _ := io.ReadAll(resp.Body)
		log.Printf("❌ Expected status 200/204, got %d: %s", resp.StatusCode, string(body))
		return false
	}

	log.Printf("✅ Flow %d deleted successfully", flowID)
	return true
}

func testFlowNotInList(flowID int) bool {
	resp, err := http.Get("http://localhost:8080/api/flows")
	if err != nil {
		log.Printf("❌ Failed to get flow list: %v", err)
		return false
	}
	defer resp.Body.Close()

	var flows []Flow
	json.NewDecoder(resp.Body).Decode(&flows)

	for _, f := range flows {
		if f.ID == flowID {
			log.Printf("❌ Flow %d still in list after deletion", flowID)
			return false
		}
	}

	log.Printf("✅ Flow %d correctly removed from list", flowID)
	return true
}
