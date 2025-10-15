import csv  # this is for reading csv or txt
import cmath # for logarithm
from bigtree import Node #this is just for making tree and print it on console not for making decision
from graphviz import Digraph # for visualizing decision_tree
import sys

"""
Keeping the possible states of the subtree structures that should occur as a list.
Ascendant structure prevents the main columns that come before the subtree from being used again.
"""
class TreeNodes:
    def __init__(self, sub_list,ascendant):
        self.sub_list = sub_list
        self.ascendant = ascendant

"""
Using the Graphviz library, the tree coloring function is a structure in which those 
in the pattern are colored red and those not in the pattern are colored black. 
Also, the expressions in the column are indicated with square brackets.
"""
def create_colored_tree(root, highlighted_path_objects):
    dot = Digraph()
    added_gv_nodes = set()

    def add_nodes_edges(node):
        node_gv_id = str(id(node))

        color = 'red' if node in highlighted_path_objects else 'black'

        if node_gv_id not in added_gv_nodes:
            label_node=node.name
            for items in title:
                if items==node.name:
                    label_node="["+node.name+"]"
            dot.node(node_gv_id, label=label_node, color=color)
            added_gv_nodes.add(node_gv_id)

        for child in node.children:
            child_gv_id = str(id(child))
            dot.edge(node_gv_id, child_gv_id)
            add_nodes_edges(child)

    add_nodes_edges(root)
    return dot

"""
The concept expressed by Dictionary is actually the output states of 
the instantaneous tree fragment specifically encountered.
"""
def entropy(dictionary):
    total = 0
    answer=0
    for a in dictionary:
        total+=dictionary[a]
    for b in dictionary:
        p=dictionary[b]/total
        answer+=-p*cmath.log(p,2)
    answer=round(answer.real,3)
    return answer
"""
A parameter called values_dict is used to specify how many output states there are. 
part of the gain calculation is the entropy calculation of values_dict, 
and when calculating from main_dict you need the sum of the values of the elements in values_dict.
The variable parameter called 
main_dict is the dictionary structure created for the first column we have in each calculation.
"""
def information_gain(main_dict,values_dict):
    total=0
    answer=0
    for a in values_dict:
        total+=values_dict[a]
    for b in main_dict:
        p=0
        for c in main_dict[b]:
            p+=main_dict[b][c]/total
        print(b,":",main_dict[b],"->Entropy:",entropy(main_dict[b]))
        answer+=-(p*entropy(main_dict[b]))
    answer+=entropy(values_dict)
    answer = round(answer.real, 3)
    return answer

"""
The main purpose of the tree_maker function is to branch the tree structure, whether it is the main root or any part, 
the variable called top holds the index value of that column to determine which one will be the root for that part.
The variable called Final List creates a list that holds how many of which variables are specific to that column. 
The using_list parameter is a structure where the current possible states are sent from the main list.
set_of_attr is the parameter sent to prevent the previously used columns from being used again.
values_dict is the parameter needed to send to information gain.
"""
def tree_maker(using_list,set_of_attr,values_dict):
    temp_result = 0
    top = ""
    final_list=[]
    for k in range(len(title) - 1):

        temp_dict = {}
        if k not in set_of_attr:
            for j in range(len(using_list)):
                if using_list[j][k] not in temp_dict:
                    temp_dict[using_list[j][k]] = {}
                if using_list[j][len(title) - 1] not in temp_dict[using_list[j][k]]:
                    temp_dict[using_list[j][k]][using_list[j][len(title) - 1]] = 0
                temp_dict[using_list[j][k]][using_list[j][len(title) - 1]] += 1
            print("For:",title[k])
            result = information_gain(temp_dict, values_dict)
            print("Gain:",result,"\n")
            if attribute_possibilities:# In the first case, to get the keys of the columns
                attribute_keys.append(list(temp_dict.keys()))
            if result > temp_result:
                temp_result = result
                final_list = list(temp_dict.items())
                top = k
    return [top,final_list]

"""
This is the function that uses the list created with tree_maker. 
The main purpose here is to create a custom list for each element in the list created 
with tree_maker and append it to the structure called Tree_Nodes_List. 
To give an example, let the results from tree_maker be as follows:(sunny:{yes:3,no:2},overcast:{no:4}) in this case, 
the main purpose of this function is to pull only the sunny ones from the list structure above itself 
and keep it in a list.If we need to talk about the ascendant part, 
the main purpose is to keep it in order not to reuse the previous one.
"""
def child_list_maker(tree_maker_list,using_list,ascendant_set):
    temp_set=ascendant_set.copy()
    for j in range(len(tree_maker_list[1])):
        sub_list=[]
        for k in range(len(using_list)):
            if tree_maker_list[1][j][0]==using_list[k][tree_maker_list[0]]:
                sub_list.append(using_list[k])
        temp_set.add(tree_maker_list[0])
        element=TreeNodes(sub_list,temp_set)
        Tree_Nodes_List.append(element)

"""
According to the received inputs,
the tree roaming algorithm is printed as png at the same time while roaming with the algorithm.
"""
def traverse_tree(list_of_inputs):
    temp=root
    highlight_nodes = [temp]
    main_flag=True
    temp_child=root
    while main_flag:
        for k in range(len(list_of_inputs)):
            tree_flag=False
            if list_of_inputs[k][0]==temp.name:

               for l in range(len(temp.children)):
                   if temp.children[l].name==list_of_inputs[k][1]:
                       temp=temp.children[l]
                       highlight_nodes.append(temp)
                       tree_flag=True
                       break
               if not tree_flag:
                   main_flag=False
                   print("One of the data you have entered is incorrect.")
                   break
               if temp.children:
                   if len(temp.children) > 1 and not temp.children[0].children:
                       temp_child=temp.children[0]
                       backup=temp
                       for c in range(len(temp.children)):
                          temp=temp.children[c]
                          highlight_nodes.append(temp)
                          temp=backup
                       temp=temp_child
                   else:
                       temp = temp.children[0]
                       highlight_nodes.append(temp)
            if not temp.children: #IMPORTANT FACTOR:If tree includes some conflict factors,the tree can not decide which way it should go to
               answer=temp.name
               main_flag=False
               print("--------FINAL ANSWER--------")
               if temp_child is not root:print("The final result can not be determined because of inputs")
               else:print(title[len(title) - 1] + ":" + answer)
               dot = create_colored_tree(root, highlight_nodes)
               dot.render('decision_tree', format='png', cleanup=True,view=True)
               break

csv_list = []
file_input = input("Please enter file name: ")
print()

"""
The main purpose of the try expect block is to exit directly from the code 
using the sys library in the absence of the file name entered as input.
Since the delimiter is not completely clear when reading, the delimiter is determined 
by reading a certain amount at first and then using the sniffer function of the csv library, 
the delimiter is determined and the split operation is performed accordingly.
All expressions are lowered except title to make it case sensitive.
"""

try:
    with open(file_input, newline='') as csv_file:

        sample = csv_file.read(1024)
        csv_file.seek(0)
        try:
            sniffer = csv.Sniffer()
            dialect = sniffer.sniff(sample)
        except csv.Error:
            dialect = csv.get_dialect('excel')

        csv_read = csv.reader(csv_file, dialect)
        title = next(csv_read)
        count=0
        for row in csv_read:
            lower_row = [elem.lower() for elem in row]
            csv_list.append(lower_row)

except FileNotFoundError:
    print(f"Error: File '{file_input}' not found.")
    sys.exit(1)

"""
Until the while loop, the main purpose is to find the root and find the first elements to be found 
under the root and assign them to dictionary_values_list, Tree_Nodes_List and node_list.dictionary_values_list 
keeps the output probabilities of the elements in dictionary format, while Tree_Nodes_List keeps 
sub_lists and ascendants. node_list is a structure that is kept in order to be able to use the parents 
properly when using the tree structure specified in the Node library. these 3 structures work synchronously 
with each other and coincide at the same time specific to that value.
"""

output_dict={}
for i in range(len(csv_list)):
    if csv_list[i][len(title) - 1] not in output_dict:
        output_dict[csv_list[i][len(title) - 1]]=0
    output_dict[csv_list[i][len(title) - 1]]+=1
print("Finding for Root:",entropy(output_dict)," ",output_dict,"\n")

attribute_possibilities=True
attribute_keys=[]

root_list=tree_maker(csv_list,set(),output_dict)

attribute_possibilities=False

dictionary_values_list=root_list[1]
print("Highest Gain is:",title[root_list[0]], "Split based on:",title[root_list[0]], "\n")
node_list=[]
root=Node(title[root_list[0]])
for i in range(len(dictionary_values_list)):
    node=Node(dictionary_values_list[i][0],parent=root)
    node_list.append(node)

Tree_Nodes_List=[]
child_list_maker(root_list,csv_list,set())
counter=0 # using for node_list to synchronized with dictionary_values_list and Tree_Nodes_List
"""
The main purpose of the while loop is to create the rest of the loop and 
the whole tree structure after the root is found.
At first, the entropy is checked, if the entropy is 0, 
this is the pure state and the pure state means that only the leaf node is left and 
the leaf node is determined immediately and the process is finished.
If the entropy is not zero, go to the tree_maker structure and calculate the gain and find the highest gain.
To make a footnote in the structure specified here, if the value returned by tree_maker is 0. 
Index is not an integer, then either all ascendants are full or all information gain values are 0. 
This is a conflict situation and the problem is solved by adding both situations as leaf node. 
If this footnote is not encountered, go to child_list_maker and create new sub_list and ascendants. 
In order for node_list, dictionary_values_list and Tree_Nodes_List to work synchronously, 
the first elements of dictionary_values_list and Tree_Nodes_List are popped and queue structure is made. 
Node_list looks for the next element.
"""
while len(dictionary_values_list)>0:

  entropy_answer=entropy(dictionary_values_list[0][1])
  print("----- Subtree for:",dictionary_values_list[0][0],":",dictionary_values_list[0][1],"Parent:",node_list[counter].parent.name,"Entropy:",entropy_answer,"-----\n")
  if entropy_answer==0:

      key = next(iter(dictionary_values_list[0][1]))
      print("All",key,"-> Leaf Node\n")
      node = Node(key,parent=node_list[counter])
  else:
      answer_list=tree_maker(Tree_Nodes_List[0].sub_list,Tree_Nodes_List[0].ascendant,dictionary_values_list[0][1])
      if isinstance(answer_list[0],str):
          print("In this case,The Decision Tree can not decide which parameter will be divided")
          for key in dictionary_values_list[0][1].keys():
              node = Node(key, parent=node_list[counter])
      else:
          print("Highest Gain is:", title[answer_list[0]], "Split based on:", title[answer_list[0]], "\n")
          node = Node(title[answer_list[0]], parent=node_list[counter])
          dictionary_values_list.extend(answer_list[1])
          child_list_maker(answer_list, Tree_Nodes_List[0].sub_list, Tree_Nodes_List[0].ascendant)
          for item in answer_list[1]:
              node2 = Node(item[0], parent=node)
              node_list.append(node2)
  Tree_Nodes_List.pop(0)
  dictionary_values_list.pop(0)
  counter=counter+1

"""
Finally after the tree has been made,decision_tree is printed on console.
Then the user enters some inputs and 
According to inputs the traverse_tree will be determined and will be shown as png type.
"""
root.hshow()
input_flag=True
while input_flag:
    input_list = []
    for i in range(len(title) - 1):
        data = input(f"Please enter {title[i]} {attribute_keys[i]} type: ")
        data = data.lower()
        input_list.append([title[i], data])
    print("\n")
    traverse_tree(input_list)
    data=input("\nDo you want to continue? (Press y or Y for continue,Press anything else to quit): ")
    data=data.lower()
    print()
    if not data =="y":
        input_flag=False