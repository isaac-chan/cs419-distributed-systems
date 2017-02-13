#hadoop installation dir
HADOOP_DIR=/home/isaac/Downloads/hadoop-2.7.3

#project location dir
PROJECT_DIR=/home/isaac/IdeaProjects/cs4191

#remove previous output
rm -f -d $PROJECT_DIR/output/ --recursive

#start node
$HADOOP_DIR/sbin/start-all.sh

#load input text file into HDFS
$HADOOP_DIR/bin/hadoop fs -put $PROJECT_DIR/input/input.txt /input.txt

#run word count JAR
$HADOOP_DIR/bin/hadoop jar /$PROJECT_DIR/cs419-1-1.0-SNAPSHOT.jar wordCount /input.txt /output/

#get output from HDFS
$HADOOP_DIR/bin/hadoop fs -get /output/ $PROJECT_DIR/

#print out the output
cat $PROJECT_DIR/output/part-r-00000

#clean HDFS
$HADOOP_DIR/bin/hadoop fs -rm -r -skipTrash /input.txt
$HADOOP_DIR/bin/hadoop fs -rm -r -skipTrash /output

#stop node
$HADOOP_DIR/sbin/stop-all.sh
