import java.io.IOException;
import java.io.DataOutputStream;

import java.util.regex.Pattern;
import java.util.Arrays;

import org.apache.hadoop.conf.Configuration;
import org.apache.hadoop.io.IntWritable;
import org.apache.hadoop.io.Text;
import org.apache.hadoop.fs.Path;
import org.apache.hadoop.mapreduce.Job;
import org.apache.hadoop.mapreduce.Mapper;
import org.apache.hadoop.mapreduce.Reducer;
import org.apache.hadoop.mapreduce.lib.input.TextInputFormat;
import org.apache.hadoop.mapreduce.lib.output.FileOutputFormat;
import org.apache.hadoop.mapreduce.TaskAttemptContext;
import org.apache.hadoop.mapreduce.RecordWriter;
import org.apache.hadoop.fs.FileSystem;
import org.apache.hadoop.fs.FSDataOutputStream;

public class BoughtTogether {

    /* count the total number of pairs */
    public enum COUNTERS {
        TOTALPAIRS
    };

    /* generate combinations of items that were bought together */
    public static class LineMapper
        extends Mapper<Object, Text, Text, IntWritable> {

        private final static Pattern itemSeparator = Pattern.compile(", *");
        private final static IntWritable one = new IntWritable(1);

        public void map(Object key, Text value, Context context)
                throws IOException, InterruptedException {

            /* input lines of case-sensitive items separated by commas */
            String line = value.toString();

            /* split comma-separated items */
            String[] items = itemSeparator.split(line);

            /* sort the items to prevent duplicate pairs (x,y) and (y,x) */
            Arrays.sort(items);

            /* output combinations */
            for (int i=0; i<items.length; ++i) {
                String first = items[i];
                for (int j=i+1; j<items.length; ++j) {
                    String second = items[j];
                    Text combination = new Text("(" + first + ", " + second + ")");
                    context.write(combination, one);
                }
            }
        }
    }

    /* count the combinations */
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
            context.getCounter(COUNTERS.TOTALPAIRS).increment(1);

            /* write the count */
            context.write(new Text(combination), count);
        }
    }

    /* custom formatted file output */
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

    /* custom formatted record writer */
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
            out.writeBytes("Total Pairs: " + context.getCounter(COUNTERS.TOTALPAIRS).getValue() + "\n");
            out.close();
        }

        /* print pairs */
        @Override
        public void write(Text key, IntWritable value)
                throws IOException, InterruptedException {
            out.writeBytes(key.toString() + " " + value.get() + "\n");
        }
    }

    /* entry point */
    public static void main(String[] args)
            throws Exception {

        Configuration configuration = new Configuration();

        Job job = Job.getInstance(configuration, "bought together");
        job.setJarByClass(BoughtTogether.class);
        job.setMapperClass(LineMapper.class);
        job.setReducerClass(CombinationReducer.class);
        job.setOutputKeyClass(Text.class);
        job.setOutputValueClass(IntWritable.class);
        job.setOutputFormatClass(FormattedFileOutputFormat.class);

        TextInputFormat.addInputPath(job, new Path(args[0]));
        FileOutputFormat.setOutputPath(job, new Path(args[1]));

        System.exit(job.waitForCompletion(true) ? 0 : 1);
    }
}
