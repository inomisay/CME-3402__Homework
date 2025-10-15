# Decision Tree Classifier from Scratch (CME 3402)

This repository contains the submission for the **CME 3402: Concepts of Programming Languages** course assignment. The project is a from-scratch implementation of the **ID3 (Iterative Dichotomiser 3) decision tree algorithm** in three different programming languages: **Python, Java, and Go**.

The primary goal is to build a classifier without relying on any external machine learning libraries, focusing on the manual calculation of **entropy** and **information gain** to construct the tree.

## ✨ Features

-   **ID3 Algorithm from Scratch**: Core logic for calculating entropy and information gain is implemented manually to determine the best attribute for each split.
-   **Multi-Language Implementation**: The same algorithm is implemented in Python, Java, and Go to compare different programming paradigms.
-   **Flexible Data Parsing**: Reads and processes both `.csv` and `.txt` files. The implementation includes heuristics to auto-detect delimiters like commas, semicolons, or tabs.
-   **Interactive Prediction CLI**: After building the tree, the application enters an interactive mode where users can input new data and receive a class prediction.
-   **Detailed Calculation Output**: The program displays the step-by-step entropy and information gain calculations for each node, providing insight into the tree-building process.
-   **Tree Visualization**:
    -   A clear, text-based representation of the final tree is printed to the console.
    -   Graphical visualizations are generated as `.png` images using Graphviz.
-   **Conflict Handling**: The implementations detect records with identical features but conflicting outcomes and provide options for resolution.

---

## 🚀 Languages & Implementations

This project was a collaborative effort, with each team member responsible for one language implementation.

<table>
  <thead>
    <tr>
      <th>Language</th>
      <th>Developer</th>
      <th>Description</th>
    </tr>
  </thead>
  <tbody>
    <tr>
      <td>🐍 Python</td>
      <td>Ege Yıldırım</td>
      <td>This implementation leverages Python's dynamic typing and data structures. It uses the <code>bigtree</code> library for console visualization and <code>graphviz</code> for generating PNG outputs of the decision path.</td>
    </tr>
    <tr>
      <td>☕ Java</td>
      <td>Sara Mahyanbakhshayesh</td>
      <td>A robust, object-oriented implementation that models the tree components (nodes, dataset, etc.) as classes. It generates a <code>.dot</code> file that can be rendered into an image using Graphviz.</td>
    </tr>
    <tr>
      <td>🐹 Go (Golang)</td>
      <td>Yasamin Valishariatpanahi</td>
      <td>A concurrent and efficient implementation that relies on Go's strong standard library for file parsing and processing. It also exports a <code>.dot</code> file for visualization.</td>
    </tr>
  </tbody>
</table>

---

## ⚙️ How to Run

### Prerequisites
* **Python**: Python 3, `bigtree` and `graphviz` libraries.
* **Java**: JDK 11 or higher.
* **Go**: Go compiler.
* **Graphviz**: Required to convert `.dot` files into images.

### Running the Programs
1.  **Python**
    ```bash
    # Navigate to the Python project directory
    python <your_script_name>.py
    ```

2.  **Java**
    ```bash
    # Navigate to the Java project directory
    javac Main.java
    java Main
    ```

3.  **Go**
    ```bash
    # Navigate to the Go project directory
    go run main.go
    ```
   

---

## 📋 Usage

For all implementations, the program will first prompt you to enter the name of a dataset file.

1.  The program reads the data and displays the step-by-step calculations for finding the root node and subsequent branches.
2.  Once the tree is built, a text-based version is printed to the console.
3.  The program then enters an interactive loop, prompting you for attribute values for a new instance.
4.  After you provide the inputs, it will traverse the tree and output the final prediction.
5.  For the **Java** and **Go** versions, `.dot` files representing the tree are saved. You can convert them to an image using Graphviz:
    ```bash
    dot -Tpng decisionTree.dot -o decisionTree.png
    ```
   

## 📊 Datasets

The algorithm has been tested on the following datasets provided for the assignment:
* `weather.csv`
* `contact_lenses.csv`
* `breast_cancer.csv`

## 👥 Authors

-   **Yasamin Valishariatpanahi** - *Go Implementation*
-   **Sara Mahyanbakhshayesh** - *Java Implementation*
-   **Ege Yıldırım** - *Python Implementation*
