import java.io.IOException;

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

public class BoughtTogether {

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

            /* write the count */
            context.write(new Text(combination), count);
        }
    }

    /* entry point */
    public static void main(String[] args)
        throws Exception {

        Configuration configuration = new Configuration();

        Job job = Job.getInstance(configuration, "bought together");
        job.setJarByClass(BoughtTogether.class);
        job.setMapperClass(LineMapper.class);
        job.setCombinerClass(CombinationReducer.class);
        job.setReducerClass(CombinationReducer.class);
        job.setOutputKeyClass(Text.class);
        job.setOutputValueClass(IntWritable.class);

        TextInputFormat.addInputPath(job, new Path(args[0]));
        FileOutputFormat.setOutputPath(job, new Path(args[1]));

        System.exit(job.waitForCompletion(true) ? 0 : 1);
    }
}
