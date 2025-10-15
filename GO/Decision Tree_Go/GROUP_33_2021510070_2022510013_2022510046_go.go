package main

import (
	"encoding/csv"  // Handles reading CSV files
	"fmt"           // Prints to terminal
	"log"           // Logs errors
	"math"          // Provides mathematical functions like log2 in entropy formula
	"os"            // Handles file operations
	"path/filepath" // Handles file path operations (joining directory and filename safely across OSes)
	"slices"        // Provides functions for slice operations like cloning and sorting
	"strings"       // Handles string operations
	"time"          // Provides current date/time utilities (for timestamping prediction filenames)
)

var outputMode string             // "terminal", "file", "none"
var outputBuilder strings.Builder // Buffer to write to file if needed

// ANSI color codes (used only if outputMode == "terminal")
const (
	green  = "\033[32m"
	yellow = "\033[33m"
	blue   = "\033[34m"
	reset  = "\033[0m"
)

func colorize(text string, colorCode string) string {
	if outputMode == "terminal" {
		return colorCode + text + reset
	}
	return text
}

// Represents a node in the Decision Tree
type DecisionTreeNode struct {
	SplitAttr     string                       // Attribute to split on
	Branches      map[string]*DecisionTreeNode // Child nodes for each value
	FinalDecision string                       // Leaf node result
	IsLeafNode    bool                         // Whether this node is a leaf or not
}

// Reads CSV or TXT file and returns rows and headers
func readCSVFile(filePath string) ([][]string, []string) {
	// Read entire file content
	content, err := os.ReadFile(filePath)
	if err != nil {
		log.Fatal("Failed to read file:", err)
	}

	// Detect delimiter
	text := string(content)
	lines := strings.Split(text, "\n")
	if len(lines) < 2 {
		log.Fatal("CSV/TXT file must contain at least a header and one data row.")
	}

	sample := lines[0]
	delimiters := []rune{',', ';', '\t', '.', '|'}
	bestDelimiter := ','
	maxCount := -1
	for _, delim := range delimiters {
		count := strings.Count(sample, string(delim))
		if count > maxCount {
			bestDelimiter = delim
			maxCount = count
		}
	}

	// Create CSV reader
	reader := csv.NewReader(strings.NewReader(text))
	reader.Comma = bestDelimiter
	reader.TrimLeadingSpace = true

	rows, err := reader.ReadAll()
	if err != nil {
		log.Fatal("Failed to parse CSV content:", err)
	}

	// Trim all cells in all rows
	for i := range rows {
		for j := range rows[i] {
			rows[i][j] = strings.TrimSpace(rows[i][j])
		}
	}

	headers := rows[0]
	return rows[1:], headers
}

func entropyPrint(s string) {
	switch outputMode {
	case "terminal":
		fmt.Print(s)
	case "file":
		outputBuilder.WriteString(s)
	case "none":
		// do nothing
	}
}

func askEntropyOutputPreference() {
	fmt.Println("---------------------------------------------------------")
	fmt.Println("How would you like to view entropy and gain calculations?")
	fmt.Println("1. Show in terminal")
	fmt.Println("2. Save to file (entropy_output.txt)")
	fmt.Println("3. Skip showing calculations")
	fmt.Print("Enter choice (1/2/3): ")

	var choice string
	fmt.Scanln(&choice)

	switch choice {
	case "1":
		outputMode = "terminal"
	case "2":
		outputMode = "file"
	case "3":
		outputMode = "none"
	default:
		fmt.Println("Invalid input. Defaulting to terminal.")
		outputMode = "terminal"
	}
}

// Calculates entropy for a dataset
func calculateEntropy(data [][]string) float64 {
	outcomeCounts := make(map[string]int)
	for _, row := range data {
		outcome := row[len(row)-1]
		outcomeCounts[outcome]++
	}
	totalRecords := float64(len(data))
	entropy := 0.0

	entropyPrint("\n📊 Entropy Calculation\n")
	entropyPrint(strings.Repeat("═", 70) + "\n")
	entropyPrint(fmt.Sprintf("%-12s | %-6s | %-11s | %-13s\n", "Outcome", "Count", "Probability", "Contribution"))
	entropyPrint(strings.Repeat("─", 70) + "\n")

	for outcome, count := range outcomeCounts {
		prob := float64(count) / totalRecords
		contrib := -prob * math.Log2(prob)
		entropy += contrib
		entropyPrint(fmt.Sprintf("%-12s | %-6d | %-11.4f | %-13.4f\n", outcome, count, prob, contrib))
	}

	entropyPrint(strings.Repeat("─", 70) + "\n")
	entropyPrint(fmt.Sprintf("🔹 Total Entropy = %.4f\n", entropy))
	return entropy
}

// Computes information gain for a given attribute
func computeInformationGain(data [][]string, attrIndex int, headers []string) float64 {
	attrLabel := headers[attrIndex]
	entropyPrint("\n" + strings.Repeat("═", 70) + "\n")
	entropyPrint(fmt.Sprintf("📊 Information Gain for attribute: %s\n", attrLabel))
	entropyPrint(strings.Repeat("═", 70) + "\n")

	originalEntropy := calculateEntropy(data)
	entropyPrint(fmt.Sprintf("%sEntropy before split:%s %.4f\n", colorize("🔸 ", yellow), reset, originalEntropy))

	subsets := make(map[string][][]string)
	for _, row := range data {
		val := row[attrIndex]
		subsets[val] = append(subsets[val], row)
	}
	total := float64(len(data))
	weightedEntropy := 0.0

	entropyPrint("\n📎 Attribute Value Splits\n")
	entropyPrint(strings.Repeat("-", 70) + "\n")
	entropyPrint(fmt.Sprintf("%-12s | %-6s | %-9s | %-10s | %-15s\n", "Value", "Count", "Weight", "Entropy", "Contribution"))
	entropyPrint(strings.Repeat("-", 70) + "\n")

	for val, subset := range subsets {
		subsetEntropy := calculateEntropy(subset)
		weight := float64(len(subset)) / total
		contribution := weight * subsetEntropy
		weightedEntropy += contribution
		entropyPrint(fmt.Sprintf("%-12s | %-6d | %-9.4f | %-10.4f | %-15.4f\n", val, len(subset), weight, subsetEntropy, contribution))
	}

	infoGain := originalEntropy - weightedEntropy
	entropyPrint(strings.Repeat("-", 70) + "\n")
	entropyPrint(fmt.Sprintf("%sInformation Gain = %.4f - %.4f = %.4f%s\n", colorize("✅ ", green), originalEntropy, weightedEntropy, infoGain, reset))

	return infoGain
}

// Gets the most frequently occurring decision outcome
func getMostCommonDecision(data [][]string) string {
	counts := make(map[string]int) // Map to count frequency of each decision outcome

	for _, row := range data {
		// Get the decision outcome from the last column
		decisionOutcome := row[len(row)-1]
		counts[decisionOutcome]++ // Increment count for this outcome
	}

	max := 0
	var mostCommon string

	// Find the decision outcome with the highest count
	for decisionOutcome, count := range counts {
		if count > max {
			max = count
			mostCommon = decisionOutcome
		}
	}

	return mostCommon // Return the most frequently occurring decision
}

// Checks if all decision outcomes in the dataset are the same
func allDecisionsSame(data [][]string) bool {
	// Take the decision outcome from the first row as the reference
	firstDecision := data[0][len(data[0])-1]

	// Compare the decision outcome of each row with the first one
	for _, row := range data {
		currentDecision := row[len(row)-1]
		if currentDecision != firstDecision {
			return false // Found a different outcome → not all are the same
		}
	}

	return true // All outcomes are the same
}

// Returns the conflicting keys and user's preference: true = use "Can't decide", false = use most common
func detectConflicts(data [][]string, headers []string) (map[string]bool, bool) {
	conflictMap := make(map[string]map[string][]int)
	conflictingKeys := make(map[string]bool)

	for i, row := range data {
		key := strings.Join(row[:len(row)-1], "|")
		outcome := row[len(row)-1]

		if _, exists := conflictMap[key]; !exists {
			conflictMap[key] = make(map[string][]int)
		}
		conflictMap[key][outcome] = append(conflictMap[key][outcome], i+2)
	}

	conflictFound := false
	for key, outcomes := range conflictMap {
		if len(outcomes) > 1 {
			conflictFound = true
			conflictingKeys[key] = true
			fmt.Println("⚠️ Conflict detected:")
			fmt.Printf(" - Features: %s\n", strings.ReplaceAll(key, "|", ","))
			for outcome, lines := range outcomes {
				fmt.Printf("   - Outcome '%s' at lines: %v\n", outcome, lines)
			}
		}
	}

	// Ask user preference only if conflicts exist
	useCantDecide := false
	if conflictFound {
		fmt.Println("\nThere are conflicts in the data (same inputs → different outcomes).")
		fmt.Println("1. Use 'Can't decide' for conflicting inputs")
		fmt.Println("2. Automatically choose the most common decision")
		fmt.Print("Choose how to handle them (1/2): ")

		var choice string
		fmt.Scanln(&choice)
		useCantDecide = (choice == "1")
	}

	return conflictingKeys, useCantDecide
}

// Recursively constructs a decision tree using the ID3 algorithm.
// It selects the best attribute to split on based on information gain and creates
// branches for each unique attribute value.
func buildDecisionTree(data [][]string, headers []string, used []bool, conflictingKeys map[string]bool, useCantDecide bool) *DecisionTreeNode {
	// Base case: If all decisions are the same, return a leaf node with that decision.
	if allDecisionsSame(data) {
		return &DecisionTreeNode{FinalDecision: data[0][len(data[0])-1], IsLeafNode: true}
	}

	// Check if this feature combination has a conflict
	conflictKey := strings.Join(data[0][:len(data[0])-1], "|")
	if conflictingKeys[conflictKey] {
		if useCantDecide {
			return &DecisionTreeNode{FinalDecision: "Can't decide", IsLeafNode: true}
		}
		return &DecisionTreeNode{FinalDecision: getMostCommonDecision(data), IsLeafNode: true}
	}

	bestGain := -1.0
	bestIndex := -1

	// Find the attribute with the highest information gain that hasn't been used yet.
	for i := range headers[:len(headers)-1] { // Exclude the outcome column
		if !used[i] {
			gain := computeInformationGain(data, i, headers)
			if gain > bestGain {
				bestGain = gain
				bestIndex = i
			}
		}
	}

	// If no attribute provides gain, return a leaf node with the most common decision.
	if bestIndex == -1 {
		return &DecisionTreeNode{FinalDecision: getMostCommonDecision(data), IsLeafNode: true}
	}

	// Create a new decision node with the selected attribute.
	node := &DecisionTreeNode{
		SplitAttr: headers[bestIndex],
		Branches:  make(map[string]*DecisionTreeNode),
	}

	// Create a copy of the 'used' slice and mark the current attribute as used.
	usedCopy := slices.Clone(used)
	usedCopy[bestIndex] = true

	// Partition data by the values of the selected attribute.
	partitions := make(map[string][][]string)
	for _, row := range data {
		val := row[bestIndex]
		partitions[val] = append(partitions[val], row)
	}

	// Recursively build subtrees for each partitioned subset.
	for value, subset := range partitions {
		node.Branches[value] = buildDecisionTree(subset, headers, usedCopy, conflictingKeys, useCantDecide)
	}

	return node
}

// Prints the DecisionTree
func printTree(node *DecisionTreeNode, prefix string) {
	if node.IsLeafNode {
		fmt.Println(prefix + green + "📌 Decision: ✅ " + node.FinalDecision + reset)
		return
	}

	fmt.Println(prefix + blue + "📦 Attribute: " + node.SplitAttr + reset)

	branchCount := len(node.Branches)
	i := 0
	for val, child := range node.Branches {
		connector := "├──"
		childPrefix := prefix + "  │   "
		if i == branchCount-1 {
			connector = "└──"
			childPrefix = prefix + "      "
		}
		fmt.Printf("%s  %s %s[%s]%s\n", prefix, connector, yellow, val, reset)
		printTree(child, childPrefix)
		i++
	}
}

// Exports the tree to Graphviz DOT file to visualize a Decision Tree
func exportTreeToDot(tree *DecisionTreeNode, datasetFilename string) {
	var builder strings.Builder

	builder.WriteString(`digraph DecisionTree {
  fontname="Helvetica,Arial,sans-serif";
  labelfontname="Georgia";
  node [fontname="Helvetica", style=filled, fontcolor=black];
  edge [fontname="Helvetica", penwidth=2];
  rankdir=TB;
  bgcolor="white";
  label="Decision Tree";
  labelloc=top;
  labeljust=center;
  fontsize=24;
  nodesep=0.7;
  ranksep=0.8;
`)

	nodeCounter := 0
	conditionCounter := 0

	var walk func(*DecisionTreeNode) int
	walk = func(node *DecisionTreeNode) int {
		currentID := nodeCounter
		nodeCounter++

		if node.IsLeafNode {
			// Green box for leaf (final decision)
			builder.WriteString(fmt.Sprintf(
				`  node%d [label="%s", shape=box, style="rounded,filled", fillcolor="#b3f3b3", color="#2e8b57", penwidth=2];`+"\n",
				currentID, node.FinalDecision))
		} else {
			// Yellow box for decision attribute
			builder.WriteString(fmt.Sprintf(
				`  node%d [label="%s", shape=box, style="rounded,filled", fillcolor="#fef0b3", color="#e6ac00", penwidth=2];`+"\n",
				currentID, node.SplitAttr))
		}

		for val, child := range node.Branches {
			childID := walk(child)

			condID := 10000 + conditionCounter
			conditionCounter++

			// Blue ellipse for condition label
			builder.WriteString(fmt.Sprintf(
				`  node%d [label="%s", shape=ellipse, fillcolor="#eaf4ff", color="#6495ed", fontcolor="#1e3f66", penwidth=1.6];`+"\n",
				condID, val))

			// Connect attribute → condition → child
			builder.WriteString(fmt.Sprintf("  node%d -> node%d [color=gray50];\n", currentID, condID))
			builder.WriteString(fmt.Sprintf("  node%d -> node%d [color=gray50];\n", condID, childID))
		}
		return currentID
	}

	walk(tree)
	builder.WriteString("}\n")

	// File output
	treeFolder := "decision_tree"
	os.MkdirAll(treeFolder, 0755)
	baseName := strings.TrimSuffix(filepath.Base(datasetFilename), filepath.Ext(datasetFilename))
	dotPath := filepath.Join(treeFolder, fmt.Sprintf("%s_decisionTree.dot", baseName))
	pngPath := filepath.Join(treeFolder, fmt.Sprintf("%s_decisionTree.png", baseName))

	err := os.WriteFile(dotPath, []byte(builder.String()), 0644)
	if err != nil {
		log.Fatal("Failed to write DOT file:", err)
	}

	fmt.Println("📦 DOT file saved at:", dotPath)
	fmt.Println("For high-resolution PNG:")
	fmt.Printf("👉 Use: dot -Tpng -Gdpi=300 %s -o %s\n", dotPath, pngPath)
	fmt.Println("If the Tree is Large:")
	fmt.Printf("👉 Use: dot -Tpng -Gdpi=300 -Gscale=2 %s -o %s\n", dotPath, pngPath)
	fmt.Println("SVG Format:")
	fmt.Printf("👉 Use: dot -Tsvg %s -o %s\n", dotPath, pngPath)
}

// Predicts output label for given input
// It traverses the decision tree to predict the output label for a given input.
// It recursively follows the branches based on the input attributes until it reaches a leaf node.
func predict(tree *DecisionTreeNode, headers []string, input map[string]string) string {
	// If the current node is a leaf, return its final decision label.
	if tree.IsLeafNode {
		return tree.FinalDecision
	}

	// Get the attribute value from input corresponding to the current split attribute.
	val := input[tree.SplitAttr]

	// Check if there is a branch for this attribute value.
	child, exists := tree.Branches[val]
	if !exists {
		// If no branch exists for this value, return "Unknown" as the prediction.
		return "Unknown"
	}

	// Recursively predict using the subtree for the matched attribute value.
	return predict(child, headers, input)
}

func getUniqueValues(data [][]string, attrIndex int) []string {
	uniqueMap := make(map[string]bool)
	for _, row := range data {
		if attrIndex < len(row) {
			val := row[attrIndex]
			uniqueMap[val] = true
		}
	}
	var result []string
	for val := range uniqueMap {
		result = append(result, val)
	}
	return result
}

// Exports the tree to Graphviz DOT file to visualize User Prediciton's Decision Tree
func exportPredictionPath(tree *DecisionTreeNode, input map[string]string, filename string) {
	var builder strings.Builder

	// Graph styling
	builder.WriteString(`digraph PredictionPath {
  fontname="Helvetica,Arial,sans-serif";
  labelfontname="Georgia";
  node [fontname="Helvetica", style=filled, fontcolor=black];
  edge [fontname="Helvetica", color=gray50, fontcolor=gray30, penwidth=1.6];
  rankdir=TB;
  bgcolor="white";
  label="Prediction Path";
  labelloc=top;
  labeljust=center;
  fontsize=22;
  nodesep=0.7;
  ranksep=0.9;
`)

	nodeCounter := 0

	// Recursive traversal
	var walk func(*DecisionTreeNode, string, bool)
	walk = func(node *DecisionTreeNode, id string, highlight bool) {
		if node.IsLeafNode {
			style := `shape=box, style="rounded,filled", fillcolor="#b3f3b3", color="#2e8b57", penwidth=2`
			if highlight {
				style = `shape=box, style="rounded,filled", fillcolor="#b3f3b3", color=red, penwidth=2.4`
			}
			builder.WriteString(fmt.Sprintf("  %s [label=\"%s\", %s];\n", id, node.FinalDecision, style))
			return
		}

		style := `shape=box, style="rounded,filled", fillcolor="#fef0b3", color="#e6ac00", penwidth=2`
		if highlight {
			style = `shape=box, style="rounded,filled", fillcolor="#fef0b3", color=red, penwidth=2.4`
		}
		builder.WriteString(fmt.Sprintf("  %s [label=\"%s\", %s];\n", id, node.SplitAttr, style))

		for val, child := range node.Branches {
			condID := fmt.Sprintf("cond%d", nodeCounter)
			nodeCounter++
			childID := fmt.Sprintf("node%d", nodeCounter)
			nodeCounter++

			isPath := highlight && input[node.SplitAttr] == val

			condStyle := `shape=ellipse, style=filled, fillcolor="#eaf4ff", color="#6495ed", fontcolor="#1e3f66", penwidth=1.6`
			if isPath {
				condStyle = `shape=ellipse, style=filled, fillcolor="#eaf4ff", color=red, fontcolor="#1e3f66", penwidth=2.4`
			}
			builder.WriteString(fmt.Sprintf("  %s [label=\"%s\", %s];\n", condID, val, condStyle))

			builder.WriteString(fmt.Sprintf("  %s -> %s;\n", id, condID))
			builder.WriteString(fmt.Sprintf("  %s -> %s;\n", condID, childID))

			// Recurse
			walk(child, childID, isPath)
		}
	}

	// Start traversal
	walk(tree, "node0", true)

	builder.WriteString("}\n")

	// Save output
	folder := "prediction_paths"
	os.MkdirAll(folder, 0755)
	filePath := filepath.Join(folder, filename)
	err := os.WriteFile(filePath, []byte(builder.String()), 0644)
	if err != nil {
		log.Fatal("❌ Failed to write prediction path DOT file:", err)
	}

	fmt.Println("🌳 Prediction path DOT file saved at:", filePath)
	fmt.Printf("👉 Use: dot -Tpng -Gdpi=300 %s -o %s.png\n", filePath, strings.TrimSuffix(filePath, ".dot"))
}

func main() {
	// Prompt user to enter CSV/TXT file name
	fmt.Print("Enter CSV/TXT file name: ")
	var fileName string
	fmt.Scanln(&fileName) // weather.csv, contact_lenses.csv, breast_cancer.csv

	// Read dataset and headers from the provided CSV file
	data, headers := readCSVFile(fileName)

	// Detect and log conflicting rows with identical inputs but different outcomes
	conflictingKeys, useCantDecide := detectConflicts(data, headers)

	// Ask the user how they want entropy output handled
	askEntropyOutputPreference()

	// Track which attributes are used
	used := make([]bool, len(headers)-1)

	// Build decision tree
	tree := buildDecisionTree(data, headers, used, conflictingKeys, useCantDecide)

	// Output entropy details if file mode was selected
	if outputMode == "file" {
		err := os.WriteFile("entropy_output.txt", []byte(outputBuilder.String()), 0644)
		if err != nil {
			fmt.Println("❌ Failed to write entropy_output.txt:", err)
		} else {
			fmt.Println("✅ Entropy details saved to entropy_output.txt")
		}
	}

	// Print the tree and export
	fmt.Println("\n=== Decision Tree Structure ===\n")
	printTree(tree, "")
	exportTreeToDot(tree, fileName)

	// Build attribute values for input validation
	attrValues := make(map[string][]string)
	for i := 0; i < len(headers)-1; i++ {
		attrValues[headers[i]] = getUniqueValues(data, i)
	}

	// Loop for prediction
	for {
		fmt.Println("\nPlease enter input values for prediction (type 'exit' at any prompt to quit):")
		input := make(map[string]string)

		for i := 0; i < len(headers)-1; i++ {
			attr := headers[i]
			for {
				fmt.Printf("  Enter value for %s (%s): ", attr, strings.Join(attrValues[attr], ", "))
				var val string
				fmt.Scanln(&val)

				if strings.ToLower(val) == "exit" {
					fmt.Println("Exiting prediction. Goodbye!")
					return
				}

				isValid := false
				for _, validVal := range attrValues[attr] {
					if val == validVal {
						isValid = true
						break
					}
				}

				if isValid {
					input[attr] = val
					break
				} else {
					fmt.Println("  ❌ Invalid input. Please choose one of the listed valid values.")
				}
			}
		}

		// Predict
		result := predict(tree, headers, input)
		fmt.Printf("\n>>> Predicted Decision (%s): %s\n", headers[len(headers)-1], result)

		// Export prediction path visualization
		timestamp := time.Now().Format("20060102_150405")
		dotFilename := fmt.Sprintf("prediction_%s.dot", timestamp)
		exportPredictionPath(tree, input, dotFilename)
	}
}
