#include <stdio.h>
#include <stdlib.h>

#include "common.h"
#include "compiler.h"
#include "debug.h"
#include "vm.h"

VM vm;

static void reset_stack()
{
    vm.stack_top = vm.stack;
}

void init_vm()
{
    reset_stack();
}

void free_vm()
{
}

void push(Value value)
{
    // The book didn't check for this... perhaps he'll do the check somewhere
    // else.
    if (vm.stack_top - vm.stack >= STACK_MAX) {
        printf("Stack overflow");
        exit(1);
    }

    *vm.stack_top = value;
    vm.stack_top++;
}

Value pop()
{
    // Just like the push, the book didn't really check for bad situations...
    // I'm guessing he'll do something in a few chapters, hopefully!
    if (vm.stack_top == &vm.stack[0]) {
        printf("Trying to pop from an empty stack!");
        exit(1);
    }

    vm.stack_top--;
    return *vm.stack_top;
}

static InterpretResult run()
{
#define READ_BYTE() (*vm.ip++)
#define READ_CONSTANT() (vm.chunk->constants.values[READ_BYTE()])
#define BINARY_OP(op)     \
    do {                  \
        double b = pop(); \
        double a = pop(); \
        push(a op b);     \
    } while (false)

    for (;;)
    {
#ifdef DEBUG_TRACE_EXECUTION
        printf("          ");
        for (Value *slot = vm.stack; slot < vm.stack_top; slot++)
        {
            printf("[ ");
            print_value(*slot);
            printf(" ]");
        }
        printf("\n");
        disassemble_instruction(vm.chunk, (int)(vm.ip - vm.chunk->code));
#endif
        uint8_t instruction;
        switch (instruction = READ_BYTE()) {
        case OP_CONSTANT: {
            Value constant = READ_CONSTANT();
            push(constant);
        } break;
        case OP_ADD: {
            BINARY_OP(+);
        } break;
        case OP_SUBTRACT: {
            BINARY_OP(-);
        } break;
        case OP_MULTIPLY: {
            BINARY_OP(*);
        } break;
        case OP_DIVIDE: {
            BINARY_OP(/);
        } break;
        case OP_NEGATE: {
            push(-pop());
        } break;
        case OP_RETURN: {
            print_value(pop());
            printf("\n");
            return INTERPRET_OK;
        }
        }
    }

#undef READ_BYTE
#undef READ_CONSTANT
#undef BINARY_OP
}

InterpretResult interpret(const char *source)
{
    compile(source);
    return INTERPRET_OK;
}
