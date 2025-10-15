import javax.imageio.ImageIO;
import javax.swing.*;
import java.awt.*;
import java.awt.image.BufferedImage;
import java.io.*;                  // For reading the CSV file
import java.util.*;                // For using List, Map, Scanner, etc.
import java.util.List;

public class DecisionTree {

    // Dataset loaded into memory as a list of rows (string arrays)
    static List<String[]> data = new ArrayList<>();
    static List<String> headers; // Column names
    static Map<String, Set<String>> attributeValues = new HashMap<>();
    static List<String> chosenPath = new ArrayList<>();

    public static void main(String[] args) throws IOException, InterruptedException {
        //weather.csv
        //contact_lenses.csv
        //breast_cancer.csv
        Scanner sc = new Scanner(System.in);

        // Ask for the base name of the dataset (without extension)
        System.out.println(" ");
        System.out.print("Enter dataset base name (e.g., weather, contact_lenses, or breast_cancer): ");
        String baseName = sc.nextLine().trim();

        // Ask for the file format
        String format = "";
        while (true) {
            System.out.print("Enter file format (csv or txt): ");
            format = sc.nextLine().trim().toLowerCase();
            if (format.equals("csv") || format.equals("txt")) {
                break;
            } else {
                System.out.println("Invalid format. Please enter 'csv' or 'txt'.");
            }
        }

        // Construct the filename with the correct extension
        String fileName = baseName + "." + format;

        System.out.println("📂 Loading file: " + fileName + "\n ");

        // Load the dataset from the constructed filename
        try {
            readDataFile(fileName);
        } catch (IOException e) {
            System.err.println("❌ Error reading the file: " + e.getMessage());
            return; // Exit or handle it accordingly
        }


        // Create a list of indices of feature columns (excluding the label column)
        List<Integer> features = new ArrayList<>();
        for (int i = 0; i < headers.size() - 1; i++) {
            features.add(i);
        }

        boolean hasConflicts = detectConflictingDuplicates(data);

        if (hasConflicts) {
            Scanner scanner = new Scanner(System.in);
            System.out.print("❓ Conflicts found. Do you want to continue building the tree? (yes/no): ");
            String answer = scanner.nextLine().trim().toLowerCase();

            if (!answer.equals("yes")) {
                System.out.println("👋 Program exited by user.");
                System.exit(0); // Exit the program
            }
        }


        // Build the decision tree using ID3 algorithm
        //Node root = buildTree(data, features);
        Node root = buildTree(data, features);

        // Print the decision tree in a readable format
        System.out.println("\n################### DECISION TREE (TEXT-BASED) IS ###################");
        //printTree(root, 0);
        printTree2(root, "", true);

        // Visual Tree
        try {
            // Ask user if they want to generate image
            label:
            while (true) {
                System.out.print("\n⚡ Do you want to generate and view an image of the tree? (yes/no, or press 'x' to exit): ");
                String answer = sc.nextLine().trim().toLowerCase();

                switch (answer) {
                    case "x":
                        System.out.println("\uD83D\uDD1A Exiting...");
                        return; // Exit the program

                    case "yes":
                        exportToGraphviz(root, "tree.dot");
                        generatePNGFromDot("tree.dot", "tree.png");
                        //showImageInWindow("tree.png");

                        openImageWithDefaultViewer("tree.png");// Opens and keeps the image window open
                        break label; // Exit the loop after showing the image
                    case "no":
                        System.out.println("Image generation skipped.");
                        break label; // Exit the loop
                    default:
                        System.out.println("Invalid input. Please enter 'yes', 'no', or 'x'.");
                        break;
                }
            }
        } catch (Exception e) {
            System.err.println("\uD83D\uDEE0 Graphviz Error: " + e.getMessage());
        }

        // Allow user to make predictions using the tree
        System.out.println("\n###################------------------------------------ PREDICTION IS ------------------------------------###################");
        System.out.println("⚠️ Please read carefully: your answer must be one of the options listed in parentheses.");
        System.out.println("   If not, you may get an unknown answer and need to try again to find the correct one.\n");
        System.out.println("-------------------------------------------------------------------------------------------------------------------------------");

        while (true) {
            System.out.println("Write (exit) to exit the program or fill it with proper word.");
            List<String> inputs = new ArrayList<>();
            for (int i = 0; i < headers.size() - 1; i++) {
                //System.out.print(headers.get(i) + ": ");
                String attr = headers.get(i);
                Set<String> values = attributeValues.get(attr);
                System.out.print(attr + " " + values + ": ");

                String input = sc.nextLine();
                if (input.equalsIgnoreCase("exit")) {
                    System.out.println("\uD83D\uDD1A Exiting...");
                    System.exit(0); // Fully exit the program
                }
                inputs.add(input);
            }
            // Clearing the chosenPath before calling predict() on every iteration
            chosenPath.clear();
            // Perform prediction
            String prediction = predict(root, inputs);

            if(prediction.equals("Unknown")){
                System.out.println("Prediction: " + prediction);
                System.out.println("\uD83D\uDD04 No prediction found.\nPlease read carefully: your answer must be one of the options listed in parentheses.");
            }else{
                System.out.println("\uD83D\uDCCA Prediction: " + prediction);
                if (!prediction.equalsIgnoreCase("Unknown")) {
                    exportToGraphviz2(root, "highlighted_tree.dot", chosenPath);
                    generatePNGFromDot("highlighted_tree.dot", "highlighted_tree.png");
                    openImageWithDefaultViewer("highlighted_tree.png");
                }
            }
            System.out.println("-------------------------------------------------------------------------------------------------------------------------------");
        }
    }
    // Function to read a CSV and text files and store the headers and rows
    static void readDataFile(String fileName) throws IOException {
        BufferedReader br = new BufferedReader(new FileReader(fileName));
        String headerLine = br.readLine();

        // Detect delimiter: Check for comma, semicolon, tab (\t), and space (in that order)
        String delimiter;
        if (headerLine.contains(",")) {
            delimiter = ",";
        } else if (headerLine.contains(";")) {
            delimiter = ";";
        } else if (headerLine.contains("\t")) {
            delimiter = "\t";
        } else if (headerLine.contains(" ")) {
            delimiter = " ";
        } else if (headerLine.contains(".")) {
            delimiter = ".";
        }else {
            delimiter = ","; // Fallback
        }

        // Process headers
        headers = Arrays.asList(headerLine.split(delimiter));

        // Initialize map for unique attribute values
        attributeValues.clear();
        for (String header : headers) {
            attributeValues.put(header, new HashSet<>());
        }

        data.clear(); // Clear previous data

        String line;
        while ((line = br.readLine()) != null) {
            String[] row = line.split(delimiter);
            data.add(row);

            for (int i = 0; i < headers.size() && i < row.length; i++) {
                String value = row[i].trim();
                attributeValues.get(headers.get(i)).add(value);
            }
        }

        br.close();
    }

    // Checks for conflicting duplicates in the dataset
    static boolean detectConflictingDuplicates(List<String[]> data) {
        Map<String, String> featureToLabel = new HashMap<>();
        Map<String, Integer> featureToLine = new HashMap<>();
        boolean conflictFound = false;

        for (int i = 0; i < data.size(); i++) {
            String[] row = data.get(i);
            String label = row[row.length - 1];
            String featureKey = String.join(",", Arrays.copyOf(row, row.length - 1));
            int lineNumber = i + 2; // Line numbers start from 1

            if (featureToLabel.containsKey(featureKey)) {
                String existingLabel = featureToLabel.get(featureKey);
                int firstLine = featureToLine.get(featureKey);
                if (!existingLabel.equals(label)) {
                    System.out.println("⚠️ Conflict detected:");
                    System.out.println(" - Features: " + featureKey);
                    System.out.println(" - Label '" + existingLabel + "' at line " + firstLine);
                    System.out.println(" - Label '" + label + "' at line " + lineNumber);
                    conflictFound = true;
                }
            } else {
                featureToLabel.put(featureKey, label);
                featureToLine.put(featureKey, lineNumber);
            }
        }

        if (!conflictFound) {
            System.out.println("✅ No conflicting duplicates found.");
        }

        return conflictFound;
    }

    // Recursive function to build the decision tree
    static Node buildTree(List<String[]> subset, List<Integer> features) {
        Map<String, Integer> labelCount = countLabels(subset); // Count class label occurrences

        // Step-wise reporting
        System.out.println("\n==============================");
        System.out.println("📦 Step: Examining Subset (" + subset.size() + " records)");
        System.out.println("Class Distribution:");
        labelCount.forEach((label, count) -> System.out.println(" - " + label + ": " + count));

        // If only one class label exists, return it as a leaf node
        if (labelCount.size() == 1) {
            String result = labelCount.keySet().iterator().next();
            System.out.println("✅ Pure subset. Creating Leaf Node: " + result);
            return new Node(null, result);
        }

        // If no more features to split, return majority class as leaf
        if (features.isEmpty()) {
            String majority = majorityLabel(labelCount);
            System.out.println("⚠️ No remaining features. Creating Leaf Node with Majority Class: " + majority);
            return new Node(null, majority);
        }

        // Step 1: Calculate entropy of current dataset
        double baseEntropy = entropy(subset);
        System.out.println("\nStep 1: Entropy of Current Subset");
        System.out.println("Entropy(S) = " + String.format("%.4f", baseEntropy));

        // Step 2: Information Gain for each feature
        System.out.println("\nStep 2: Information Gain for Each Feature");
        int bestFeature = -1;
        double bestGain = -1;

        for (int feature : features) {
            double gain = informationGain(subset, feature, baseEntropy);
            System.out.printf(" - %s: Gain = %.4f\n", headers.get(feature), gain);
            if (gain > bestGain) {
                bestGain = gain;
                bestFeature = feature;
            }
        }

        // If no feature provides information gain, return majority class
        if (bestFeature == -1) {
            String majority = majorityLabel(labelCount);
            System.out.println("⚠️ No feature improves entropy. Creating Leaf Node with Majority Class: " + majority);
            return new Node(null, majority);
        }

        // Step 3: Splitting based on best feature
        String bestAttr = headers.get(bestFeature);
        System.out.println("\nStep 3: Best Feature to Split: " + bestAttr);
        System.out.println("→ Splitting on: " + bestAttr + " (Gain = " + String.format("%.4f", bestGain) + ")");

        Node node = new Node(bestAttr);
        Map<String, List<String[]>> splits = splitByFeature(subset, bestFeature);
        List<Integer> newFeatures = new ArrayList<>(features);
        newFeatures.remove(Integer.valueOf(bestFeature));

        // Step 4: Recurse on each split
        for (String value : splits.keySet()) {
            System.out.println("\n🔸 Creating Branch: " + bestAttr + " = " + value);
            node.children.put(value, buildTree(splits.get(value), newFeatures));
        }

        return node;
    }

    // Calculate entropy of a dataset
    static double entropy(List<String[]> subset) {
        Map<String, Integer> labelCount = countLabels(subset);
        double entropy = 0;
        int total = subset.size();

        // Compute entropy formula: -sum(p * log2(p))
        for (int count : labelCount.values()) {
            double p = (double) count / total;
            entropy -= p * (Math.log(p) / Math.log(2)); // Negative point calculated her!
        }
        return entropy;
    }

    // Calculate information gain of splitting on a feature
    static double informationGain(List<String[]> subset, int feature, double baseEntropy) {
        Map<String, List<String[]>> splits = splitByFeature(subset, feature);
        List<List<String[]>> groups = new ArrayList<>(splits.values());
        double newEntropy = 0;
        int total = subset.size();

        // Calculate weighted entropy after split using classic for loop
        for (List<String[]> group : groups) {
            double weight = (double) group.size() / total;
            newEntropy += weight * entropy(group);
        }
        return baseEntropy - newEntropy;
    }

    // Split dataset based on the values of a feature
    static Map<String, List<String[]>> splitByFeature(List<String[]> subset, int feature) {
        Map<String, List<String[]>> splits = new HashMap<>();
        for (String[] row : subset) {
            splits.computeIfAbsent(row[feature], k -> new ArrayList<>()).add(row);
        }
        return splits;
    }

    // Count the frequency of each label in the dataset
    static Map<String, Integer> countLabels(List<String[]> subset) {
        Map<String, Integer> count = new HashMap<>();
        for (String[] row : subset) {
            String label = row[row.length - 1]; // Output column (last)
            count.put(label, count.getOrDefault(label, 0) + 1);
        }
        return count;
    }

    // Return the most common label (majority class)
    static String majorityLabel(Map<String, Integer> labelCount) {
        return labelCount.entrySet().stream().max(Map.Entry.comparingByValue()).get().getKey();
    }

    // Improved tree printer with clean branching visuals
    static void printTree2(Node node, String prefix, boolean isLast) {
        if (node.isLeaf()) {
            System.out.println(prefix + (isLast ? "└── " : "├── ") + "--> " + node.label);
            return;
        }

        if (!prefix.isEmpty()) {
            System.out.println(prefix + (isLast ? "└── " : "├── ") + node.attribute);
        } else {
            // Root node
            System.out.println(node.attribute);
        }

        List<Map.Entry<String, Node>> entries = new ArrayList<>(node.children.entrySet());
        for (int i = 0; i < entries.size(); i++) {
            String value = entries.get(i).getKey();
            Node child = entries.get(i).getValue();
            boolean childIsLast = (i == entries.size() - 1);

            if (child.isLeaf()) {
                System.out.println(prefix + (isLast ? "    " : "│   ") + (childIsLast ? "└── " : "├── ") + value + " --> " + child.label);
            } else {
                System.out.println(prefix + (isLast ? "    " : "│   ") + (childIsLast ? "└── " : "├── ") + value);
                printTree2(child, prefix + (isLast ? "    " : "│   ") + (childIsLast ? "    " : "│   "), true);
            }
        }
    }

    // Predict the answer
    static String predict(Node node, List<String> input) {
        while (!node.isLeaf()) {
            // Find which attribute this node tests
            String attribute = node.attribute;

            // Find the index of this attribute in the input
            int index = headers.indexOf(attribute);
            if (index == -1) return "Unknown"; // Attribute not found

            // Normalize input value
            String value = input.get(index).trim().toLowerCase();

            // Now search for a matching child key manually
            boolean found = false;
            for (String key : node.children.keySet()) {
                if (key.trim().toLowerCase().equals(value)) {
                    // Add attribute=value to chosen path
                    chosenPath.add(attribute + "=" + key);
                    node = node.children.get(key);
                    found = true;
                    break;
                }
            }

            if (!found) return "Unknown";
        }

        // Return the final label at the leaf node
        return node.label;
    }

    // Graph functions:
    //######################################################################################################################################
    // Show image 1
    // Exports the decision tree to a Graphviz DOT file for visualization
    static void exportToGraphviz(Node node, String filePath) throws IOException {
        BufferedWriter writer = new BufferedWriter(new FileWriter(filePath));
        writer.write("digraph Tree {\n");
        writer.write("rankdir=TB;\n"); // Top-down layout
        writer.write("node [shape=box, style=filled, color=black, fontname=\"Segoe UI Variable Text Semibold\", fontsize=16];\n");
        generateDotContent(node, writer, new int[]{0}, "N0");
        writer.write("}\n");
        writer.close();
    }

    // Recursively generates the DOT content by traversing the tree nodes and writing node/edge definitions
    //used to generate the full tree without highlighting
    static void generateDotContent(Node node, BufferedWriter writer, int[] id, String name) throws IOException {
        String label = node.isLeaf() ? node.label : node.attribute;
        String fill = node.isLeaf() ? "#ccffcc" : "#ffe0b3"; // pale green / pale orange
        writer.write(name + " [label=\"" + label + "\", fillcolor=\"" + fill + "\"];\n");

        for (Map.Entry<String, Node> entry : node.children.entrySet()) {
            id[0]++;
            String childName = "N" + id[0];
            writer.write(name + " -> " + childName + " [label=\"" + entry.getKey() + "\", fontsize=14];\n");
            generateDotContent(entry.getValue(), writer, id, childName);
        }
    }

    // Uses the Graphviz 'dot' tool to convert a DOT file into a high-resolution PNG image
    static void generatePNGFromDot(String dotFile, String pngFile) throws IOException, InterruptedException {
        //ProcessBuilder pb = new ProcessBuilder("dot", "-Tpng", dotFile, "-o", pngFile);
        ProcessBuilder pb = new ProcessBuilder("dot", "-Tpng", "-Gdpi=300", dotFile, "-o", pngFile);

        Process process = pb.start();
        int exitCode = process.waitFor();
        if (exitCode != 0) {
            throw new RuntimeException("Graphviz conversion failed. Check if 'dot' is installed and in your PATH.");
        }
    }

    // Showing And Opening Image
    static void openImageWithDefaultViewer(String imagePath) throws IOException {
        File imageFile = new File(imagePath);
        if (!Desktop.isDesktopSupported()) {
            System.err.println("Desktop not supported. Cannot open image automatically.");
            return;
        }

        Desktop desktop = Desktop.getDesktop();
        if (imageFile.exists()) {
            desktop.open(imageFile); // Open with system's default image viewer
        } else {
            System.err.println("Image file does not exist: " + imagePath);
        }
    }

    //######################################################################################################################################
    // Show image 2 after running and getting input
    static void exportToGraphviz2(Node root, String filePath, List<String> selections) throws IOException {
        BufferedWriter writer = new BufferedWriter(new FileWriter(filePath));
        writer.write("digraph Tree {\n");
        writer.write("rankdir=TB;\n");
        writer.write("bgcolor=white;\n");
        writer.write("node [shape=box, style=filled, color=black, fillcolor=white, fontname=\"Segoe UI\", fontsize=16];\n");

        Map<Node, String> nodeIds = new HashMap<>();
        int[] counter = new int[]{0};
        assignNodeIds(root, "N0", nodeIds, counter);

        Set<String> visitedNodes = new HashSet<>();
        Set<String> visitedEdges = new HashSet<>();
        followSelectedPath(root, selections, nodeIds, visitedNodes, visitedEdges);

        generateDotContent(root, writer, nodeIds, visitedNodes, visitedEdges);
        writer.write("}\n");
        writer.close();
    }

    // Assign unique IDs to each node
    static void assignNodeIds(Node node, String id, Map<Node, String> nodeIds, int[] counter) {
        nodeIds.put(node, id);
        for (Map.Entry<String, Node> entry : node.children.entrySet()) {
            counter[0]++;
            assignNodeIds(entry.getValue(), "N" + counter[0], nodeIds, counter);
        }
    }

    // Follow path from selections and mark nodes and edges to highlight
    static void followSelectedPath(Node node, List<String> selections, Map<Node, String> nodeIds,
                                   Set<String> visitedNodes, Set<String> visitedEdges) {
        Node current = node;
        visitedNodes.add(nodeIds.get(current));

        for (String sel : selections) {
            String[] parts = sel.split("=", 2); // split into at most 2 parts
            if (parts.length < 2) continue;     // skip if not proper key=value format

            String attr = parts[0];
            String value = parts[1];

            if (!current.attribute.equals(attr)) {
                continue;  // skip if not matching attribute yet
            }

            Node child = current.children.get(value);
            if (child == null) return;

            String fromId = nodeIds.get(current);
            String toId = nodeIds.get(child);

            visitedNodes.add(toId);
            visitedEdges.add(fromId + "->" + toId);

            current = child;

            if (current.isLeaf()) break;
        }

    }

    // Draw nodes and edges
    //designed to highlight a specific decision path based on user input selections.
    static void generateDotContent(Node node, BufferedWriter writer, Map<Node, String> nodeIds,
                                   Set<String> visitedNodes, Set<String> visitedEdges) throws IOException {
        String nodeId = nodeIds.get(node);
        boolean isPath = visitedNodes.contains(nodeId);
        String fill = isPath ? (node.isLeaf() ? "#ccffcc" : "#ffd699") : "white";

        writer.write(String.format("%s [label=\"%s\", fillcolor=\"%s\"];\n", nodeId,
                node.isLeaf() ? node.label : node.attribute, fill));

        for (Map.Entry<String, Node> entry : node.children.entrySet()) {
            Node child = entry.getValue();
            String childId = nodeIds.get(child);
            boolean edgeIsPath = visitedEdges.contains(nodeId + "->" + childId);

            writer.write(String.format("%s -> %s [label=\"%s\", color=\"%s\"];\n", nodeId, childId,
                    entry.getKey(), edgeIsPath ? "orange" : "black"));

            generateDotContent(child, writer, nodeIds, visitedNodes, visitedEdges);
        }
    }
    //######################################################################################################################################

    // Node class represents a decision node or a leaf node in the tree
    static class Node {
        String attribute;                        // The attribute used to split data at this node
        Map<String, Node> children = new LinkedHashMap<>(); // Child nodes for each possible attribute value
        String label;                            // If it's a leaf, this is the predicted class label

        Node(String attribute) {
            this.attribute = attribute;          // Constructor for decision node
        }

        Node(String attribute, String label) {
            this.attribute = attribute;          // Constructor for leaf node with class label
            this.label = label;
        }

        boolean isLeaf() {
            return label != null;                // Returns true if it's a leaf node
        }
    }
}
