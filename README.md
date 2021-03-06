# calc
A command line calculator.

## Installation
```bash
git clone https://github.com/pointlander/calc.git
cd calc
go install
```

## Usage
```bash
calc
```

## Language
```
e <- sp e1 !.
e1 <- e2 ( add e2
         / minus e2
         )*
e2 <- e3 ( multiply e3
         / divide e3
         / modulus e3
         )*
e3 <- e4 ( exponentiation e4
         )*
e4 <- minus value
    / value
value <- matrix
       / imaginary
       / number
       / exp1
       / exp2
       / natural
       / pi
       / prec
       / simplify
       / derivative
       / log
       / sqrt
       / cos
       / sin
       / tan
       / variable
       / sub
variable <- [A-Za-z]+ sp
matrix <- '[' sp (e1 / row)+ ']' sp
imaginary <- < [-]? [0-9]+ ([.] [0-9]*)? > 'i' sp
number <- < [-]? [0-9]+ ([.] [0-9]*)? > sp
exp1 <- 'exp' open e1 close
exp2 <- 'e^' value
natural <- 'e' sp
pi <- 'pi' sp
prec <- 'prec' open e1 close
simplify <- 'simplify' open e1 close
derivative <- 'derivative' open e1 close
log <- 'log' open e1 close
sqrt <- 'sqrt' open e1 close
cos <- 'cos' open e1 close
sin <- 'sin' open e1 close
tan <- 'tan' open e1 close
sub <- open e1 close
add <- '+' sp
minus <- '-' sp
multiply <- '*' sp
divide <- '/' sp
modulus <- '%' sp
exponentiation <- '^' sp
open <- '(' sp
close <- ')' sp
sp <- ( ' ' / '\t' )*
row <- ';' sp
```
