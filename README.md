# MapReduce Project

Isaac Chan and Taylor Alexander Brown

CS 419 Distributed Systems

Oregon State University

Winter 2017

## Configuration

Configure Hadoop as a [single node cluster](http://hadoop.apache.org/docs/current/hadoop-project-dist/hadoop-common/SingleCluster.html).

## wordCount

## BoughtTogether

Given input files containing lists of rewards card transactions, compute how many times each pair of items are bought together.

### Usage

    $ cd BoughtTogether/
    $ export PROJECT_DIR=`pwd`
    $ export HADOOP_DIR=/path/to/hadoop-2.7.3
    $ $HADOOP_DIR/bin/hdfs namenode -format
    $ ./run.sh

### Example Output

    $ cat output/part-r-00000
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
