# MapReduce Project

Isaac Chan and Taylor Alexander Brown

CS 419 Distributed Systems

Oregon State University

Winter 2017

## Configuration

Configure Hadoop as a [single node cluster](http://hadoop.apache.org/docs/current/hadoop-project-dist/hadoop-common/SingleCluster.html) with pseudo-distributed operation.

Export path to Hadoop path for use in run scripts:

    $ export HADOOP_DIR=/path/to/hadoop-2.7.3

Format the Hadoop Distributed File System:

    $ $HADOOP_DIR/bin/hdfs namenode -format

## wordCount

Given an input file with words, sort the words and ouput the number of times the word occurs in the file.

### Usage

    $ cd sample_wordCount/
    $ export PROJECT_DIR=`pwd`
    $ ./run.sh

### Example Output

    $ cat output/part-r-00000
    1 cow     1
    2 jumps   1
    3 moon    1
    4 over    1
    5 the     2
    Unique Words: 5
    Total WOrds: 6

## BoughtTogether

Given input files containing lists of rewards card transactions, compute how many times each pair of items are bought together.

### Usage

    $ cd BoughtTogether/
    $ export PROJECT_DIR=`pwd`
    $ ./run.sh

### Example Output

    $ cat output/part-r-00000
    (Apples, Bananas) 1
    (Apples, BeavMoo Milk) 1
    (Apples, Best Bread) 1
    (Bananas, BeavMoo Milk) 1
    (Bananas, Best Bread) 1
    (BeavMoo Milk, Best Bread) 2
    (BeavMoo Milk, Fluffy Pizza) 1
    (BeavMoo Milk, Whitey Toothpaste) 1
    (Best Bread, Fluffy Pizza) 1
    (Best Bread, Whitey Toothpaste) 1
    (Fluffy Pizza, Whitey Toothpaste) 1
    Total Pairs: 12
