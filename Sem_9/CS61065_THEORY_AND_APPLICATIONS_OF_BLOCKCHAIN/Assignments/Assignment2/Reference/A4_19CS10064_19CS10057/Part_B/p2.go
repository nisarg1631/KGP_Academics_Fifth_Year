//  Assignment No. - 4 Part - B
//  Hyperledger Fabric
//  CS61065 - Theory and Applications of Blockchain
//  Semester - 7 (Autumn 2022-23)
//  Group Members - Vanshita Garg (19CS10064) & Shristi Singh (19CS10057)

// Import dependencies and define smart contract
package main

import (
    "fmt"
    "log"
    "encoding/json"
    "github.com/hyperledger/fabric-contract-api-go/contractapi"
	"math"
)

type SmartContract struct {
    contractapi.Contract
}

// struct treenode with json tags
type TreeNode struct {
	Val int `json:"val"`
	Left *TreeNode `json:"left"`
	Right *TreeNode `json:"right"`
}	

// struct Binary search tree with json tags
type MyBST struct {
	PrimaryKey string `json:"primaryKey"`
	Root *TreeNode `json:"root"`
}

// go function for inserting a value in the tree
func (s *SmartContract) Insert(ctx contractapi.TransactionContextInterface, val int) error {
	// check if any BST already exists in the ledger using ReadMyBST()
	tree, err := ReadMyBST(ctx)
	// if err != nil, this means that no BST exists in the ledger
	// if no BST exists in the ledger, create a new BST and insert the value in the tree
	if err != nil {
		tree = &MyBST{PrimaryKey: "MyBST", Root: &TreeNode{Val: val}}
		treeBytes, err := json.Marshal(tree)
		if err != nil {
			return err
		}
		err = ctx.GetStub().PutState(tree.PrimaryKey, treeBytes)
		if err != nil {
			return err
		}
		return nil
	}
	// if yes, insert the value in the tree using UpdateMyBST()
	if tree != nil {
		// use UpdateMyBST() to insert the value in the tree
		err = UpdateMyBST(ctx, val, tree, 0)
		if err != nil {
			return err
		}
	}
	return nil
}

// go function Delete(ctx contractapi.TransactionContextInterface, val int) error
func (s *SmartContract) Delete(ctx contractapi.TransactionContextInterface, val int) error {
	// check if any BST already exists in the ledger using ReadMyBST()
	tree, err := ReadMyBST(ctx)
	if err != nil {
		return err
	}
	// if yes, delete the value from the tree using UpdateMyBST()
	if tree != nil {
		// update the tree
		err = UpdateMyBST(ctx, val, tree, 1)
		if err != nil {
			return err
		}	
	}
	return nil
}

// go function Preorder(ctx contractapi.TransactionContextInterface) (string, error)
func (s *SmartContract) Preorder(ctx contractapi.TransactionContextInterface) (string, error) {
	// check if any BST already exists in the ledger using ReadMyBST()
	tree, err := ReadMyBST(ctx)
	if err != nil {
		return "", err
	}
	// if yes, return the preorder traversal of the tree using preorder()
	if tree != nil {
		// get the preorder traversal of the tree using preorder()
		preorderTraversal := preorder(tree.Root)
		// convert the preorder traversal to a comma separated string and return it
		return fmt.Sprint(preorderTraversal), nil
	}
	return "", nil
}

// go fucntion Inorder(ctx contractapi.TransactionContextInterface) (string, error)
func (s *SmartContract) Inorder(ctx contractapi.TransactionContextInterface) (string, error) {
	// check if any BST already exists in the ledger using ReadMyBST()
	tree, err := ReadMyBST(ctx)
	if err != nil {
		return "", err
	}
	// if yes, return the inorder traversal of the tree using inorder()
	if tree != nil {
		// get the inorder traversal of the tree using inorder()
		inorderTraversal := inorder(tree.Root)
		// convert the inorder traversal to a comma separated string and return it
		return fmt.Sprint(inorderTraversal), nil
	}
	return "", nil
}

// go function TreeHeight(ctx contractapi.TransactionContextInterface) (string, error)
func (s *SmartContract) TreeHeight(ctx contractapi.TransactionContextInterface) (string, error) {
	// check if any BST already exists in the ledger using ReadMyBST()
	tree, err := ReadMyBST(ctx)
	if err != nil {
		return "0", err
	}
	// if yes, return the height of the tree using treeHeight()
	if tree != nil {
		// get the height of the tree using treeHeight()
		height := heightOfTree(tree.Root)
		// convert the height to a string and return it
		return fmt.Sprint(height), nil
	}
	return "0", nil
}


// go function InsertValue() to insert a new value in the tree as per the rules of BST.
func (tree *MyBST) InsertValue(val int) error {
	if tree.Root == nil {
		tree.Root = &TreeNode{Val: val}
		return nil
	}
	// do not insert if the value is already present in the tree
	if tree.SearchValue(val) {
		return fmt.Errorf("Value already present in the tree")
	}
	curr := tree.Root
	for curr != nil {
		if val < curr.Val {
			if curr.Left == nil {
				curr.Left = &TreeNode{Val: val}
				return nil
			}
			curr = curr.Left
		} else {
			if curr.Right == nil {
				curr.Right = &TreeNode{Val: val}
				return nil 
			}
			curr = curr.Right
		}
	}
	return nil
}

// go function SearchValue() to search a value in the tree.
func (tree *MyBST) SearchValue(val int) bool {
	curr := tree.Root
	for curr != nil {
		if val == curr.Val {
			return true
		} else if val < curr.Val {
			curr = curr.Left
		} else {
			curr = curr.Right
		}
	}
	return false
}

// write a DeleteValue() function that deletes a value from the tree input as *MyBST
func (tree *MyBST) DeleteValue(val int) error {
	// if the tree is empty, return an error
	if tree.Root == nil {
		return fmt.Errorf("Tree is empty")
	}
	// if the value to be deleted is not present in the tree, return an error
	if !tree.SearchValue(val) {
		return fmt.Errorf("Value not present in the tree")
	}
	// if the value to be deleted is the root node, call the deleteRootNode() function
	if tree.Root.Val == val {
		tree.Root = deleteRootNode(tree.Root)
		return nil
	}
	// if the value to be deleted is not the root node, call the deleteNode() function
	tree.Root = deleteNode(tree.Root, val)
	return nil
}

// write a deleteRootNode() function that deletes the root node of the tree input as *TreeNode
func deleteRootNode(root *TreeNode) *TreeNode {
	// if the root node has no children, return nil
	if root.Left == nil && root.Right == nil {
		return nil
	}
	// if the root node has only one child, return the child
	if root.Left == nil {
		return root.Right
	}
	if root.Right == nil {
		return root.Left
	}
	// if the root node has two children, find the inorder successor of the root node
	// and replace the root node with the inorder successor
	inorderSuccessor := findInorderSuccessor(root.Right)
	root.Val = inorderSuccessor.Val
	root.Right = deleteNode(root.Right, inorderSuccessor.Val)
	return root
}

// findInorderSuccessor() function to find the inorder successor of a node input as *TreeNode
func findInorderSuccessor(root *TreeNode) *TreeNode {
	for root.Left != nil {
		root = root.Left
	}
	return root
}

// go function deleteNode() to delete a node from the tree.
func deleteNode(root *TreeNode, val int) *TreeNode {
	if root == nil {
		return nil
	}
	if val < root.Val {
		root.Left = deleteNode(root.Left, val)
	} else if val > root.Val {
		root.Right = deleteNode(root.Right, val)
	} else {
		if root.Left == nil {
			return root.Right
		} else if root.Right == nil {
			return root.Left
		}
	}
	return root
}

// go function for preoder traversal of the tree
func preorder(root *TreeNode) []int {
	if root == nil {
		return []int{}
	}
	return append(append([]int{root.Val}, preorder(root.Left)...), preorder(root.Right)...)
}

// go function for inorder traversal of the tree
func inorder(root *TreeNode) []int {
	if root == nil {
		return []int{}
	}
	return append(append(inorder(root.Left), root.Val), inorder(root.Right)...)
}

// go function heightOfTree() to calculate the height of the tree
func heightOfTree(root *TreeNode) int {
	if root == nil {
		return 0	
	}
	ret := 1 + int(math.Max(float64(heightOfTree(root.Left)), float64(heightOfTree(root.Right))))
	return ret
}


// go function for UpdateMyBST(ctx contractapi.TransactionContextInterface, val int, bst *MyBST, operation int) error
func UpdateMyBST(ctx contractapi.TransactionContextInterface, val int, bst *MyBST, operation int) error {
	if operation == 0 {
		err := bst.InsertValue(val)
		if err != nil {
			return err
		}
	} else if operation == 1 {
		err := bst.DeleteValue(val)
		if err != nil {
			return err
		}
	} else {
		return fmt.Errorf("invalid operation")
	}
	treeBytes, err := json.Marshal(bst)
	if err != nil {
		return err
	}
	// Delete the old tree from the state
	err = ctx.GetStub().DelState(bst.PrimaryKey)
	if err != nil {
		return err
	}
	// Put the updated tree in the state if the tree is not empty
	if bst.Root != nil {
		err = ctx.GetStub().PutState(bst.PrimaryKey, treeBytes)
		if err != nil {	
			return err
		}
	}
	return nil
}

// go function for ReadMyBST(ctx contractapi.TransactionContextInterface) (*MyBST, error)
func ReadMyBST(ctx contractapi.TransactionContextInterface) (*MyBST, error) {
	// retrieve the tree from the ledger using the ctx.GetStub.GetStateByRange() function
	iterator, err := ctx.GetStub().GetStateByRange("", "")
	// read the first key-value pair from the iterator
	result, err := iterator.Next()
	// read the first BST from the iterato
	if err != nil {
		// iterator.Next() returns an error if there are no more key-value pairs to read
		return nil, fmt.Errorf("no tree found")
	}
	tree := MyBST{}
	err = json.Unmarshal(result.Value, &tree)
	if err != nil {
		return nil, err
	}
	return &tree, nil
}

// go function MyBSTExists(ctx contractapi.TransactionContextInterface, key string) (bool, error)
func MyBSTExists(ctx contractapi.TransactionContextInterface, key string) (bool, error) {
	// return true if the BST entry exists in the ledger
	treeBytes, err := ctx.GetStub().GetState(key)
	if err != nil {
		return false, err
	}
	return treeBytes != nil, nil
}

// function main
func main() {
    // create a new smart contract
    smartContract, err := contractapi.NewChaincode(new(SmartContract))
    if err != nil {
        log.Panicf("Error creating smart contract: %v", err)
    }
    // start the smart contract
    if err := smartContract.Start(); err != nil {
        log.Panicf("Error starting smart contract: %v", err)
    }
}
