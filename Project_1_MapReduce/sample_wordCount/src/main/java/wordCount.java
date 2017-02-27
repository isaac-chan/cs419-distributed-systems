/**
 * Created by isaac on 2/3/17.
 */

import java.io.DataOutputStream;
import java.io.IOException;
import java.util.*;

import org.apache.hadoop.conf.Configuration;
import org.apache.hadoop.fs.FSDataOutputStream;
import org.apache.hadoop.fs.FileSystem;
import org.apache.hadoop.fs.Path;
import org.apache.hadoop.io.IntWritable;
import org.apache.hadoop.io.Text;
import org.apache.hadoop.mapred.Counters;
import org.apache.hadoop.mapreduce.*;
import org.apache.hadoop.mapreduce.lib.input.FileInputFormat;
import org.apache.hadoop.mapreduce.lib.output.FileOutputFormat;

public class wordCount {

    public enum COUNTERS {
        TOTALWORDS,
        REDUCERRUNS
    };

    public static class TokenizerMapper
            extends Mapper<Object, Text, Text, IntWritable> {

        private final static IntWritable one = new IntWritable(1);
        private Text word = new Text();

        //private List<String> uniquewords = new ArrayList<String>();

        public void map(Object key, Text value, Context context
        ) throws IOException, InterruptedException {
            StringTokenizer itr = new StringTokenizer(value.toString().toLowerCase().replaceAll("[^a-zA-Z ]", ""));

            while (itr.hasMoreTokens()) {
                word.set(itr.nextToken());
                context.write(word, one);

                //uniquewords.add(word.toString());
            }
            //Set<String> uniquewordsSet = new HashSet<String>(uniquewords);

            //IntWritable num_unique = new IntWritable(uniquewordsSet.size());

            //context.write(new Text("Unique Words: "), num_unique);
            //context.write(new Text("Total Words: "), new IntWritable(((int) context.getCounter(COUNTERS.TOTALWORDS).getValue())));
        }
    }

    public static class CombinationReducer
            extends Reducer<Text, IntWritable, Text, IntWritable> {

        private IntWritable count = new IntWritable();

        public void reduce(Text combination, Iterable<IntWritable> values, Context context)
                throws IOException, InterruptedException {

            /* count the combinations */
            int sum = 0;
            for (IntWritable val : values) {
                sum += val.get();
            }
            count.set(sum);
            context.getCounter(COUNTERS.REDUCERRUNS).increment(1);
            context.getCounter(COUNTERS.TOTALWORDS).increment(sum);

            String serial = String.valueOf(context.getCounter(COUNTERS.REDUCERRUNS).getValue());
            /* write the count */

            context.write(new Text(serial + " " + combination), count);
        }
    }

    public static class FormattedFileOutputFormat
            extends FileOutputFormat<Text, IntWritable> {

        @Override
        public RecordWriter<Text, IntWritable> getRecordWriter(TaskAttemptContext arg0)
                throws IOException, InterruptedException {

            Path path = FileOutputFormat.getOutputPath(arg0);
            Path fullPath = new Path(path, "part-r-00000");
            FileSystem fs = path.getFileSystem(arg0.getConfiguration());
            FSDataOutputStream fileOut = fs.create(fullPath, arg0);

            return new FormattedRecordWriter(fileOut);
        }
    }

    public static class FormattedRecordWriter
            extends RecordWriter<Text, IntWritable> {

        private DataOutputStream out;

        public FormattedRecordWriter(DataOutputStream stream) {
            out = stream;
        }

        /* print total number of pairs and close */
        @Override
        public void close(TaskAttemptContext context)
                throws IOException, InterruptedException {
            out.writeBytes("Unique Words: " + context.getCounter(COUNTERS.REDUCERRUNS).getValue() + "\n");
            out.writeBytes("Total Words: " + context.getCounter(COUNTERS.TOTALWORDS).getValue() + "\n");
            out.close();
        }

        /* print pairs */
        @Override
        public void write(Text key, IntWritable value)
                throws IOException, InterruptedException {
            out.writeBytes(key.toString() + " " + value.get() + "\n");
        }
    }

    public static void main(String[] args) throws Exception {
        Configuration conf = new Configuration();
        Job job = Job.getInstance(conf, "word count");
        job.setJarByClass(wordCount.class);
        job.setMapperClass(TokenizerMapper.class);

        job.setReducerClass(CombinationReducer.class);
        job.setOutputKeyClass(Text.class);
        job.setOutputValueClass(IntWritable.class);
        job.setOutputFormatClass(FormattedFileOutputFormat.class);

        FileInputFormat.addInputPath(job, new Path(args[0]));
        FileOutputFormat.setOutputPath(job, new Path(args[1]));

        System.exit(job.waitForCompletion(true) ? 0 : 1);
    }
}
