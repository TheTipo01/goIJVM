package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

var (
	// The stack
	stack = make([]int16, 0)
	// Map for knowing at what line a given label is
	labelMap = make(map[string]int)
	// Map for storing variables
	variables = make(map[string]int16)
	// Instruction counter
	ic int
)

func main() {
	// Checks if file is specified
	if len(os.Args) < 2 {
		fmt.Println("No file specified! Drag and drop one on the executable!\nPress any key to exit")

		var dummy string
		_, _ = fmt.Scanf("%s", &dummy)
		return
	}

	// Read the file
	data, err := ioutil.ReadFile(os.Args[1])
	if err != nil {
		fmt.Println("Error while reading file:", err, "\nPress any key to exit")

		var dummy string
		_, _ = fmt.Scanf("%s", &dummy)
		return
	}

	// Replaces windows file endings with unix one and split the program on newline
	program := strings.Split(strings.ReplaceAll(string(data), "\r\n", "\n"), "\n")

	// Initializes labels
	initializeLabelMap(program)

	// Handle every instruction
	for ic >= 0 && ic < len(program) {
		handleInstruction(program[ic])
	}

	// Exits
	fmt.Println("Program execution terminated. Press a key to exit")
	var dummy string
	_, _ = fmt.Scanf("%s", &dummy)
}

func handleInstruction(is string) {
	// Note: INVOKEVIRTUAL, IRETURN and LDC_W are not implemented

	// WIDE is also not implemented for simplicity, as the entire stack is of 16bit
	// It is also stripped from the instruction, if present

	divided := strings.Split(strings.TrimPrefix(is, "WIDE "), " ")

	switch strings.ToUpper(divided[0]) {
	// Push a byte onto stack
	case "BIPUSH":
		// Tries to parse the number
		n, err := strconv.ParseUint(divided[1], 10, 16)
		if err == nil {
			stack = append(stack, int16(n))
		} else {
			fmt.Printf("\nError on line %d: %s", ic+1, err.Error())
		}

		ic++
		break

	// Copy top word on stack and push onto stack
	case "DUP":
		// Checks if we have at least one element in the stack
		if len(divided) > 0 {
			stack = append(stack, stack[len(stack)-1])
		} else {
			fmt.Printf("\nError on line %d: DUP needs at least one element in the stack!", ic+1)
		}

		ic++
		break

	// Unconditional jump
	case "GOTO":
		// Checks if we have more than two element in the stack
		if len(divided) > 1 {
			// Checks if label exists and jumps to it
			if _, ok := labelMap[divided[1]]; ok {
				ic = labelMap[divided[1]]
			} else {
				fmt.Printf("\nError on line %d: label doesn't exist", ic+1)
			}
		} else {
			fmt.Printf("\nError on line %d: GOTO needs a label separated by a space", ic+1)
		}

		ic++
		break

	// Halt the simulator
	case "HALT":
		ic = -1
		break

	// Pop two words from stack; push their sum
	case "IADD":
		// Checks if we have more than two element in the stack
		if len(stack) > 1 {
			// Computes sum of the top two element in the stack
			sum := stack[len(stack)-1] + stack[len(stack)-2]
			// Resizes the stack
			stack = stack[:len(stack)-2]
			// Adds the result to the stack
			stack = append(stack, sum)
		} else {
			fmt.Printf("\nError on line %d: IADD needs at least two element in the stack", ic+1)
		}

		ic++
		break

	// Pop two words from stack; push Boolean AND
	case "IAND":
		// Checks if we have more than two element in the stack
		if len(stack) > 1 {
			// Computes OR of the top two element in the stack
			and := stack[len(stack)-1] | stack[len(stack)-2]
			// Resizes the stack
			stack = stack[:len(stack)-2]
			// Adds the result to the stack
			stack = append(stack, and)
		} else {
			fmt.Printf("\nError on line %d: IAND needs at least two element in the stack", ic+1)
		}

		ic++
		break

	// Pop word from stack and branch if it is zero
	case "IFEQ":
		// Checks if we have more than one in the stack
		if len(stack) > 0 {
			// Checks the condition
			if stack[len(stack)-1] == 0 {
				// Pop the element
				stack = stack[:len(stack)-1]

				// Check if there's a label
				if len(divided) > 1 {
					// And if the label exist
					if _, ok := labelMap[divided[1]]; ok {
						ic = labelMap[divided[1]]
					} else {
						fmt.Printf("\nError on line %d: label doesn't exist", ic+1)
					}
				} else {
					fmt.Printf("\nError on line %d: IFEQ needs a label separated by a space", ic+1)
				}
			}
		} else {
			fmt.Printf("\nError on line %d: IFEQ needs at least one element in the stack", ic+1)
		}

		break

	// Pop word from stack and branch if it is less than zero
	case "IFLT":
		// Check if we have more than 1 element in the stack
		if len(stack) > 0 {
			// Checks the condition
			if stack[len(stack)-1] < 0 {
				// Pop the element
				stack = stack[:len(stack)-1]

				// Check if there's a label
				if len(divided) > 1 {
					// And if the label exist
					if _, ok := labelMap[divided[1]]; ok {
						ic = labelMap[divided[1]]
					} else {
						fmt.Printf("\nError on line %d: label doesn't exist", ic+1)
					}
				} else {
					fmt.Printf("\nError on line %d: IFLT needs a label separated by a space", ic+1)
				}
			}
		} else {
			fmt.Printf("\nError on line %d: IFLT needs at least one element in the stack", ic+1)
		}

		break

	// Pop two words from stack and branch if they are equal
	case "IF_ICMPEQ":
		// Check if we have more than 2 element in the stack
		if len(stack) > 1 {
			// Checks the condition
			if stack[len(stack)-1] == stack[len(stack)-2] {

				// Check if there's a label
				if len(divided) > 1 {
					// And if the label exist
					if _, ok := labelMap[divided[1]]; ok {
						ic = labelMap[divided[1]]
					} else {
						fmt.Printf("\nError on line %d: label doesn't exist", ic+1)
					}
				} else {
					fmt.Printf("\nError on line %d: IF_ICMPEQ needs a label separated by a space", ic+1)
				}
			}

			// Pop the two element
			stack = stack[:len(stack)-2]
		} else {
			fmt.Printf("\nError on line %d: IF_ICMPEQ needs at least one element in the stack", ic+1)
		}

		break

	// Add a constant value to a local variable
	case "IINC":
		// Checks if the variable name and byte exists
		if len(divided) > 2 {
			n, err := strconv.ParseUint(divided[1], 10, 16)
			if err == nil {
				variables[divided[1]] = int16(n)
			} else {
				fmt.Printf("\nError on line %d: %s", ic+1, err.Error())
			}
		} else {
			fmt.Printf("\nError on line %d: IINC needs a variable name and a constant, seperated by spaces", ic+1)
		}

		ic++
		break

	// Push local variable onto stack
	case "ILOAD":
		// Checks if there's the variable name
		if len(divided) > 1 {
			if _, ok := variables[divided[1]]; ok {
				stack = append(stack, variables[divided[1]])
			} else {
				fmt.Printf("\nError on line %d: variable doesn't exist", ic+1)
			}
		} else {
			fmt.Printf("\nError on line %d: ILOAD needs a variable name seperated by a space", ic+1)
		}

		ic++
		break

	// Reads a character from the keyboard buffer and pushes it onto the stack. If no character is available, 0 is pushed
	case "IN":
		// Scan a character
		var c rune
		_, _ = fmt.Scanf("%c", &c)

		// Add it to the stack
		stack = append(stack, int16(c))

		ic++
		break

	// Pop word from stack and store in local variable
	case "ISTORE":
		// Checks if we have at least one element in the stack
		if len(stack) > 0 {
			// And if there's a label specified
			if len(divided) > 1 {
				// And if the label exist
				if _, ok := variables[divided[1]]; ok {
					// Store last item in the stack as a variable
					variables[divided[1]] = stack[len(stack)-1]
					// And removes word from the stack
					stack = stack[:len(stack)-1]
				} else {
					fmt.Printf("\nError on line %d: label doesn't exist", ic+1)
				}
			} else {
				fmt.Printf("\nError on line %d: ISTORE needs a variable name seperated by a space", ic+1)
			}
		} else {
			fmt.Printf("\nError on line %d: ISTORE needs at least one element in the stack", ic+1)
		}

		ic++
		break

	// Pop two words from stack; subtract the second to top word from top word, push the answer;
	case "ISUB":
		// Checks if we have at least two element in the stack
		if len(stack) > 1 {
			sub := stack[len(stack)-1] - stack[len(stack)-2]
			stack = stack[:len(stack)-2]
			stack = append(stack, sub)
		} else {
			fmt.Printf("\nError on line %d: ISUB needs at least two element in the stack", ic+1)
		}

		ic++
		break

	// Do nothing
	case "NOP":
		ic++
		break

	// Pop word off stack and print it to standard out
	case "OUT":
		// Checks if we have at least one element in the stack
		if len(stack) > 0 {
			fmt.Println("OUT:", stack[len(stack)-1])
			stack = stack[:len(stack)-1]
		} else {
			fmt.Printf("\nError on line %d: OUT needs at least one element in the stack", ic+1)
		}

		ic++
		break

	// Delete word from top of stack
	case "POP":
		// Checks if we have at least one element in the stack
		if len(stack) > 0 {
			stack = stack[:len(stack)-1]
		} else {
			fmt.Printf("\nError on line %d: POP needs at least one element in the stack", ic+1)
		}

		ic++
		break

	// Swap the two top words on the stack
	case "SWAP":
		// Checks if we have at least two element in the stack
		if len(stack) > 1 {
			// Swaps
			stack[len(stack)-1], stack[len(stack)-2] = stack[len(stack)-2], stack[len(stack)-1]
		} else {
			fmt.Printf("\nError on line %d: SWAP needs at least two element in the stack", ic+1)
		}

		ic++
		break

	// Prints status of the stack and the variables
	case "DEBUG":
		if len(stack) > 0 {
			fmt.Printf("\n\nStack:\n")
			for i, s := range stack {
				fmt.Printf("%d: %d\n", i, s)
			}

			fmt.Printf("\n\n")
		}

		if len(variables) > 0 {
			fmt.Printf("\n\nVariables:\n")
			for k, v := range variables {
				fmt.Printf("%s: %d\n", k, v)
			}

			fmt.Printf("\n\n")
		}

		ic++
		break

	default:
		fmt.Printf("\nWarning on line %d: unknown instruction", ic+1)
	}
}

// Checks for labels in the program
func initializeLabelMap(program []string) {
	for i, line := range program {
		line = strings.TrimSpace(line)

		if strings.HasSuffix(line, ":") {
			labelMap[strings.TrimSuffix(line, ".")] = i
		}
	}
}
