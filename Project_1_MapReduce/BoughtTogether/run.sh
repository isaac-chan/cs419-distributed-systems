#hadoop installation dir
if [[ -z "$HADOOP_DIR" ]] ; then
    HADOOP_DIR=/home/isaac/Downloads/hadoop-2.7.3
fi

#project location dir
if [[ -z "$PROJECT_DIR" ]] ; then
    PROJECT_DIR=/home/isaac/cs419-distributed-systems/Project_1_MapReduce/BoughtTogether
fi

#build JAR
cd $PROJECT_DIR;mvn clean package

#remove previous output
rm -f -d $PROJECT_DIR/output/ --recursive

#start node
$HADOOP_DIR/sbin/start-dfs.sh

#load input text file into HDFS
$HADOOP_DIR/bin/hadoop fs -put $PROJECT_DIR/input/P2t2.txt /file.txt

#run word count JAR
$HADOOP_DIR/bin/hadoop jar $PROJECT_DIR/target/cs419-1-1.0-SNAPSHOT.jar BoughtTogether /file.txt /output/

#get output from HDFS
$HADOOP_DIR/bin/hadoop fs -get /output/ $PROJECT_DIR/

#print out the output
cat $PROJECT_DIR/output/part-r-00000

#clean HDFS
$HADOOP_DIR/bin/hadoop fs -rm -r -skipTrash /file.txt
$HADOOP_DIR/bin/hadoop fs -rm -r -skipTrash /output

#stop node
$HADOOP_DIR/sbin/stop-dfs.sh
