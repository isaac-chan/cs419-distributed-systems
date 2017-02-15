# MapReduce Project

Isaac Chan and Taylor Alexander Brown

CS 419 Distributed Systems

Oregon State University

Winter 2017

## wordCount

Given an input file with words, sort the words and ouput the number of times the word occurs in the file.

### Usage
Run using: cs419-distributed-systems/sample_wordCount/run.sh

### Example Output
    Total Words:   6
    Unique Words:  5
    cow    1
    jumps  1
    moon   1
    over   1
    the    2

## BoughtTogether

Given input files containing lists of rewards card transactions, compute how many times each pair of items are bought together.

### Usage

    $ cd BoughtTogether/
    $ hadoop com.sun.tools.javac.Main BoughtTogether.java
    $ jar cf bt.jar BoughtTogether*.class
    $ hadoop jar bt.jar BoughtTogether bought-together-input/ bought-together-output/

### Example Output

    $ less bought-together-output/part-r-00000
    Apples, Bananas 1
    Apples, BeavMoo Milk    1
    Apples, Best Bread      1
    Bananas, BeavMoo Milk   1
    Bananas, Best Bread     1
    BeavMoo Milk, Best Bread        2
    BeavMoo Milk, Fluffy Pizza      1
    BeavMoo Milk, Whitey Toothpaste 1
    Best Bread, Fluffy Pizza        1
    Best Bread, Whitey Toothpaste   1
    Fluffy Pizza, Whitey Toothpaste 1
